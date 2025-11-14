package http

import (
	"net/http"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/rs/cors"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/extension/optionutil"
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
// Example: `WithAllowOriginVaryRequestFunc(func(r *http.Request, origin string) (bool, []string) { return r.Header.Get("X-Custom") == "true", nil })`
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

	// ServerRegistrar holds the service registration instance.
	ServerRegistrar rtservice.ServerRegistrar

	// ServerMiddlewares holds a map of named server middlewares.
	ServerMiddlewares map[string]middleware.Middleware
}

// FromServerOptions creates a new HTTP ServerOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromServerOptions(opts []options.Option) *ServerOptions {
	// Apply HTTP server-specific options first
	o := optionutil.NewT[ServerOptions](opts...)
	// Removed: o.Container = rtcontainer.FromOptions(opts)
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

// WithServerMiddlewares adds a map of named server middlewares to the options.
func WithServerMiddlewares(mws map[string]middleware.Middleware) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		if o.ServerMiddlewares == nil {
			o.ServerMiddlewares = make(map[string]middleware.Middleware)
		}
		for name, mw := range mws {
			o.ServerMiddlewares[name] = mw
		}
	})
}

// ClientOptions is a container for HTTP client-specific options.
type ClientOptions struct {
	// ClientOptions allows passing native Kratos HTTP client dial options.
	ClientOptions []transhttp.ClientOption

	// ClientMiddlewares holds a map of named client middlewares.
	ClientMiddlewares map[string]middleware.Middleware

	// Discoveries holds a map of named service discovery clients.
	Discoveries map[string]registry.Discovery
}

// FromClientOptions creates a new HTTP ClientOptions struct by applying a slice of functional options.
// It also initializes and includes the common service-level options, ensuring they are applied only once.
func FromClientOptions(opts []options.Option) *ClientOptions {
	// Apply HTTP client-specific options first
	o := optionutil.NewT[ClientOptions](opts...)
	// Removed: o.Container = rtcontainer.FromOptions(opts)
	return o
}

// WithClientOptions appends Kratos HTTP client options.
func WithClientOptions(opts ...transhttp.ClientOption) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		o.ClientOptions = append(o.ClientOptions, opts...)
	})
}

// WithClientMiddlewares adds a map of named client middlewares to the options.
func WithClientMiddlewares(mws map[string]middleware.Middleware) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		if o.ClientMiddlewares == nil {
			o.ClientMiddlewares = make(map[string]middleware.Middleware)
		}
		for name, mw := range mws {
			o.ClientMiddlewares[name] = mw
		}
	})
}

// WithDiscoveries adds a map of named service discovery clients to the options.
func WithDiscoveries(discoveries map[string]registry.Discovery) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		if o.Discoveries == nil {
			o.Discoveries = make(map[string]registry.Discovery)
		}
		for name, d := range discoveries {
			o.Discoveries[name] = d
		}
		
	})
}
