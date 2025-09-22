package bootstrap

import (
	"fmt"
	"sync"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/goexts/generic/configure"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/internal/decoder"

	"github.com/origadmin/runtime/log"                      // Re-import runtime/log
	runtimeRegistry "github.com/origadmin/runtime/registry" // Placeholder for actual registry package
)

// DefaultBootstrapPath the default bootstrap configuration file path
const DefaultBootstrapPath = "configs/bootstrap.toml"

// ComponentFactoryFunc defines the signature for a function that creates a component.
// It receives the full ConfigDecoder and the specific component's configuration map.
// The returned interface{} should be the actual component instance.
type ComponentFactoryFunc func(decoder interfaces.ConfigDecoder, config map[string]any) (interface{}, error)

var (
	componentFactories     = make(map[string]ComponentFactoryFunc)
	componentFactoriesLock sync.RWMutex
)

// RegisterComponentFactory registers a ComponentFactoryFunc for a given component type.
// This allows for dynamic creation of components based on configuration.
// It is safe for concurrent use.
func RegisterComponentFactory(componentType string, factory ComponentFactoryFunc) {
	componentFactoriesLock.Lock()
	defer componentFactoriesLock.Unlock()
	if _, exists := componentFactories[componentType]; exists {
		kratoslog.Warnf("Component factory for type '%s' already registered. Overwriting.", componentType)
	}
	componentFactories[componentType] = factory
}

// NewProvider is the one-stop-shop function that loads configuration, creates all components,
// and returns them via the ComponentProvider interface, along with a cleanup function.
func NewProvider(bootstrapPath string, options ...Option) (interfaces.ComponentProvider, func(), error) {
	// 1. REUSE: Call NewDecoder to get the decoder and cleanup function.
	decoder, cleanup, err := NewDecoderWithOptions(bootstrapPath, options...)
	if err != nil {
		return nil, nil, err // Decoder creation failed, nothing to clean up yet.
	}

	// --- Logger (with graceful fallback) ---
	var logger kratoslog.Logger // Type remains kratoslog.Logger
	loggerCfg, err := decoder.DecodeLogger()
	// Initialize a temporary helper for early logging, before the main logger is fully configured.
	// This helper will be replaced by the fully configured one later.
	var tempHelper *kratoslog.Helper = kratoslog.NewHelper(kratoslog.DefaultLogger)

	if err != nil {
		// Log a warning and create a default logger if config is missing or invalid
		tempHelper.Warnf("Failed to decode logger config, using default: %v", err)
		logger = log.NewLogger(nil) // Use runtime/log.NewLogger for default
	} else {
		logger = log.NewLogger(loggerCfg) // Use runtime/log.NewLogger for configured
		// log.NewLogger does not return an error, so no need to check err here.
	}

	// Now that the main logger is determined, create the helper for consistent logging.
	helper := kratoslog.NewHelper(logger)

	// --- Discoveries & Registrars (with graceful fallback) ---
	registrars := make(map[string]registry.Registrar)
	discoveries := make(map[string]registry.Discovery)
	var defaultRegistrar registry.Registrar

	var registriesCfg struct {
		DefaultRegistry string                            `json:"default" yaml:"default"`
		Registries      map[string]*discoveryv1.Discovery `json:"registries" yaml:"registries"`
	}
	// CORRECTED: Decode the specific "registries" section, not the entire config root.
	err = decoder.Decode("registries", &registriesCfg)
	if err != nil {
		// Distinguish between config not found and parsing error if possible, for now, treat as not found for graceful fallback.
		helper.Infof("No service registry configuration found or failed to decode (%v). Running in local mode.", err)
		// Continue with empty maps and nil defaultRegistrar
	} else {
		for name, registryCfg := range registriesCfg.Registries {
			if registryCfg == nil || registryCfg.GetType() == "" || registryCfg.GetType() == "none" {
				helper.Infof("Skipping registry '%s' due to missing or 'none' type.", name)
				continue
			}

			helper.Infof("Initializing service registry and discovery '%s' with type: %s", name, registryCfg.GetType())
			r, err := runtimeRegistry.NewRegistrar(registryCfg) // Assuming runtimeRegistry.NewRegistrar returns Kratos type
			if err != nil {
				cleanup()
				return nil, nil, fmt.Errorf("failed to create registrar for '%s': %w", name, err)
			}
			d, err := runtimeRegistry.NewDiscovery(registryCfg) // Assuming runtimeRegistry.NewDiscovery returns Kratos type
			if err != nil {
				cleanup()
				return nil, nil, fmt.Errorf("failed to create discovery for '%s': %w", name, err)
			}
			registrars[name] = r
			discoveries[name] = d


			if name == registriesCfg.DefaultRegistry {
				defaultRegistrar = r
			}
		}

		// Handle case where default registry is specified but not found
		if registriesCfg.DefaultRegistry != "" && defaultRegistrar == nil {
			cleanup()
			return nil, nil, fmt.Errorf("default registry '%s' not found in configured registries", registriesCfg.DefaultRegistry)
		}
	}

	// --- Dynamic Components (Scan and create from [components] section) ---
	components := make(map[string]interface{})
	components["logger"] = logger
	components["discoveries"] = discoveries
	components["registrars"] = registrars
	if defaultRegistrar != nil {
		components["defaultRegistrar"] = defaultRegistrar
	}

	var rawGenericComponents map[string]map[string]any
	err = decoder.Decode("components", &rawGenericComponents)
	if err != nil {
		helper.Infof("No generic components configuration found or failed to decode (%v). Skipping dynamic component creation.", err)
	} else {
		componentFactoriesLock.RLock()
		defer componentFactoriesLock.RUnlock()

		for compName, compConfig := range rawGenericComponents {
			compType, ok := compConfig["type"].(string)
			if !ok || compType == "" {
				helper.Warnf("Component '%s' has no 'type' field or it's not a string. Skipping.", compName)
				continue
			}

			factory, found := componentFactories[compType]
			if !found {
				helper.Warnf("No factory registered for component type '%s' (component '%s'). Skipping.", compType, compName)
				continue
			}

			helper.Infof("Creating dynamic component '%s' of type '%s'.", compName, compType)
			compInstance, factoryErr := factory(decoder, compConfig)
			if factoryErr != nil {
				helper.Errorf("Failed to create component '%s' of type '%s': %v. Skipping.", compName, compType, factoryErr)
				continue
			}
			components[compName] = compInstance
		}
	}

	// 3. ASSEMBLE: Create the concrete provider struct.
	provider := &componentProviderImpl{
		logger:           logger,
		discoveries:      discoveries,
		registrars:       registrars,
		defaultRegistrar: defaultRegistrar,
		components:       components,
	}

	// 4. RETURN: Return the provider and the original cleanup function.
	return provider, cleanup, nil
}

