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
)

// structuredConfigImpl implements the interfaces.StructuredConfig interface.
// It wraps a generic interfaces.Config and provides type-safe, path-based decoding methods.
type structuredConfigImpl struct {
	interfaces.Config // Embed the generic config interface
	paths             map[constant.ComponentKey]string
	cache             sync.Map // Cache for decoded configurations
}

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

// DecodedConfig returns the underlying generic configuration.
func (c *structuredConfigImpl) DecodedConfig() any {
	return c.Config
}

// Statically assert that structuredConfigImpl implements the full StructuredConfig interface.
var _ interfaces.StructuredConfig = (*structuredConfigImpl)(nil)

// decodeAndCache implements a robust decoding logic with caching for single-type components.
// It first checks the `paths` map for a pre-discovered path. If not found,
// it falls back to using the componentKey directly as the path.
// If the component is already in the cache, it returns the cached value.
// Otherwise, it decodes the component, stores it in the cache, and then returns it.
func (c *structuredConfigImpl) decodeAndCache(componentKey constant.ComponentKey, valuePtr any) (any, error) {
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

// DecodeApp decodes and returns the App configuration.
func (c *structuredConfigImpl) DecodeApp() (*appv1.App, error) {
	val, err := c.decodeAndCache(constant.ConfigApp, &appv1.App{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	appConfig := val.(*appv1.App)
	if appConfig.Id == "" && appConfig.Name == "" {
		return nil, nil
	}
	return appConfig, nil
}

// DecodeLogger decodes and returns the Logger configuration.
func (c *structuredConfigImpl) DecodeLogger() (*loggerv1.Logger, error) {
	val, err := c.decodeAndCache(constant.ComponentLogger, &loggerv1.Logger{})
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

// DecodeData decodes and returns the Data configuration.
func (c *structuredConfigImpl) DecodeData() (*datav1.Data, error) {
	val, err := c.decodeAndCache(constant.ComponentData, &datav1.Data{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Data), nil
}

// DecodeCaches decodes and returns the Caches configuration.
func (c *structuredConfigImpl) DecodeCaches() (*datav1.Caches, error) {
	val, err := c.decodeAndCache(constant.ComponentCaches, &datav1.Caches{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Caches), nil
}

// DecodeDatabases decodes and returns the Databases configuration.
func (c *structuredConfigImpl) DecodeDatabases() (*datav1.Databases, error) {
	val, err := c.decodeAndCache(constant.ComponentDatabases, &datav1.Databases{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Databases), nil
}

// DecodeObjectStores decodes and returns the ObjectStores configuration.
func (c *structuredConfigImpl) DecodeObjectStores() (*datav1.ObjectStores, error) {
	val, err := c.decodeAndCache(constant.ComponentObjectStores, &datav1.ObjectStores{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.ObjectStores), nil
}

// DecodeDefaultDiscovery decodes and returns the default discovery name.
func (c *structuredConfigImpl) DecodeDefaultDiscovery() (string, error) {
	var defaultRegistry string
	val, err := c.decodeAndCache(constant.ComponentDefaultRegistry, &defaultRegistry)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	return *val.(*string), nil
}

// multiFormatDecoder is a generic helper to decode configuration components
// that can be represented as a map, a slice, or a direct struct.
type multiFormatDecoder[ItemPtr any, ResultPtr any] struct {
	c             *structuredConfigImpl
	componentKey  constant.ComponentKey
	sliceToResult func([]ItemPtr) ResultPtr
	isResultEmpty func(ResultPtr) bool
	setNameOnItem func(name string, item ItemPtr)
}

// decode executes the decoding process.
func (d *multiFormatDecoder[ItemPtr, ResultPtr]) decode() (ResultPtr, error) {
	var zero ResultPtr
	// 1. Check cache for the final, converted result.
	if cachedVal, ok := d.c.cache.Load(d.componentKey); ok {
		if cachedVal == nil {
			return zero, nil
		}
		return cachedVal.(ResultPtr), nil
	}

	path, ok := d.c.paths[d.componentKey]
	if !ok {
		path = string(d.componentKey)
	} else if path == "" {
		d.c.cache.Store(d.componentKey, nil)
		return zero, nil
	}

	// 2. Attempt to decode as a map.
	var m map[string]ItemPtr
	err := d.c.Decode(path, &m)
	if err == nil {
		if len(m) > 0 {
			items := make([]ItemPtr, 0, len(m))
			for name, item := range m {
				if d.setNameOnItem != nil {
					d.setNameOnItem(name, item)
				}
				items = append(items, item)
			}
			result := d.sliceToResult(items)
			d.c.cache.Store(d.componentKey, result)
			return result, nil
		}
	} else if errors.Is(err, kratosconfig.ErrNotFound) {
		d.c.cache.Store(d.componentKey, nil)
		return zero, nil
	}

	// 3. Attempt to decode as a slice.
	var s []ItemPtr
	if err := d.c.Decode(path, &s); err == nil {
		if len(s) > 0 {
			result := d.sliceToResult(s)
			d.c.cache.Store(d.componentKey, result)
			return result, nil
		}
	}

	// 4. Attempt to decode as a direct struct.
	var direct ResultPtr
	if err := d.c.Decode(path, &direct); err == nil {
		if d.isResultEmpty == nil || !d.isResultEmpty(direct) {
			d.c.cache.Store(d.componentKey, direct)
			return direct, nil
		}
	}

	// 5. If all attempts fail or result in empty config, cache nil.
	d.c.cache.Store(d.componentKey, nil)
	return zero, nil
}

// DecodeDiscoveries decodes the Discoveries configuration using the multi-format decoder.
func (c *structuredConfigImpl) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	decoder := &multiFormatDecoder[*discoveryv1.Discovery, *discoveryv1.Discoveries]{
		c:            c,
		componentKey: constant.ComponentRegistries,
		sliceToResult: func(items []*discoveryv1.Discovery) *discoveryv1.Discoveries {
			if len(items) == 0 {
				return nil
			}
			return &discoveryv1.Discoveries{Configs: items}
		},
		isResultEmpty: func(r *discoveryv1.Discoveries) bool {
			return r == nil || len(r.Configs) == 0
		},
		setNameOnItem: func(name string, item *discoveryv1.Discovery) {
			if item.Name == "" {
				item.Name = name
			}
		},
	}
	return decoder.decode()
}

// DecodeMiddlewares decodes the Middlewares configuration using the multi-format decoder.
func (c *structuredConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	decoder := &multiFormatDecoder[*middlewarev1.Middleware, *middlewarev1.Middlewares]{
		c:            c,
		componentKey: constant.ComponentMiddlewares,
		sliceToResult: func(items []*middlewarev1.Middleware) *middlewarev1.Middlewares {
			if len(items) == 0 {
				return nil
			}
			return &middlewarev1.Middlewares{Configs: items}
		},
		isResultEmpty: func(r *middlewarev1.Middlewares) bool {
			return r == nil || len(r.Configs) == 0
		},
		setNameOnItem: func(name string, item *middlewarev1.Middleware) {
			if item.Name == "" {
				item.Name = name
			}
		},
	}
	return decoder.decode()
}

// DecodeServers decodes the Servers configuration using the multi-format decoder.
func (c *structuredConfigImpl) DecodeServers() (*transportv1.Servers, error) {
	decoder := &multiFormatDecoder[*transportv1.Server, *transportv1.Servers]{
		c:            c,
		componentKey: constant.ComponentServers,
		sliceToResult: func(items []*transportv1.Server) *transportv1.Servers {
			if len(items) == 0 {
				return nil
			}
			return &transportv1.Servers{Configs: items}
		},
		isResultEmpty: func(r *transportv1.Servers) bool {
			return r == nil || len(r.Configs) == 0
		},
		setNameOnItem: func(name string, item *transportv1.Server) {
			if item.Name == "" {
				item.Name = name
			}
		},
	}
	return decoder.decode()
}

// DecodeClients decodes the Clients configuration using the multi-format decoder.
func (c *structuredConfigImpl) DecodeClients() (*transportv1.Clients, error) {
	decoder := &multiFormatDecoder[*transportv1.Client, *transportv1.Clients]{
		c:            c,
		componentKey: constant.ComponentClients,
		sliceToResult: func(items []*transportv1.Client) *transportv1.Clients {
			if len(items) == 0 {
				return nil
			}
			return &transportv1.Clients{Configs: items}
		},
		isResultEmpty: func(r *transportv1.Clients) bool {
			return r == nil || len(r.Configs) == 0
		},
		setNameOnItem: func(name string, item *transportv1.Client) {
			if item.Name == "" {
				item.Name = name
			}
		},
	}
	return decoder.decode()
}
