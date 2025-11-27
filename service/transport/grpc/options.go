package grpc

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	kgprc "github.com/go-kratos/kratos/v2/transport/grpc"
	grpcx "google.golang.org/grpc"

	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/service/transport"
)

// ServerOptions is a container for gRPC server-specific options.
type ServerOptions struct {
	// ServerOptions allows passing native Kratos gRPC server options.
	ServerOptions []kgprc.ServerOption

	Context   context.Context
	Registrar transport.GRPCRegistrar

	// ServerMiddlewares holds a map of named server middlewares.
	ServerMiddlewares map[string]middleware.Middleware
}

// FromServerOptions creates a new gRPC ServerOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromServerOptions(opts []options.Option) *ServerOptions {
	// Apply gRPC server-specific options first
	o := optionutil.NewT[ServerOptions](opts...)
	if o.Context == nil {
		o.Context = context.Background()
	}
	return o
}

// WithServerOptions appends Kratos gRPC server options.
func WithServerOptions(opts ...kgprc.ServerOption) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.ServerOptions = append(o.ServerOptions, opts...)
	})
}

// WithRegistrar sets the gRPC registrar to use for service registration.
func WithRegistrar(registrar transport.GRPCRegistrar) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.Registrar = registrar
	})
}

// WithContext sets the context.Context for the service.
func WithContext(ctx context.Context) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.Context = ctx
	})
}

// WithContextRegistry sets both the context.Context and ServerRegistrar for the service.
func WithContextRegistry(ctx context.Context, registrar transport.GRPCRegistrar) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.Context = ctx
		o.Registrar = registrar
	})
}

// WithServerMiddlewares adds a map of named server middlewares to the options.
func WithServerMiddlewares(mws map[string]middleware.Middleware) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		if o.ServerMiddlewares == nil {
			o.ServerMiddlewares = make(map[string]middleware.Middleware)
		}
		for name, mw := range mws {
			o.ServerMiddlewares[name] = mw
		}
	})
}

// ClientOptions is a container for gRPC client-specific options.
type ClientOptions struct {
	// DialOptions allows passing native gRPC client dial options.
	DialOptions []grpcx.DialOption

	// ClientMiddlewares holds a map of named client middlewares.
	ClientMiddlewares map[string]middleware.Middleware

	// Discoveries holds a map of named service discovery clients.
	Discoveries map[string]registry.Discovery
}

// FromClientOptions creates a new gRPC ClientOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromClientOptions(opts []options.Option) *ClientOptions {
	// Apply gRPC client-specific options first
	o := optionutil.NewT[ClientOptions](opts...)
	return o
}

// WithDialOptions appends native gRPC client dial options.
func WithDialOptions(opts ...grpcx.DialOption) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		o.DialOptions = append(o.DialOptions, opts...)
	})
}

// WithClientMiddlewares adds a map of named client middlewares to the options.
func WithClientMiddlewares(mws map[string]middleware.Middleware) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		if o.ClientMiddlewares == nil {
			o.ClientMiddlewares = make(map[string]middleware.Middleware)
		}
		for name, mw := range mws {
			o.ClientMiddlewares[name] = mw
		}
	})
}

// WithDiscoveries adds a map of named service discovery clients to the options.
func WithDiscoveries(discoveries map[string]registry.Discovery) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		if o.Discoveries == nil {
			o.Discoveries = make(map[string]registry.Discovery)
		}
		for name, d := range discoveries {
			o.Discoveries[name] = d
		}
	})
}
