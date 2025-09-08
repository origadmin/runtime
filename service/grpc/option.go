package grpc

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/framework/runtime/service"
)

// optionsKey is a private key type to avoid collisions in context.
type (
	serverOptionsKey struct{}
	clientOptionsKey struct{}
)

// WithServerOption is an option to add a grpc.ServerOption to the context.
func WithServerOption(opt ...grpc.ServerOption) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		opts, _ := o.Context.Value(serverOptionsKey{}).([]grpc.ServerOption)
		o.Context = context.WithValue(o.Context, serverOptionsKey{}, append(opts, opt...))
	}
}

// FromServerOptions returns the collected grpc.ServerOption from the service.Options' Context.
func FromServerOptions(o *service.Options) []grpc.ServerOption {
	if o == nil || o.Context == nil {
		return nil
	}
	opts, _ := o.Context.Value(serverOptionsKey{}).([]grpc.ServerOption)
	return opts
}

// WithClientOption is an option to add a grpc.ClientOption to the context.
func WithClientOption(opt ...grpc.ClientOption) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		opts, _ := o.Context.Value(clientOptionsKey{}).([]grpc.ClientOption)
		o.Context = context.WithValue(o.Context, clientOptionsKey{}, append(opts, opt...))
	}
}

// FromClientOptions returns the collected grpc.ClientOption from the service.Options' Context.
func FromClientOptions(o *service.Options) []grpc.ClientOption {
	if o == nil || o.Context == nil {
		return nil
	}
	opts, _ := o.Context.Value(clientOptionsKey{}).([]grpc.ClientOption)
	return opts
}
