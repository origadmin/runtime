package http

import (
	"net/http"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/rs/cors"

	rtcontainer "github.com/origadmin/runtime/container"
	"github.com/origadmin/runtime/extension/optionutil"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	rtservice "github.com/origadmin/runtime/service"
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
	// ServerOptions allows passing native Kratos HTTP server options.
	ServerOptions []transhttp.ServerOption

	// CorsOptions allows applying advanced, code-based configurations to the CORS handler.
	// These options will be applied *after* the base configuration from the proto file.
	CorsOptions []CorsOption

	// Options holds common service-level configurations.
	Options []options.Option

	// Container holds the application container instance.
	Container interfaces.Container

	// ServerRegistrar holds the service registration instance.
	ServerRegistrar rtservice.ServerRegistrar
}

// FromServerOptions creates a new HTTP ServerOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromServerOptions(opts []options.Option) *ServerOptions {
	// Apply HTTP server-specific options first
	o := optionutil.NewT[ServerOptions](opts...)
	o.Options = opts
	o.Container = rtcontainer.FromOptions(opts)
	o.ServerRegistrar = rtservice.ServerRegistrarFromOptions(opts)
	return o
}

// WithServerOptions appends Kratos HTTP server options.
func WithServerOptions(opts ...transhttp.ServerOption) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.ServerOptions = append(o.ServerOptions, opts...)
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
	// ClientOptions allows passing native Kratos HTTP client dial options.
	ClientOptions []transhttp.ClientOption

	// Options holds common service-level configurations.
	Options []options.Option

	// Container holds the application container instance.
	Container interfaces.Container
}

// FromClientOptions creates a new HTTP ClientOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromClientOptions(opts []options.Option) *ClientOptions {
	// Apply HTTP client-specific options first
	o := optionutil.NewT[ClientOptions](opts...)
	o.Options = opts
	o.Container = rtcontainer.FromOptions(opts)

	return o
}

// WithClientOptions appends Kratos HTTP client options.
func WithClientOptions(opts ...transhttp.ClientOption) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		o.ClientOptions = append(o.ClientOptions, opts...)
	})
}
