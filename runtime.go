package runtime

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/storage"
)

// Runtime defines the application's runtime environment, providing convenient access to core components.
// It encapsulates an interfaces.Container and is the primary object that applications will interact with.
type Runtime struct {
	result bootstrap.Result
}

// New is the core constructor for a Runtime instance.
// It takes a fully initialized bootstrap result, which is guaranteed to be non-nil.
func New(result bootstrap.Result) *Runtime {
	if result == nil {
		panic("bootstrap.Result cannot be nil when creating a new Runtime")
	}
	return &Runtime{result: result}
}

// NewFromBootstrap is a convenience constructor that simplifies application startup.
// It encapsulates the entire process of calling bootstrap.New and then runtime.New.
// It accepts bootstrap.Option parameters directly, allowing the user to configure the bootstrap process.
func NewFromBootstrap(bootstrapPath string, opts ...bootstrap.Option) (*Runtime, error) {
	bootstrapResult, err := bootstrap.New(bootstrapPath, opts...)
	if err != nil {
		return nil, err
	}

	rt := New(bootstrapResult)
	return rt, nil
}

// WithAppInfo returns an Option that sets the application information.
func WithAppInfo(info *AppInfo) bootstrap.Option {
	return bootstrap.WithAppInfo((*interfaces.AppInfo)(info))
}

// AppInfo returns the application's configured information (ID, name, version, metadata).
func (r *Runtime) AppInfo() *interfaces.AppInfo {
	return r.result.AppInfo()
}

// Config returns the configuration decoder, allowing access to raw configuration values.
func (r *Runtime) Config() interfaces.Config {
	return r.result.Config()
}

// StructuredConfig returns the structured configuration decoder, which provides path-based access to configuration values.
// This is typically used for initializing components like servers and clients that require specific configuration blocks.
func (r *Runtime) StructuredConfig() interfaces.StructuredConfig {
	return r.result.StructuredConfig()
}

// Logger returns the configured Kratos logger.
func (r *Runtime) Logger() log.Logger {
	return r.result.Logger()
}

// Container returns the underlying dependency injection container. This method is primarily for advanced use cases
// where direct access to the container is necessary. For most common operations, prefer using the specific
// accessor methods provided by the Runtime (e.g., Logger(), Config(), Component()).
func (r *Runtime) Container() interfaces.Container {
	return r.result.Container()
}

// NewApp creates a new Kratos application instance.
// It wires together the runtime's configured components (like the default registrar) with the provided transport servers.
// It now accepts additional Kratos options for more flexible configuration.
func (r *Runtime) NewApp(servers []transport.Server, options ...kratos.Option) *kratos.App {
	opts := []kratos.Option{
		kratos.Logger(r.Logger()),
		kratos.Server(servers...),
	}
	info := (*AppInfo)(r.AppInfo())
	opts = append(opts, info.Options()...)

	if registrar := r.DefaultRegistrar(); registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}

	opts = append(opts, options...)
	return kratos.New(opts...)
}

// Component retrieves a generic, user-defined component by its registered name.
func (r *Runtime) Component(name string) (interface{}, bool) {
	return r.result.Container().Component(name)
}

// DefaultRegistrar returns the default service registrar, used for service self-registration.
// It may be nil if no default registry is configured.
func (r *Runtime) DefaultRegistrar() registry.Registrar {
	return r.result.Container().DefaultRegistrar()
}

// Discovery returns a service discovery component by its configured name.
func (r *Runtime) Discovery(name string) (registry.Discovery, bool) {
	return r.result.Container().Discovery(name)
}

// Registrar returns a service registrar component by its configured name.
func (r *Runtime) Registrar(name string) (registry.Registrar, bool) {
	return r.result.Container().Registrar(name)
}

// Storage returns the configured storage provider.
func (r *Runtime) Storage() storage.Provider {
	return r.result.Container().StorageProvider()
}

// Cleanup executes the cleanup function for all resources acquired during bootstrap.
// This should be called via defer right after the Runtime is created.
func (r *Runtime) Cleanup() {
	// The cleanup function itself can be nil if no resources need cleaning.
	if cleanup := r.result.Cleanup(); cleanup != nil {
		cleanup()
	}
}
