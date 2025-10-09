package http

import (
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
	"github.com/origadmin/runtime/service"
)

// ServerOptions is a container for HTTP server-specific options.
type ServerOptions struct {
	// ServiceOptions holds common service-level configurations.
	ServiceOptions *service.Options

	// HttpServerOptions allows passing native Kratos HTTP server options.
	HttpServerOptions []transhttp.ServerOption
}

// FromServerOptions creates a new HTTP ServerOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromServerOptions(opts []options.Option) *ServerOptions {
	o := &ServerOptions{}
	// Apply HTTP server-specific options first
	optionutil.Apply(o, opts...)

	// Initialize and include common service-level options if not already set.
	// This prevents redundant application of common options.
	if o.ServiceOptions == nil {
		o.ServiceOptions = service.FromOptions(opts)
	}

	return o
}

// WithHttpServerOptions appends Kratos HTTP server options.
func WithHttpServerOptions(opts ...transhttp.ServerOption) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.HttpServerOptions = append(o.HttpServerOptions, opts...)
	})
}

// ClientOptions is a container for HTTP client-specific options.
type ClientOptions struct {
	// ServiceOptions holds common service-level configurations.
	ServiceOptions *service.Options

	// HttpClientOptions allows passing native Kratos HTTP client options.
	HttpClientOptions []transhttp.ClientOption
}

// FromClientOptions creates a new HTTP ClientOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromClientOptions(opts []options.Option) *ClientOptions {
	o := &ClientOptions{}
	// Apply HTTP client-specific options first
	optionutil.Apply(o, opts...)

	// Initialize and include common service-level options if not already set.
	// This prevents redundant application of common options.
	if o.ServiceOptions == nil {
		o.ServiceOptions = service.FromOptions(opts)
	}

	return o
}

// WithHttpClientOptions appends Kratos HTTP client options.
func WithHttpClientOptions(opts ...transhttp.ClientOption) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		o.HttpClientOptions = append(o.HttpClientOptions, opts...)
	})
}
