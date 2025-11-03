package http

import (
	"net/http"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/rs/cors"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
	"github.com/origadmin/runtime/service"
)

// CorsOption defines a functional option for configuring advanced CORS settings in code.
// It allows developers to override or enhance the base configuration provided by the proto file.
type CorsOption func(*cors.Options)

// WithAllowOriginFunc sets the `AllowOriginFunc` on the underlying `cors.Options`.
// This provides dynamic, request-based control over allowed origins, which cannot be defined in a static proto file.
// Example: `WithAllowOriginFunc(func(origin string) bool { return origin == "https://example.com" })`
func WithAllowOriginFunc(f func(origin string) bool) CorsOption {
	return func(o *cors.Options) {
		o.AllowOriginFunc = f
	}
}

// WithAllowOriginVaryRequestFunc sets the `AllowOriginVaryRequestFunc` on the underlying `cors.Options`.
// This allows dynamic origin control based on the full http.Request.
// Example: `WithAllowOriginVaryRequestFunc(func(r *http.Request, origin string) bool { return r.Header.Get("X-Custom") == "true" })`
func WithAllowOriginVaryRequestFunc(f func(r *http.Request, origin string) (bool, []string)) CorsOption {
	return func(o *cors.Options) {
		o.AllowOriginVaryRequestFunc = f
	}
}

// ServerOptions is a container for HTTP server-specific options.
type ServerOptions struct {
	// ServiceOptions holds common service-level configurations.
	ServiceOptions *service.Options

	// HttpServerOptions allows passing native Kratos HTTP server options.
	HttpServerOptions []transhttp.ServerOption

	// CorsOptions allows applying advanced, code-based configurations to the CORS handler.
	// These options will be applied *after* the base configuration from the proto file.
	CorsOptions []CorsOption

	// Options holds common service-level configurations.
	Options []options.Option
}

// FromServerOptions creates a new HTTP ServerOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromServerOptions(opts []options.Option) *ServerOptions {
	// Apply HTTP server-specific options first
	o := optionutil.NewT[ServerOptions](opts...)
	o.Options = opts
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

// WithCorsOptions appends advanced, code-based CORS configurations.
// These will be applied on top of the CORS settings from the proto configuration.
func WithCorsOptions(opts ...CorsOption) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.CorsOptions = append(o.CorsOptions, opts...)
	})
}

// ClientOptions is a container for HTTP client-specific options.
type ClientOptions struct {
	// ServiceOptions holds common service-level configurations.
	ServiceOptions *service.Options

	// HttpClientOptions allows passing native Kratos HTTP client options.
	HttpClientOptions []transhttp.ClientOption

	// Options holds common service-level configurations.
	Options []options.Option
}

// FromClientOptions creates a new HTTP ClientOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromClientOptions(opts []options.Option) *ClientOptions {
	// Apply HTTP client-specific options first
	o := optionutil.NewT[ClientOptions](opts...)
	o.Options = opts

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
