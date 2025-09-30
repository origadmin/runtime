package bootstrap

import (
	"fmt"
	"time"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goexts/generic/configure"

	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/bootstrap/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	bootstrapConfig "github.com/origadmin/runtime/bootstrap/internal/config"
	"github.com/origadmin/runtime/bootstrap/internal/container"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/toolkits/errors"
)

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

// bootstrapperImpl implements the interfaces.Bootstrapper interface.
type bootstrapperImpl struct {
	provider interfaces.Container
	config   interfaces.Config
	cleanup  func()
}

// Provider implements interfaces.Bootstrapper.
func (b *bootstrapperImpl) Provider() interfaces.Container {
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

// New creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading and component initialization.
// It now returns the interfaces.Bootstrapper interface.
func New(bootstrapPath string, opts ...Option) (interfaces.Bootstrapper, error) {
	// 1. Apply provider-level Options
	providerOpts := configure.Apply(&Options{}, opts)

	// AppInfo is a mandatory input for creating a valid provider.
	// Check if appInfo is nil OR if it's not valid (e.g., empty ID, Name, Version).
	baseAppInfo := providerOpts.appInfo // This is from WithAppInfo() option

	// Provide a default AppInfo if none was given via options
	if baseAppInfo == nil {
		baseAppInfo = &interfaces.AppInfo{}
	}
	// Set StartTime if not already set (it's a runtime value, not from config)
	if baseAppInfo.StartTime.IsZero() {
		baseAppInfo.StartTime = time.Now()
	}

	// 2. Create the configuration decoder, which returns the interfaces.Config.
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
	// 3. Decode AppInfo from the configuration and merge/overwrite
	// 4. Merge AppInfo from config using the type-safe decoder interface
	if decoder, ok := cfg.(interfaces.AppConfigDecoder); ok {
		if configAppInfo, err := decoder.DecodeApp(); err == nil && configAppInfo != nil {
			if configAppInfo.Id != "" {
				baseAppInfo.ID = configAppInfo.Id
			}
			if configAppInfo.Name != "" {
				baseAppInfo.Name = configAppInfo.Name
			}
			if configAppInfo.Version != "" {
				baseAppInfo.Version = configAppInfo.Version
			}
			if configAppInfo.Env != "" {
				baseAppInfo.Env = configAppInfo.Env
			}
			if len(configAppInfo.Metadata) > 0 {
				baseAppInfo.Metadata = configAppInfo.Metadata
			}
		}
	}
	if baseAppInfo.ID == "" || baseAppInfo.Name == "" || baseAppInfo.Version == "" {
		return nil, errors.New("app info (ID, Name, Version) is required and must be valid after merging with config")
	}

	// 3. Create the component provider implementation.
	// This will hold all the initialized components.
	p := container.NewContainer(baseAppInfo, cfg)

	// 4. Initialize core components by consuming the config.
	// This is where the magic happens: logger, registries, etc., are created.
	if err := p.Initialize(cfg); err != nil {
		// Even if initialization fails, we should still call the cleanup function.
		cleanup()
		return nil, fmt.Errorf("failed to initialize components: %w", err)
	}

	// 5. Initialize user-defined components registered via WithComponent.
	for key, factory := range providerOpts.componentFactories {
		instance, err := factory(cfg, p)
		if err != nil {
			cleanup()
			return nil, fmt.Errorf("failed to create component '%s' using factory: %w", key, err)
		}
		p.RegisterComponent(key, instance)
	}

	// 7. Return the provider, the config, and the final cleanup function.

	return &bootstrapperImpl{
		provider: p,
		config:   cfg,
		cleanup:  cleanup,
	}, nil
}
