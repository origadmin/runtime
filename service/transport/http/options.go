package http

import (
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// Empty keys
var (
	// serverOptionsKey is the context key for HTTP server options
	serverOptionsKey = optionutil.Key[[]transhttp.ServerOption]{}
	// clientOptionsKey is the context key for HTTP client options
	clientOptionsKey = optionutil.Key[[]transhttp.ClientOption]{}
)

type httpServerOptions struct {
	ServerOptions []transhttp.ServerOption
}

type httpClientOptions struct {
	ClientOptions []transhttp.ClientOption
}

// WithServerOption adds HTTP server options to the context.
func WithServerOption(opts ...transhttp.ServerOption) options.Option {
	return optionutil.Update(func(o *httpServerOptions) {
		o.ServerOptions = append(o.ServerOptions, opts...)
	})
}

// WithClientOption adds HTTP client options to the context.
func WithClientOption(opts ...transhttp.ClientOption) options.Option {
	return optionutil.Update(func(o *httpClientOptions) {
		o.ClientOptions = append(o.ClientOptions, opts...)
	})
}

// FromServerOptions retrieves HTTP server options from the service.Options.
func FromServerOptions(opts ...options.Option) []transhttp.ServerOption {
	var o httpServerOptions
	optionutil.Apply(&o, opts...)
	return o.ServerOptions
}

// FromClientOptions retrieves HTTP client options from the service.Options.
func FromClientOptions(opts ...options.Option) []transhttp.ClientOption {
	var o httpClientOptions
	optionutil.Apply(&o, opts...)
	return o.ClientOptions
}
