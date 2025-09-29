package bootstrap

import (
	"errors"
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goexts/generic/configure"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/bootstrap/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	bootstrapConfig "github.com/origadmin/runtime/bootstrap/internal/config"
	"github.com/origadmin/runtime/bootstrap/internal/provider"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/interfaces"
)

// componentFactoryRegistryImpl implements interfaces.ComponentFactoryRegistry.
type componentFactoryRegistryImpl struct{}

// GetFactory retrieves a component factory by its type.
func (r *componentFactoryRegistryImpl) GetFactory(componentType string) (interfaces.ComponentFactoryFunc, bool) {
	return getFactory(componentType)
}

// LoadBootstrapConfig loads the bootstrapv1.Bootstrap definition from a local bootstrap configuration file.
// This function is the first step in the configuration process.
// The temporary config used to load the sources is closed internally.
func LoadBootstrapConfig(bootstrapPath string) (*bootstrapv1.Bootstrap, error) {
	// Create a temporary Kratos config instance to load the bootstrap.yaml file.
	// We assume runtimeconfig.NewFileSource exists and creates a file source.
	bootConfig := kratosconfig.New(
		kratosconfig.WithSource(file.NewSource(bootstrapPath)),
	)
	// Defer closing the config and handle its error
	defer func() {
		if err := bootConfig.Close(); err != nil {
			// Log the error, as we can't return it from a deferred function
			log.Errorf("failed to close temporary bootstrap config: %v", err)
		}
	}()

	// CRITICAL FIX: Load the config after creating it
	if err := bootConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load temporary bootstrap config: %w", err)
	}

	var bc bootstrapv1.Bootstrap
	if err := bootConfig.Scan(&bc); err != nil { // Reverted to direct Scan
		return nil, fmt.Errorf("failed to scan bootstrap config from %s: %w", bootstrapPath, err)
	}

	return &bc, nil
}

// NewDecoder creates a new configuration decoder instance.
// It orchestrates the entire configuration decoding process, including path resolution and source merging.
// The returned interfaces.Config is ready to be consumed by NewProvider or other tools.
func NewDecoder(bootstrapPath string, opts ...DecoderOption) (interfaces.Config, error) {
	// 1. Apply options
	decoderOpts := &decoderOptions{}
	for _, o := range opts {
		o(decoderOpts)
	}

	// If a custom config is provided, use it directly.
	if decoderOpts.customConfig != nil {
		return decoderOpts.customConfig, nil
	}

	// If a direct Kratos config is provided, use it to create the interfaces.Config.
	if decoderOpts.kratosConfig != nil {
		// If a transformer is provided, use it to transform the Kratos config.
		if decoderOpts.configTransformer != nil {
			cfg, err := decoderOpts.configTransformer.Transform(decoderOpts.kratosConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to create interfaces.Config from kratosConfig using transformer: %w", err)
			}
			return cfg, nil
		}
		// Otherwise, use the default implementation.
		return bootstrapConfig.NewConfigImpl(decoderOpts.kratosConfig, nil), nil // Pass nil for paths as they are not relevant here
	}

	// 2. Load BootstrapConfig from file
	bootstrapCfg, err := LoadBootstrapConfig(bootstrapPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load bootstrap config: %w", err)
	}

	// 3. Merge paths with a clear priority:
	//    1. DefaultComponentPaths (lowest)
	//    2. WithDefaultPaths option
	//    3. Paths from bootstrap.yaml (highest)

	// Start with a copy of the public default map
	finalPaths := make(map[string]string)
	for component, path := range constant.DefaultComponentPaths {
		finalPaths[component] = path
	}

	// Apply paths from WithDefaultPaths option
	if decoderOpts.defaultPaths != nil {
		for component, path := range decoderOpts.defaultPaths {
			finalPaths[component] = path
		}
	}

	// 4. Create the final Kratos config.Config from all sources
	var sources *sourcev1.Sources
	if bootstrapCfg != nil {
		sources = &sourcev1.Sources{Sources: bootstrapCfg.GetSources()}
	}
	// If no sources are found, it means bootstrap.yaml was empty or invalid, but we can still proceed
	// if a config factory is provided to handle the empty Kratos config.
	if sources == nil {
		sources = &sourcev1.Sources{}
	}

	finalKratosConfig, err := runtimeconfig.NewConfig(sources, decoderOpts.configOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create final kratos config: %w", err)
	}

	// CRITICAL FIX: Load the configuration after creating it
	if err := finalKratosConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load final kratos config: %w", err)
	}

	// 5. Create the interfaces.Config implementation with the final merged paths
	// If a config transformer is provided, use it to transform the finalKratosConfig.
	if decoderOpts.configTransformer != nil {
		cfg, err := decoderOpts.configTransformer.Transform(finalKratosConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create interfaces.Config from finalKratosConfig using transformer: %w", err)
		}
		return cfg, nil
	}

	// Otherwise, use the default implementation.
	// We assume bootstrapConfig.NewConfigImpl exists in the internal package.
	decoder := bootstrapConfig.NewConfigImpl(finalKratosConfig, finalPaths)

	return decoder, nil
}

