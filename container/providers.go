package container

import (
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
)

// RegistryProvider provides access to Discovery and Registrar components.
type RegistryProvider interface {
	Discoveries() (map[string]registry.KDiscovery, error)
	Discovery(name string) (registry.KDiscovery, error)
	Registrars() (map[string]registry.KRegistrar, error)
	Registrar(name string) (registry.KRegistrar, error)
	DefaultRegistrar(globalDefaultName string) (registry.KRegistrar, error) // Modified to accept globalDefaultName
	RegisterDiscovery(name string, discovery registry.KDiscovery)
	RegisterRegistrar(name string, registrar registry.KRegistrar)
	SetOptions(opts ...options.Option) // Add SetOptions method
}

type ServerMiddlewareProvider interface {
	ServerMiddlewares() (map[string]middleware.KMiddleware, error)
	ServerMiddleware(name string) (middleware.KMiddleware, error)
	RegisterServerMiddleware(name string, middleware middleware.KMiddleware)
}

type ClientMiddlewareProvider interface {
	ClientMiddlewares() (map[string]middleware.KMiddleware, error)
	ClientMiddleware(name string) (middleware.KMiddleware, error)
	RegisterClientMiddleware(name string, middleware middleware.KMiddleware)
}

// MiddlewareProvider provides access to ServerMiddleware components.
type MiddlewareProvider interface {
	ServerMiddlewareProvider
	ClientMiddlewareProvider
	SetOptions(opts ...options.Option) // Add SetOptions method
}

// CacheProvider provides access to Cache components.
type CacheProvider interface {
	Caches() (map[string]storage.Cache, error)
	Cache(name string) (storage.Cache, error)
	DefaultCache(globalDefaultName string) (storage.Cache, error) // Modified to accept globalDefaultName
	RegisterCache(name string, cache storage.Cache)
	SetOptions(opts ...options.Option) // Add SetOptions method
}

// DatabaseProvider provides access to Database components.
type DatabaseProvider interface {
	Databases() (map[string]storage.Database, error)
	Database(name string) (storage.Database, error)
	DefaultDatabase(globalDefaultName string) (storage.Database, error) // Modified to accept globalDefaultName
	RegisterDatabase(name string, db storage.Database)
	SetOptions(opts ...options.Option) // Add SetOptions method
}

// ObjectStoreProvider provides access to ObjectStore components.
type ObjectStoreProvider interface {
	ObjectStores() (map[string]storage.ObjectStore, error)
	ObjectStore(name string) (storage.ObjectStore, error)
	DefaultObjectStore(globalDefaultName string) (storage.ObjectStore, error) // Modified to accept globalDefaultName
	RegisterObjectStore(name string, store storage.ObjectStore)
	SetOptions(opts ...options.Option) // Add SetOptions method
}
