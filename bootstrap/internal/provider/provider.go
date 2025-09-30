package provider

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	// appv1 "github.com/origadmin/runtime/api/gen/go/app/v1" // Removed: No longer directly used here
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	"github.com/origadmin/runtime/middleware"

	"github.com/origadmin/runtime/bootstrap/constant"
	"github.com/origadmin/runtime/interfaces" // Ensure this is imported for interfaces.AppInfo and ComponentFactoryRegistry
	runtimeLog "github.com/origadmin/runtime/log"
	runtimeRegistry "github.com/origadmin/runtime/registry"
)

// componentProviderImpl is the default implementation of the interfaces.ComponentProvider interface.
// It holds all the initialized components for the application runtime.
type componentProviderImpl struct {
	appInfo                  interfaces.AppInfo // Modified: Now stores interfaces.AppInfo
	logger                   log.Logger
	discoveries              map[string]registry.Discovery
	registrars               map[string]registry.Registrar
	defaultRegistrar         registry.Registrar
	config                   interfaces.Config // Added: Store the configuration decoder
	components               map[string]interface{}
	serverMiddlewaresMap     map[string]middleware.KMiddleware
	clientMiddlewaresMap     map[string]middleware.KMiddleware
	componentFactoryRegistry interfaces.ComponentFactoryRegistry // Added: Store the component factory registry
}

func (p *componentProviderImpl) ServerMiddleware(name middleware.Name) (middleware.KMiddleware, bool) {
	mw, ok := p.serverMiddlewaresMap[string(name)]
	return mw, ok
}

func (p *componentProviderImpl) ClientMiddleware(name middleware.Name) (middleware.KMiddleware, bool) {
	mw, ok := p.clientMiddlewaresMap[string(name)]
	return mw, ok
}

// Statically assert that componentProviderImpl implements the interface.
var _ interfaces.ComponentProvider = (*componentProviderImpl)(nil)

// NewComponentProvider creates a new, uninitialized component provider.
// It now accepts the interfaces.AppInfo, interfaces.Config, and interfaces.ComponentFactoryRegistry instances.
func NewComponentProvider(appInfo interfaces.AppInfo, cfg interfaces.Config, factoryRegistry interfaces.ComponentFactoryRegistry) *componentProviderImpl {
	return &componentProviderImpl{
		appInfo:                  appInfo,
		config:                   cfg, // Store the config instance
		components:               make(map[string]interface{}),
		componentFactoryRegistry: factoryRegistry,
	}
}

// RegisterComponent adds a user-defined component to the provider's internal registry.
// This is intended to be called by the bootstrap process after the component has been decoded.
func (p *componentProviderImpl) RegisterComponent(name string, comp interface{}) {
	// Ensure the map is initialized.
	if p.components == nil {
		p.components = make(map[string]interface{})
	}

	// Check for duplicates, as this likely indicates a configuration error.
	if _, exists := p.components[name]; exists {
		p.Logger().Log(log.LevelWarn, "msg", "overwriting an existing component registration", "name", name)
	}

	p.components[name] = comp
	p.Logger().Log(log.LevelInfo, "msg", "registered component", "name", name)
}

// InitComponents consumes the configuration and initializes all core and generic components.
// This is the main logic hub for component creation.
func (p *componentProviderImpl) InitComponents(cfg interfaces.Config) error {
	// 1. Initialize Logger with graceful fallback.
	if err := p.initLogger(cfg); err != nil {
		// Even if the logger fails to initialize from config, a fallback is created.
		// We log the error but do not stop the bootstrap process.
		log.Errorf("failed to initialize logger component, error: %v", err)
	}
	helper := log.NewHelper(p.Logger())
	// 2. Initialize Registries and Discoveries with graceful fallback.
	if err := p.initRegistries(cfg); err != nil {
		// Log the error but continue, as local mode is the fallback.
		helper.Errorf("failed to initialize registries component, error: %v", err)
	}

	if err := p.initMiddlewares(cfg); err != nil {
		// Log the error but continue.
		helper.Errorf("failed to initialize middlewares component, error: %v", err)
	}

	// 3. Initialize generic components from the [components] config section.
	if err := p.initGenericComponents(cfg); err != nil {
		// Log the error but continue.
		helper.Errorf("failed to initialize generic components, error: %v", err)
	}

	return nil
}

// initLogger handles the initialization of the logger component.
func (p *componentProviderImpl) initLogger(cfg interfaces.Config) error {
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
		p.Logger().Log(log.LevelWarn, "msg", "failed to decode logger component", "error", err)
		return err
	}

	// 4. Create the logger instance. NewLogger handles the nil config gracefully.
	logger := runtimeLog.NewLogger(loggerCfg)

	// 5. Set the logger for the provider and globally for the Kratos framework.
	p.logger = logger
	runtimeLog.SetLogger(p.logger)

	return nil
}

