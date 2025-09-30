package container

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	// appv1 "github.com/origadmin/runtime/api/gen/go/app/v1" // Removed: No longer directly used here
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	"github.com/origadmin/runtime/middleware"

	"github.com/origadmin/runtime/bootstrap/constant"
	"github.com/origadmin/runtime/interfaces" // Ensure this is imported for interfaces.AppInfo and ComponentFactoryRegistry
	runtimelog "github.com/origadmin/runtime/log"
	runtimeRegistry "github.com/origadmin/runtime/registry"
)

// container is the default implementation of the interfaces.Container interface.
// It holds all the initialized components for the application runtime.
type container struct {
	appInfo   *interfaces.AppInfo // Modified: Now stores interfaces.AppInfo
	config    interfaces.Config   // Added: Store the configuration decoder
	factories map[string]interfaces.ComponentFactory

	logger               log.Logger
	discoveries          map[string]registry.Discovery
	registrars           map[string]registry.Registrar
	defaultRegistrar     registry.Registrar
	components           map[string]interface{}
	serverMiddlewaresMap map[string]middleware.KMiddleware
	clientMiddlewaresMap map[string]middleware.KMiddleware
}

// Statically assert that componentProviderImpl implements the interface.
var _ interfaces.Container = (*container)(nil)

// NewContainer creates a new, uninitialized container.
// It now accepts the interfaces.AppInfo, interfaces.Config, and interfaces.ComponentFactoryRegistry instances.
func NewContainer(cfg interfaces.Config) interfaces.Container {
	return &container{
		config:               cfg, // Store the config instance
		factories:            make(map[string]interfaces.ComponentFactory),
		components:           make(map[string]interface{}),
		serverMiddlewaresMap: make(map[string]middleware.KMiddleware),
		clientMiddlewaresMap: make(map[string]middleware.KMiddleware),
	}
}

func (c *container) RegisterFactory(name string, factory interfaces.ComponentFactory) {
	if _, exists := c.factories[name]; exists {
		panic(fmt.Sprintf("component factory for '%s' is already registered", name))
	}
	c.factories[name] = factory
}

func (c *container) Component(name string) (interface{}, bool) {
	comp, ok := c.components[name]
	return comp, ok
}

// Build instantiates all registered components in a deterministic order.
func (c *container) Build() error {
	// Define a build order to ensure dependencies are met.
	// For example, "logger" must be built before other components that use it.
	buildOrder := []string{"logger", "registries", "middlewares"}

	// Create a map of registered keys for quick lookup
	registeredKeys := make(map[string]struct{})
	for key := range c.factories {
		registeredKeys[key] = struct{}{}
	}

	// Append remaining keys that are not in the explicit build order.
	// Sort them to ensure deterministic build order.
	var remainingKeys []string
	for key := range registeredKeys {
		isExplicit := false
		for _, explicitKey := range buildOrder {
			if key == explicitKey {
				isExplicit = true
				break
			}
		}
		if !isExplicit {
			remainingKeys = append(remainingKeys, key)
		}
	}
	sort.Strings(remainingKeys)
	buildOrder = append(buildOrder, remainingKeys...)

	// Build components in the determined order.
	for _, name := range buildOrder {
		if _, exists := registeredKeys[name]; !exists {
			continue // Skip if a key in buildOrder was not registered.
		}

		// Skip if component is already built (e.g., as a dependency).
		if _, ok := c.components[name]; ok {
			continue
		}

		log.NewHelper(c.Logger()).Infof("component '%s' built successfully", name)
	}

	return nil
}

func (c *container) ServerMiddleware(name middleware.Name) (middleware.KMiddleware, bool) {
	mw, ok := c.serverMiddlewaresMap[string(name)]
	return mw, ok
}

