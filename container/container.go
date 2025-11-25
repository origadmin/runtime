package container

import (
	"fmt"
	"sort"
	"sync"

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
// The configuration of the container is finalized upon its creation.
type Container interface {
	Registry(opts ...options.Option) (RegistryProvider, error)
	Middleware(opts ...options.Option) (MiddlewareProvider, error)
	Cache(opts ...options.Option) (CacheProvider, error)
	Database(opts ...options.Option) (DatabaseProvider, error)
	ObjectStore(opts ...options.Option) (ObjectStoreProvider, error)
	Components(opts ...options.Option) (map[string]interfaces.Component, error)
	Component(name string) (interfaces.Component, error)
	RegisterComponent(name string, comp interfaces.Component)
	Logger() runtimelog.Logger
	AppInfo() interfaces.AppInfo

	DefaultCache() (storageiface.Cache, error)
	DefaultDatabase() (storageiface.Database, error)
	DefaultObjectStore() (storageiface.ObjectStore, error)
	DefaultRegistrar() (runtimeregistry.KRegistrar, error)
}

// containerImpl implements the Container interface.
type containerImpl struct {
	mu                  sync.RWMutex
	config              interfaces.StructuredConfig
	logger              runtimelog.Logger
	appInfo             interfaces.AppInfo
	opts                []options.Option // Stored options from New, used for sub-provider initialization
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
	// Global default names from options, highest priority
	globalDefaultCacheName       string
	globalDefaultDatabaseName    string
	globalDefaultObjectStoreName string
	globalDefaultRegistrarName   string
}

// New creates a new Container instance. All configuration should be passed here.
func New(config interfaces.StructuredConfig, opts ...options.Option) Container {
	co := optionutil.NewT[containerOptions](opts...)

	var baseLogger runtimelog.Logger
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		runtimelog.NewHelper(runtimelog.DefaultLogger).Warnf("failed to decode logger config, using default logger: %v", err)
		loggerConfig = nil
	}
	baseLogger = runtimelog.NewLogger(loggerConfig)

	enrichedLogger := baseLogger
	if co.appInfo != nil {
		// Removed kratos tracing specific fields, use runtimelog's With for general context
		enrichedLogger = runtimelog.With(baseLogger,
			"service.name", co.appInfo.Name(),
			"service.version", co.appInfo.Version(),
			"service.id", co.appInfo.ID(),
		)
	}

	return &containerImpl{
		config:              config,
		logger:              enrichedLogger,
		appInfo:             co.appInfo,
		opts:                opts, // Store the initial options for sub-provider initialization
		componentFactories:  co.componentFactories,
		cachedComponents:    make(map[string]interfaces.Component),
		middlewareProvider:  middleware.NewProvider(enrichedLogger),
		cacheProvider:       cache.NewProvider(enrichedLogger),
		databaseProvider:    database.NewProvider(enrichedLogger),
		registryProvider:    registry.NewProvider(enrichedLogger),
		objectStoreProvider: objectstore.NewProvider(enrichedLogger),
		// Initialize global default names from options
		globalDefaultCacheName:       co.defaultCacheName,
		globalDefaultDatabaseName:    co.defaultDatabaseName,
		globalDefaultObjectStoreName: co.defaultObjectStoreName,
		globalDefaultRegistrarName:   co.defaultRegistrarName,
	}
}

func (c *containerImpl) AppInfo() interfaces.AppInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.appInfo
}

func (c *containerImpl) Logger() runtimelog.Logger {
	return c.logger
}

func (c *containerImpl) initializeComponents(opts ...options.Option) {
	logHelper := runtimelog.NewHelper(c.logger)
	factories := make([]ComponentFactory, 0, len(c.componentFactories))
	for _, factory := range c.componentFactories {
		factories = append(factories, factory)
	}
	sort.Slice(factories, func(i, j int) bool {
		return factories[i].Priority() < factories[j].Priority()
	})

	// Combine initial options with any specific options for Components()
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
	c.mu.Lock()
	defer c.mu.Unlock()
	maps := make(map[string]interfaces.Component, len(c.cachedComponents))
	for k, v := range c.cachedComponents {
		maps[k] = v
	}
	return maps, nil
}

func (c *containerImpl) Component(name string) (interfaces.Component, error) {
	// Ensure all components are initialized before trying to access one.
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
		runtimelog.NewHelper(c.logger).Warnf("component with name '%s' is being overwritten", name)
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

// DefaultRegistrar returns the default registrar, prioritizing global options, then config.
func (c *containerImpl) DefaultRegistrar() (runtimeregistry.KRegistrar, error) {
	provider, err := c.Registry() // Ensure provider is initialized
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	globalName := c.globalDefaultRegistrarName
	c.mu.RUnlock()

	return provider.DefaultRegistrar(globalName)
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

// DefaultCache returns the default cache, prioritizing global options, then config.
func (c *containerImpl) DefaultCache() (storageiface.Cache, error) {
	provider, err := c.Cache() // Ensure provider is initialized
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	globalName := c.globalDefaultCacheName
	c.mu.RUnlock()

	return provider.DefaultCache(globalName)
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

// DefaultDatabase returns the default database, prioritizing global options, then config.
func (c *containerImpl) DefaultDatabase() (storageiface.Database, error) {
	provider, err := c.Database() // Ensure provider is initialized
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	globalName := c.globalDefaultDatabaseName
	c.mu.RUnlock()

	return provider.DefaultDatabase(globalName)
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

// DefaultObjectStore returns the default object store, prioritizing global options, then config.
func (c *containerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	provider, err := c.ObjectStore() // Ensure provider is initialized
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	globalName := c.globalDefaultObjectStoreName
	c.mu.RUnlock()

	return provider.DefaultObjectStore(globalName)
}
