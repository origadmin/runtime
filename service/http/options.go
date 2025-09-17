package http

import (
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/runtime/optionutil"
	"github.com/origadmin/runtime/service"
)

// Option keys
var (
	// serverOptionsKey is the context key for HTTP server options
	serverOptionsKey = optionutil.Key[[]transhttp.ServerOption]{}
	// clientOptionsKey is the context key for HTTP client options
	clientOptionsKey = optionutil.Key[[]transhttp.ClientOption]{}
)

// WithServerOption adds HTTP server options to the context.
func WithServerOption(opts ...transhttp.ServerOption) service.Option {
	if len(opts) == 0 {
		return func(*service.Options) {}
	}
	return func(o *service.Options) {
		optionutil.Append(o, serverOptionsKey, opts...)
	}
}

// WithClientOption adds HTTP client options to the context.
func WithClientOption(opts ...transhttp.ClientOption) service.Option {
	if len(opts) == 0 {
		return func(*service.Options) {}
	}
	return func(o *service.Options) {
		optionutil.Append(o, clientOptionsKey, opts...)
	}
}

// FromServerOptions retrieves HTTP server options from the service.Options.
func FromServerOptions(o *service.Options) []transhttp.ServerOption {
	return optionutil.Slice(o, serverOptionsKey)
}

// FromClientOptions retrieves HTTP client options from the service.Options.
func FromClientOptions(o *service.Options) []transhttp.ClientOption {
	return optionutil.Slice(o, clientOptionsKey)
}
