package config

import (
	"errors"
	"fmt"
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

// NewStructured creates a new structured config implementation.
// Currently switched to the Eager implementation.
func NewStructured(cfg interfaces.ConfigLoader, paths map[constant.ComponentKey]string) interfaces.StructuredConfig {
	if paths == nil {
		paths = make(map[constant.ComponentKey]string)
	}
	// Switch implementation here:
	// return newLazyConfig(cfg, paths) // Old implementation
	return newEagerConfig(cfg, paths) // New implementation
}

// =============================================================================
// Eager Implementation (New)
// =============================================================================

// eagerConfigImpl implements interfaces.StructuredConfig with eager loading.
// All configurations are decoded at initialization time.
type eagerConfigImpl struct {
	loader interfaces.ConfigLoader

	app             *appv1.App
	data            *datav1.Data
	logger          *loggerv1.Logger
	caches          *datav1.Caches
	databases       *datav1.Databases
	objectStores    *datav1.ObjectStores
	defaultRegistry string
	discoveries     *discoveryv1.Discoveries
	middlewares     *middlewarev1.Middlewares
	servers         *transportv1.Servers
	clients         *transportv1.Clients
}

func newEagerConfig(cfg interfaces.ConfigLoader, paths map[constant.ComponentKey]string) *eagerConfigImpl {
	impl := &eagerConfigImpl{
		loader: cfg,
	}

	// Helper to resolve path
	getPath := func(key constant.ComponentKey) (string, bool) {
		path, ok := paths[key]
		if !ok {
			return string(key), true
		}
		if path == "" {
			return "", false
		}
		return path, true
	}

	// Helper for simple struct decoding
	decodeStruct := func(key constant.ComponentKey, target any) {
		path, enabled := getPath(key)
		if !enabled {
			return
		}
		err := cfg.Decode(path, target)
		if err != nil && !errors.Is(err, kratosconfig.ErrNotFound) {
			panic(fmt.Sprintf("failed to decode config for %s (path: %s): %v", key, path, err))
		}
	}

	// 1. Decode Simple Components
	impl.app = &appv1.App{}
	decodeStruct(constant.ConfigApp, impl.app)
	if impl.app.Id == "" && impl.app.Name == "" {
		impl.app = nil
	}

	impl.logger = &loggerv1.Logger{}
	decodeStruct(constant.ComponentLogger, impl.logger)
	if impl.logger.Name == "" && len(impl.logger.Level) == 0 {
		impl.logger = nil
	}

	impl.data = &datav1.Data{}
	decodeStruct(constant.ComponentData, impl.data)

	impl.caches = &datav1.Caches{}
	decodeStruct(constant.ComponentCaches, impl.caches)

	impl.databases = &datav1.Databases{}
	decodeStruct(constant.ComponentDatabases, impl.databases)

	impl.objectStores = &datav1.ObjectStores{}
	decodeStruct(constant.ComponentObjectStores, impl.objectStores)

	decodeStruct(constant.ComponentDefaultRegistry, &impl.defaultRegistry)

	// 2. Decode Complex Components
	decodeMulti := func(key constant.ComponentKey, decodeFunc func(interfaces.ConfigLoader, string) (any, error)) any {
		path, enabled := getPath(key)
		if !enabled {
			return nil
		}
		res, err := decodeFunc(cfg, path)
		if err != nil {
			panic(fmt.Sprintf("failed to decode multi-format config for %s (path: %s): %v", key, path, err))
		}
		return res
	}

	impl.discoveries = decodeMulti(constant.ComponentRegistries, func(l interfaces.ConfigLoader, p string) (any, error) {
		return decodeDiscoveries(l, p)
	}).(*discoveryv1.Discoveries)

	impl.middlewares = decodeMulti(constant.ComponentMiddlewares, func(l interfaces.ConfigLoader, p string) (any, error) {
		return decodeMiddlewares(l, p)
	}).(*middlewarev1.Middlewares)

	impl.servers = decodeMulti(constant.ComponentServers, func(l interfaces.ConfigLoader, p string) (any, error) {
		return decodeServers(l, p)
	}).(*transportv1.Servers)

	impl.clients = decodeMulti(constant.ComponentClients, func(l interfaces.ConfigLoader, p string) (any, error) {
		return decodeClients(l, p)
	}).(*transportv1.Clients)

	return impl
}

func (c *eagerConfigImpl) DecodedConfig() any { return c.loader }

// Getters for Eager Implementation (Always return nil error)

func (c *eagerConfigImpl) DecodeApp() (*appv1.App, error)              { return c.app, nil }
func (c *eagerConfigImpl) DecodeLogger() (*loggerv1.Logger, error)     { return c.logger, nil }
func (c *eagerConfigImpl) DecodeData() (*datav1.Data, error)           { return c.data, nil }
func (c *eagerConfigImpl) DecodeCaches() (*datav1.Caches, error)       { return c.caches, nil }
func (c *eagerConfigImpl) DecodeDatabases() (*datav1.Databases, error) { return c.databases, nil }
func (c *eagerConfigImpl) DecodeObjectStores() (*datav1.ObjectStores, error) {
	return c.objectStores, nil
}
func (c *eagerConfigImpl) DecodeDefaultDiscovery() (string, error)     { return c.defaultRegistry, nil }
func (c *eagerConfigImpl) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	return c.discoveries, nil
}
func (c *eagerConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	return c.middlewares, nil
}
func (c *eagerConfigImpl) DecodeServers() (*transportv1.Servers, error) { return c.servers, nil }
func (c *eagerConfigImpl) DecodeClients() (*transportv1.Clients, error) { return c.clients, nil }

