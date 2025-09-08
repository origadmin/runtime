package http

import (
	"context"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/framework/runtime/service"
)

// optionsKey is a private key type to avoid collisions in context.
type (
	httpServerOptionsKey struct{}
	httpClientOptionsKey struct{}
)

// WithServerOption is an option to add a transhttp.ServerOption to the context.
func WithServerOption(opt ...transhttp.ServerOption) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		opts, _ := o.Context.Value(httpServerOptionsKey{}).([]transhttp.ServerOption)
		o.Context = context.WithValue(o.Context, httpServerOptionsKey{}, append(opts, opt...))
	}
}

// WithClientOption is an option to add a transhttp.ClientOption to the context.
func WithClientOption(opt ...transhttp.ClientOption) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		opts, _ := o.Context.Value(httpClientOptionsKey{}).([]transhttp.ClientOption)
		o.Context = context.WithValue(o.Context, httpClientOptionsKey{}, append(opts, opt...))
	}
}

// FromServerOptions returns the collected transhttp.ServerOption from the service.Options' Context.
func FromServerOptions(o *service.Options) []transhttp.ServerOption {
	if o == nil || o.Context == nil {
		return nil
	}
	opts, _ := o.Context.Value(httpServerOptionsKey{}).([]transhttp.ServerOption)
	return opts
}

// FromClientOptions returns the collected transhttp.ClientOption from the service.Options' Context.
func FromClientOptions(o *service.Options) []transhttp.ClientOption {
	if o == nil || o.Context == nil {
		return nil
	}
	opts, _ := o.Context.Value(httpClientOptionsKey{}).([]transhttp.ClientOption)
	return opts
}
