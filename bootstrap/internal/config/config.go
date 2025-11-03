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
	"github.com/origadmin/toolkits/decode"
)

// structuredConfigImpl implements the interfaces.StructuredConfig interface.
// It wraps a generic interfaces.Config and provides type-safe, path-based decoding methods.
type structuredConfigImpl struct {
	interfaces.Config // Embed the generic config interface
	paths             map[string]string
}

// Statically assert that structuredConfigImpl implements the full StructuredConfig interface.
var _ interfaces.StructuredConfig = (*structuredConfigImpl)(nil)

// --- Reusable Converters ---

var (
	serverConverter = decode.NewConverter(
		func(name string, item *transportv1.Server) (*transportv1.Server, bool) {
			if item.Name == "" {
				item.Name = name
			}
			return item, true
		},
		func(items []*transportv1.Server) *transportv1.Servers {
			if items == nil {
				return nil
			}
			return &transportv1.Servers{Configs: items}
		},
	)

	clientConverter = decode.NewConverter(
		func(name string, item *transportv1.Client) (*transportv1.Client, bool) {
			if item.Name == "" {
				item.Name = name
			}
			return item, true
		},
		func(items []*transportv1.Client) *transportv1.Clients {
			if items == nil {
				return nil
			}
			return &transportv1.Clients{Configs: items}
		},
	)

	discoveryConverter = decode.NewConverter(
		func(name string, item *discoveryv1.Discovery) (*discoveryv1.Discovery, bool) {
			if item.Name == "" {
				item.Name = name
			}
			return item, true
		},
		func(items []*discoveryv1.Discovery) *discoveryv1.Discoveries {
			if items == nil {
				return nil
			}
			return &discoveryv1.Discoveries{Configs: items}
		},
	)

	middlewareConverter = decode.NewConverter(
		func(name string, item *middlewarev1.Middleware) (*middlewarev1.Middleware, bool) {
			if item.Name == "" {
				item.Name = name
			}
			return item, true
		},
		func(items []*middlewarev1.Middleware) *middlewarev1.Middlewares {
			if items == nil {
				return nil
			}
			return &middlewarev1.Middlewares{Configs: items}
		},
	)
)

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
	var m map[string]*discoveryv1.Discovery
	if err := c.decodeComponent(constant.ComponentRegistries, &m); err == nil && len(m) > 0 {
		return discoveryConverter.FromMap(m), nil
	}

	var ds []*discoveryv1.Discovery
	if err := c.decodeComponent(constant.ComponentRegistries, &ds); err == nil {
		return discoveryConverter.FromSlice(ds), nil
	}

	var s discoveryv1.Discoveries
	if err := c.decodeComponent(constant.ComponentRegistries, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *structuredConfigImpl) DecodeDefaultDiscovery() (string, error) {
	defaultRegistry := ""
	if err := c.decodeComponent(constant.ComponentDefaultRegistry, &defaultRegistry); err != nil {
		return "", err
	}
	if defaultRegistry == "" {
		return "", nil
	}
	return defaultRegistry, nil
}

// DecodeMiddlewares implements the MiddlewareConfigDecoder interface.
func (c *structuredConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	var m map[string]*middlewarev1.Middleware
	if err := c.decodeComponent(constant.ComponentMiddlewares, &m); err == nil && len(m) > 0 {
		return middlewareConverter.FromMap(m), nil
	}

	var ms []*middlewarev1.Middleware
	if err := c.decodeComponent(constant.ComponentMiddlewares, &ms); err == nil {
		return middlewareConverter.FromSlice(ms), nil
	}

	var s middlewarev1.Middlewares
	if err := c.decodeComponent(constant.ComponentMiddlewares, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// DecodeServers implements the ServiceConfigDecoder interface.
func (c *structuredConfigImpl) DecodeServers() (*transportv1.Servers, error) {
	var m map[string]*transportv1.Server
	if err := c.decodeComponent(constant.ComponentServers, &m); err == nil && len(m) > 0 {
		return serverConverter.FromMap(m), nil
	}

	var ss []*transportv1.Server
	if err := c.decodeComponent(constant.ComponentServers, &ss); err == nil {
		return serverConverter.FromSlice(ss), nil
	}

	var s transportv1.Servers
	if err := c.decodeComponent(constant.ComponentServers, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// DecodeClients implements the ServiceConfigDecoder interface.
func (c *structuredConfigImpl) DecodeClients() (*transportv1.Clients, error) {
	var m map[string]*transportv1.Client
	if err := c.decodeComponent(constant.ComponentClients, &m); err == nil && len(m) > 0 {
		return clientConverter.FromMap(m), nil
	}
	var cs []*transportv1.Client
	if err := c.decodeComponent(constant.ComponentClients, &cs); err == nil {
		return clientConverter.FromSlice(cs), nil
	}

	var s transportv1.Clients
	if err := c.decodeComponent(constant.ComponentClients, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
