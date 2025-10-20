package bootstrap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/bootstrap/internal/container"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	// "github.com/origadmin/toolkits/errors" // REMOVED: imported and not used
)

// New creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading and component initialization.
// It now returns the interfaces.Bootstrapper interface.
func New(bootstrapPath string, opts ...Option) (interfaces.Bootstrapper, error) {
	// 1. Apply options to get access to WithAppInfo.
	// Assuming options have been flattened as per our discussion.
	providerOpts := FromOptions(opts...)

	// 2. Load configuration first. This is a critical change in the flow.
	cfg, sc, err := LoadConfig(bootstrapPath, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	cleanup := func() {
		if cfg != nil {
			cfg.Close()
		}
	}

	// 3. Decode AppInfo from the configuration.
	configAppInfo, err := sc.DecodeApp()
	// It's okay if this fails (e.g., 'app' key not in config), so we don't return the error immediately.
	// A hard error would prevent using WithAppInfo as the only source.
	if err != nil {
		cleanup()
		log.Debugf("failed to decode app info from config, will rely on WithAppInfo option: %v", err)
	}

	// 4. Merge AppInfo from options (as base) and config (as override).
	// Start with the AppInfo provided via the WithAppInfo option. It can be nil.
	finalAppInfo := providerOpts.appInfo
	if finalAppInfo == nil {
		// If no AppInfo was provided via options, create a new one to populate from config.
		finalAppInfo = &interfaces.AppInfo{}
	}

	// Merge values from the config. Config values take precedence.
	if configAppInfo != nil {
		if configAppInfo.Id != "" {
			finalAppInfo.ID = configAppInfo.Id
		}
		if configAppInfo.Name != "" {
			finalAppInfo.Name = configAppInfo.Name
		}
		if configAppInfo.Version != "" {
			finalAppInfo.Version = configAppInfo.Version
		}
		if configAppInfo.Env != "" {
			finalAppInfo.Env = configAppInfo.Env
		}
		if len(configAppInfo.Metadata) > 0 {
			finalAppInfo.Metadata = configAppInfo.Metadata
		}
	}

	// 5. Set runtime values and validate.
	if finalAppInfo.ID == "" || finalAppInfo.Name == "" || finalAppInfo.Version == "" {
		cleanup()
		return nil, runtimeerrors.NewStructured("bootstrap", "app info (ID, Name, Version) is required but was not found in config or WithAppInfo option").WithCaller()
	}

	// 3. Create the component provider implementation.
	// This will hold all the initialized components.
	builder := container.NewBuilder(providerOpts.componentFactories).WithConfig(sc)

	// 4. Initialize core components by consuming the config.
	// This is where the magic happens: logger, registries, etc., are created.
	c, err := builder.Build()
	if err != nil {
		cleanup()
		// Even if initialization fails, we should still call the cleanup function.
		return nil, fmt.Errorf("failed to initialize components: %w", err)
	}

	// 7. Return the container, the config, and the final cleanup function.
	return &bootstrapperImpl{
		config:    cfg,
		appInfo:   finalAppInfo,
		container: c,
		cleanup:   cleanup,
	}, nil
}
