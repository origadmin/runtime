package container

import (
	"github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
)

// RegistryProvider provides access to Discovery and Registrar components.
// The provider is responsible for creating and managing the lifecycle of these components.
type RegistryProvider interface {
	Discoveries() (map[string]registry.KDiscovery, error)
	Discovery(name string) (registry.KDiscovery, error)
	Registrars() (map[string]registry.KRegistrar, error)
	Registrar(name string) (registry.KRegistrar, error)
	DefaultRegistrar(globalDefaultName string) (registry.KRegistrar, error)
	RegisterDiscovery(name string, discovery registry.KDiscovery)
	RegisterRegistrar(name string, registrar registry.KRegistrar)
}

// ServerMiddlewareProvider provides access to server-side middleware components.
type ServerMiddlewareProvider interface {
	ServerMiddlewares() (map[string]middleware.KMiddleware, error)
	ServerMiddleware(name string) (middleware.KMiddleware, error)
	RegisterServerMiddleware(name string, mw middleware.KMiddleware)
}

// ClientMiddlewareProvider provides access to client-side middleware components.
type ClientMiddlewareProvider interface {
	ClientMiddlewares() (map[string]middleware.KMiddleware, error)
	ClientMiddleware(name string) (middleware.KMiddleware, error)
	RegisterClientMiddleware(name string, mw middleware.KMiddleware)
}

// MiddlewareProvider provides access to both server and client middleware components.
type MiddlewareProvider interface {
	ServerMiddlewareProvider
	ClientMiddlewareProvider
}

// CacheProvider provides access to Cache components.
type CacheProvider interface {
	Caches() (map[string]storage.Cache, error)
	Cache(name string) (storage.Cache, error)
	DefaultCache(globalDefaultName string) (storage.Cache, error)
	RegisterCache(name string, cache storage.Cache)
}

// DatabaseProvider provides access to Database components.
type DatabaseProvider interface {
	Databases() (map[string]storage.Database, error)
	Database(name string) (storage.Database, error)
	DefaultDatabase(globalDefaultName string) (storage.Database, error)
	RegisterDatabase(name string, db storage.Database)
}

// ObjectStoreProvider provides access to ObjectStore components.
type ObjectStoreProvider interface {
	ObjectStores() (map[string]storage.ObjectStore, error)
	ObjectStore(name string) (storage.ObjectStore, error)
	DefaultObjectStore(globalDefaultName string) (storage.ObjectStore, error)
	RegisterObjectStore(name string, store storage.ObjectStore)
}
