package container

import (
	"fmt"

	"github.com/origadmin/runtime/container/internal/cache"
	"github.com/origadmin/runtime/container/internal/database"
	"github.com/origadmin/runtime/container/internal/middleware"
	"github.com/origadmin/runtime/container/internal/objectstore"
	"github.com/origadmin/runtime/container/internal/registry"
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
	runtimeregistry "github.com/origadmin/runtime/registry"
)

// Container defines the interface for accessing various runtime components.
// It supports lazy initialization of components via factories and dependency injection.
type Container interface {
	// Registry returns the service registry provider.
	Registry(opts ...options.Option) (RegistryProvider, error)
	// Middleware returns the middleware provider.
	Middleware(opts ...options.Option) (MiddlewareProvider, error)
	// Cache returns the cache provider.
	Cache(opts ...options.Option) (CacheProvider, error)
	// Database returns the database provider.
	Database(opts ...options.Option) (DatabaseProvider, error)
	// ObjectStore returns the object storage provider.
	ObjectStore(opts ...options.Option) (ObjectStoreProvider, error)

	// Component retrieves a generic component by its name, using a lazy-initialization strategy.
	// If the component is not yet created, the container will use its registered factory to create it.
	// The provided options are only used during the creation of the component.
	// This method is concurrency-safe.
	Component(name string, opts ...options.Option) (interfaces.Component, error)

	// RegisterComponent registers a pre-built component instance.
	RegisterComponent(name string, comp interfaces.Component)

	// RegisterFactory registers a factory for creating a component on demand.
	RegisterFactory(name string, factory ComponentFactory)

	// HasComponent checks if a component instance or factory has been registered.
	HasComponent(name string) bool

	// RegisteredComponents returns a slice of names for all registered components and factories.
	RegisteredComponents() []string

	// Logger returns the container's configured logger.
	Logger() runtimelog.Logger
	// AppInfo returns the application's metadata.
	AppInfo() interfaces.AppInfo

	// DefaultCache returns the default cache instance.
	DefaultCache() (storageiface.Cache, error)
	// DefaultDatabase returns the default database instance.
	DefaultDatabase() (storageiface.Database, error)
	// DefaultObjectStore returns the default object storage instance.
	DefaultObjectStore() (storageiface.ObjectStore, error)
	// DefaultRegistrar returns the default service registrar instance.
	DefaultRegistrar() (runtimeregistry.KRegistrar, error)
}

// containerImpl implements the Container interface.
// It simplifies the provider access by delegating initialization safety to the providers themselves.
type containerImpl struct {
	config   interfaces.StructuredConfig
	logger   runtimelog.Logger
	initOpts *containerOptions // initOpts now holds AppInfo

	// Concurrency-safe store for generic components and their factories.
	componentStore *componentStore

	// Built-in provider instances.
	registryProvider    *registry.Provider
	middlewareProvider  *middleware.Provider
	cacheProvider       *cache.Provider
	databaseProvider    *database.Provider
	objectStoreProvider *objectstore.Provider
}

// New creates a new, concurrency-safe Container instance.
// It receives the final AppInfo from the bootstrap process.
func New(config interfaces.StructuredConfig, opts ...options.Option) Container {
	initOpts := optionutil.NewT[containerOptions](opts...)

	var baseLogger runtimelog.Logger
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		runtimelog.NewHelper(runtimelog.DefaultLogger).Warnf("failed to decode logger config, using default logger: %v", err)
		loggerConfig = nil
	}
	baseLogger = runtimelog.NewLogger(loggerConfig)

	enrichedLogger := baseLogger
	if initOpts.appInfo != nil { // Use appInfo from initOpts for logger enrichment
		enrichedLogger = runtimelog.With(baseLogger,
			"service.name", initOpts.appInfo.Name(),
			"service.version", initOpts.appInfo.Version(),
			"service.id", initOpts.appInfo.ID(),
		)
	}

	c := &containerImpl{
		config:         config,
		logger:         enrichedLogger,
		initOpts:       initOpts, // initOpts now holds AppInfo
		componentStore: newComponentStore(enrichedLogger),
	}

	// Eagerly create provider shells, but defer configuration.
	c.registryProvider = registry.NewProvider(c.logger)
	c.middlewareProvider = middleware.NewProvider(c.logger)
	c.cacheProvider = cache.NewProvider(c.logger)
	c.databaseProvider = database.NewProvider(c.logger)
	c.objectStoreProvider = objectstore.NewProvider(c.logger)

	return c
}

// AppInfo returns the application's metadata.
func (c *containerImpl) AppInfo() interfaces.AppInfo {
	if c.initOpts != nil {
		return c.initOpts.appInfo
	}
	return nil // Or a default empty AppInfo
}

// Logger returns the container's configured logger instance.
func (c *containerImpl) Logger() runtimelog.Logger {
	return c.logger
}

// Component retrieves a generic component by name, using a lazy-initialization strategy.
func (c *containerImpl) Component(name string, opts ...options.Option) (interfaces.Component, error) {
	// 1. Check for an existing instance. If found, return it immediately.
	// Options are ignored for already-created singleton instances.
	if comp, ok := c.componentStore.GetInstance(name); ok {
		return comp, nil
	}

	// 2. Check for a factory.
	factory, ok := c.componentStore.GetFactory(name)
	if !ok {
		return nil, fmt.Errorf("component '%s' not found", name)
	}

	// 3. Create the component using the factory. Options are only used here.
	comp, err := factory.NewComponent(c.config, c, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create component '%s' from factory: %w", name, err)
	}

	// 4. Store the new instance for future calls (singleton).
	c.componentStore.RegisterInstance(name, comp)

	return comp, nil
}

