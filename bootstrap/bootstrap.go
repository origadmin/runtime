package bootstrap

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/toolkits/errors"
)

// New creates a new component provider, which is the main entry point for application startup.
// It orchestrates the entire process of configuration loading and component initialization.
// It now returns the interfaces.Bootstrapper interface.
func New(bootstrapPath string, opts ...Option) (interfaces.Bootstrapper, error) {
	// 1. Apply options to get access to WithAppInfo.
	// Assuming options have been flattened as per our discussion.
	providerOpts := FromOptions(opts...)

	// 2. Load configuration first. This is a critical change in the flow.
	cfg, err := LoadConfig(bootstrapPath, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// The cleanup function will be built up, starting with closing the config.
	cleanup := func() {
		if err := cfg.Close(); err != nil {
			log.Errorf("failed to close config: %v", err)
		}
	}

	// 3. Decode AppInfo from the configuration.
	configAppInfo, err := cfg.DecodeApp()
	// It's okay if this fails (e.g., 'app' key not in config), so we don't return the error immediately.
	// A hard error would prevent using WithAppInfo as the only source.
	if err != nil {
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
	if finalAppInfo.StartTime.IsZero() {
		finalAppInfo.StartTime = time.Now()
	}

	if finalAppInfo.ID == "" || finalAppInfo.Name == "" || finalAppInfo.Version == "" {
		cleanup() // Call cleanup as we are failing.
		return nil, errors.New("app info (ID, Name, Version) is required but was not found in config or WithAppInfo option")
	}

	//// 3. Create the component provider implementation.
	//// This will hold all the initialized components.
	//p := container.NewContainer(cfg)
	//
	//// 4. Initialize core components by consuming the config.
	//// This is where the magic happens: logger, registries, etc., are created.
	//if err := p.Initialize(cfg); err != nil {
	//	// Even if initialization fails, we should still call the cleanup function.
	//	cleanup()
	//	return nil, fmt.Errorf("failed to initialize components: %w", err)
	//}
	//
	//// 5. Initialize user-defined components registered via WithComponent.
	//for key, factory := range providerOpts.componentFactories {
	//	instance, err := factory(cfg, p)
	//	if err != nil {
	//		cleanup()
	//		return nil, fmt.Errorf("failed to create component '%s' using factory: %w", key, err)
	//	}
	//	p.RegisterComponent(key, instance)
	//}

	// 7. Return the container, the config, and the final cleanup function.
	return &bootstrapperImpl{
		appInfo: finalAppInfo,
		config:  cfg,
		cleanup: cleanup,
	}, nil
}
