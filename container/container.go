package container

import (
	"errors"
	"fmt"
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

	Components() (map[string]interfaces.Component, error)
	Component(name string) (interfaces.Component, error)
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

// Components loads components from config on first call, then returns the cache.
// This version uses a map-based configuration, which is more intuitive.
func (c *containerImpl) Components() (map[string]interfaces.Component, error) {
	var allErrors error
	logHelper := log.NewHelper(c.logger)

	c.onceComponents.Do(func() {
		// Decode the 'components' key into a map where keys are component names.
		componentConfigs := make(map[string]interfaces.StructuredConfig)
		if err := c.config.ScanKey("components", &componentConfigs); err != nil {
			logHelper.Errorf("failed to decode 'components' map: %v", err)
			allErrors = errors.Join(allErrors, fmt.Errorf("failed to decode 'components' map: %w", err))
			return
		}

		if len(componentConfigs) == 0 {
			logHelper.Info("no component configurations found under 'components' key")
			return
		}

		c.componentMu.Lock()
		defer c.componentMu.Unlock()

		for name, compConfig := range componentConfigs {
			// The framework needs to know the 'type' to find the factory.
			var typeHolder struct {
				Type string `json:"type"`
			}
			if err := compConfig.Scan(&typeHolder); err != nil {
				err = fmt.Errorf("failed to scan 'type' for component '%s': %w", name, err)
				logHelper.Error(err)
				allErrors = errors.Join(allErrors, err)
				continue
			}

			if typeHolder.Type == "" {
				err := fmt.Errorf("component '%s' is missing a 'type' field", name)
				logHelper.Error(err)
				allErrors = errors.Join(allErrors, err)
				continue
			}

			factory, ok := c.componentFactories[typeHolder.Type]
			if !ok {
				err := fmt.Errorf("component factory for type '%s' (component '%s') not found", typeHolder.Type, name)
				logHelper.Error(err)
				allErrors = errors.Join(allErrors, err)
				continue
			}

			// The factory gets its own structured config and the container instance.
			// The factory is responsible for parsing what it needs from the config.
			comp, err := factory.NewComponent(compConfig, c, c.opts...)
			if err != nil {
				err = fmt.Errorf("failed to create component '%s' with type '%s': %w", name, typeHolder.Type, err)
				logHelper.Error(err)
				allErrors = errors.Join(allErrors, err)
				continue
			}
			c.cachedComponents[name] = comp
		}
	})

	c.componentMu.RLock()
	defer c.componentMu.RUnlock()
	// Return a copy to prevent modification of the internal map
	maps := make(map[string]interfaces.Component, len(c.cachedComponents))
	for k, v := range c.cachedComponents {
		maps[k] = v
	}

	return maps, allErrors
}

// Component returns a specific component by name, loading all if not already loaded.
func (c *containerImpl) Component(name string) (interfaces.Component, error) {
	// This ensures the lazy-loading from config is triggered if it hasn't been already.
	if _, err := c.Components(); err != nil {
		// A specific component might have been created successfully even if others failed.
		// We proceed to check the cache regardless of the aggregate error.
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
	c.cachedComponents[name] = comp
}
