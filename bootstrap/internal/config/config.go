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

// DecodeApp implements the AppConfigDecoder interface.
func (c *structuredConfigImpl) DecodeApp() (*appv1.App, error) {
	path, ok := c.paths[constant.ConfigApp]
	if !ok {
		// If no path is specified, try decoding from the root.
		path = "app"
	}
	appConfig := new(appv1.App)
	if err := c.Decode(path, appConfig); err != nil {
		return nil, err
	}
	return appConfig, nil
}

// DecodeLogger implements the LoggerConfigDecoder interface.
func (c *structuredConfigImpl) DecodeLogger() (*loggerv1.Logger, error) {
	path, ok := c.paths[constant.ComponentLogger]
	if !ok {
		return nil, errors.New("logger component path not configured")
	}

	loggerConfig := new(loggerv1.Logger)
	if err := c.Decode(path, loggerConfig); err != nil {
		return nil, err
	}
	return loggerConfig, nil
}

// DecodeDiscoveries implements the DiscoveriesConfigDecoder interface.
func (c *structuredConfigImpl) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	path, ok := c.paths[constant.ComponentRegistries]
	if !ok {
		return nil, errors.New("registries component path not configured")
	}

	var discoveries map[string]*discoveryv1.Discovery
	if err := c.Decode(path, &discoveries); err != nil {
		return nil, err
	}

	return discoveries, nil
}

// DecodeMiddleware implements the MiddlewareConfigDecoder interface.
func (c *structuredConfigImpl) DecodeMiddleware() (*middlewarev1.Middlewares, error) {
	path, ok := c.paths[constant.ComponentMiddlewares]
	if !ok {
		return nil, errors.New("middlewares component path not configured")
	}

	var middlewares *middlewarev1.Middlewares
	if err := c.Decode(path, &middlewares); err != nil {
		return nil, err
	}

	return middlewares, nil
}
