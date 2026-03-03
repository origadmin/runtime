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
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/constant"
)

// NewStructured creates a new structured config implementation.
func NewStructured(cfg runtimeconfig.KConfig, paths map[constant.ComponentKey]string) contracts.StructuredConfig {
	if paths == nil {
		paths = make(map[constant.ComponentKey]string)
	}
	return newEagerConfig(cfg, paths)
}

// scanConfig is a helper to perform standard Kratos scanning.
func scanConfig(c runtimeconfig.KConfig, path string, target any) error {
	if path == "" {
		return c.Scan(target)
	}
	return c.Value(path).Scan(target)
}

// =============================================================================
// Eager Implementation (Deprecated)
// =============================================================================

type eagerConfigImpl struct {
	loader runtimeconfig.KConfig

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

func newEagerConfig(cfg runtimeconfig.KConfig, paths map[constant.ComponentKey]string) *eagerConfigImpl {
	impl := &eagerConfigImpl{
		loader: cfg,
	}

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

	decodeStruct := func(key constant.ComponentKey, target any) {
		path, enabled := getPath(key)
		if !enabled {
			return
		}
		var err error
		if path == "" {
			err = cfg.Scan(target)
		} else {
			err = cfg.Value(path).Scan(target)
		}
		if err != nil && !errors.Is(err, kratosconfig.ErrNotFound) {
			panic(fmt.Sprintf("failed to decode config for %s: %v", key, err))
		}
	}

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

	decodeMulti := func(key constant.ComponentKey, decodeFunc func(runtimeconfig.KConfig, string) (any, error)) any {
		path, enabled := getPath(key)
		if !enabled {
			return nil
		}
		res, err := decodeFunc(cfg, path)
		if err != nil {
			panic(err)
		}
		return res
	}

	impl.discoveries, _ = decodeMulti(constant.ComponentRegistries, func(c runtimeconfig.KConfig, p string) (any, error) {
		return decodeDiscoveries(c, p)
	}).(*discoveryv1.Discoveries)
	impl.middlewares, _ = decodeMulti(constant.ComponentMiddlewares, func(c runtimeconfig.KConfig, p string) (any, error) {
		return decodeMiddlewares(c, p)
	}).(*middlewarev1.Middlewares)
	impl.servers, _ = decodeMulti(constant.ComponentServers, func(c runtimeconfig.KConfig, p string) (any, error) {
		return decodeServers(c, p)
	}).(*transportv1.Servers)
	impl.clients, _ = decodeMulti(constant.ComponentClients, func(c runtimeconfig.KConfig, p string) (any, error) {
		return decodeClients(c, p)
	}).(*transportv1.Clients)

	return impl
}

func (c *eagerConfigImpl) DecodedConfig() any { return c.loader }

func (c *eagerConfigImpl) DecodeApp() (*appv1.App, error)              { return c.app, nil }
func (c *eagerConfigImpl) DecodeLogger() (*loggerv1.Logger, error)     { return c.logger, nil }
func (c *eagerConfigImpl) DecodeData() (*datav1.Data, error)           { return c.data, nil }
func (c *eagerConfigImpl) DecodeCaches() (*datav1.Caches, error)       { return c.caches, nil }
func (c *eagerConfigImpl) DecodeDatabases() (*datav1.Databases, error) { return c.databases, nil }
func (c *eagerConfigImpl) DecodeObjectStores() (*datav1.ObjectStores, error) {
	return c.objectStores, nil
}
func (c *eagerConfigImpl) DecodeDefaultDiscovery() (string, error) { return c.defaultRegistry, nil }
func (c *eagerConfigImpl) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	return c.discoveries, nil
}
func (c *eagerConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	return c.middlewares, nil
}
func (c *eagerConfigImpl) DecodeServers() (*transportv1.Servers, error) { return c.servers, nil }
func (c *eagerConfigImpl) DecodeClients() (*transportv1.Clients, error) { return c.clients, nil }

// Shared Helpers (Keep compilation)

func decodeDiscoveries(c runtimeconfig.KConfig, path string) (*discoveryv1.Discoveries, error) {
	var result discoveryv1.Discoveries
	if err := c.Value(path).Scan(&result); err == nil && len(result.Configs) > 0 {
		return &result, nil
	}
	return &discoveryv1.Discoveries{}, nil
}

func decodeMiddlewares(c runtimeconfig.KConfig, path string) (*middlewarev1.Middlewares, error) {
	var result middlewarev1.Middlewares
	if err := c.Value(path).Scan(&result); err == nil && len(result.Configs) > 0 {
		return &result, nil
	}
	return &middlewarev1.Middlewares{}, nil
}

func decodeServers(c runtimeconfig.KConfig, path string) (*transportv1.Servers, error) {
	var result transportv1.Servers
	if err := c.Value(path).Scan(&result); err == nil && len(result.Configs) > 0 {
		return &result, nil
	}
	return &transportv1.Servers{}, nil
}

func decodeClients(c runtimeconfig.KConfig, path string) (*transportv1.Clients, error) {
	var result transportv1.Clients
	if err := c.Value(path).Scan(&result); err == nil && len(result.Configs) > 0 {
		return &result, nil
	}
	return &transportv1.Clients{}, nil
}

// Lazy Implementation (Legacy, keep compilation)
type lazyConfigImpl struct {
	runtimeconfig.KConfig
	paths map[constant.ComponentKey]string
	cache sync.Map
}

func (c *lazyConfigImpl) DecodedConfig() any                                    { return c.KConfig }
func (c *lazyConfigImpl) DecodeApp() (*appv1.App, error)                        { return nil, nil }
func (c *lazyConfigImpl) DecodeLogger() (*loggerv1.Logger, error)               { return nil, nil }
func (c *lazyConfigImpl) DecodeData() (*datav1.Data, error)                     { return nil, nil }
func (c *lazyConfigImpl) DecodeCaches() (*datav1.Caches, error)                 { return nil, nil }
func (c *lazyConfigImpl) DecodeDatabases() (*datav1.Databases, error)           { return nil, nil }
func (c *lazyConfigImpl) DecodeObjectStores() (*datav1.ObjectStores, error)     { return nil, nil }
func (c *lazyConfigImpl) DecodeDefaultDiscovery() (string, error)               { return "", nil }
func (c *lazyConfigImpl) DecodeDiscoveries() (*discoveryv1.Discoveries, error)  { return nil, nil }
func (c *lazyConfigImpl) DecodeMiddlewares() (*middlewarev1.Middlewares, error) { return nil, nil }
func (c *lazyConfigImpl) DecodeServers() (*transportv1.Servers, error)          { return nil, nil }
func (c *lazyConfigImpl) DecodeClients() (*transportv1.Clients, error)          { return nil, nil }
