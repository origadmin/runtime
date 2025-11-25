package container

import (
	"fmt"
	"sync"

	"github.com/origadmin/runtime/container/internal/cache"
	"github.com/origadmin/runtime/container/internal/database"
	"github.com/origadmin/runtime/container/internal/middleware"
	"github.com/origadmin/runtime/container/internal/objectstore"
	"github.com/origadmin/runtime/container/internal/registry"
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/constant" // Import constant package
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
	Component(name string, opts ...options.Option) (interfaces.Component, error) // Modified to accept options
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
	mu       sync.RWMutex
	config   interfaces.StructuredConfig
	logger   runtimelog.Logger
	appInfo  interfaces.AppInfo
	initOpts *containerOptions // Options used during initialization.

	// Caches for initialized providers and components.
	// Using sync.Map for concurrent read/write.
	providers  sync.Map // Stores initialized Provider instances (e.g., RegistryProvider)
	components sync.Map // Stores initialized generic Component instances

	// Map to store sync.Once for each provider, ensuring initFunc is called only once.
	providerOnce sync.Map
}

// New creates a new Container instance.
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
	if initOpts.appInfo != nil {
		enrichedLogger = runtimelog.With(baseLogger,
			"service.name", initOpts.appInfo.Name(),
			"service.version", initOpts.appInfo.Version(),
			"service.id", initOpts.appInfo.ID(),
		)
	}

	return &containerImpl{
		config:       config,
		logger:       enrichedLogger,
		appInfo:      initOpts.appInfo,
		initOpts:     initOpts,
		providerOnce: sync.Map{}, // Initialize the sync.Map for sync.Once
	}
}

func (c *containerImpl) AppInfo() interfaces.AppInfo {
	return c.appInfo
}

func (c *containerImpl) Logger() runtimelog.Logger {
	return c.logger
}

