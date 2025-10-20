package config

import (
	"errors"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
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

// decodeComponent implements a simple and robust decoding logic.
// It no longer contains any fallback logic. It trusts the `paths` map provided by the bootstrap package.
func (c *structuredConfigImpl) decodeComponent(componentKey string, value any) error {
	path, ok := c.paths[componentKey]

	// If the key is not in the paths map, or the path is explicitly empty, it's considered disabled or not configured.
	if !ok || path == "" {
		return nil // This is not an error.
	}

	// Attempt to decode using the provided path.
	err := c.Decode(path, value)
	if err != nil {
		// If the error is specifically a "not found" error, it's not a fatal issue.
		if errors.Is(err, kratosconfig.ErrNotFound) {
			return nil
		}
		// Any other error (e.g., parsing) is a real problem.
		return err
	}
	return nil
}

// DecodeApp implements the AppConfigDecoder interface.
func (c *structuredConfigImpl) DecodeApp() (*appv1.App, error) {
	appConfig := new(appv1.App)
	if err := c.decodeComponent(constant.ConfigApp, appConfig); err != nil {
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
	if err := c.decodeComponent(constant.ComponentLogger, loggerConfig); err != nil {
		return nil, err
	}
	if loggerConfig.Name == "" && len(loggerConfig.Level) == 0 {
		return nil, nil
	}
	return loggerConfig, nil
}

// DecodeDiscoveries implements the DiscoveriesConfigDecoder interface.
func (c *structuredConfigImpl) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	var list []*discoveryv1.Discovery
	var m map[string]*discoveryv1.Discovery

	// Try to decode as a map first.
	if err := c.decodeComponent(constant.ComponentRegistries, &m); err == nil && len(m) > 0 {
		for name, item := range m {
			if item.Name == "" {
				item.Name = name
			}
			list = append(list, item)
		}
		return &discoveryv1.Discoveries{Discoveries: list}, nil
	}

	// If map fails or is empty, try to decode as a list.
	if err := c.decodeComponent(constant.ComponentRegistries, &list); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return &discoveryv1.Discoveries{Discoveries: list}, nil
}

// DecodeMiddlewares implements the MiddlewareConfigDecoder interface.
func (c *structuredConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	val := new(middlewarev1.Middlewares)
	if err := c.decodeComponent(constant.ComponentMiddlewares, val); err != nil {
		return nil, err
	}
	if val == nil || len(val.Middlewares) == 0 {
		return nil, nil
	}
	return val, nil
}

// DecodeServers implements the ServiceConfigDecoder interface.
// It intelligently handles both map and slice formats for server configurations.
func (c *structuredConfigImpl) DecodeServers() (*transportv1.Servers, error) {
	var list []*transportv1.Server
	var m map[string]*transportv1.Server

	// Try to decode as a map first.
	if err := c.decodeComponent(constant.ComponentServers, &m); err == nil && len(m) > 0 {
		for name, item := range m {
			if item.Name == "" {
				item.Name = name
			}
			list = append(list, item)
		}
		return &transportv1.Servers{Servers: list}, nil
	}

	// If map fails or is empty, try to decode as a list.
	if err := c.decodeComponent(constant.ComponentServers, &list); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return &transportv1.Servers{Servers: list}, nil
}

// DecodeClients implements the ServiceConfigDecoder interface.
// It intelligently handles both map and slice formats for client configurations.
func (c *structuredConfigImpl) DecodeClients() (*transportv1.Clients, error) {
	var list []*transportv1.Client
	var m map[string]*transportv1.Client

	// Try to decode as a map first.
	if err := c.decodeComponent(constant.ComponentClients, &m); err == nil && len(m) > 0 {
		for name, item := range m {
			if item.Name == "" {
				item.Name = name
			}
			list = append(list, item)
		}
		return &transportv1.Clients{Clients: list}, nil
	}

	// If map fails or is empty, try to decode as a list.
	if err := c.decodeComponent(constant.ComponentClients, &list); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return &transportv1.Clients{Clients: list}, nil
}
