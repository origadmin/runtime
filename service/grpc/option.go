package grpc

import (
	"context"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"

	"github.com/origadmin/framework/runtime/service"
)

// optionsKey is a private key type to avoid collisions in context.
type (
	grpcServerOptionsKey struct{}
	grpcClientOptionsKey struct{}
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

// WithClientOption is an option to add a native grpc.DialOption to the context.
func WithClientOption(opt ...grpc.DialOption) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		opts, _ := o.Context.Value(grpcClientOptionsKey{}).([]grpc.DialOption)
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

// FromClientOptions returns the collected native grpc.DialOption from the service.Options' Context.
func FromClientOptions(o *service.Options) []grpc.DialOption {
	if o == nil || o.Context == nil {
		return nil
	}
	opts, _ := o.Context.Value(grpcClientOptionsKey{}).([]grpc.DialOption)
	return opts
}