// =============================================================================
// Lazy Implementation (Old)
// =============================================================================

// lazyConfigImpl implements interfaces.StructuredConfig with lazy loading and caching.
type lazyConfigImpl struct {
	interfaces.ConfigLoader
	paths map[constant.ComponentKey]string
	cache sync.Map
}

func newLazyConfig(cfg interfaces.ConfigLoader, paths map[constant.ComponentKey]string) *lazyConfigImpl {
	return &lazyConfigImpl{
		ConfigLoader: cfg,
		paths:        paths,
		cache:        sync.Map{},
	}
}

func (c *lazyConfigImpl) DecodedConfig() any { return c.ConfigLoader }

func (c *lazyConfigImpl) decodeAndCache(componentKey constant.ComponentKey, valuePtr any) (any, error) {
	if cachedVal, ok := c.cache.Load(componentKey); ok {
		return cachedVal, nil
	}

	path, ok := c.paths[componentKey]
	if !ok {
		path = string(componentKey)
	} else if path == "" {
		return nil, nil
	}

	err := c.Decode(path, valuePtr)
	if err != nil {
		if errors.Is(err, kratosconfig.ErrNotFound) {
			c.cache.Store(componentKey, nil)
			return nil, nil
		}
		return nil, err
	}

	c.cache.Store(componentKey, valuePtr)
	return valuePtr, nil
}

