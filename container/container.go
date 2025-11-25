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
type Container interface {
	Registry(opts ...options.Option) (RegistryProvider, error)
	Middleware(opts ...options.Option) (MiddlewareProvider, error)
	Cache(opts ...options.Option) (CacheProvider, error)
	Database(opts ...options.Option) (DatabaseProvider, error)
	ObjectStore(opts ...options.Option) (ObjectStoreProvider, error)
	Component(name string, opts ...options.Option) (interfaces.Component, error)
	RegisterComponent(name string, comp interfaces.Component)
	Logger() runtimelog.Logger
	AppInfo() interfaces.AppInfo

	DefaultCache() (storageiface.Cache, error)
	DefaultDatabase() (storageiface.Database, error)
	DefaultObjectStore() (storageiface.ObjectStore, error)
	DefaultRegistrar() (runtimeregistry.KRegistrar, error)
}

// containerImpl implements the Container interface.
// WARNING: This implementation is NOT CONCURRENCY-SAFE for operations that modify
// or lazily initialize internal state after the New function returns.
// Concurrent access to methods like Component, RegisterComponent, and the provider
// accessors (Registry, Cache, etc.) without external synchronization will lead to
// race conditions, data corruption, or unpredictable behavior.
type containerImpl struct {
	config   interfaces.StructuredConfig
	logger   runtimelog.Logger
	appInfo  interfaces.AppInfo
	initOpts *containerOptions // Options used during initialization.

	// Eagerly initialized concrete provider instances
	registryProvider    *registry.Provider
	middlewareProvider  *middleware.Provider
	cacheProvider       *cache.Provider
	databaseProvider    *database.Provider
	objectStoreProvider *objectstore.Provider

	// Flags to track if structural configuration has been applied to providers
	registryConfigured    bool
	middlewareConfigured  bool
	cacheConfigured       bool
	databaseConfigured    bool
	objectStoreConfigured bool

	// Components map (still not concurrency-safe)
	components map[string]interfaces.Component
}

// New creates a new Container instance and eagerly initializes all known provider instances.
// Configuration and options are applied lazily upon first access to each provider.
// This function is expected to be called once during application startup.
// It returns an error if any provider instance fails to be created.
// The returned Container instance is NOT concurrency-SAFE for subsequent method calls.
func New(config interfaces.StructuredConfig, opts ...options.Option) (Container, error) {
	initOpts := optionutil.NewT[containerOptions](opts...)

	var baseLogger runtimelog.Logger
	loggerConfig, err := config.DecodeLogger()
	if err != nil {
		runtimelog.NewHelper(runtimelog.DefaultLogger).Warnf("failed to decode logger config, using default logger: %v", err)
		loggerConfig = nil
	}
	baseLogger = runtimelog.NewLogger(loggerConfig)

	enrichedLogger := baseLogger
	if initOpts.appInfo != nil {
		enrichedLogger = runtimelog.With(baseLogger,
			"service.name", initOpts.appInfo.Name(),
			"service.version", initOpts.appInfo.Version(),
			"service.id", initOpts.appInfo.ID(),
		)
	}

	c := &containerImpl{
		config:     config,
		logger:     enrichedLogger,
		appInfo:    initOpts.appInfo,
		initOpts:   initOpts,
		components: make(map[string]interfaces.Component),
	}

	// --- Eagerly create concrete provider instances (without configuration or options) ---

	// Registry Provider
	c.registryProvider = registry.NewProvider(c.logger)

	// Middleware Provider
	c.middlewareProvider = middleware.NewProvider(c.logger)

	// Cache Provider
	c.cacheProvider = cache.NewProvider(c.logger)

	// Database Provider
	c.databaseProvider = database.NewProvider(c.logger)

	// ObjectStore Provider
	c.objectStoreProvider = objectstore.NewProvider(c.logger)

	return c, nil
}

func (c *containerImpl) AppInfo() interfaces.AppInfo {
	return c.appInfo
}

func (c *containerImpl) Logger() runtimelog.Logger {
	return c.logger
}

