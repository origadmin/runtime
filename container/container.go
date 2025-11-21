package container

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/container/internal/cache"
	"github.com/origadmin/runtime/container/internal/component"
	"github.com/origadmin/runtime/container/internal/database"
	"github.com/origadmin/runtime/container/internal/middleware"
	"github.com/origadmin/runtime/container/internal/objectstore"
	"github.com/origadmin/runtime/container/internal/registry"
	"github.com/origadmin/runtime/interfaces"
	containerInterfaces "github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
)

// Container defines the interface for accessing various runtime components.
type Container interface {
	Registry() (RegistryProvider, error)
	Middleware() (MiddlewareProvider, error)
	Cache() (CacheProvider, error)
	Database() (DatabaseProvider, error)
	ObjectStore() (ObjectStoreProvider, error)
	Component() (ComponentProvider, error)
}

// containerImpl implements the interfaces.Container interface with lazy loading and caching.
type containerImpl struct {
	config containerInterfaces.StructuredConfig
	logger log.Logger
	opts   []options.Option

	componentFactories map[string]ComponentFactory

	middlewareProvider  *middleware.Provider
	cacheProvider       *cache.Provider
	databaseProvider    *database.Provider
	registryProvider    *registry.Provider
	objectStoreProvider *objectstore.Provider
	componentProvider   *component.Provider
}

// NewContainer creates a new Container instance with the given configuration.
func NewContainer(config containerInterfaces.StructuredConfig, componentFactories map[string]ComponentFactory, opts ...options.Option) Container {
	// Instantiate Logger internally
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		// Log the error using a default logger, then proceed with default logger
		log.DefaultLogger.Log(log.LevelError, "failed to decode logger config, using default logger: %v", err)
		return &containerImpl{
			config:              config,
			logger:              log.DefaultLogger,
			opts:                opts,
			middlewareProvider:  middleware.NewProvider(log.DefaultLogger, opts),
			cacheProvider:       cache.NewProvider(log.DefaultLogger, opts),
			databaseProvider:    database.NewProvider(log.DefaultLogger, opts),
			registryProvider:    registry.NewProvider(log.DefaultLogger, opts),
			objectStoreProvider: objectstore.NewProvider(log.DefaultLogger, opts),
			componentProvider:   component.NewProvider(log.DefaultLogger, opts, componentFactories),
		}
	}

	logger := runtimelog.NewLogger(loggerConfig)

	impl := &containerImpl{
		config:              config,
		logger:              logger,
		opts:                opts,
		middlewareProvider:  middleware.NewProvider(logger, opts),
		cacheProvider:       cache.NewProvider(logger, opts),
		databaseProvider:    database.NewProvider(logger, opts),
		registryProvider:    registry.NewProvider(logger, opts),
		objectStoreProvider: objectstore.NewProvider(logger, opts),
		componentProvider:   component.NewProvider(logger, opts, componentFactories),
	}

	return impl
}

// Registry implements Container.
func (c *containerImpl) Registry() (RegistryProvider, error) {
	discoveries, err := c.config.DecodeDiscoveries()
	if err != nil {
		return nil, err
	}
	c.registryProvider.SetConfig(discoveries)
	return c.registryProvider, nil
}

// Middleware implements Container.
func (c *containerImpl) Middleware() (MiddlewareProvider, error) {
	middlewares, err := c.config.DecodeMiddlewares()
	if err != nil {
		return nil, err
	}
	c.middlewareProvider.SetConfig(middlewares)
	return c.middlewareProvider, nil
}

// Cache implements Container.
func (c *containerImpl) Cache() (CacheProvider, error) {
	caches, err := c.config.DecodeCaches()
	if err != nil {
		return nil, err
	}
	c.cacheProvider.SetConfig(caches)
	return c.cacheProvider, nil
}

// Database implements Container.
func (c *containerImpl) Database() (DatabaseProvider, error) {
	databases, err := c.config.DecodeDatabases()
	if err != nil {
		return nil, err
	}
	c.databaseProvider.SetConfig(databases)
	return c.databaseProvider, nil
}

// ObjectStore implements Container.
func (c *containerImpl) ObjectStore() (ObjectStoreProvider, error) {
	filestores, err := c.config.DecodeObjectStores()
	if err != nil {
		return nil, err
	}
	c.objectStoreProvider.SetConfig(filestores)
	return c.objectStoreProvider, nil
}

// Component implements Container.
func (c *containerImpl) Component() ComponentProvider {
	return c.componentProvider
}
