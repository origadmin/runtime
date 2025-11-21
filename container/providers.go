package container

import (
	kratosmiddleware "github.com/go-kratos/kratos/v2/middleware"
	kratosregistry "github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/interfaces/storage"
)

// RegistryProvider provides access to Discovery and Registrar components.
type RegistryProvider interface {
	Discoveries() (map[string]kratosregistry.Discovery, error)
	Discovery(name string) (kratosregistry.Discovery, error)
	Registrars() (map[string]kratosregistry.Registrar, error)
	Registrar(name string) (kratosregistry.Registrar, error)
	DefaultRegistrar() (kratosregistry.Registrar, error)
	RegisterDiscovery(name string, discovery kratosregistry.Discovery)
}

type ServerMiddlewareProvider interface {
	ServerMiddlewares() (map[string]kratosmiddleware.Middleware, error)
	ServerMiddleware(name string) (kratosmiddleware.Middleware, error)
	RegisterServerMiddleware(name string, middleware kratosmiddleware.Middleware)
}

type ClientMiddlewareProvider interface {
	ClientMiddlewares() (map[string]kratosmiddleware.Middleware, error)
	ClientMiddleware(name string) (kratosmiddleware.Middleware, error)
	RegisterClientMiddleware(name string, middleware kratosmiddleware.Middleware)
}

// MiddlewareProvider provides access to ServerMiddleware components.
type MiddlewareProvider interface {
	ServerMiddlewareProvider
	ClientMiddlewareProvider
}

// CacheProvider provides access to Cache components.
type CacheProvider interface {
	Caches() (map[string]storage.Cache, error)
	Cache(name string) (storage.Cache, error)
	RegisterCache(name string, cache storage.Cache)
}

// DatabaseProvider provides access to Database components.
type DatabaseProvider interface {
	Databases() (map[string]storage.Database, error)
	Database(name string) (storage.Database, error)
	RegisterDatabase(name string, db storage.Database)
}

// ObjectStoreProvider provides access to ObjectStore components.
type ObjectStoreProvider interface {
	ObjectStores() (map[string]storage.ObjectStore, error)
	ObjectStore(name string) (storage.ObjectStore, error)
	RegisterObjectStore(name string, store storage.ObjectStore)
}
