package container

import (
	// stderrors "errors" // REMOVED: imported and not used
	"fmt"
	"reflect"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/data/storage"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces" // Ensure this is imported for interfaces.AppInfo and ComponentFactoryRegistry
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
	runtimeMiddleware "github.com/origadmin/runtime/middleware" // Import runtime/middleware package, but only for internal use.
	runtimeRegistry "github.com/origadmin/runtime/registry"
)

// container is the default implementation of the interfaces.Container interface.
// It holds all the initialized components for the application runtime.
type container struct {
	discoveries          map[string]registry.Discovery
	registrars           map[string]registry.Registrar
	defaultRegistrar     registry.Registrar
	components           map[string]interface{}
	serverMiddlewaresMap map[string]middleware.Middleware
	clientMiddlewaresMap map[string]middleware.Middleware
	storageProvider      storageiface.Provider
}

func (c *container) StorageProvider() storageiface.Provider {
	return c.storageProvider
}

// Statically assert that componentProviderImpl implements the interface.
var _ interfaces.Container = (*container)(nil)

// Builder is a builder for creating and initializing a Container.
// It provides a fluent API for configuring and building a container instance.
type Builder struct {
	container *container
	logger    runtimelog.Logger
	config    interfaces.StructuredConfig
	err       error // Track errors during building
	factories map[string]interfaces.ComponentFactory
}

func defaultInstance[T any](
	registrars map[string]T,
	defaultKey string,
	activeKey string) T {
	if defaultKey == "" {
		defaultKey = interfaces.GlobalDefaultKey
	}
	if activeKey == "" {
		activeKey = defaultKey
	}

	if len(registrars) == 1 {
		for _, registrar := range registrars {
			return registrar
		}
	} else if len(registrars) > 1 {
		if d, ok := registrars[activeKey]; ok {
			return d
		}
		if d, ok := registrars[defaultKey]; ok {
			return d
		}
		for _, registrar := range registrars {
			return registrar
		}
	}
	var zero T
	return zero
}

// NewContainer creates a new, uninitialized container.
// Deprecated: Use NewBuilder for more flexible container initialization.
func NewContainer() interfaces.Container {
	return &container{
		components:           make(map[string]interface{}),
		serverMiddlewaresMap: make(map[string]middleware.Middleware), // Corrected type
		clientMiddlewaresMap: make(map[string]middleware.Middleware), // Corrected type
	}
}

// NewBuilder creates a new Builder instance.
// The builder can be used to configure and build a container with the provided config.
func NewBuilder(componentFactories map[string]interfaces.ComponentFactory) *Builder {
	return &Builder{
		factories: componentFactories,
		container: &container{
			components:           make(map[string]interface{}),
			discoveries:          make(map[string]registry.Discovery),
			registrars:           make(map[string]registry.Registrar),
			serverMiddlewaresMap: make(map[string]middleware.Middleware), // Corrected type
			clientMiddlewaresMap: make(map[string]middleware.Middleware), // Corrected type
		},
	}
}

// WithFactory registers a component factory with the container.
func (b *Builder) WithFactory(name string, factory interfaces.ComponentFactory) *Builder {
	if b.err != nil {
		return b
	}
	if _, exists := b.factories[name]; exists {
		b.err = fmt.Errorf("component factory for '%s' is already registered", name)
		return b
	}
	b.factories[name] = factory
	return b
}

// WithConfig sets the configuration that will be used when building the container.
func (b *Builder) WithConfig(cfg interfaces.StructuredConfig) *Builder {
	if b.err != nil {
		return b
	}
	b.config = cfg
	return b
}

// WithLogger sets a custom logger for the container.
func (b *Builder) WithLogger(logger log.Logger) *Builder {
	if b.err != nil {
		return b
	}
	b.logger = logger
	return b
}

// WithComponent adds a pre-initialized component to the container.
func (b *Builder) WithComponent(name string, component interface{}) *Builder {
	if b.err != nil {
		return b
	}
	b.container.components[name] = component
	return b
}

