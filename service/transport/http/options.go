package http

import (
	"context"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/rs/cors"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/helpers/optionutil"
	"github.com/origadmin/runtime/service/transport"
)

// CorsOption defines a functional option for configuring advanced CORS settings in code.
type CorsOption func(*cors.Options)

// WithAllowOriginFunc sets the AllowOriginFunc.
func WithAllowOriginFunc(f func(origin string) bool) CorsOption {
	return func(o *cors.Options) {
		o.AllowOriginFunc = f
	}
}

// ServerOptions is a container for HTTP server-specific options.
type ServerOptions struct {
	// ServerOptions allows passing native Kratos HTTP server options.
	ServerOptions []transhttp.ServerOption

	// CorsOptions allows applying advanced, code-based configurations.
	CorsOptions []CorsOption

	// Registrar is the HTTP registrar to use.
	Registrar transport.HTTPRegistrar

	// ServerMiddlewares holds a map of named server middlewares.
	ServerMiddlewares map[string]middleware.Middleware
	Context           context.Context
}

// FromServerOptions creates a new HTTP ServerOptions struct.
func FromServerOptions(opts []options.Option) *ServerOptions {
	o := optionutil.NewT[ServerOptions](opts...)
	if o.Context == nil {
		o.Context = context.Background()
	}
	return o
}

// WithServerOptions appends Kratos HTTP server options.
func WithServerOptions(opts ...transhttp.ServerOption) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.ServerOptions = append(o.ServerOptions, opts...)
	})
}

// WithRegistrar sets the HTTP registrar.
func WithRegistrar(registrar transport.HTTPRegistrar) options.Option {
	return optionutil.Update(func(o *ServerOptions) {
		o.Registrar = registrar
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

// FromClientOptions creates a new HTTP ClientOptions struct.
func FromClientOptions(opts []options.Option) *ClientOptions {
	o := optionutil.NewT[ClientOptions](opts...)
	return o
}

// WithClientOptions appends Kratos HTTP client options.
func WithClientOptions(opts ...transhttp.ClientOption) options.Option {
	return optionutil.Update(func(o *ClientOptions) {
		o.ClientOptions = append(o.ClientOptions, opts...)
	})
}

// WithClientMiddlewares adds a map of named client middlewares.
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

// WithDiscoveries adds a map of named service discovery clients.
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
