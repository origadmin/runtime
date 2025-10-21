package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"

	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
)

// mergeAppInfo handles the logic of merging application information from configuration
// and options, and then validates the final result.
func mergeAppInfo(sc interfaces.StructuredConfig, providerOptions *ProviderOptions) (*interfaces.AppInfo, error) {
	// 1. Decode AppInfo from the configuration.
	configAppInfo, err := sc.DecodeApp()
	// It's okay if this fails (e.g., 'app' key not in config), so we don't return the error immediately.
	// A hard error would prevent using WithAppInfo as the only source.
	if err != nil {
		log.Debugf("failed to decode app info from config, will rely on WithAppInfo option: %v", err)
	}

	// 2. Merge AppInfo from options (as base) and config (as override).
	// Start with the AppInfo provided via the WithAppInfo option. It can be nil.
	appInfo := providerOptions.appInfo
	if appInfo == nil {
		// If no AppInfo was provided via options, create a new one to populate from config.
		appInfo = &interfaces.AppInfo{}
	}

	// Merge values from the config. Config values take precedence.
	if configAppInfo != nil {
		if configAppInfo.Id != "" {
			appInfo.ID = configAppInfo.Id
		}
		if configAppInfo.Name != "" {
			appInfo.Name = configAppInfo.Name
		}
		if configAppInfo.Version != "" {
			appInfo.Version = configAppInfo.Version
		}
		if configAppInfo.Env != "" {
			appInfo.Env = configAppInfo.Env
		}
		if len(configAppInfo.Metadata) > 0 {
			appInfo.Metadata = configAppInfo.Metadata
		}
	}

	// 3. Validate the final merged AppInfo.
	if appInfo.ID == "" || appInfo.Name == "" || appInfo.Version == "" {
		return nil, runtimeerrors.NewStructured("bootstrap", "app info (ID, Name, Version) is required but was not found in config or WithAppInfo option").WithCaller()
	}

	return appInfo, nil
}