func (c *lazyConfigImpl) DecodeApp() (*appv1.App, error) {
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

func (c *lazyConfigImpl) DecodeLogger() (*loggerv1.Logger, error) {
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

func (c *lazyConfigImpl) DecodeData() (*datav1.Data, error) {
	val, err := c.decodeAndCache(constant.ComponentData, &datav1.Data{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Data), nil
}

func (c *lazyConfigImpl) DecodeCaches() (*datav1.Caches, error) {
	val, err := c.decodeAndCache(constant.ComponentCaches, &datav1.Caches{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Caches), nil
}

func (c *lazyConfigImpl) DecodeDatabases() (*datav1.Databases, error) {
	val, err := c.decodeAndCache(constant.ComponentDatabases, &datav1.Databases{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.Databases), nil
}

func (c *lazyConfigImpl) DecodeObjectStores() (*datav1.ObjectStores, error) {
	val, err := c.decodeAndCache(constant.ComponentObjectStores, &datav1.ObjectStores{})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return val.(*datav1.ObjectStores), nil
}

func (c *lazyConfigImpl) DecodeDefaultDiscovery() (string, error) {
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

func (c *lazyConfigImpl) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
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

func (c *lazyConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
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

func (c *lazyConfigImpl) DecodeServers() (*transportv1.Servers, error) {
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

func (c *lazyConfigImpl) DecodeClients() (*transportv1.Clients, error) {
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

// =============================================================================
// Shared Helpers (Used by both implementations)
// =============================================================================

func decodeDiscoveries(c interfaces.ConfigLoader, path string) (*discoveryv1.Discoveries, error) {
	var m map[string]*discoveryv1.Discovery
	if err := c.Decode(path, &m); err == nil && len(m) > 0 {
		items := make([]*discoveryv1.Discovery, 0, len(m))
		for name, item := range m {
			if item.Name == "" {
				item.Name = name
			}
			items = append(items, item)
		}
		return &discoveryv1.Discoveries{Configs: items}, nil
	}

	var s []*discoveryv1.Discovery
	if err := c.Decode(path, &s); err == nil && len(s) > 0 {
		return &discoveryv1.Discoveries{Configs: s}, nil
	}
	return nil, nil
}

func decodeMiddlewares(c interfaces.ConfigLoader, path string) (*middlewarev1.Middlewares, error) {
	var m map[string]*middlewarev1.Middleware
	if err := c.Decode(path, &m); err == nil && len(m) > 0 {
		items := make([]*middlewarev1.Middleware, 0, len(m))
		for name, item := range m {
			if item.Name == "" {
				item.Name = name
			}
			items = append(items, item)
		}
		return &middlewarev1.Middlewares{Configs: items}, nil
	}

	var s []*middlewarev1.Middleware
	if err := c.Decode(path, &s); err == nil && len(s) > 0 {
		return &middlewarev1.Middlewares{Configs: s}, nil
	}
	return nil, nil
}

func decodeServers(c interfaces.ConfigLoader, path string) (*transportv1.Servers, error) {
	var m map[string]*transportv1.Server
	if err := c.Decode(path, &m); err == nil && len(m) > 0 {
		items := make([]*transportv1.Server, 0, len(m))
		for name, item := range m {
			if item.Name == "" {
				item.Name = name
			}
			items = append(items, item)
		}
		return &transportv1.Servers{Configs: items}, nil
	}

	var s []*transportv1.Server
	if err := c.Decode(path, &s); err == nil && len(s) > 0 {
		return &transportv1.Servers{Configs: s}, nil
	}
	return nil, nil
}

func decodeClients(c interfaces.ConfigLoader, path string) (*transportv1.Clients, error) {
	var m map[string]*transportv1.Client
	if err := c.Decode(path, &m); err == nil && len(m) > 0 {
		items := make([]*transportv1.Client, 0, len(m))
		for name, item := range m {
			if item.Name == "" {
				item.Name = name
			}
			items = append(items, item)
		}
		return &transportv1.Clients{Configs: items}, nil
	}

	var s []*transportv1.Client
	if err := c.Decode(path, &s); err == nil && len(s) > 0 {
		return &transportv1.Clients{Configs: s}, nil
	}
	return nil, nil
}

// multiFormatDecoder is a generic helper for the Lazy implementation.
type multiFormatDecoder[ItemPtr any, ResultPtr any] struct {
	c             *lazyConfigImpl
	componentKey  constant.ComponentKey
	sliceToResult func([]ItemPtr) ResultPtr
	isResultEmpty func(ResultPtr) bool
	setNameOnItem func(name string, item ItemPtr)
}

func (d *multiFormatDecoder[ItemPtr, ResultPtr]) decode() (ResultPtr, error) {
	var zero ResultPtr
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

	var s []ItemPtr
	if err := d.c.Decode(path, &s); err == nil {
		if len(s) > 0 {
			result := d.sliceToResult(s)
			d.c.cache.Store(d.componentKey, result)
			return result, nil
		}
	}

	var direct ResultPtr
	if err := d.c.Decode(path, &direct); err == nil {
		if d.isResultEmpty == nil || !d.isResultEmpty(direct) {
			d.c.cache.Store(d.componentKey, direct)
			return direct, nil
		}
	}

	d.c.cache.Store(d.componentKey, nil)
	return zero, nil
}
