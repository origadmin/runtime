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
	Registry() (RegistryProvider, error)
	Middleware() (MiddlewareProvider, error)
	Cache() (CacheProvider, error)
	Database() (DatabaseProvider, error)
	ObjectStore() (ObjectStoreProvider, error)

	// Components returns all initialized components.
	// The first call will trigger the initialization process.
	Components() (map[string]interfaces.Component, error)
	// Component returns a specific component by name.
	// The first call will trigger the initialization process.
	Component(name string) (interfaces.Component, error)
	// RegisterComponent allows for dynamic, programmatic registration of components.
	// This is the primary mechanism for factories to add created components to the container.
	RegisterComponent(name string, comp interfaces.Component)
}

// containerImpl implements the Container interface with lazy loading and caching.
type containerImpl struct {
	config interfaces.StructuredConfig
	logger log.Logger
	opts   []options.Option

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
func NewContainer(config interfaces.StructuredConfig, componentFactories map[string]ComponentFactory, opts ...options.Option) Container {
	// Instantiate Logger internally
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		// Log the error using a default logger, then proceed with default logger
		log.DefaultLogger.Log(log.LevelError, "failed to decode logger config, using default logger: %v", err)
		logger := log.DefaultLogger
		return &containerImpl{
			config:              config,
			logger:              logger,
			opts:                opts,
			componentFactories:  componentFactories,
			cachedComponents:    make(map[string]interfaces.Component),
			middlewareProvider:  middleware.NewProvider(logger, opts),
			cacheProvider:       cache.NewProvider(logger, opts),
			databaseProvider:    database.NewProvider(logger, opts),
			registryProvider:    registry.NewProvider(logger, opts),
			objectStoreProvider: objectstore.NewProvider(logger, opts),
		}
	}

	logger := runtimelog.NewLogger(loggerConfig)

	impl := &containerImpl{
		config:              config,
		logger:              logger,
		opts:                opts,
		componentFactories:  componentFactories,
		cachedComponents:    make(map[string]interfaces.Component),
		middlewareProvider:  middleware.NewProvider(logger, opts),
		cacheProvider:       cache.NewProvider(logger, opts),
		databaseProvider:    database.NewProvider(logger, opts),
		registryProvider:    registry.NewProvider(logger, opts),
		objectStoreProvider: objectstore.NewProvider(logger, opts),
	}

	return impl
}

// initializeComponents is the core logic for component initialization.
// The container acts as a coordinator, having zero knowledge of the configuration structure.
// It sorts the registered factories by priority and then executes them.
func (c *containerImpl) initializeComponents() {
	logHelper := log.NewHelper(c.logger)

	// 1. Create a slice of factories to be sorted.
	factories := make([]ComponentFactory, 0, len(c.componentFactories))
	for _, factory := range c.componentFactories {
		factories = append(factories, factory)
	}

	// 2. Sort the factories by priority (lower first).
	sort.Slice(factories, func(i, j int) bool {
		return factories[i].Priority() < factories[j].Priority()
	})

	// 3. Iterate through the sorted factories and execute them.
	// The container does not parse any configuration. It passes the root config
	// to each factory, and the factory is responsible for finding its own config,
	// creating components, and registering them back to the container.
	for _, factory := range factories {
		logHelper.Infof("executing component factory with priority %d", factory.Priority())

		// The container IGNORES the returned component from NewComponent.
		// The factory is expected to register all created components via container.RegisterComponent.
		if _, err := factory.NewComponent(c.config, c, c.opts...); err != nil {
			c.initErr = fmt.Errorf("error executing factory with priority %d: %w", factory.Priority(), err)
			logHelper.Errorf("halting component initialization due to error: %v", c.initErr)
			return // Stop on first error to prevent dependency issues.
		}
	}
}

// Components returns all initialized components.
func (c *containerImpl) Components() (map[string]interfaces.Component, error) {
	c.onceComponents.Do(c.initializeComponents)
	if c.initErr != nil {
		return nil, c.initErr
	}

	c.componentMu.RLock()
	defer c.componentMu.RUnlock()
	// Return a copy to prevent modification of the internal map
	maps := make(map[string]interfaces.Component, len(c.cachedComponents))
	for k, v := range c.cachedComponents {
		maps[k] = v
	}
	return maps, nil
}

// Component returns a specific component by name.
func (c *containerImpl) Component(name string) (interfaces.Component, error) {
	c.onceComponents.Do(c.initializeComponents)
	if c.initErr != nil {
		return nil, c.initErr
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
