package container

import (
	"fmt"
	"sort"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/container/internal/cache"
	"github.com/origadmin/runtime/container/internal/database"
	"github.com/origadmin/runtime/container/internal/middleware"
	"github.com/origadmin/runtime/container/internal/objectstore"
	"github.com/origadmin/runtime/container/internal/registry"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
)

// Container defines the interface for accessing various runtime components.
// It elevates Component management to a first-class concern.
type Container interface {
	Registry(opts ...options.Option) (RegistryProvider, error)
	Middleware(opts ...options.Option) (MiddlewareProvider, error)
	Cache(opts ...options.Option) (CacheProvider, error)
	Database(opts ...options.Option) (DatabaseProvider, error)
	ObjectStore(opts ...options.Option) (ObjectStoreProvider, error)
	// Components returns all initialized components.
	// The first call will trigger the initialization process.
	Components(opts ...options.Option) (map[string]interfaces.Component, error)
	// Component returns a specific component by name.
	// The first call will trigger the initialization process.
	Component(name string) (interfaces.Component, error)
	// RegisterComponent allows for dynamic, programmatic registration of components.
	// This is the primary mechanism for factories to add created components to the container.
	RegisterComponent(name string, comp interfaces.Component)

	// Logger returns the configured logger instance.
	Logger() log.Logger
}

// containerImpl implements the Container interface with lazy loading and caching.
type containerImpl struct {
	config interfaces.StructuredConfig
	logger log.Logger
	opts   []options.Option // Static options from NewContainer

	// Component-related fields are now directly in the container
	componentFactories map[string]ComponentFactory
	cachedComponents   map[string]interfaces.Component
	initErr            error
	onceComponents     sync.Once
	componentMu        sync.RWMutex

	// Providers for other services
	middlewareProvider  *middleware.Provider
	cacheProvider       *cache.Provider
	databaseProvider    *database.Provider
	registryProvider    *registry.Provider
	objectStoreProvider *objectstore.Provider
}

// NewContainer creates a new Container instance with the given configuration.
func NewContainer(config interfaces.StructuredConfig, opts ...options.Option) Container {
	componentFactories := ComponentFactoryFromOptions(opts...)
	var logger log.Logger
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		// Log the error using a default logger, then proceed with default logger
		log.DefaultLogger.Log(log.LevelError, "failed to decode logger config, using default logger: %v", err)
		// Set loggerConfig to nil so runtimelog.NewLogger returns DefaultLogger
		loggerConfig = nil
	}

	logger = runtimelog.NewLogger(loggerConfig)

	impl := &containerImpl{
		config:              config,
		logger:              logger,
		opts:                opts,
		componentFactories:  componentFactories,
		cachedComponents:    make(map[string]interfaces.Component),
		middlewareProvider:  middleware.NewProvider(logger),
		cacheProvider:       cache.NewProvider(logger),
		databaseProvider:    database.NewProvider(logger),
		registryProvider:    registry.NewProvider(logger),
		objectStoreProvider: objectstore.NewProvider(logger),
	}

	return impl
}

// initializeComponents is the core logic for component initialization.
func (c *containerImpl) initializeComponents(opts ...options.Option) {
	logHelper := log.NewHelper(c.logger)

	factories := make([]ComponentFactory, 0, len(c.componentFactories))
	for _, factory := range c.componentFactories {
		factories = append(factories, factory)
	}

	sort.Slice(factories, func(i, j int) bool {
		return factories[i].Priority() < factories[j].Priority()
	})

	// Merge static opts from NewContainer and dynamic opts from Components() call
	finalOpts := append(append([]options.Option{}, c.opts...), opts...)

	for _, factory := range factories {
		logHelper.Infof("executing component factory with priority %d", factory.Priority())
		if _, err := factory.NewComponent(c.config, c, finalOpts...); err != nil {
			c.initErr = fmt.Errorf("error executing factory with priority %d: %w", factory.Priority(), err)
			logHelper.Errorf("halting component initialization due to error: %v", c.initErr)
			return
		}
	}
}