// bootstrapperImpl implements the interfaces.Bootstrapper interface.
type bootstrapperImpl struct {
	provider interfaces.ComponentProvider
	config   interfaces.Config
	cleanup  func()
}

// Provider implements interfaces.Bootstrapper.
func (b *bootstrapperImpl) Provider() interfaces.ComponentProvider {
	return b.provider
}

// Config implements interfaces.Bootstrapper.
func (b *bootstrapperImpl) Config() interfaces.Config {
	return b.config
}

// Cleanup implements interfaces.Bootstrapper.
func (b *bootstrapperImpl) Cleanup() func() {
	return b.cleanup
}

// NewProvider creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading and component initialization.
// It now returns the interfaces.Bootstrapper interface.
func NewProvider(bootstrapPath string, opts ...Option) (interfaces.Bootstrapper, error) {
	// 1. Apply provider-level options
	providerOpts := configure.Apply(&options{}, opts)

	// AppInfo is a mandatory input for creating a valid provider.
	// Check if appInfo is nil OR if it's not valid (e.g., empty ID, Name, Version).
	appInfo := providerOpts.appInfo
	if appInfo.ID == "" || appInfo.Name == "" || appInfo.Version == "" {
		return nil, errors.New("app info is required and must be valid")
	}

	// 2. Create the configuration decoder, passing through any decoder options.
	cfg, err := NewDecoder(bootstrapPath, providerOpts.decoderOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create config decoder: %w", err)
	}

	// The cleanup function will be built up, starting with closing the config.
	cleanup := func() {
		if err := cfg.Close(); err != nil {
			// Use a global logger to report cleanup errors.
			// This is a best-effort action.
			log.Errorf("failed to close config: %v", err)
		}
	}

	// 3. Create the component provider implementation.
	// This will hold all the initialized components.
	componentFactoryRegistry := &componentFactoryRegistryImpl{}
	p := provider.NewComponentProvider(providerOpts.appInfo, cfg, componentFactoryRegistry)

	// 4. Initialize core components by consuming the config.
	// This is where the magic happens: logger, registries, etc., are created.
	if err := p.InitComponents(cfg); err != nil {
		// Even if initialization fails, we should still call the cleanup function.
		cleanup()
		return nil, fmt.Errorf("failed to initialize components: %w", err)
	}

	// 5. Initialize user-defined components registered via WithComponent.
	for _, comp := range providerOpts.componentsToConfigure {
		if err := cfg.Decode(comp.Key, comp.Target); err != nil {
			cleanup()
			return nil, fmt.Errorf("failed to decode component '%s': %w", comp.Key, err)
		}
		// Register the populated struct as a component.
		p.RegisterComponent(comp.Key, comp.Target)
	}

	// 7. Return the provider, the config, and the final cleanup function.
	return &bootstrapperImpl{
		provider: p,
		config:   cfg,
		cleanup:  cleanup,
	}, nil
}