// initRegistries handles the initialization of the service discovery and registration components.
func (p *componentProviderImpl) initRegistries(cfg interfaces.Config) error {
	var registriesBlock struct {
		Default     string                            `json:"default"`
		Discoveries map[string]*discoveryv1.Discovery `json:"discoveries"`
	}

	// For registries, we use the generic Decode to get both the 'default' key and the map.
	err := cfg.Decode(constant.ComponentRegistries, &registriesBlock)

	// Graceful Fallback: If there's an error or no registries are configured, run in local mode.
	if err != nil || registriesBlock.Discoveries == nil {
		p.Logger().Log(log.LevelInfo, "msg", "no registries configured or failed to decode, running in local mode", "error", err)
		p.discoveries = make(map[string]registry.Discovery)
		p.registrars = make(map[string]registry.Registrar)
		return nil // Not a fatal error
	}

	p.discoveries = make(map[string]registry.Discovery, len(registriesBlock.Discoveries))
	p.registrars = make(map[string]registry.Registrar, len(registriesBlock.Discoveries))

	for name, discoveryCfg := range registriesBlock.Discoveries {
		// Create Discovery
		d, err := runtimeRegistry.NewDiscovery(discoveryCfg)
		if err != nil {
			p.Logger().Log(log.LevelWarn, "msg", "failed to create discovery", "name", name, "error", err)
			continue // Skip this one
		}
		p.discoveries[name] = d

		// Create Registrar
		r, err := runtimeRegistry.NewRegistrar(discoveryCfg)
		if err != nil {
			p.Logger().Log(log.LevelWarn, "msg", "failed to create registrar", "name", name, "error", err)
			continue // Skip this one
		}
		p.registrars[name] = r
	}

	// Set the default registrar
	if registriesBlock.Default != "" {
		if r, ok := p.registrars[registriesBlock.Default]; ok {
			p.defaultRegistrar = r
			p.Logger().Log(log.LevelInfo, "msg", "default registrar set", "name", registriesBlock.Default)
		} else {
			p.Logger().Log(log.LevelWarn, "msg", "default registrar not found", "name", registriesBlock.Default)
		}
	}

	return nil
}

func (p *componentProviderImpl) initMiddlewares(cfg interfaces.Config) error {
	p.serverMiddlewaresMap = make(map[string]middleware.KMiddleware)
	p.clientMiddlewaresMap = make(map[string]middleware.KMiddleware)
	v, ok := cfg.(interfaces.MiddlewareConfigDecoder)
	if ok {
		middlewares, err := v.DecodeMiddleware()
		if err != nil {
			return fmt.Errorf("failed to decode middlewares: %w", err)
		}
		for _, mc := range middlewares.GetMiddlewares() {
			if mc.GetEnabled() {
				mclient, ok := middleware.NewClient(mc)
				if !ok {
					continue
				}
				mserver, ok := middleware.NewServer(mc)
				if !ok {
					continue
				}
				p.serverMiddlewaresMap[mc.GetType()] = mserver
				p.clientMiddlewaresMap[mc.GetType()] = mclient
			}
		}
	}
	return nil
}

// initGenericComponents handles the initialization of user-defined components.
func (p *componentProviderImpl) initGenericComponents(cfg interfaces.Config) error {
	var componentsMap map[string]map[string]interface{}
	if err := cfg.Decode(constant.ComponentComponents, &componentsMap); err != nil {
		// If the components key doesn't exist, it's not an error, just means no generic components.
		return nil
	}

	for name, compCfg := range componentsMap {
		// The 'type' field is mandatory for finding the factory.
		compType, ok := compCfg["type"].(string)
		if !ok || compType == "" {
			p.Logger().Log(log.LevelWarn, "msg", "component type is missing or not a string, skipping", "name", name)
			continue
		}

		// Get the factory for this component type.
		factory, found := p.componentFactoryRegistry.GetFactory(compType)
		if !found {
			p.Logger().Log(log.LevelWarn, "msg", "component factory not found, skipping", "type", compType, "name", name)
			continue
		}

		// Create the component instance.
		instance, err := factory(cfg, compCfg)
		if err != nil {
			p.Logger().Log(log.LevelWarn, "msg", "failed to create component instance", "name", name, "type", compType, "error", err)
			continue
		}

		// Store the created component.
		p.components[name] = instance
		p.Logger().Log(log.LevelInfo, "msg", "initialized generic component", "name", name, "type", compType)
	}

	return nil
}

// AppInfo implements the interfaces.ComponentProvider interface.
func (p *componentProviderImpl) AppInfo() interfaces.AppInfo { // Modified: Now returns interfaces.AppInfo
	return p.appInfo
}

// Logger implements the interfaces.ComponentProvider interface.
func (p *componentProviderImpl) Logger() log.Logger {
	// Ensure a logger always exists, even if initialization failed.
	if p.logger == nil {
		p.logger = log.NewStdLogger(os.Stderr)
	}
	return p.logger
}

// Discoveries implements the interfaces.ComponentProvider interface.
func (p *componentProviderImpl) Discoveries() map[string]registry.Discovery {
	return p.discoveries
}

// Registrars implements the interfaces.ComponentProvider interface.
func (p *componentProviderImpl) Registrars() map[string]registry.Registrar {
	return p.registrars
}

// DefaultRegistrar implements the interfaces.ComponentProvider interface.
func (p *componentProviderImpl) DefaultRegistrar() registry.Registrar {
	return p.defaultRegistrar
}

// Config implements the interfaces.ComponentProvider interface.
func (p *componentProviderImpl) Config() interfaces.Config {
	return p.config
}

// Component implements the interfaces.ComponentProvider interface.
func (p *componentProviderImpl) Component(name string) (interface{}, bool) {
	c, ok := p.components[name]
	return c, ok
}
