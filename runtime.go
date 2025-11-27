package runtime

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1" // Import appv1
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/bootstrap/internal/config" // For NewAppInfo
	"github.com/origadmin/runtime/container"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/constant"
	"github.com/origadmin/runtime/interfaces/options"
)

// App defines the application's runtime environment.
type App struct {
	result         bootstrap.Result
	container      container.Container
	appInfo        interfaces.AppInfo // Final AppInfo after bootstrap
	mu             sync.RWMutex
	globalOpts     []options.Option
	componentOpts  map[string][]options.Option
	containerOpts  []options.Option
	initialAppInfo interfaces.AppInfo // Initial AppInfo from New()
}

// New creates a new, partially initialized App instance with essential metadata and container configurations.
// It accepts functional options to allow for pre-configuring the container.
// The App is not fully functional until Load is called.
func New(name, version string, opts ...Option) (*App, error) {
	initialAppInfo := newAppInfo(name, version)

	app := &App{
		initialAppInfo: initialAppInfo, // Store initial AppInfo
		componentOpts:  make(map[string][]options.Option),
		containerOpts:  make([]options.Option, 0),
	}

	// Apply functional options, which will populate containerOpts
	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

// Load reads the configuration from the given path, completes the App initialization,
// and prepares it for running.
func (r *App) Load(path string, bootOpts ...bootstrap.Option) error {
	// 1. Bootstrap from configuration file.
	res, err := bootstrap.New(path, bootOpts...)
	if err != nil {
		return fmt.Errorf("failed to bootstrap configuration from '%s': %w", path, err)
	}
	r.result = res

	// 2. Merge AppInfo: Code-defined values for Name and Version take precedence.
	// Config-defined values for Env, ID, and Metadata supplement or override.
	finalAppInfo := mergeAppInfo(r.initialAppInfo, convertProtoAppToAppInfo(res.App()))
	r.appInfo = finalAppInfo // Store the final merged AppInfo

	// 3. Create the container.
	ctnOpts := r.containerOpts                                                 // Use collected container options
	r.container = container.New(res.StructuredConfig(), r.appInfo, ctnOpts...) // Pass final AppInfo to container

	return nil
}

// WarmUp attempts to initialize all configured providers and generic components.
// This method should be called after all configurations have been added to the App.
// It returns an error if any component fails to initialize, allowing for early error detection.
func (r *App) WarmUp() error {
	var initErrors []error

	// Known providers
	if _, err := r.RegistryProvider(); err != nil {
		initErrors = append(initErrors, fmt.Errorf("failed to warm up RegistryProvider: %w", err))
	}
	if _, err := r.DatabaseProvider(); err != nil {
		initErrors = append(initErrors, fmt.Errorf("failed to warm up DatabaseProvider: %w", err))
	}
	if _, err := r.CacheProvider(); err != nil {
		initErrors = append(initErrors, fmt.Errorf("failed to warm up CacheProvider: %w", err))
	}
	if _, err := r.ObjectStoreProvider(); err != nil {
		initErrors = append(initErrors, fmt.Errorf("failed to warm up ObjectStoreProvider: %w", err))
	}
	if _, err := r.MiddlewareProvider(); err != nil {
		initErrors = append(initErrors, fmt.Errorf("failed to warm up MiddlewareProvider: %w", err))
	}

	// Generic components (iterate through all configured ones)
	r.mu.RLock()
	configuredComponentNames := make([]string, 0, len(r.componentOpts))
	for name := range r.componentOpts {
		configuredComponentNames = append(configuredComponentNames, name)
	}
	r.mu.RUnlock()

	for _, name := range configuredComponentNames {
		// Skip known providers as they are already handled above
		if name == string(constant.ComponentRegistries) ||
			name == string(constant.ComponentDatabases) ||
			name == string(constant.ComponentCaches) ||
			name == string(constant.ComponentObjectStores) ||
			name == string(constant.ComponentMiddlewares) {
			continue
		}
		if _, err := r.Component(name); err != nil {
			initErrors = append(initErrors, fmt.Errorf("failed to warm up generic component '%s': %w", name, err))
		}
	}

	if len(initErrors) > 0 {
		return fmt.Errorf("runtime app warm-up failed with %d errors: %v", len(initErrors), initErrors)
	}

	return nil
}

// AddGlobalOptions adds options that will be applied to all providers.
func (r *App) AddGlobalOptions(opts ...options.Option) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.globalOpts = append(r.globalOpts, opts...)
}

// AddComponentOptions provides a generic way to add pre-configured options for any named component.
func (r *App) AddComponentOptions(name string, opts ...options.Option) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.componentOpts[name] = append(r.componentOpts[name], opts...)
}

// ConfigureRegistry adds pre-configured options for the registry provider.
func (r *App) ConfigureRegistry(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentRegistries), opts...)
}

// ConfigureDatabase adds pre-configured options for the database provider.
func (r *App) ConfigureDatabase(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentDatabases), opts...)
}

// ConfigureCache adds pre-configured options for the cache provider.
func (r *App) ConfigureCache(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentCaches), opts...)
}

// ConfigureObjectStore adds pre-configured options for the object store provider.
func (r *App) ConfigureObjectStore(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentObjectStores), opts...)
}