// Build initializes and returns the Container using the provided config.
// It returns the container and any error that occurred during initialization.
func (b *Builder) Build(opts ...options.Option) (interfaces.Container, error) {
	if b.err != nil {
		return nil, b.err
	}

	if b.config == nil {
		return nil, runtimeerrors.NewStructured("bootstrap", "config must be provided using WithConfig()").WithCaller()
	}

	// 1. Initialize Logger with graceful fallback.
	if err := b.initLogger(); err != nil {
		// Even if the logger fails to initialize from config, a fallback is created.
		// We log the error but do not stop the bootstrap process.
		log.Errorf("failed to initialize logger component, error: %v", err) // This uses the temporary logger if c.logger is not yet set
	}
	helper := log.NewHelper(b.logger) // Now c.Logger() should be initialized or fallbacked

	// 2. Initialize Registries and Discoveries with graceful fallback.
	if err := b.initRegistries(opts...); err != nil {
		// Log the error but continue, as local mode is the fallback.
		helper.Errorf("failed to initialize registries component, error: %v", err)
	}

	if err := b.initMiddlewares(opts...); err != nil {
		// Log the error but continue.
		helper.Errorf("failed to initialize middlewares component, error: %v", err)
	}

	// 3. Initialize generic components from the [components] config section.
	if err := b.initGenericComponents(opts...); err != nil {
		// Log the error but continue.
		helper.Errorf("failed to initialize generic components, error: %v", err)
	}

	// Ensure logger is set
	if b.logger == nil {
		b.logger = log.DefaultLogger
	}

	return b.container, nil
}

// MustBuild is like Build but panics if an error occurs.
// This is useful for initialization code that should fail fast.
func (b *Builder) MustBuild() interfaces.Container {
	c, err := b.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build container: %v", err))
	}
	return c
}

// Initialize consumes the configuration and initializes all core and generic components.
// This is the main logic hub for component creation.
func (b *Builder) Initialize(cfg interfaces.Config) error {

	return nil
}

// initLogger handles the initialization of the logger component.
func (b *Builder) initLogger() error {
	loggerCfg, err := b.config.DecodeLogger()
	if err != nil {
		return fmt.Errorf("failed to decode logger config: %w", err)
	}
	// 4. Create the logger instance. NewLogger handles the nil config gracefully.
	logger := runtimelog.NewLogger(loggerCfg)

	// 5. Set the logger for the provider and globally for the Kratos framework.
	b.logger = logger
	return nil
}

func (b *Builder) buildStorageProvider() error {
	// 6. Create the storage provider instance. NewStorageProvider handles the nil config gracefully.
	storageProvider, err := storage.New(b.config)
	if err != nil {
		return fmt.Errorf("failed to create storage provider: %w", err)
	}
	// 7. Set the storage provider for the container.
	b.container.storageProvider = storageProvider
	return nil
}

// initRegistries handles the initialization of the service discovery and registration components.
func (b *Builder) initRegistries(opts ...options.Option) error {
	helper := log.NewHelper(b.Logger()) // Use log.Helper

	discoveriesCfg, err := b.config.DecodeDiscoveries()
	// Graceful Fallback: If there's an error or no registries are configured, run in local mode.
	if err != nil {
		helper.Infow("msg", "no registries configured or failed to decode, running in local mode", "error", err)
		return nil // Not a fatal error
	}

	if discoveriesCfg == nil || len(discoveriesCfg.GetConfigs()) == 0 {
		helper.Infow("msg", "no registries configured, running in local mode")
		b.container.discoveries = make(map[string]registry.Discovery)
		b.container.registrars = make(map[string]registry.Registrar)
		return nil // Not a fatal error
	}

	discoveries := discoveriesCfg.GetConfigs()
	b.container.discoveries = make(map[string]registry.Discovery, len(discoveries))
	b.container.registrars = make(map[string]registry.Registrar, len(discoveries))

	for _, discoveryCfg := range discoveries {
		key := discoveryCfg.GetType()
		if discoveryCfg.GetName() != "" {
			key = discoveryCfg.GetName()
		}
		// Create Discovery
		d, err := runtimeRegistry.NewDiscovery(discoveryCfg, opts...)
		if err != nil {
			helper.Warnw("msg", "failed to create discovery", "key", key, "error", err)
			continue // Skip this one
		}

		// Create Registrar
		r, err := runtimeRegistry.NewRegistrar(discoveryCfg, opts...)
		if err != nil {
			helper.Warnw("msg", "failed to create registrar", "key", key, "error", err)
			continue // Skip this one
		}

		b.container.discoveries[key] = d
		b.container.registrars[key] = r
	}

	b.container.defaultRegistrar = defaultInstance(b.container.registrars, discoveriesCfg.GetDefault(), discoveriesCfg.GetActive())
	if b.container.defaultRegistrar == nil {
		helper.Warnw("msg", "no default registrar set")
	}

	return nil
}

