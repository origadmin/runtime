package runtime

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/container"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/storage"
)

// App defines the application's runtime environment, providing convenient access to core components.
// It encapsulates an interfaces.Container and is the primary object that applications will interact with.
type App struct {
	result    bootstrap.Result
	container container.Container
}

// New is the core constructor for a App instance.
// It takes a fully initialized bootstrap result, which is guaranteed to be non-nil.
func New(result bootstrap.Result) *App {
	if result == nil {
		panic("bootstrap.Result cannot be nil when creating a new App")
	}
	return &App{
		result:    result,
		container: container.New(result.StructuredConfig()),
	}
}

// NewFromBootstrap is a convenience constructor that simplifies application startup.
// It encapsulates the entire process of calling bootstrap.New and then runtime.New.
// It accepts bootstrap.Option parameters directly, allowing the user to configure the bootstrap process.
func NewFromBootstrap(bootstrapPath string, opts ...bootstrap.Option) (*App, error) {
	bootstrapResult, err := bootstrap.New(bootstrapPath, opts...)
	if err != nil {
		return nil, err
	}

	rt := New(bootstrapResult)
	return rt, nil
}

// Config returns the configuration decoder, allowing access to raw configuration values.
func (r *App) Config() interfaces.Config {
	return r.result.Config()
}

// StructuredConfig returns the structured configuration decoder, which provides path-based access to configuration values.
// This is typically used for initializing components like servers and clients that require specific configuration blocks.
func (r *App) StructuredConfig() interfaces.StructuredConfig {
	return r.result.StructuredConfig()
}

// Logger returns the configured Kratos logger.
func (r *App) Logger() log.Logger {
	return r.container.Logger()
}

// Container returns the underlying dependency injection container. This method is primarily for advanced use cases
// where direct access to the container is necessary. For most common operations, prefer using the specific
// accessor methods provided by the App (e.g., Logger(), Config(), Component()).
//
// The returned container is a new instance with the provided options applied, allowing for scoped configuration.
// The original container remains unchanged.
func (r *App) Container(opts ...options.Option) container.Container {
	// Return a new container with the provided options applied
	return r.container.WithOptions(opts...)
}

// NewApp creates a new Kratos application instance.
// It wires together the runtime's configured components (like the default registrar) with the provided transport servers.
// It now accepts additional Kratos options for more flexible configuration.
func (r *App) NewApp(servers []transport.Server, options ...kratos.Option) *kratos.App {
	opts := []kratos.Option{
		kratos.Logger(r.Logger()),
		kratos.Server(servers...),
	}
	info := r.AppInfo()
	opts = append(opts, info.Options()...)

	if registrar := r.DefaultRegistrar(); registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}

	opts = append(opts, options...)
	return kratos.New(opts...)
}

// Component retrieves a generic, user-defined component by its registered name.
func (r *App) Component(name string) (interface{}, error) {
	return r.container.Component(name)
}

// DefaultRegistrar returns the default service registrar, used for service self-registration.
// It may be nil if no default registry is configured.
func (r *App) DefaultRegistrar() (registry.Registrar, error) {
	reg, err := r.container.Registry()
	if err != nil {
		return nil, err
	}
	return reg.DefaultRegistrar()
}

// RegistryProvider returns the service registry provider, used for service discovery.
// It may be nil if no default registry is configured.
func (r *App) RegistryProvider(opts ...options.Option) (container.RegistryProvider, error) {
	return r.container.Registry(opts...)
}

// Registrar returns a service registrar component by its configured name.
func (r *App) Registrar(name string) (registry.Registrar, bool) {
	return r.container.Registrar(name)
}

// Storage returns the configured storage provider.
func (r *App) Storage() storage.Provider {
	return r.container.StorageProvider()
}

// AppInfo returns the application's metadata.
func (r *App) AppInfo() *interfaces.AppInfo {
	app := r.result.AppInfo()
	if app == nil {
		// Return a default AppInfo if not set
		return &interfaces.AppInfo{
			ID:        "unknown",
			Name:      "unknown",
			Version:   "v0.0.0",
			Env:       "dev",
			StartTime: time.Now(),
			Metadata:  make(map[string]string),
		}
	}

	// Convert protobuf AppInfo to interfaces.AppInfo
	info := &interfaces.AppInfo{
		ID:        app.Id,
		Name:      app.Name,
		Version:   app.Version,
		Env:       app.Environment,
		StartTime: app.StartTime.AsTime(),
		Metadata:  app.Metadata,
	}

	return info
}

// Options returns the Kratos options based on the application's metadata.
func (info *interfaces.AppInfo) Options() []kratos.Option {
	opts := []kratos.Option{
		kratos.ID(info.ID),
		kratos.Name(info.Name),
		kratos.Version(info.Version),
		kratos.Metadata(info.Metadata),
	}

	// Add environment if set
	if info.Env != "" {
		opts = append(opts, kratos.Env(info.Env))
	}

	return opts
}