// ConfigureMiddleware adds pre-configured options for the middleware provider.
func (r *App) ConfigureMiddleware(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentMiddlewares), opts...)
}

// getMergedOptions safely retrieves and merges global and component-specific options.
func (r *App) getMergedOptions(name string) []options.Option {
	r.mu.RLock()
	defer r.mu.RUnlock()
	global := r.globalOpts
	specific := r.componentOpts[name]
	// The order is important: specific options should be able to override global ones if the underlying
	// implementation processes them last. We place global options first.
	return append(append([]options.Option{}, global...), specific...)
}

// Config returns the configuration decoder.
func (r *App) Config() interfaces.Config {
	return r.result.Config()
}

// StructuredConfig returns the structured configuration decoder.
func (r *App) StructuredConfig() interfaces.StructuredConfig {
	return r.result.StructuredConfig()
}

// Logger returns the configured Kratos logger.
func (r *App) Logger() log.Logger {
	return r.container.Logger()
}

// Container returns the underlying dependency injection container.
func (r *App) Container() container.Container {
	return r.container
}

// NewApp creates a new Kratos application instance.
func (r *App) NewApp(servers []transport.Server, options ...kratos.Option) *kratos.App {
	info := r.appInfo

	md := info.Metadata()
	if md == nil {
		md = make(map[string]string)
	}
	if info.Env() != "" {
		md["env"] = info.Env()
	}

	opts := []kratos.Option{
		kratos.Logger(r.Logger()),
		kratos.Server(servers...),
		kratos.ID(info.ID()),
		kratos.Name(info.Name()),
		kratos.Version(info.Version()),
		kratos.Metadata(md),
	}

	if registrar, _ := r.DefaultRegistrar(); registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}

	opts = append(opts, options...)
	return kratos.New(opts...)
}

// Component retrieves a generic, user-defined component by its registered name.
func (r *App) Component(name string) (interface{}, error) {
	opts := r.getMergedOptions(name)
	return r.container.Component(name, opts...)
}

// DefaultRegistrar returns the default service registrar.
func (r *App) DefaultRegistrar() (registry.Registrar, error) {
	return r.container.DefaultRegistrar()
}

// RegistryProvider returns the service registry provider, using pre-configured options.
func (r *App) RegistryProvider() (container.RegistryProvider, error) {
	opts := r.getMergedOptions(string(constant.ComponentRegistries))
	return r.container.Registry(opts...)
}

// DatabaseProvider returns the database provider, using pre-configured options.
func (r *App) DatabaseProvider() (container.DatabaseProvider, error) {
	opts := r.getMergedOptions(string(constant.ComponentDatabases))
	return r.container.Database(opts...)
}

// CacheProvider returns the cache provider, using pre-configured options.
func (r *App) CacheProvider() (container.CacheProvider, error) {
	opts := r.getMergedOptions(string(constant.ComponentCaches))
	return r.container.Cache(opts...)
}

// ObjectStoreProvider returns the object store provider, using pre-configured options.
func (r *App) ObjectStoreProvider() (container.ObjectStoreProvider, error) {
	opts := r.getMergedOptions(string(constant.ComponentObjectStores))
	return r.container.ObjectStore(opts...)
}

// MiddlewareProvider returns the middleware provider, using pre-configured options.
func (r *App) MiddlewareProvider() (container.MiddlewareProvider, error) {
	opts := r.getMergedOptions(string(constant.ComponentMiddlewares))
	return r.container.Middleware(opts...)
}

// AppInfo returns the application's metadata as an interface.
func (r *App) AppInfo() interfaces.AppInfo {
	return r.appInfo
}

// mergeAppInfo combines AppInfo from runtime and config, with runtime Name and Version taking precedence.
func mergeAppInfo(runtime, config interfaces.AppInfo) interfaces.AppInfo {
	if config == nil {
		return runtime
	}

	// Start with the runtime AppInfo as the base.
	finalID := runtime.ID()
	finalEnv := runtime.Env()
	finalMetadata := runtime.Metadata()

	// If config provides values, they take precedence for these fields.
	if config.ID() != "" {
		finalID = config.ID()
	}
	if config.Env() != "" {
		finalEnv = config.Env()
	}
	if config.Metadata() != nil {
		if finalMetadata == nil {
			finalMetadata = make(map[string]string)
		}
		for k, v := range config.Metadata() {
			finalMetadata[k] = v
		}
	}

	// Create a new AppInfo instance with the merged values.
	// Name and Version are always taken from the runtime AppInfo.
	return config.NewAppInfo(
		runtime.Name(),
		runtime.Version(),
		finalID,
		finalEnv,
		finalMetadata,
	)
}

// convertProtoAppToAppInfo converts a protobuf App message to an interfaces.AppInfo.
func convertProtoAppToAppInfo(protoApp *appv1.App) interfaces.AppInfo {
	if protoApp == nil {
		return nil
	}
	return config.NewAppInfo(
		protoApp.GetName(),
		protoApp.GetVersion(),
		protoApp.GetId(),
		protoApp.GetEnv(),
		protoApp.GetMetadata(),
	)
}