// Component retrieves a generic component by name.
// WARNING: This method is NOT CONCURRENCY-SAFE. Concurrent calls may lead to race conditions.
func (c *containerImpl) Component(name string, opts ...options.Option) (interfaces.Component, error) {
	comp, ok := c.components[name]
	if ok {
		// If the component implements SetOptions, apply them
		if so, ok := comp.(interface{ SetOptions(...options.Option) }); ok {
			so.SetOptions(opts...)
		}
		return comp, nil
	}
	// This is a placeholder for a generic component factory mechanism.
	// In a real scenario, you would look up a factory for 'name' and use it.
	// For now, we'll return an error, indicating that generic components need factories.
	return nil, fmt.Errorf("component '%s' not found and no factory is registered for it", name)
}

// RegisterComponent registers an already created component instance.
// WARNING: This method is NOT CONCURRENCY-SAFE. Concurrent calls may lead to race conditions.
func (c *containerImpl) RegisterComponent(name string, comp interfaces.Component) {
	if _, loaded := c.components[name]; loaded {
		runtimelog.NewHelper(c.logger).Warnf("component with name '%s' is being overwritten", name)
	}
	c.components[name] = comp
}

// Registry retrieves the RegistryProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) Registry(opts ...options.Option) (RegistryProvider, error) {
	if !c.registryConfigured {
		discoveries, err := c.config.DecodeDiscoveries()
		if err != nil {
			return nil, fmt.Errorf("failed to decode discoveries config for RegistryProvider: %w", err)
		}
		c.registryProvider.SetConfig(discoveries)
		c.registryConfigured = true
	}
	c.registryProvider.SetOptions(opts...)
	return c.registryProvider, nil
}

// DefaultRegistrar retrieves the default registrar from the RegistryProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) DefaultRegistrar() (runtimeregistry.KRegistrar, error) {
	// This call will trigger configuration and options setting for registryProvider
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

// Middleware retrieves the MiddlewareProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) Middleware(opts ...options.Option) (MiddlewareProvider, error) {
	if !c.middlewareConfigured {
		middlewares, err := c.config.DecodeMiddlewares()
		if err != nil {
			return nil, fmt.Errorf("failed to decode middlewares config for MiddlewareProvider: %w", err)
		}
		c.middlewareProvider.SetConfig(middlewares)
		c.middlewareConfigured = true
	}
	c.middlewareProvider.SetOptions(opts...)
	return c.middlewareProvider, nil
}

// Cache retrieves the CacheProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) Cache(opts ...options.Option) (CacheProvider, error) {
	if !c.cacheConfigured {
		caches, err := c.config.DecodeCaches()
		if err != nil {
			return nil, fmt.Errorf("failed to decode caches config for CacheProvider: %w", err)
		}
		c.cacheProvider.SetConfig(caches)
		c.cacheConfigured = true
	}
	c.cacheProvider.SetOptions(opts...)
	return c.cacheProvider, nil
}

// DefaultCache retrieves the default cache from the CacheProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) DefaultCache() (storageiface.Cache, error) {
	// This call will trigger configuration and options setting for cacheProvider
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

// Database retrieves the DatabaseProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) Database(opts ...options.Option) (DatabaseProvider, error) {
	if !c.databaseConfigured {
		databases, err := c.config.DecodeDatabases()
		if err != nil {
			return nil, fmt.Errorf("failed to decode databases config for DatabaseProvider: %w", err)
		}
		c.databaseProvider.SetConfig(databases)
		c.databaseConfigured = true
	}
	c.databaseProvider.SetOptions(opts...)
	return c.databaseProvider, nil
}

// DefaultDatabase retrieves the default database from the DatabaseProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) DefaultDatabase() (storageiface.Database, error) {
	// This call will trigger configuration and options setting for databaseProvider
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

// ObjectStore retrieves the ObjectStoreProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) ObjectStore(opts ...options.Option) (ObjectStoreProvider, error) {
	if !c.objectStoreConfigured {
		filestores, err := c.config.DecodeObjectStores()
		if err != nil {
			return nil, fmt.Errorf("failed to decode object stores config for ObjectStoreProvider: %w", err)
		}
		c.objectStoreProvider.SetConfig(filestores)
		c.objectStoreConfigured = true
	}
	c.objectStoreProvider.SetOptions(opts...)
	return c.objectStoreProvider, nil
}

// DefaultObjectStore retrieves the default object store from the ObjectStoreProvider.
// WARNING: This method is NOT CONCURRENCY-SAFE.
func (c *containerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	// This call will trigger configuration and options setting for objectStoreProvider
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
