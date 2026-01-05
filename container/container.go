package container

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

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
type Container interface {
	Registry(opts ...options.Option) (RegistryProvider, error)
	Middleware(opts ...options.Option) (MiddlewareProvider, error)
	Cache(opts ...options.Option) (CacheProvider, error)
	Database(opts ...options.Option) (DatabaseProvider, error)
	ObjectStore(opts ...options.Option) (ObjectStoreProvider, error)
	Component(name string, opts ...options.Option) (interfaces.Component, error)
	RegisterComponent(name string, comp interfaces.Component)
	RegisterFactory(name string, factory ComponentFactory)
	HasComponent(name string) bool
	RegisteredComponents() []string
	Logger() runtimelog.Logger
	AppInfo() interfaces.AppInfo
	DefaultCache() (storageiface.Cache, error)
	DefaultDatabase() (storageiface.Database, error)
	DefaultObjectStore() (storageiface.ObjectStore, error)
	DefaultRegistrar() (runtimeregistry.KRegistrar, error)
}

// containerImpl implements the Container interface.
type containerImpl struct {
	config   interfaces.StructuredConfig
	logger   runtimelog.Logger
	helper   *runtimelog.Helper
	initOpts *containerOptions

	componentStore *componentStore

	registryProvider    *registry.Provider
	middlewareProvider  *middleware.Provider
	cacheProvider       *cache.Provider
	databaseProvider    *database.Provider
	objectStoreProvider *objectstore.Provider
}

// New creates a new, concurrency-safe Container instance.
func New(config interfaces.StructuredConfig, opts ...options.Option) Container {
	initOpts := optionutil.NewT[containerOptions](opts...)

	// Use the logger from options or create a new one from config.
	baseLogger := initOpts.Logger
	if baseLogger == nil {
		loggerConfig, err := config.DecodeLogger()
		if err != nil {
			// Use default logger to report the error, then proceed with default.
			runtimelog.NewHelper(runtimelog.DefaultLogger).Warnf("failed to decode logger config, using default logger: %v", err)
			loggerConfig = nil
		}
		baseLogger = runtimelog.NewLogger(loggerConfig)
	}

	enrichedLogger := baseLogger
	if initOpts.appInfo != nil {
		enrichedLogger = runtimelog.With(baseLogger,
			"service.name", initOpts.appInfo.Name(),
			"service.version", initOpts.appInfo.Version(),
			"service.id", initOpts.appInfo.ID(),
		)
	}

	helper := runtimelog.NewHelper(log.With(enrichedLogger, "module", "container"))
	helper.Debug("Runtime container initializing...")

	c := &containerImpl{
		config:         config,
		logger:         enrichedLogger,
		helper:         helper,
		initOpts:       initOpts,
		componentStore: newComponentStore(enrichedLogger),
	}

	c.registryProvider = registry.NewProvider(c.logger)
	c.middlewareProvider = middleware.NewProvider(c.logger)
	c.cacheProvider = cache.NewProvider(c.logger)
	c.databaseProvider = database.NewProvider(c.logger)
	c.objectStoreProvider = objectstore.NewProvider(c.logger)

	helper.Debug("Runtime container initialized successfully.")
	return c
}

func (c *containerImpl) AppInfo() interfaces.AppInfo {
	if c.initOpts != nil {
		return c.initOpts.appInfo
	}
	return nil
}

func (c *containerImpl) Logger() runtimelog.Logger {
	return c.logger
}

func (c *containerImpl) Component(name string, opts ...options.Option) (interfaces.Component, error) {
	if comp, ok := c.componentStore.GetInstance(name); ok {
		c.helper.Debugf("Component '%s' retrieved from cache.", name)
		return comp, nil
	}

	factory, ok := c.componentStore.GetFactory(name)
	if !ok {
		return nil, fmt.Errorf("component '%s' not found", name)
	}

	c.helper.Debugf("Component '%s' not found in cache, creating from factory...", name)
	comp, err := factory.NewComponent(c.config, c, opts...)
	if err != nil {
		c.helper.Errorf("Failed to create component '%s' from factory: %v", name, err)
		return nil, fmt.Errorf("failed to create component '%s' from factory: %w", name, err)
	}

	c.componentStore.RegisterInstance(name, comp)
	c.helper.Infof("Component '%s' created and registered successfully.", name)
	return comp, nil
}

