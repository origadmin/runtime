package container

import (
	"fmt"
	"sort"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/origadmin/runtime/container/internal/cache"
	"github.com/origadmin/runtime/container/internal/database"
	"github.com/origadmin/runtime/container/internal/middleware"
	"github.com/origadmin/runtime/container/internal/objectstore"
	"github.com/origadmin/runtime/container/internal/registry"
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
)

// Container defines the interface for accessing various runtime components.
type Container interface {
	Registry(opts ...options.Option) (RegistryProvider, error)
	Middleware(opts ...options.Option) (MiddlewareProvider, error)
	Cache(opts ...options.Option) (CacheProvider, error)
	Database(opts ...options.Option) (DatabaseProvider, error)
	ObjectStore(opts ...options.Option) (ObjectStoreProvider, error)
	Components(opts ...options.Option) (map[string]interfaces.Component, error)
	Component(name string) (interfaces.Component, error)
	RegisterComponent(name string, comp interfaces.Component)
	Logger() log.Logger
	AppInfo() interfaces.AppInfo
	WithOptions(opts ...options.Option) Container
}

// containerImpl implements the Container interface.
type containerImpl struct {
	mu                  sync.RWMutex
	config              interfaces.StructuredConfig
	logger              log.Logger
	appInfo             interfaces.AppInfo
	opts                []options.Option
	componentFactories  map[string]ComponentFactory
	cachedComponents    map[string]interfaces.Component
	componentsErr       error
	componentsOnce      sync.Once
	middlewareProvider  *middleware.Provider
	middlewareOnce      sync.Once
	middlewareErr       error
	cacheProvider       *cache.Provider
	cacheOnce           sync.Once
	cacheErr            error
	databaseProvider    *database.Provider
	databaseOnce        sync.Once
	databaseErr         error
	registryProvider    *registry.Provider
	registryOnce        sync.Once
	registryErr         error
	objectStoreProvider *objectstore.Provider
	objectStoreOnce     sync.Once
	objectStoreErr      error
}

// New creates a new Container instance.
func New(config interfaces.StructuredConfig, opts ...options.Option) Container {
	co := optionutil.NewT[containerOptions](opts...)

	var baseLogger log.Logger
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		log.NewHelper(log.DefaultLogger).Warnf("failed to decode logger config, using default logger: %v", err)
		loggerConfig = nil
	}
	baseLogger = runtimelog.NewLogger(loggerConfig)

	enrichedLogger := baseLogger
	if co.appInfo != nil {
		enrichedLogger = log.With(baseLogger,
			"service.name", co.appInfo.Name(),
			"service.version", co.appInfo.Version(),
			"service.id", co.appInfo.ID(),
			"trace.id", tracing.TraceID(),
			"span.id", tracing.SpanID(),
		)
	}

	return &containerImpl{
		config:              config,
		logger:              enrichedLogger,
		appInfo:             co.appInfo,
		opts:                opts,
		componentFactories:  co.componentFactories,
		cachedComponents:    make(map[string]interfaces.Component),
		middlewareProvider:  middleware.NewProvider(enrichedLogger),
		cacheProvider:       cache.NewProvider(enrichedLogger),
		databaseProvider:    database.NewProvider(enrichedLogger),
		registryProvider:    registry.NewProvider(enrichedLogger),
		objectStoreProvider: objectstore.NewProvider(enrichedLogger),
	}
}

func (c *containerImpl) WithOptions(opts ...options.Option) Container {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.opts = append(c.opts, opts...)
	co := optionutil.NewT[containerOptions](c.opts...)

	if c.componentFactories == nil {
		c.componentFactories = make(map[string]ComponentFactory)
	}
	for name, factory := range co.componentFactories {
		c.componentFactories[name] = factory
	}
	if co.appInfo != nil {
		c.appInfo = co.appInfo
	}
	return c
}

func (c *containerImpl) AppInfo() interfaces.AppInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.appInfo
}

func (c *containerImpl) Logger() log.Logger {
	return c.logger
}