// RegisterComponent registers a pre-built component instance.
func (c *containerImpl) RegisterComponent(name string, comp interfaces.Component) {
	c.componentStore.RegisterInstance(name, comp)
}

// RegisterFactory registers a factory for creating a component on demand.
func (c *containerImpl) RegisterFactory(name string, factory ComponentFactory) {
	c.componentStore.RegisterFactory(name, factory)
}

// HasComponent checks if a component instance or factory has been registered.
func (c *containerImpl) HasComponent(name string) bool {
	return c.componentStore.Has(name)
}

// RegisteredComponents returns a slice of names for all registered components and factories.
func (c *containerImpl) RegisteredComponents() []string {
	return c.componentStore.List()
}

// Registry returns the service registry provider.
func (c *containerImpl) Registry(opts ...options.Option) (RegistryProvider, error) {
	discoveries, err := c.config.DecodeDiscoveries()
	if err != nil {
		runtimelog.NewHelper(c.logger).Errorf("failed to decode discoveries config for RegistryProvider: %v", err)
	}
	c.registryProvider.Initialize(discoveries, opts...)
	return c.registryProvider, nil
}

// Middleware returns the middleware provider.
func (c *containerImpl) Middleware(opts ...options.Option) (MiddlewareProvider, error) {
	middlewares, err := c.config.DecodeMiddlewares()
	if err != nil {
		runtimelog.NewHelper(c.logger).Errorf("failed to decode middlewares config for MiddlewareProvider: %v", err)
	}
	c.middlewareProvider.Initialize(middlewares, opts...)
	return c.middlewareProvider, nil
}

// Cache returns the cache provider.
func (c *containerImpl) Cache(opts ...options.Option) (CacheProvider, error) {
	caches, err := c.config.DecodeCaches()
	if err != nil {
		runtimelog.NewHelper(c.logger).Errorf("failed to decode caches config for CacheProvider: %v", err)
	}
	c.cacheProvider.Initialize(caches, opts...)
	return c.cacheProvider, nil
}

// Database returns the database provider.
func (c *containerImpl) Database(opts ...options.Option) (DatabaseProvider, error) {
	databases, err := c.config.DecodeDatabases()
	if err != nil {
		runtimelog.NewHelper(c.logger).Errorf("failed to decode databases config for DatabaseProvider: %v", err)
	}
	c.databaseProvider.Initialize(databases, opts...)
	return c.databaseProvider, nil
}

// ObjectStore returns the object storage provider.
func (c *containerImpl) ObjectStore(opts ...options.Option) (ObjectStoreProvider, error) {
	filestores, err := c.config.DecodeObjectStores()
	if err != nil {
		runtimelog.NewHelper(c.logger).Errorf("failed to decode object stores config for ObjectStoreProvider: %v", err)
	}
	c.objectStoreProvider.Initialize(filestores, opts...)
	return c.objectStoreProvider, nil
}

// DefaultRegistrar returns the default service registrar instance.
func (c *containerImpl) DefaultRegistrar() (runtimeregistry.KRegistrar, error) {
	registryProvider, err := c.Registry()
	if err != nil {
		return nil, err
	}
	registrar, err := registryProvider.DefaultRegistrar(c.initOpts.defaultRegistrarName)
	if err != nil {
		return nil, fmt.Errorf("default registrar '%s' not found or failed to retrieve: %w", c.initOpts.defaultRegistrarName, err)
	}
	return registrar, nil
}

// DefaultCache returns the default cache instance.
func (c *containerImpl) DefaultCache() (storageiface.Cache, error) {
	cacheProvider, err := c.Cache()
	if err != nil {
		return nil, err
	}
	cacheInstance, err := cacheProvider.DefaultCache(c.initOpts.defaultCacheName)
	if err != nil {
		return nil, fmt.Errorf("default cache '%s' not found or failed to retrieve: %w", c.initOpts.defaultCacheName, err)
	}
	return cacheInstance, nil
}

// DefaultDatabase returns the default database instance.
func (c *containerImpl) DefaultDatabase() (storageiface.Database, error) {
	databaseProvider, err := c.Database()
	if err != nil {
		return nil, err
	}
	dbInstance, err := databaseProvider.DefaultDatabase(c.initOpts.defaultDatabaseName)
	if err != nil {
		return nil, fmt.Errorf("default database '%s' not found or failed to retrieve: %w", c.initOpts.defaultDatabaseName, err)
	}
	return dbInstance, nil
}

// DefaultObjectStore returns the default object storage instance.
func (c *containerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	objectStoreProvider, err := c.ObjectStore()
	if err != nil {
		return nil, err
	}
	osInstance, err := objectStoreProvider.DefaultObjectStore(c.initOpts.defaultObjectStoreName)
	if err != nil {
		return nil, fmt.Errorf("default object store '%s' not found or failed to retrieve: %w", c.initOpts.defaultObjectStoreName, err)
	}
	return osInstance, nil
}