func (c *containerImpl) RegisterComponent(name string, comp interfaces.Component) {
	c.helper.Infof("Registering pre-built component: %s", name)
	c.componentStore.RegisterInstance(name, comp)
}

func (c *containerImpl) RegisterFactory(name string, factory ComponentFactory) {
	c.helper.Infof("Registering component factory: %s", name)
	c.componentStore.RegisterFactory(name, factory)
}

func (c *containerImpl) HasComponent(name string) bool {
	return c.componentStore.Has(name)
}

func (c *containerImpl) RegisteredComponents() []string {
	return c.componentStore.List()
}

func (c *containerImpl) Registry(opts ...options.Option) (RegistryProvider, error) {
	c.helper.Debug("Accessing RegistryProvider...")
	discoveries, err := c.config.DecodeDiscoveries()
	if err != nil {
		c.helper.Errorf("Failed to decode discoveries config for RegistryProvider: %v", err)
	}
	c.registryProvider.Initialize(discoveries, opts...)
	return c.registryProvider, nil
}

func (c *containerImpl) Middleware(opts ...options.Option) (MiddlewareProvider, error) {
	c.helper.Debug("Accessing MiddlewareProvider...")
	middlewares, err := c.config.DecodeMiddlewares()
	if err != nil {
		c.helper.Errorf("Failed to decode middlewares config for MiddlewareProvider: %v", err)
	}
	c.middlewareProvider.Initialize(middlewares, opts...)
	return c.middlewareProvider, nil
}

func (c *containerImpl) Cache(opts ...options.Option) (CacheProvider, error) {
	c.helper.Debug("Accessing CacheProvider...")
	caches, err := c.config.DecodeCaches()
	if err != nil {
		c.helper.Errorf("Failed to decode caches config for CacheProvider: %v", err)
	}
	c.cacheProvider.Initialize(caches, opts...)
	return c.cacheProvider, nil
}

func (c *containerImpl) Database(opts ...options.Option) (DatabaseProvider, error) {
	c.helper.Debug("Accessing DatabaseProvider...")
	databases, err := c.config.DecodeDatabases()
	if err != nil {
		c.helper.Errorf("Failed to decode databases config for DatabaseProvider: %v", err)
	}
	c.databaseProvider.Initialize(databases, opts...)
	return c.databaseProvider, nil
}

func (c *containerImpl) ObjectStore(opts ...options.Option) (ObjectStoreProvider, error) {
	c.helper.Debug("Accessing ObjectStoreProvider...")
	filestores, err := c.config.DecodeObjectStores()
	if err != nil {
		c.helper.Errorf("Failed to decode object stores config for ObjectStoreProvider: %v", err)
	}
	c.objectStoreProvider.Initialize(filestores, opts...)
	return c.objectStoreProvider, nil
}

func (c *containerImpl) DefaultRegistrar() (runtimeregistry.KRegistrar, error) {
	c.helper.Debugf("Retrieving default registrar: '%s'", c.initOpts.defaultRegistrarName)
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

func (c *containerImpl) DefaultCache() (storageiface.Cache, error) {
	c.helper.Debugf("Retrieving default cache: '%s'", c.initOpts.defaultCacheName)
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

func (c *containerImpl) DefaultDatabase() (storageiface.Database, error) {
	c.helper.Debugf("Retrieving default database: '%s'", c.initOpts.defaultDatabaseName)
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

func (c *containerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	c.helper.Debugf("Retrieving default object store: '%s'", c.initOpts.defaultObjectStoreName)
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