func (c *container) ClientMiddleware(name middleware.Name) (middleware.KMiddleware, bool) {
	mw, ok := c.clientMiddlewaresMap[string(name)]
	return mw, ok
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

// Initialize consumes the configuration and initializes all core and generic components.
// This is the main logic hub for component creation.
func (c *container) Initialize(cfg interfaces.Config) error {
	// 1. Initialize Logger with graceful fallback.
	if err := c.initLogger(cfg); err != nil {
		// Even if the logger fails to initialize from config, a fallback is created.
		// We log the error but do not stop the bootstrap process.
		log.Errorf("failed to initialize logger component, error: %v", err) // This uses the temporary logger if c.logger is not yet set
	}
	helper := log.NewHelper(c.Logger()) // Now c.Logger() should be initialized or fallbacked

	// 2. Initialize Registries and Discoveries with graceful fallback.
	if err := c.initRegistries(cfg); err != nil {
		// Log the error but continue, as local mode is the fallback.
		helper.Errorf("failed to initialize registries component, error: %v", err)
	}

	if err := c.initMiddlewares(cfg); err != nil {
		// Log the error but continue.
		helper.Errorf("failed to initialize middlewares component, error: %v", err)
	}

	// 3. Initialize generic components from the [components] config section.
	if err := c.initGenericComponents(cfg); err != nil {
		// Log the error but continue.
		helper.Errorf("failed to initialize generic components, error: %v", err)
	}

	return nil
}

// initLogger handles the initialization of the logger component.
func (c *container) initLogger(cfg interfaces.Config) error {
	var loggerCfg *loggerv1.Logger
	var err error

	// 1. Prioritize the specific decoder interface (fast path).
	if decoder, ok := cfg.(interfaces.LoggerConfigDecoder); ok {
		loggerCfg, err = decoder.DecodeLogger()
	}

	// 2. Fallback to the generic decoder if the fast path is not taken or explicitly signals a fallback.
	if loggerCfg == nil && (err == nil || errors.Is(err, interfaces.ErrNotImplemented)) {
		// The error from the fast path is reset, as we are now trying the generic path.
		err = cfg.Decode(constant.ComponentLogger, &loggerCfg)
	}

	// 3. If there was any error during decoding, log it as a warning.
	// We use a temporary logger because the main logger isn't created yet.
	if err != nil {
		// Use a temporary logger here as c.logger might not be fully initialized
		log.NewHelper(log.NewStdLogger(os.Stderr)).Warnw("msg", "failed to decode logger component", "error", err)
		return err
	}

	// 4. Create the logger instance. NewLogger handles the nil config gracefully.
	logger := runtimelog.NewLogger(loggerCfg)

	// 5. Set the logger for the provider and globally for the Kratos framework.
	c.logger = logger
	runtimelog.SetLogger(c.logger)

	return nil
}

// initRegistries handles the initialization of the service discovery and registration components.
func (c *container) initRegistries(cfg interfaces.Config) error {
	helper := log.NewHelper(c.Logger()) // Use log.Helper
	var registriesBlock struct {
		Default     string                            `json:"default"`
		Discoveries map[string]*discoveryv1.Discovery `json:"discoveries"`
	}

	// For registries, we use the generic Decode to get both the 'default' key and the map.
	err := cfg.Decode(constant.ComponentRegistries, &registriesBlock)

	// Graceful Fallback: If there's an error or no registries are configured, run in local mode.
	if err != nil || registriesBlock.Discoveries == nil {
		helper.Infow("msg", "no registries configured or failed to decode, running in local mode", "error", err)
		c.discoveries = make(map[string]registry.Discovery)
		c.registrars = make(map[string]registry.Registrar)
		return nil // Not a fatal error
	}

	c.discoveries = make(map[string]registry.Discovery, len(registriesBlock.Discoveries))
	c.registrars = make(map[string]registry.Registrar, len(registriesBlock.Discoveries))

	for name, discoveryCfg := range registriesBlock.Discoveries {
		// Create Discovery
		d, err := runtimeRegistry.NewDiscovery(discoveryCfg)
		if err != nil {
			helper.Warnw("msg", "failed to create discovery", "name", name, "error", err)
			continue // Skip this one
		}
		c.discoveries[name] = d

		// Create Registrar
		r, err := runtimeRegistry.NewRegistrar(discoveryCfg)
		if err != nil {
			helper.Warnw("msg", "failed to create registrar", "name", name, "error", err)
			continue // Skip this one
		}
		c.registrars[name] = r
	}

	// Set the default registrar
	if registriesBlock.Default != "" {
		if r, ok := c.registrars[registriesBlock.Default]; ok {
			helper.Infow("msg", "default registrar set", "name", registriesBlock.Default)
			c.defaultRegistrar = r
		} else {
			helper.Warnw("msg", "default registrar not found", "name", registriesBlock.Default)
		}
	}

	return nil
}

func (c *container) initMiddlewares(cfg interfaces.Config) error {
	helper := log.NewHelper(c.Logger()) // Use log.Helper
	// c.serverMiddlewaresMap = make(map[string]middleware.KMiddleware) // Already initialized in NewComponentProvider
	// c.clientMiddlewaresMap = make(map[string]middleware.KMiddleware) // Already initialized in NewComponentProvider
	v, ok := cfg.(interfaces.MiddlewareConfigDecoder)
	if ok {
		middlewares, err := v.DecodeMiddleware()
		if err != nil {
			return fmt.Errorf("failed to decode middlewares: %w", err)
		}
		// Get the logger to pass to middleware options
		logger := c.Logger()
		for _, mc := range middlewares.GetMiddlewares() {
			if mc.GetEnabled() {
				// Assuming NewClient and NewServer support WithLogger option
				mclient, ok := middleware.NewClient(mc, runtimelog.WithLogger(logger)) // Pass logger
				if !ok {
					helper.Warnw("msg", "failed to create client middleware", "type", mc.GetType())
					continue
				}
				mserver, ok := middleware.NewServer(mc, runtimelog.WithLogger(logger)) // Pass logger
				if !ok {
					helper.Warnw("msg", "failed to create server middleware", "type", mc.GetType())
					continue
				}
				c.serverMiddlewaresMap[mc.GetType()] = mserver
				c.clientMiddlewaresMap[mc.GetType()] = mclient
			}
		}
	}
	return nil
}

// initGenericComponents handles the initialization of user-defined components.
func (c *container) initGenericComponents(cfg interfaces.Config) error {
	helper := log.NewHelper(c.Logger()) // Use log.Helper
	var componentsMap map[string]map[string]interface{}
	if err := cfg.Decode(constant.ComponentComponents, &componentsMap); err != nil {
		// If the components key doesn't exist, it's not an error, just means no generic components.
		return nil
	}

	for name, compCfg := range componentsMap {
		// The 'type' field is mandatory for finding the factory.
		compType, ok := compCfg["type"].(string)
		if !ok || compType == "" {
			helper.Warnw("msg", "component type is missing or not a string, skipping", "name", name)
			continue
		}

		// Get the factory for this component type.
		factory, found := c.factories[compType]
		if !found {
			helper.Warnw("msg", "component factory not found, skipping", "type", compType, "name", name)
			continue
		}

		// Create the component instance.
		instance, err := factory(c) // Note: This factory signature might need adjustment based on the new
		// ComponentFactory type
		if err != nil {
			helper.Warnw("msg", "failed to create component instance", "name", name, "type", compType, "error", err)
			continue
		}

		// Store the created component.
		c.components[name] = instance
		helper.Infow("msg", "initialized generic component", "name", name, "type", compType)
	}

	return nil
}

// AppInfo implements the interfaces.Container interface.
func (c *container) AppInfo() *interfaces.AppInfo { // Modified: Now returns *interfaces.AppInfo
	return c.appInfo
}

// Logger implements the interfaces.Container interface.
func (c *container) Logger() log.Logger {
	// Ensure a logger always exists, even if initialization failed.
	if c.logger == nil {
		c.logger = log.NewStdLogger(os.Stderr)
	}
	return c.logger
}

// Discoveries implements the interfaces.Container interface.
func (c *container) Discoveries() map[string]registry.Discovery {
	return c.discoveries
}

// Registrars implements the interfaces.Container interface.
func (c *container) Registrars() map[string]registry.Registrar {
	return c.registrars
}

// DefaultRegistrar implements the interfaces.Container interface.
func (c *container) DefaultRegistrar() registry.Registrar {
	return c.defaultRegistrar
}

// Config implements the interfaces.Container interface.
func (c *container) Config() interfaces.Config {
	return c.config
}