// NewDecoder initializes the configuration and returns a ready-to-use ConfigDecoder.
// This is a convenience wrapper around NewDecoderWithOptions, allowing direct passing of runtimeconfig.Option.
func NewDecoder(bootstrapPath string, configOptions ...runtimeconfig.Option) (interfaces.ConfigDecoder, func(), error) {
	return NewDecoderWithOptions(bootstrapPath, WithConfigOptions(configOptions...))
}

// NewDecoderWithOptions is the core function that loads all configuration sources,
// creates a ConfigDecoder using the specified or default provider, and returns it.
func NewDecoderWithOptions(bootstrapPath string, options ...Option) (interfaces.ConfigDecoder, func(), error) {
	// 1. Apply options to get DecoderProvider, etc.
	opts := configure.Apply(&Options{
		// Set the default decoder provider here.
		DecoderProvider: decoder.DefaultDecoderProvider,
	}, options)

	// 2. Load the bootstrap file to find out about other config sources.
	bootstrapKratosOptions := []kratosconfig.Option{
		kratosconfig.WithSource(file.NewSource(bootstrapPath)),
	}
	if opts.BootstrapOptions != nil {
		bootstrapKratosOptions = append(bootstrapKratosOptions, opts.BootstrapOptions...)
	}

	bootstrapConf := kratosconfig.New(bootstrapKratosOptions...)
	if err := bootstrapConf.Load(); err != nil {
		return nil, nil, fmt.Errorf("failed to load bootstrap config file %s: %w", bootstrapPath, err)
	}
	defer bootstrapConf.Close() // Close the bootstrap config source after we're done with it.

	// 3. Unmarshal the bootstrap config into a `sourcev1.Sources` struct.
	var sourcesDef sourcev1.Sources
	if err := bootstrapConf.Scan(&sourcesDef); err != nil {
		return nil, nil, fmt.Errorf("failed to scan bootstrap config: %w", err)
	}

	// 4. Validate the sources definition.
	if err := validateSources(&sourcesDef); err != nil {
		return nil, nil, fmt.Errorf("invalid sources definition: %w", err)
	}

	// 5. Create the final merged config from all sources.
	finalConfig, err := runtimeconfig.NewConfig(&sourcesDef, opts.ConfigOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create final config from sources: %w", err)
	}

	if err := finalConfig.Load(); err != nil {
		// If loading fails, we must close the config source we just created.
		_ = finalConfig.Close()
		return nil, nil, fmt.Errorf("failed to load final config: %w", err)
	}

	// 6. Use the decoder provider to create the ConfigDecoder.
	configDecoder, err := opts.DecoderProvider.GetConfigDecoder(finalConfig)
	if err != nil {
		_ = finalConfig.Close()
		return nil, nil, fmt.Errorf("failed to get config decoder: %w", err)
	}

	// 7. Create the cleanup function that closes the final config.
	cleanup := func() {
		if err := finalConfig.Close(); err != nil {
			// This error typically occurs during shutdown, so we print it
			// as a warning instead of causing a panic.
			fmt.Printf("[runtime/bootstrap] warning: failed to close config source: %v\n", err)
		}
	}

	return configDecoder, cleanup, nil
}

// validateSources validates the effectiveness of configuration source definitions
func validateSources(sources *sourcev1.Sources) error {
	if sources == nil {
		return fmt.Errorf("sources cannot be nil")
	}
	if len(sources.Sources) == 0 {
		return fmt.Errorf("no configuration sources defined")
	}

	// Check required fields for each configuration source
	for i, source := range sources.Sources {
		if source == nil {
			return fmt.Errorf("source #%d cannot be nil", i)
		}
		if source.Name == "" {
			return fmt.Errorf("source #%d must have a name", i)
		}
		if source.Type == "" {
			return fmt.Errorf("source #%d must have a type", i)
		}
	}

	return nil
}
