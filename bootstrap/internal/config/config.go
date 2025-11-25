package config

import (
	"errors"
	"sync"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/constant"
	"github.com/origadmin/toolkits/decode"
)

// structuredConfigImpl implements the interfaces.StructuredConfig interface.
// It wraps a generic interfaces.Config and provides type-safe, path-based decoding methods.
type structuredConfigImpl struct {
	interfaces.Config // Embed the generic config interface
	paths             map[constant.ComponentKey]string
	cache             sync.Map // Cache for decoded configurations
}

// DecodeCaches decodes and returns the Caches configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeCaches() (*datav1.Caches, error) {
	val, err := c.decodeComponent(constant.ComponentCaches, &datav1.Caches{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Caches), nil
}

// DecodeDatabases decodes and returns the Databases configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeDatabases() (*datav1.Databases, error) {
	val, err := c.decodeComponent(constant.ComponentDatabases, &datav1.Databases{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Databases), nil
}

// DecodeObjectStores decodes and returns the ObjectStores configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeObjectStores() (*datav1.ObjectStores, error) {
	val, err := c.decodeComponent(constant.ComponentObjectStores, &datav1.ObjectStores{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.ObjectStores), nil
}

// DecodedConfig returns the underlying generic configuration.
func (c *structuredConfigImpl) DecodedConfig() any {
	return c.Config
}

// Statically assert that structuredConfigImpl implements the full StructuredConfig interface.
var _ interfaces.StructuredConfig = (*structuredConfigImpl)(nil)

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
func NewStructured(cfg interfaces.Config, paths map[constant.ComponentKey]string) interfaces.StructuredConfig {
	if paths == nil {
		paths = make(map[constant.ComponentKey]string)
	}
	return &structuredConfigImpl{
		Config: cfg,
		paths:  paths,
		cache:  sync.Map{}, // Initialize the cache
	}
}

// decodeComponent implements a robust decoding logic with caching.
// It first checks the `paths` map for a pre-discovered path. If not found,
// it falls back to using the componentKey directly as the path.
// This provides both flexibility (via the paths map) and convention-based simplicity.
// If the component is already in the cache, it returns the cached value.
// Otherwise, it decodes the component, stores it in the cache, and then returns it.
func (c *structuredConfigImpl) decodeComponent(componentKey constant.ComponentKey, valuePtr any) (any, error) {
	// Check cache first
	if cachedVal, ok := c.cache.Load(componentKey); ok {
		return cachedVal, nil
	}

	path, ok := c.paths[componentKey]

	// If the key is not in the paths map, fall back to using the component key itself as the path.
	// This supports convention-over-configuration.
	if !ok {
		path = string(componentKey)
	} else if path == "" {
		// If the path is explicitly set to empty, it's considered disabled.
		return nil, nil
	}

	// Attempt to decode using the provided path.
	err := c.Decode(path, valuePtr)
	if err != nil {
		// If the error is specifically a "not found" error, it's not a fatal issue.
		if errors.Is(err, kratosconfig.ErrNotFound) {
			c.cache.Store(componentKey, nil) // Cache nil for not found to avoid repeated lookups
			return nil, nil
		}
		// Any other error (e.g., parsing) is a real problem.
		return nil, err
	}

	// Store the decoded value in the cache
	c.cache.Store(componentKey, valuePtr)
	return valuePtr, nil
}

// DecodeData decodes and returns the Data configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeData() (*datav1.Data, error) {
	val, err := c.decodeComponent(constant.ComponentData, &datav1.Data{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Data), nil
}

