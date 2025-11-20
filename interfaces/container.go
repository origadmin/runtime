package interfaces

// Container defines the interface for accessing various runtime components.
type Container interface {
	Registry() (RegistryProvider, error)
	Middleware() (MiddlewareProvider, error)
	Cache() (CacheProvider, error)
	Database() (DatabaseProvider, error)
	ObjectStore() (ObjectStoreProvider, error)
	Component() (ComponentProvider, error)
}