func (b *Builder) initMiddlewares(opts ...options.Option) error {
	helper := log.NewHelper(b.Logger()) // Use log.Helper

	middlewares, err := b.config.DecodeMiddlewares()
	if err != nil {
		return fmt.Errorf("failed to decode middlewares: %w", err)
	}
	b.container.serverMiddlewaresMap = make(map[string]middleware.Middleware)
	b.container.clientMiddlewaresMap = make(map[string]middleware.Middleware)
	if middlewares == nil || len(middlewares.GetConfigs()) == 0 {
		helper.Infow("msg", "no middlewares configured, skipping")
		return nil
	}
	// Get the logger to pass to middleware options
	logger := b.Logger()
	opts = append(opts, runtimelog.WithLogger(logger))
	for _, mc := range middlewares.GetConfigs() {
		if mc.GetEnabled() {
			key := mc.GetType()
			if mc.GetName() != "" {
				key = mc.GetName()
			}
			// Assuming NewClient and NewServer support WithLogger option
			mclient, ok := runtimeMiddleware.NewClient(mc, opts...) // Use runtimeMiddleware
			if !ok {
				helper.Warnw("msg", "failed to create client middleware", "key", key)
				continue
			}
			mserver, ok := runtimeMiddleware.NewServer(mc, opts...) // Use runtimeMiddleware
			if !ok {
				helper.Warnw("msg", "failed to create server middleware", "key", key)
				continue
			}

			b.container.serverMiddlewaresMap[key] = mserver // Store kratos middleware.Middleware
			b.container.clientMiddlewaresMap[key] = mclient // Store kratos middleware.Middleware
		}
	}
	return nil
}

// initGenericComponents handles the initialization of user-defined components.
func (b *Builder) initGenericComponents(opts ...options.Option) error {
	helper := log.NewHelper(b.Logger()) // Use log.Helper

	for name, factory := range b.factories {
		_, ok := b.container.components[name]
		if ok {
			continue
		}
		comp, err := factory(b.config, b.container)
		if err != nil {
			helper.Warnw("msg", "failed to initialize generic component", "name", name, "error", err)
			continue
		}
		b.container.RegisterComponent(name, comp)
		helper.Infow("msg", "initialized generic component", "name", name, "type", reflect.TypeOf(comp))
	}

	return nil
}

func (b *Builder) Logger() log.Logger {
	return b.logger
}

// ServerMiddlewares implements the interfaces.Container interface.
func (c *container) ServerMiddlewares() map[string]middleware.Middleware {
	return c.serverMiddlewaresMap
}

func (c *container) ServerMiddleware(name string) (middleware.Middleware, bool) {
	mw, ok := c.serverMiddlewaresMap[name]
	return mw, ok
}

// ClientMiddlewares implements the interfaces.Container interface.
func (c *container) ClientMiddlewares() map[string]middleware.Middleware {
	return c.clientMiddlewaresMap
}

func (c *container) ClientMiddleware(name string) (middleware.Middleware, bool) {
	mw, ok := c.clientMiddlewaresMap[name]
	return mw, ok
}

func (c *container) Component(name string) (interface{}, bool) {
	comp, ok := c.components[name]
	return comp, ok
}

// RegisterComponent adds a user-defined component to the provider's internal registry.
// This is intended to be called by the bootstrap process after the component has been decoded.
func (c *container) RegisterComponent(name string, comp interface{}) {
	helper := log.NewHelper(c.Logger()) // Use log.Helper
	// Ensure the map is initialized.
	if c.components == nil {
		c.components = make(map[string]interface{})
	}

	// Check for duplicates, as this likely indicates a configuration error.
	if _, exists := c.components[name]; exists {
		helper.Warnw("msg", "overwriting an existing component registration", "name", name)
	}

	c.components[name] = comp
	helper.Infow("msg", "registered component", "name", name)
}

// Logger implements the interfaces.Container interface.
func (c *container) Logger() log.Logger {
	panic("It should not be called here, panic tells you this error and needs to check your code for problems")
	// Ensure a logger always exists, even if initialization failed.
	//if c.logger == nil {
	//	c.logger = log.DefaultLogger
	//}
	//return c.logger
}

// Discoveries implements the interfaces.Container interface.
func (c *container) Discoveries() map[string]registry.Discovery {
	return c.discoveries
}

// Discovery implements the interfaces.Container interface.
func (c *container) Discovery(name string) (registry.Discovery, bool) {
	d, ok := c.discoveries[name]
	return d, ok
}

// Registrars implements the interfaces.Container interface.
func (c *container) Registrars() map[string]registry.Registrar {
	return c.registrars
}

// Registrar implements the interfaces.Container interface.
func (c *container) Registrar(name string) (registry.Registrar, bool) {
	r, ok := c.registrars[name]
	return r, ok
}

// DefaultRegistrar implements the interfaces.Container interface.
func (c *container) DefaultRegistrar() registry.Registrar {
	return c.defaultRegistrar
}
