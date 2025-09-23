package runtime

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	// appv1 "github.com/origadmin/runtime/api/gen/go/app/v1" // Removed: No longer directly used here
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces" // This import is already here and now used for interfaces.AppInfo and interfaces.Config
)

// Runtime defines the application's runtime environment, providing convenient access to core components.
// It encapsulates an interfaces.ComponentProvider and is the primary object that applications will interact with.
type Runtime struct {
	provider interfaces.ComponentProvider
}

// New is the core constructor for a Runtime instance.
// It takes a fully initialized ComponentProvider, typically created by bootstrap.NewProvider.
func New(provider interfaces.ComponentProvider) *Runtime {
	return &Runtime{provider: provider}
}

// NewFromBootstrap is a convenience constructor that simplifies application startup.
// It encapsulates the entire process of calling bootstrap.NewProvider and then runtime.New.
// It accepts bootstrap.Option parameters directly, allowing the user to configure the bootstrap process.
func NewFromBootstrap(bootstrapPath string, opts ...bootstrap.Option) (*Runtime, func(), error) {
	provider, cleanup, err := bootstrap.NewProvider(bootstrapPath, opts...)
	if err != nil {
		return nil, nil, err
	}

	rt := New(provider)
	return rt, cleanup, nil
}

// AppInfo returns the application's configured information (ID, name, version, metadata).
// Modified: Now returns the interfaces.AppInfo interface.
func (r *Runtime) AppInfo() interfaces.AppInfo {
	return r.provider.AppInfo()
}

// Logger returns the configured Kratos logger.
func (r *Runtime) Logger() log.Logger {
	return r.provider.Logger()
}

// Decoder returns the configuration decoder, allowing access to raw configuration values.
// Added: Exposes the Config decoder.
func (r *Runtime) Decoder() interfaces.Config {
	return r.provider.Config()
}

// NewApp creates a new Kratos application instance.
// It wires together the runtime's configured components (like the default registrar) with the provided transport servers.
// Modified: Uses r.AppInfo().Options() for Kratos options.
func (r *Runtime) NewApp(servers ...transport.Server) *kratos.App {
	// Get Kratos options directly from the interfaces.AppInfo
	appOpts := AppInfo(r.AppInfo()).Options()

	opts := []kratos.Option{
		kratos.Logger(r.Logger()),
		kratos.Server(servers...),
	}
	opts = append(opts, appOpts...) // Append the app info options

	if registrar := r.DefaultRegistrar(); registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}

	return kratos.New(opts...)
}

// DefaultRegistrar returns the default service registrar, used for service self-registration.
// It may be nil if no default registry is configured.
func (r *Runtime) DefaultRegistrar() registry.Registrar {
	return r.provider.DefaultRegistrar()
}

// Discovery returns a service discovery component by its configured name.
func (r *Runtime) Discovery(name string) (registry.Discovery, bool) {
	disc, ok := r.provider.Discoveries()[name]
	return disc, ok
}

// Registrar returns a service registrar component by its configured name.
func (r *Runtime) Registrar(name string) (registry.Registrar, bool) {
	reg, ok := r.provider.Registrars()[name]
	return reg, ok
}

// Component retrieves a generic, user-defined component by its registered name.
func (r *Runtime) Component(name string) (interface{}, bool) {
	return r.provider.Component(name)
}
