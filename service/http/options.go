package http

import (
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/framework/runtime/service"
	"github.com/origadmin/framework/runtime/service/optionutil"
)

// Option keys
var (
	// ServerOptionsKey is the context key for HTTP server options
	ServerOptionsKey = optionutil.OptionKey[[]transhttp.ServerOption]{}
	// ClientOptionsKey is the context key for HTTP client options
	ClientOptionsKey = optionutil.OptionKey[[]transhttp.ClientOption]{}
)

// WithServerOption adds HTTP server options to the context.
func WithServerOption(opts ...transhttp.ServerOption) service.Option {
	if len(opts) == 0 {
		return func(*service.Options) {}
	}
	return func(o *service.Options) {
		optionutil.WithSliceOption(o, ServerOptionsKey, opts...)
	}
}

// WithClientOption adds HTTP client options to the context.
func WithClientOption(opts ...transhttp.ClientOption) service.Option {
	if len(opts) == 0 {
		return func(*service.Options) {}
	}
	return func(o *service.Options) {
		optionutil.WithSliceOption(o, ClientOptionsKey, opts...)
	}
}

// FromServerOptions retrieves HTTP server options from the service.Options.
func FromServerOptions(o *service.Options) []transhttp.ServerOption {
	return optionutil.GetSliceOption(o, ServerOptionsKey)
}

// FromClientOptions retrieves HTTP client options from the service.Options.
func FromClientOptions(o *service.Options) []transhttp.ClientOption {
	return optionutil.GetSliceOption(o, ClientOptionsKey)
}
