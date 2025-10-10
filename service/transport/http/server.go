package http

import (
	"fmt"
	"net/http/pprof"

	"github.com/go-kratos/kratos/v2/middleware"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service/tls"
)

// NewHTTPServer creates a new concrete HTTP server instance based on the provided configuration.
// It returns *transhttp.Server, not the generic interfaces.Server.
func NewHTTPServer(httpConfig *transportv1.HttpServerConfig, serverOpts *ServerOptions) (*transhttp.Server, error) {
	// Prepare the Kratos HTTP server options.
	var kratosOpts []transhttp.ServerOption

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if serverOpts.ServiceOptions != nil {
		c = serverOpts.ServiceOptions.Container
	}

	// Add CORS support, merging proto config with code-based options.
	if corsConfig := httpConfig.GetCors(); corsConfig != nil && corsConfig.GetEnabled() {
		corsHandler, err := NewCorsHandler(corsConfig, serverOpts.CorsOptions...)
		if err != nil {
			// Propagate the configuration error upwards.
			return nil, fmt.Errorf("failed to create CORS handler: %w", err)
		}
		if corsHandler != nil {
			kratosOpts = append(kratosOpts, transhttp.Filter(corsHandler))
		}
	}

	// Check if middlewares are configured.
	hasMiddlewaresConfigured := len(httpConfig.GetMiddlewares()) > 0

	// If middlewares are configured but no container is provided, return an error.
	// This consolidates the nil check for the container.
	if hasMiddlewaresConfigured && c == nil {
		return nil, fmt.Errorf("application container is required for server middlewares but not found in options")
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if hasMiddlewaresConfigured {
		// 'c' is guaranteed to be non-nil at this point due to the early check above.
		for _, name := range httpConfig.GetMiddlewares() {
			m, ok := c.ServerMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("server middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		// If no specific middlewares are configured, use default ones from adapter.go.
		mws = DefaultServerMiddlewares()
	}

	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transhttp.Middleware(mws...))
	}

	// Apply other server options from protobuf config
	if httpConfig.Network != "" {
		kratosOpts = append(kratosOpts, transhttp.Network(httpConfig.Network))
	}
	if httpConfig.Addr != "" {
		kratosOpts = append(kratosOpts, transhttp.Address(httpConfig.Addr))
	}
	if httpConfig.Timeout != nil {
		kratosOpts = append(kratosOpts, transhttp.Timeout(httpConfig.Timeout.AsDuration()))
	}

	// Apply TLS configuration
	// Configure TLS for server
	if tlsConfig := httpConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := tls.NewServerTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create server TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transhttp.TLSConfig(tlsCfg))
	}

	// Apply any external Kratos HTTP server options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(serverOpts.HttpServerOptions) > 0 {
		kratosOpts = append(kratosOpts, serverOpts.HttpServerOptions...)
	}

	// Create the HTTP server instance.
	srv := transhttp.NewServer(kratosOpts...)

	// Enable pprof debugging endpoints if configured.
	if httpConfig.GetEnablePprof() {
		registerPprof(srv)
	}

	return srv, nil
}

// registerPprof registers the pprof handlers with the HTTP server.
func registerPprof(srv *transhttp.Server) {
	srv.HandleFunc("/debug/pprof", pprof.Index)
	srv.HandleFunc("/debug/cmdline", pprof.Cmdline)
	srv.HandleFunc("/debug/profile", pprof.Profile)
	srv.HandleFunc("/debug/symbol", pprof.Symbol)
	srv.HandleFunc("/debug/trace", pprof.Trace)
	srv.HandleFunc("/debug/allocs", pprof.Handler("allocs").ServeHTTP)
	srv.HandleFunc("/debug/block", pprof.Handler("block").ServeHTTP)
	srv.HandleFunc("/debug/goroutine", pprof.Handler("goroutine").ServeHTTP)
	srv.HandleFunc("/debug/heap", pprof.Handler("heap").ServeHTTP)
	srv.HandleFunc("/debug/mutex", pprof.Handler("mutex").ServeHTTP)
	srv.HandleFunc("/debug/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