// DecodeApp decodes and returns the App configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeApp() (*appv1.App, error) {
	val, err := c.decodeComponent(constant.ConfigApp, &appv1.App{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	appConfig := val.(*appv1.App)
	// If the struct is still zero-valued, it means the key was not found or disabled.
	if appConfig.Id == "" && appConfig.Name == "" {
		return nil, nil
	}
	return appConfig, nil
}

// DecodeLogger decodes and returns the Logger configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeLogger() (*loggerv1.Logger, error) {
	val, err := c.decodeComponent(constant.ComponentLogger, &loggerv1.Logger{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	loggerConfig := val.(*loggerv1.Logger)
	if loggerConfig.Name == "" && len(loggerConfig.Level) == 0 {
		return nil, nil
	}
	return loggerConfig, nil
}

// DecodeDiscoveries decodes and returns the Discoveries configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	// Attempt to decode as map
	var m map[string]*discoveryv1.Discovery
	val, err := c.decodeComponent(constant.ComponentRegistries, &m)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedMap, ok := val.(*map[string]*discoveryv1.Discovery); ok && len(*decodedMap) > 0 {
			return discoveryConverter.FromMap(*decodedMap), nil
		}
	}

	// Attempt to decode as slice
	var ds []*discoveryv1.Discovery
	val, err = c.decodeComponent(constant.ComponentRegistries, &ds)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedSlice, ok := val.(*[]*discoveryv1.Discovery); ok {
			return discoveryConverter.FromSlice(*decodedSlice), nil
		}
	}

	// Attempt to decode as direct struct
	var s discoveryv1.Discoveries
	val, err = c.decodeComponent(constant.ComponentRegistries, &s)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedStruct, ok := val.(*discoveryv1.Discoveries); ok {
			return decodedStruct, nil
		}
	}
	return nil, nil
}

// DecodeDefaultDiscovery decodes and returns the default discovery name.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeDefaultDiscovery() (string, error) {
	var defaultRegistry string
	val, err := c.decodeComponent(constant.ComponentDefaultRegistry, &defaultRegistry)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	// Dereference the pointer to string before returning
	return *val.(*string), nil
}

// DecodeMiddlewares decodes and returns the Middlewares configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	// Attempt to decode as map
	var m map[string]*middlewarev1.Middleware
	val, err := c.decodeComponent(constant.ComponentMiddlewares, &m)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedMap, ok := val.(*map[string]*middlewarev1.Middleware); ok && len(*decodedMap) > 0 {
			return middlewareConverter.FromMap(*decodedMap), nil
		}
	}

	// Attempt to decode as slice
	var ms []*middlewarev1.Middleware
	val, err = c.decodeComponent(constant.ComponentMiddlewares, &ms)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedSlice, ok := val.(*[]*middlewarev1.Middleware); ok {
			return middlewareConverter.FromSlice(*decodedSlice), nil
		}
	}

	// Attempt to decode as direct struct
	var s middlewarev1.Middlewares
	val, err = c.decodeComponent(constant.ComponentMiddlewares, &s)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedStruct, ok := val.(*middlewarev1.Middlewares); ok {
			return decodedStruct, nil
		}
	}
	return nil, nil
}

// DecodeServers decodes and returns the Servers configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeServers() (*transportv1.Servers, error) {
	// Attempt to decode as map
	var m map[string]*transportv1.Server
	val, err := c.decodeComponent(constant.ComponentServers, &m)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedMap, ok := val.(*map[string]*transportv1.Server); ok && len(*decodedMap) > 0 {
			return serverConverter.FromMap(*decodedMap), nil
		}
	}

	// Attempt to decode as slice
	var ss []*transportv1.Server
	val, err = c.decodeComponent(constant.ComponentServers, &ss)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedSlice, ok := val.(*[]*transportv1.Server); ok {
			return serverConverter.FromSlice(*decodedSlice), nil
		}
	}

	// Attempt to decode as direct struct
	var s transportv1.Servers
	val, err = c.decodeComponent(constant.ComponentServers, &s)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedStruct, ok := val.(*transportv1.Servers); ok {
			return decodedStruct, nil
		}
	}
	return nil, nil
}

// DecodeClients decodes and returns the Clients configuration.
// It retrieves the configuration from the cache if available, otherwise decodes it and stores it in the cache.
func (c *structuredConfigImpl) DecodeClients() (*transportv1.Clients, error) {
	// Attempt to decode as map
	var m map[string]*transportv1.Client
	val, err := c.decodeComponent(constant.ComponentClients, &m)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedMap, ok := val.(*map[string]*transportv1.Client); ok && len(*decodedMap) > 0 {
			return clientConverter.FromMap(*decodedMap), nil
		}
	}

	// Attempt to decode as slice
	var cs []*transportv1.Client
	val, err = c.decodeComponent(constant.ComponentClients, &cs)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedSlice, ok := val.(*[]*transportv1.Client); ok {
			return clientConverter.FromSlice(*decodedSlice), nil
		}
	}

	// Attempt to decode as direct struct
	var s transportv1.Clients
	val, err = c.decodeComponent(constant.ComponentClients, &s)
	if err != nil {
		return nil, err
	}
	if val != nil {
		if decodedStruct, ok := val.(*transportv1.Clients); ok {
			return decodedStruct, nil
		}
	}
	return nil, nil
}
