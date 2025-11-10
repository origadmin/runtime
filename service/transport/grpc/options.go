package grpc

import (
	kgprc "github.com/go-kratos/kratos/v2/transport/grpc"
	grpcx "google.golang.org/grpc"

	rtcontainer "github.com/origadmin/runtime/container"
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	rtservice "github.com/origadmin/runtime/service"
)

// ServerOptions is a container for gRPC server-specific options.
type ServerOptions struct {
	// ServerOptions allows passing native Kratos gRPC server options.
	ServerOptions []kgprc.ServerOption

	// Container holds the application container instance.
	Container interfaces.Container

	// Options holds the common service-level options.
	Options []options.Option

	// ServerRegistrar holds the service registrar instance.
	ServerRegistrar rtservice.ServerRegistrar
}

// FromServerOptions creates a new gRPC ServerOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromServerOptions(opts []options.Option) *ServerOptions {
	// Apply gRPC server-specific options first
	o := optionutil.NewT[ServerOptions](opts...)
	o.Options = opts
	o.Container = rtcontainer.FromOptions(opts)
	o.ServerRegistrar = rtservice.ServerRegistrarFromOptions(opts)
	return o
}

// WithServerOptions appends Kratos gRPC server options.
func WithServerOptions(opts ...kgprc.ServerOption) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.ServerOptions = append(o.ServerOptions, opts...)
	})
}

// ClientOptions is a container for gRPC client-specific options.
type ClientOptions struct {
	// DialOptions allows passing native gRPC client dial options.
	DialOptions []grpcx.DialOption

	// Container holds the application container instance.
	Container interfaces.Container

	// Options holds the common service-level options.
	Options []options.Option
}

// FromClientOptions creates a new gRPC ClientOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromClientOptions(opts []options.Option) *ClientOptions {
	// Apply gRPC client-specific options first
	o := optionutil.NewT[ClientOptions](opts...)
	o.Options = opts
	o.Container = rtcontainer.FromOptions(opts)

	return o
}

// WithDialOptions appends native gRPC client dial options.
func WithDialOptions(opts ...grpcx.DialOption) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		o.DialOptions = append(o.DialOptions, opts...)
	})
}