// Component retrieves a generic component by name.
// It will create and cache the component on first request, using provided options.
func (c *containerImpl) Component(name string, opts ...options.Option) (interfaces.Component, error) {
	if cached, ok := c.components.Load(name); ok {
		return cached.(interfaces.Component), nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after lock
	if cached, ok := c.components.Load(name); ok {
		return cached.(interfaces.Component), nil
	}

	// This is a placeholder for a generic component factory mechanism.
	// In a real scenario, you would look up a factory for 'name' and use it.
	// For now, we'll return an error, indicating that generic components need factories.
	return nil, fmt.Errorf("component '%s' not found and no factory is registered for it", name)
}

// RegisterComponent registers an already created component instance.
func (c *containerImpl) RegisterComponent(name string, comp interfaces.Component) {
	if _, loaded := c.components.LoadOrStore(name, comp); loaded {
		runtimelog.NewHelper(c.logger).Warnf("component with name '%s' is being overwritten", name)
		c.components.Store(name, comp) // Explicitly overwrite
	}
}

// initProvider is a generic helper for creating and caching providers, ensuring initFunc is called only once.
func (c *containerImpl) initProvider(
	providerName string,
	initFunc func(logger runtimelog.Logger, finalOpts ...options.Option) (any, error),
	callOpts ...options.Option,
) (any, error) {
	// Fast path: if already cached, return immediately.
	if cached, ok := c.providers.Load(providerName); ok {
		return cached, nil
	}

	// Get or create a sync.Once for this provider name.
	// LoadOrStore returns the existing value or stores the new one.
	// The actual value is an interface{}, so we need a type assertion.
	onceVal, _ := c.providerOnce.LoadOrStore(providerName, &sync.Once{})
	once := onceVal.(*sync.Once)

	var (
		provider any
		err      error
	)

	once.Do(func() {
		// The callOpts are now considered the final, merged options from the runtime.App.
		provider, err = initFunc(c.logger, callOpts...)
		if err != nil {
			// If initFunc fails, we need to ensure subsequent calls also fail.
			// One way is to store the error, or re-panic if it's critical.
			// For simplicity and standard sync.Once behavior, we'll let it fail once.
			return
		}
		c.providers.Store(providerName, provider)
	})

	// If there was an error during initFunc, it will be returned here.
	// If provider is nil and err is not, it means initFunc failed.
	if err != nil {
		// If initFunc failed, we should remove the sync.Once from providerOnce
		// so that a subsequent call can retry initialization.
		// However, sync.Once is designed to run only once. If we want retry logic,
		// we might need a different pattern (e.g., storing the error in the cache).
		// For simplicity and standard sync.Once behavior, we'll let it fail once.
		return nil, err
	}

	// Retrieve the provider from cache after initialization.
	// This handles the case where multiple goroutines hit LoadOrStore for sync.Once
	// but only one successfully initializes. Others will wait and then Load.
	finalProvider, ok := c.providers.Load(providerName)
	if !ok {
		// This should ideally not happen if initFunc successfully stored it.
		// It might happen if initFunc returned an error and didn't store.
		return nil, fmt.Errorf("provider '%s' was not stored after initialization attempt", providerName)
	}

	return finalProvider, nil
}

func (c *containerImpl) Registry(opts ...options.Option) (RegistryProvider, error) {
	provider, err := c.initProvider(string(constant.ComponentRegistries), func(logger runtimelog.Logger, finalOpts ...options.Option) (any, error) {
		discoveries, err := c.config.DecodeDiscoveries()
		if err != nil {
			return nil, err
		}
		p := registry.NewProvider(logger, finalOpts...)
		p.SetConfig(discoveries) // Set structural config after creation
		return p, nil
	}, opts...)

	if err != nil {
		return nil, err
	}
	return provider.(RegistryProvider), nil
}

func (c *containerImpl) DefaultRegistrar() (runtimeregistry.KRegistrar, error) {
	// Attempt to get the RegistryProvider from cache.
	cachedProvider, ok := c.providers.Load(string(constant.ComponentRegistries))
	if !ok {
		return nil, fmt.Errorf("registry provider not initialized, call Registry() first")
	}
	provider := cachedProvider.(RegistryProvider)

	// Attempt to get the default registrar from the provider.
	registrar, err := provider.DefaultRegistrar(c.initOpts.defaultRegistrarName)
	if err != nil {
		return nil, fmt.Errorf("default registrar '%s' not found or failed to retrieve: %w", c.initOpts.defaultRegistrarName, err)
	}
	return registrar, nil
}

func (c *containerImpl) Middleware(opts ...options.Option) (MiddlewareProvider, error) {
	provider, err := c.initProvider(string(constant.ComponentMiddlewares), func(logger runtimelog.Logger, finalOpts ...options.Option) (any, error) {
		middlewares, err := c.config.DecodeMiddlewares()
		if err != nil {
			return nil, err
		}
		p := middleware.NewProvider(logger, finalOpts...)
		p.SetConfig(middlewares) // Set structural config after creation
		return p, nil
	}, opts...)

	if err != nil {
		return nil, err
	}
	return provider.(MiddlewareProvider), nil
}

func (c *containerImpl) Cache(opts ...options.Option) (CacheProvider, error) {
	provider, err := c.initProvider(string(constant.ComponentCaches), func(logger runtimelog.Logger, finalOpts ...options.Option) (any, error) {
		caches, err := c.config.DecodeCaches()
		if err != nil {
			return nil, err
		}
		p := cache.NewProvider(logger, finalOpts...)
		p.SetConfig(caches) // Set structural config after creation
		return p, nil
	}, opts...)

	if err != nil {
		return nil, err
	}
	return provider.(CacheProvider), nil
}

func (c *containerImpl) DefaultCache() (storageiface.Cache, error) {
	cachedProvider, ok := c.providers.Load(string(constant.ComponentCaches))
	if !ok {
		return nil, fmt.Errorf("cache provider not initialized, call Cache() first")
	}
	provider := cachedProvider.(CacheProvider)

	cacheInstance, err := provider.DefaultCache(c.initOpts.defaultCacheName)
	if err != nil {
		return nil, fmt.Errorf("default cache '%s' not found or failed to retrieve: %w", c.initOpts.defaultCacheName, err)
	}
	return cacheInstance, nil
}

func (c *containerImpl) Database(opts ...options.Option) (DatabaseProvider, error) {
	provider, err := c.initProvider(string(constant.ComponentDatabases), func(logger runtimelog.Logger, finalOpts ...options.Option) (any, error) {
		databases, err := c.config.DecodeDatabases()
		if err != nil {
			return nil, err
		}
		p := database.NewProvider(logger, finalOpts...)
		p.SetConfig(databases) // Set structural config after creation
		return p, nil
	}, opts...)

	if err != nil {
		return nil, err
	}
	return provider.(DatabaseProvider), nil
}

func (c *containerImpl) DefaultDatabase() (storageiface.Database, error) {
	cachedProvider, ok := c.providers.Load(string(constant.ComponentDatabases))
	if !ok {
		return nil, fmt.Errorf("database provider not initialized, call Database() first")
	}
	provider := cachedProvider.(DatabaseProvider)

	dbInstance, err := provider.DefaultDatabase(c.initOpts.defaultDatabaseName)
	if err != nil {
		return nil, fmt.Errorf("default database '%s' not found or failed to retrieve: %w", c.initOpts.defaultDatabaseName, err)
	}
	return dbInstance, nil
}

func (c *containerImpl) ObjectStore(opts ...options.Option) (ObjectStoreProvider, error) {
	provider, err := c.initProvider(string(constant.ComponentObjectStores), func(logger runtimelog.Logger, finalOpts ...options.Option) (any, error) {
		filestores, err := c.config.DecodeObjectStores()
		if err != nil {
			return nil, err
		}
		p := objectstore.NewProvider(logger, finalOpts...)
		p.SetConfig(filestores) // Set structural config after creation
		return p, nil
	}, opts...)

	if err != nil {
		return nil, err
	}
	return provider.(ObjectStoreProvider), nil
}

func (c *containerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	cachedProvider, ok := c.providers.Load(string(constant.ComponentObjectStores))
	if !ok {
		return nil, fmt.Errorf("object store provider not initialized, call ObjectStore() first")
	}
	provider := cachedProvider.(ObjectStoreProvider)

	osInstance, err := provider.DefaultObjectStore(c.initOpts.defaultObjectStoreName)
	if err != nil {
		return nil, fmt.Errorf("default object store '%s' not found or failed to retrieve: %w", c.initOpts.defaultObjectStoreName, err)
	}
	return osInstance, nil
}