func (c *containerImpl) initializeComponents(opts ...options.Option) {
	logHelper := log.NewHelper(c.logger)
	factories := make([]ComponentFactory, 0, len(c.componentFactories))
	for _, factory := range c.componentFactories {
		factories = append(factories, factory)
	}
	sort.Slice(factories, func(i, j int) bool {
		return factories[i].Priority() < factories[j].Priority()
	})

	finalOpts := append(append([]options.Option{}, c.opts...), opts...)
	for _, factory := range factories {
		logHelper.Infof("executing component factory with priority %d", factory.Priority())
		if _, err := factory.NewComponent(c.config, c, finalOpts...); err != nil {
			c.componentsErr = fmt.Errorf("error executing factory with priority %d: %w", factory.Priority(), err)
			logHelper.Errorf("halting component initialization due to error: %v", c.componentsErr)
			return
		}
	}
}

func (c *containerImpl) Components(opts ...options.Option) (map[string]interfaces.Component, error) {
	c.componentsOnce.Do(func() {
		c.initializeComponents(opts...)
	})
	if c.componentsErr != nil {
		return nil, c.componentsErr
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	maps := make(map[string]interfaces.Component, len(c.cachedComponents))
	for k, v := range c.cachedComponents {
		maps[k] = v
	}
	return maps, nil
}

func (c *containerImpl) Component(name string) (interfaces.Component, error) {
	if _, err := c.Components(); err != nil {
		return nil, err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	comp, ok := c.cachedComponents[name]
	if !ok {
		return nil, fmt.Errorf("component '%s' not found", name)
	}
	return comp, nil
}

func (c *containerImpl) RegisterComponent(name string, comp interfaces.Component) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.cachedComponents[name]; exists {
		log.NewHelper(c.logger).Warnf("component with name '%s' is being overwritten", name)
	}
	c.cachedComponents[name] = comp
}

func (c *containerImpl) Registry(opts ...options.Option) (RegistryProvider, error) {
	c.registryOnce.Do(func() {
		discoveries, err := c.config.DecodeDiscoveries()
		if err != nil {
			c.registryErr = err
			return
		}
		finalOpts := append(append([]options.Option{}, c.opts...), opts...)
		c.registryProvider.SetConfig(discoveries, finalOpts...)
	})
	return c.registryProvider, c.registryErr
}

func (c *containerImpl) Middleware(opts ...options.Option) (MiddlewareProvider, error) {
	c.middlewareOnce.Do(func() {
		middlewares, err := c.config.DecodeMiddlewares()
		if err != nil {
			c.middlewareErr = err
			return
		}
		finalOpts := append(append([]options.Option{}, c.opts...), opts...)
		c.middlewareProvider.SetConfig(middlewares, finalOpts...)
	})
	return c.middlewareProvider, c.middlewareErr
}

func (c *containerImpl) Cache(opts ...options.Option) (CacheProvider, error) {
	c.cacheOnce.Do(func() {
		caches, err := c.config.DecodeCaches()
		if err != nil {
			c.cacheErr = err
			return
		}
		finalOpts := append(append([]options.Option{}, c.opts...), opts...)
		c.cacheProvider.SetConfig(caches, finalOpts...)
	})
	return c.cacheProvider, c.cacheErr
}

func (c *containerImpl) Database(opts ...options.Option) (DatabaseProvider, error) {
	c.databaseOnce.Do(func() {
		databases, err := c.config.DecodeDatabases()
		if err != nil {
			c.databaseErr = err
			return
		}
		finalOpts := append(append([]options.Option{}, c.opts...), opts...)
		c.databaseProvider.SetConfig(databases, finalOpts...)
	})
	return c.databaseProvider, c.databaseErr
}

func (c *containerImpl) ObjectStore(opts ...options.Option) (ObjectStoreProvider, error) {
	c.objectStoreOnce.Do(func() {
		filestores, err := c.config.DecodeObjectStores()
		if err != nil {
			c.objectStoreErr = err
			return
		}
		finalOpts := append(append([]options.Option{}, c.opts...), opts...)
		c.objectStoreProvider.SetConfig(filestores, finalOpts...)
	})
	return c.objectStoreProvider, c.objectStoreErr
}
