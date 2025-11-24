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
)

// App defines the application's runtime environment.
type App struct {
	result    bootstrap.Result
	container container.Container
	appInfo   interfaces.AppInfo // Now an interface
}

// New is the core constructor for a App instance.
// It now accepts container options to allow for proper AppInfo injection.
func New(result bootstrap.Result, ctnOpts ...options.Option) *App {
	if result == nil {
		panic("bootstrap.Result cannot be nil when creating a new App")
	}

	// Create the container, passing the options through.
	ctn := container.New(result.StructuredConfig(), ctnOpts...)

	return &App{
		result:    result,
		container: ctn,
		appInfo:   ctn.AppInfo(), // Get the definitive AppInfo from the container.
	}
}

// NewFromBootstrap is a convenience constructor that simplifies application startup.
func NewFromBootstrap(bootstrapPath string, opts ...Option) (*App, error) {
	// 1. Apply runtime options to get bootstrap options and a potential AppInfo.
	rtOpts := &appOptions{}
	for _, o := range opts {
		o(rtOpts)
	}

	// 2. Call bootstrap.New
	bootstrapResult, err := bootstrap.New(bootstrapPath, rtOpts.bootstrapOpts...)
	if err != nil {
		return nil, err
	}

	// Note: The logic for merging AppInfo from config has been removed.
	// The AppInfo provided via runtime.WithAppInfo is now considered definitive.
	// If no AppInfo is provided, the container will operate without one (or with its own defaults).

	// 3. Create the App instance, passing the AppInfo to the container via an option.
	rt := New(bootstrapResult, container.WithAppInfo(rtOpts.appInfo))
	return rt, nil
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
func (r *App) Container(opts ...options.Option) container.Container {
	return r.container.WithOptions(opts...)
}

// NewApp creates a new Kratos application instance.
func (r *App) NewApp(servers []transport.Server, options ...kratos.Option) *kratos.App {
	info := r.AppInfo()
	if info == nil {
		panic("AppInfo not available in runtime.App, cannot create kratos.App")
	}

	// Prepare metadata, ensuring it's not nil.
	md := info.Metadata()
	if md == nil {
		md = make(map[string]string)
	}
	// Correctly inject the env as part of the metadata.
	if info.Env() != "" {
		md["env"] = info.Env()
	}

	opts := []kratos.Option{
		kratos.Logger(r.Logger()),
		kratos.Server(servers...),
		kratos.ID(info.ID()),
		kratos.Name(info.Name()),
		kratos.Version(info.Version()),
		kratos.Metadata(md), // Pass the enriched metadata.
	}

	if registrar, _ := r.DefaultRegistrar(); registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}

	opts = append(opts, options...)
	return kratos.New(opts...)
}

// Component retrieves a generic, user-defined component by its registered name.
func (r *App) Component(name string) (interface{}, error) {
	return r.container.Component(name)
}

// DefaultRegistrar returns the default service registrar.
func (r *App) DefaultRegistrar() (registry.Registrar, error) {
	reg, err := r.container.Registry()
	if err != nil {
		return nil, err
	}
	return reg.DefaultRegistrar()
}

// RegistryProvider returns the service registry provider.
func (r *App) RegistryProvider(opts ...options.Option) (container.RegistryProvider, error) {
	return r.container.Registry(opts...)
}

// AppInfo returns the application's metadata.
func (r *App) AppInfo() interfaces.AppInfo {
	return r.appInfo
}
