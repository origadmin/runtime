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

// containerImpl implements the interfaces.Container interface with lazy loading and caching.
type containerImpl struct {
	config containerInterfaces.StructuredConfig
	logger log.Logger
	opts   []options.Option

	// Generic component factories, which should be passed in.
	componentFactories containerInterfaces.ComponentFactory

	middlewareProvider  *middleware.Provider
	cacheProvider       *cache.Provider
	databaseProvider    *database.Provider
	registryProvider    *registry.Provider
	objectStoreProvider *objectstore.Provider
	componentProvider   *component.Provider
}

// NewContainer creates a new Container instance with the given configuration.
func NewContainer(config containerInterfaces.StructuredConfig, genericFactories containerInterfaces.GenericComponentFactory, opts ...options.Option) interfaces.Container {
	// Instantiate Logger internally
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		// Log the error using a default logger, then proceed with default logger
		log.DefaultLogger.Log(log.LevelError, "failed to decode logger config, using default logger: %v", err)
		return &containerImpl{
			config:                    config,
			logger:                    log.DefaultLogger,
			opts:                      opts,
			middlewareProvider:        middleware.NewProvider(log.DefaultLogger, opts),
			cacheProvider:             cache.NewProvider(log.DefaultLogger, opts),
			databaseProvider:          database.NewProvider(log.DefaultLogger, opts),
			registryProvider:          registry.NewProvider(log.DefaultLogger, opts),
			objectStoreProvider:       objectstore.NewProvider(log.DefaultLogger, opts),
			componentProvider:         component.NewProvider(log.DefaultLogger, opts, genericFactories),
			genericComponentFactories: genericFactories,
		}
	}

	logger := runtimelog.NewLogger(loggerConfig)

	impl := &containerImpl{
		config:                    config,
		logger:                    logger,
		opts:                      opts,
		middlewareProvider:        middleware.NewProvider(logger, opts),
		cacheProvider:             cache.NewProvider(logger, opts),
		databaseProvider:          database.NewProvider(logger, opts),
		registryProvider:          registry.NewProvider(logger, opts),
		objectStoreProvider:       objectstore.NewProvider(logger, opts),
		genericComponentFactories: genericFactories,
	}

	return impl
}

// Registry implements Container.
func (c *containerImpl) Registry() (containerInterfaces.RegistryProvider, error) {
	discoveries, err := c.config.DecodeDiscoveries()
	if err != nil {
		return nil, err
	}
	c.registryProvider.SetConfig(discoveries)
	return c.registryProvider, nil
}

// Middleware implements Container.
func (c *containerImpl) Middleware() (containerInterfaces.MiddlewareProvider, error) {
	middlewares, err := c.config.DecodeMiddlewares()
	if err != nil {
		return nil, err
	}
	c.middlewareProvider.SetConfig(middlewares)
	return c.middlewareProvider, nil
}

// Cache implements Container.
func (c *containerImpl) Cache() (containerInterfaces.CacheProvider, error) {
	caches, err := c.config.DecodeCaches()
	if err != nil {
		return nil, err
	}
	c.cacheProvider.SetConfig(caches)
	return c.cacheProvider, nil
}

// Database implements Container.
func (c *containerImpl) Database() (containerInterfaces.DatabaseProvider, error) {
	databases, err := c.config.DecodeDatabases()
	if err != nil {
		return nil, err
	}
	c.databaseProvider.SetConfig(databases)
	return c.databaseProvider, nil
}

// ObjectStore implements Container.
func (c *containerImpl) ObjectStore() (containerInterfaces.ObjectStoreProvider, error) {
	filestores, err := c.config.DecodeFilestores()
	if err != nil {
		return nil, err
	}
	c.objectStoreProvider.SetConfig(filestores)
	return c.objectStoreProvider, nil
}

// Component implements Container.
func (c *containerImpl) Component() (containerInterfaces.ComponentProvider, error) {
	c.componentOnce.Do(func() {
		c.cachedComponentProvider = internal.NewComponentProvider(c.config, c.logger, c.opts, c.genericComponentFactories)
	})
	return c.cachedComponentProvider, c.componentErr
}
