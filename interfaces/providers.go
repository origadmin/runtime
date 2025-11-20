package interfaces

import (
	kratosmiddleware "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
)

// RegistryProvider provides access to Discovery and Registrar components.
type RegistryProvider interface {
	Discoveries() (map[string]registry.Discovery, error)
	Discovery(name string) (registry.Discovery, error)
	Registrars() (map[string]registry.Registrar, error)
	Registrar(name string) (registry.Registrar, error)
	DefaultRegistrar() (registry.Registrar, error)
	RegisterDiscovery(name string, discovery registry.Discovery)
}

// MiddlewareProvider provides access to ServerMiddleware components.
type MiddlewareProvider interface {
	ServerMiddlewares() (map[string]kratosmiddleware.Middleware, error)
	ServerMiddleware(name string) (kratosmiddleware.Middleware, error)
	RegisterServerMiddleware(name string, middleware kratosmiddleware.Middleware)
	ClientMiddlewares() (map[string]kratosmiddleware.Middleware, error)
	ClientMiddleware(name string) (kratosmiddleware.Middleware, error)
	RegisterClientMiddleware(name string, middleware kratosmiddleware.Middleware)
}

// CacheProvider provides access to Cache components.
type CacheProvider interface {
	Caches() (map[string]Cache, error)
	Cache(name string) (Cache, error)
	RegisterCache(name string, cache Cache)
}

// DatabaseProvider provides access to Database components.
type DatabaseProvider interface {
	Databases() (map[string]Database, error)
	Database(name string) (Database, error)
	RegisterDatabase(name string, db Database)
}

// ObjectStoreProvider provides access to ObjectStore components.
type ObjectStoreProvider interface {
	ObjectStores() (map[string]ObjectStore, error)
	ObjectStore(name string) (ObjectStore, error)
}

// ComponentProvider provides access to generic Component components.
type ComponentProvider interface {
	Components() (map[string]Component, error)
	Component(name string) (Component, error)
}