// Components returns all initialized components.
func (c *containerImpl) Components(opts ...options.Option) (map[string]interfaces.Component, error) {
	c.onceComponents.Do(func() {
		c.initializeComponents(opts...) // Pass opts to the actual initialization logic
	})
	if c.initErr != nil {
		return nil, c.initErr
	}

	c.componentMu.RLock()
	defer c.componentMu.RUnlock()
	maps := make(map[string]interfaces.Component, len(c.cachedComponents))
	for k, v := range c.cachedComponents {
		maps[k] = v
	}
	return maps, nil
}

// Component returns a specific component by name.
func (c *containerImpl) Component(name string) (interfaces.Component, error) {
	// Component() call will trigger initializeComponents() if not already done
	if _, err := c.Components(); err != nil {
		return nil, err
	}

	c.componentMu.RLock()
	defer c.componentMu.RUnlock()

	comp, ok := c.cachedComponents[name]
	if !ok {
		return nil, fmt.Errorf("component '%s' not found", name)
	}
	return comp, nil
}

// RegisterComponent allows for dynamic, programmatic registration of components.
func (c *containerImpl) RegisterComponent(name string, comp interfaces.Component) {
	c.componentMu.Lock()
	defer c.componentMu.Unlock()
	if _, exists := c.cachedComponents[name]; exists {
		log.NewHelper(c.logger).Warnf("component with name '%s' is being overwritten", name)
	}
	c.cachedComponents[name] = comp
}

// Logger returns the configured logger instance.
func (c *containerImpl) Logger() log.Logger {
	return c.logger
}

// Registry implements Container.
func (c *containerImpl) Registry(opts ...options.Option) (RegistryProvider, error) {
	c.onceComponents.Do(func() {
		c.initializeComponents()
	})
	if c.initErr != nil {
		return nil, c.initErr
	}
	discoveries, err := c.config.DecodeDiscoveries()
	if err != nil {
		return nil, err
	}
	finalOpts := append(append([]options.Option{}, c.opts...), opts...)
	c.registryProvider.SetConfig(discoveries, finalOpts...)
	return c.registryProvider, nil
}

// Middleware implements Container.
func (c *containerImpl) Middleware(opts ...options.Option) (MiddlewareProvider, error) {
	c.onceComponents.Do(func() {
		c.initializeComponents()
	})
	if c.initErr != nil {
		return nil, c.initErr
	}
	middlewares, err := c.config.DecodeMiddlewares()
	if err != nil {
		return nil, err
	}
	finalOpts := append(append([]options.Option{}, c.opts...), opts...)
	c.middlewareProvider.SetConfig(middlewares, finalOpts...)
	return c.middlewareProvider, nil
}

// Cache implements Container.
func (c *containerImpl) Cache(opts ...options.Option) (CacheProvider, error) {
	c.onceComponents.Do(func() {
		c.initializeComponents()
	})
	if c.initErr != nil {
		return nil, c.initErr
	}
	caches, err := c.config.DecodeCaches()
	if err != nil {
		return nil, err
	}
	finalOpts := append(append([]options.Option{}, c.opts...), opts...)
	c.cacheProvider.SetConfig(caches, finalOpts...)
	return c.cacheProvider, nil
}

// Database implements Container.
func (c *containerImpl) Database(opts ...options.Option) (DatabaseProvider, error) {
	c.onceComponents.Do(func() {
		c.initializeComponents()
	})
	if c.initErr != nil {
		return nil, c.initErr
	}
	databases, err := c.config.DecodeDatabases()
	if err != nil {
		return nil, err
	}
	finalOpts := append(append([]options.Option{}, c.opts...), opts...)
	c.databaseProvider.SetConfig(databases, finalOpts...)
	return c.databaseProvider, nil
}

// ObjectStore implements Container.
func (c *containerImpl) ObjectStore(opts ...options.Option) (ObjectStoreProvider, error) {
	c.onceComponents.Do(func() {
		c.initializeComponents()
	})
	if c.initErr != nil {
		return nil, c.initErr
	}
	filestores, err := c.config.DecodeObjectStores()
	if err != nil {
		return nil, err
	}
	finalOpts := append(append([]options.Option{}, c.opts...), opts...)
	c.objectStoreProvider.SetConfig(filestores, finalOpts...)
	return c.objectStoreProvider, nil
}
