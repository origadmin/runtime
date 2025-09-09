package grpc

import (
	"context"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/runtime/service"
)

// optionsKey is a private key type to avoid collisions in context.
type (
	grpcServerOptionsKey struct{}
	grpcClientOptionsKey struct{} // This key will now store transgrpc.ClientOption
)

// WithServerOption is an option to add a Kratos transgrpc.ServerOption to the context.
func WithServerOption(opt ...transgrpc.ServerOption) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		opts, _ := o.Context.Value(grpcServerOptionsKey{}).([]transgrpc.ServerOption)
		o.Context = context.WithValue(o.Context, grpcServerOptionsKey{}, append(opts, opt...))
	}
}

// WithClientOption is an option to add a transgrpc.ClientOption to the context.
func WithClientOption(opt ...transgrpc.ClientOption) service.Option { // Change parameter type
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		opts, _ := o.Context.Value(grpcClientOptionsKey{}).([]transgrpc.ClientOption) // Change stored type
		o.Context = context.WithValue(o.Context, grpcClientOptionsKey{}, append(opts, opt...))
	}
}

// FromServerOptions returns the collected Kratos transgrpc.ServerOption from the service.Options' Context.
func FromServerOptions(o *service.Options) []transgrpc.ServerOption {
	if o == nil || o.Context == nil {
		return nil
	}
	opts, _ := o.Context.Value(grpcServerOptionsKey{}).([]transgrpc.ServerOption)
	return opts
}

// FromClientOptions returns the collected transgrpc.ClientOption from the service.Options' Context.
func FromClientOptions(o *service.Options) []transgrpc.ClientOption { // Change return type
	if o == nil || o.Context == nil {
		return nil
	}
	opts, _ := o.Context.Value(grpcClientOptionsKey{}).([]transgrpc.ClientOption) // Change retrieved type
	return opts
}
