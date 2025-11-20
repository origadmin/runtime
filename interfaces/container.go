package interfaces


// Container defines the interface for accessing various runtime components.
type Container interface {
	Registry() (RegistryProvider, error)
	ServerMiddleware() (ServerMiddlewareProvider, error)
	ClientMiddleware() (ClientMiddlewareProvider, error)
	Cache() (CacheProvider, error)
	Database() (DatabaseProvider, error)
	ObjectStore() (ObjectStoreProvider, error)
	Component() (ComponentProvider, error)
}