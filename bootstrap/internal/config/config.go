package config

import (
	"errors"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	appv1 "github.com/origadmin/runtime/api/gen/go/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	"github.com/origadmin/runtime/interfaces"
)

// structuredConfigImpl implements the interfaces.StructuredConfig interface.
// It wraps a generic interfaces.Config and provides type-safe, path-based decoding methods.
type structuredConfigImpl struct {
	interfaces.Config // Embed the generic config interface
	paths             map[string]string
}

// Statically assert that structuredConfigImpl implements the full StructuredConfig interface.
var _ interfaces.StructuredConfig = (*structuredConfigImpl)(nil)

// NewStructured creates a new structured config implementation.
// It takes a generic interfaces.Config and a path map to provide high-level decoding methods.
func NewStructured(cfg interfaces.Config, paths map[string]string) interfaces.StructuredConfig {
	if paths == nil {
		paths = make(map[string]string)
	}
	return &structuredConfigImpl{
		Config: cfg,
		paths:  paths,
	}
}

// decodeConfig implements the robust decoding logic for a component.
// It correctly distinguishes between a user-provided path, a user-disabled path, and a fallback path.
// It also correctly handles "key not found" errors as non-fatal.
func (c *structuredConfigImpl) decodeConfig(key, fallbackKey string, value any) error {
	path, ok := c.paths[key]

	// Case 1: The path is explicitly defined in the map.
	if ok {
		// Case 1a: The user has explicitly disabled this component by setting the path to "".
		if path == "" {
			return nil // Success, but do nothing.
		}
		// Case 1b: The user has provided a specific path. We must use it.
		// The `path` variable is already correctly set for this case.
	} else {
		// Case 2: The path is not in the map. This is the only case where we fall back.
		path = fallbackKey
	}

	// Attempt to decode using the determined path.
	err := c.Decode(path, value)
	if err != nil {
		// If the error is specifically a "not found" error from the underlying config source,
		// it means the config is optional and not present. This is not a fatal error.
		if errors.Is(err, kratosconfig.ErrNotFound) {
			return nil // Return nil error to indicate optional config is missing.
		}
		// For any other error (e.g., parsing error), return it as it's a real problem.
		return err
	}
	return nil
}

// DecodeApp implements the AppConfigDecoder interface.
func (c *structuredConfigImpl) DecodeApp() (*appv1.App, error) {
	appConfig := new(appv1.App)
	if err := c.decodeConfig(constant.ConfigApp, "app", appConfig); err != nil {
		return nil, err
	}
	// If the struct is still zero-valued, it means the key was not found or disabled.
	if appConfig.Id == "" && appConfig.Name == "" {
		return nil, nil
	}
	return appConfig, nil
}

// DecodeLogger implements the LoggerConfigDecoder interface.
func (c *structuredConfigImpl) DecodeLogger() (*loggerv1.Logger, error) {
	loggerConfig := new(loggerv1.Logger)
	if err := c.decodeConfig(constant.ComponentLogger, "logger", loggerConfig); err != nil {
		return nil, err
	}
	if loggerConfig.Name == "" && len(loggerConfig.Level) == 0 {
		return nil, nil
	}
	return loggerConfig, nil
}

// DecodeDiscoveries implements the DiscoveriesConfigDecoder interface.
func (c *structuredConfigImpl) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	var discoveries map[string]*discoveryv1.Discovery
	if err := c.decodeConfig(constant.ComponentRegistries, "discoveries", &discoveries); err != nil {
		return nil, err
	}
	if len(discoveries) == 0 {
		return nil, nil
	}
	return discoveries, nil
}

// DecodeMiddleware implements the MiddlewareConfigDecoder interface.
func (c *structuredConfigImpl) DecodeMiddleware() (*middlewarev1.Middlewares, error) {
	var middlewares *middlewarev1.Middlewares
	if err := c.decodeConfig(constant.ComponentMiddlewares, "middlewares", &middlewares); err != nil {
		return nil, err
	}
	if middlewares == nil || len(middlewares.Middlewares) == 0 {
		return nil, nil
	}
	return middlewares, nil
}
