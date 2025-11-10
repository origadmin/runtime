package grpc

import (
	kgprc "github.com/go-kratos/kratos/v2/transport/grpc"
	grpcx "google.golang.org/grpc"

	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/service"
)

// ServerOptions is a container for gRPC server-specific options.
type ServerOptions struct {
	// ServiceOptions holds common service-level configurations.
	ServiceOptions *service.Options

	// GrpcServerOptions allows passing native Kratos gRPC server options.
	GrpcServerOptions []kgprc.ServerOption
}

// FromServerOptions creates a new gRPC ServerOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromServerOptions(opts []options.Option) *ServerOptions {
	o := &ServerOptions{}
	// Apply gRPC server-specific options first
	optionutil.Apply(o, opts...)

	// Initialize and include common service-level options if not already set.
	// This prevents redundant application of common options.
	if o.ServiceOptions == nil {
		o.ServiceOptions = service.FromOptions(opts)
	}

	return o
}

// WithGrpcServerOptions appends Kratos gRPC server options.
func WithGrpcServerOptions(opts ...kgprc.ServerOption) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.GrpcServerOptions = append(o.GrpcServerOptions, opts...)
	})
}

// ClientOptions is a container for gRPC client-specific options.
type ClientOptions struct {
	// ServiceOptions holds common service-level configurations.
	ServiceOptions *service.Options

	// GrpcDialOptions allows passing native gRPC client dial options.
	GrpcDialOptions []grpcx.DialOption
}

// FromClientOptions creates a new gRPC ClientOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromClientOptions(opts []options.Option) *ClientOptions {
	o := &ClientOptions{}
	// Apply gRPC client-specific options first
	optionutil.Apply(o, opts...)

	// Initialize and include common service-level options if not already set.
	// This prevents redundant application of common options.
	if o.ServiceOptions == nil {
		o.ServiceOptions = service.FromOptions(opts)
	}

	return o
}

// WithGrpcDialOptions appends native gRPC client dial options.
func WithGrpcDialOptions(opts ...grpcx.DialOption) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		o.GrpcDialOptions = append(o.GrpcDialOptions, opts...)
	})
}
