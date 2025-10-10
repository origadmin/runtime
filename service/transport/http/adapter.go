// Copyright (c) 2024 OrigAdmin. All rights reserved.

package http

import (
	"net/http/pprof"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	corsv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/cors"
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// DefaultServerMiddlewares returns the default server middlewares
func DefaultServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}

// DefaultClientMiddlewares returns the default client middlewares
func DefaultClientMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}

// defaultCorsOptions returns default CORS options when none are provided
func defaultCorsOptions() []handlers.CORSOption {
	return []handlers.CORSOption{
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.MaxAge(86400), // 24 hours
	}
}

// corsOptionsFromConfig converts a CORS config to CORS options
func corsOptionsFromConfig(cfg *corsv1.Cors) []handlers.CORSOption {
	var options []handlers.CORSOption

	// Handle origin configuration
	if cfg.GetAllowAnyOrigin() {
		options = append(options, handlers.AllowedOrigins([]string{"*"}))
	} else if len(cfg.GetAllowedOrigins()) > 0 {
		options = append(options, handlers.AllowedOrigins(cfg.GetAllowedOrigins()))
	} else {
		// Default to allow any origin if no specific origin configuration
		options = append(options, handlers.AllowedOrigins([]string{"*"}))
	}

	// Handle method configuration
	if cfg.GetAllowAnyMethod() {
		options = append(options, handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}))
	} else if len(cfg.GetAllowedMethods()) > 0 {
		options = append(options, handlers.AllowedMethods(cfg.GetAllowedMethods()))
	} else {
		// Default methods if none specified
		options = append(options, handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}))
	}

	// Handle header configuration
	if cfg.GetAllowAnyHeader() {
		options = append(options, handlers.AllowedHeaders([]string{"*"}))
	} else if len(cfg.GetAllowedHeaders()) > 0 {
		options = append(options, handlers.AllowedHeaders(cfg.GetAllowedHeaders()))
	} else {
		// Default headers if none specified
		options = append(options, handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}))
	}

	// Add exposed headers if specified
	if len(cfg.GetExposedHeaders()) > 0 {
		options = append(options, handlers.ExposedHeaders(cfg.GetExposedHeaders()))
	}

	// Add credentials option if specified
	if cfg.GetAllowCredentials() {
		options = append(options, handlers.AllowCredentials())
	}

	// Handle max age configuration
	if cfg.GetMaxAge() > 0 {
		options = append(options, handlers.MaxAge(int(cfg.GetMaxAge())))
	} else {
		// Default max age if not specified
		options = append(options, handlers.MaxAge(86400)) // 24 hours
	}

	return options
}

// NewCorsHandler creates a new CORS handler for Kratos HTTP servers
func NewCorsHandler(cfg *corsv1.Cors) transhttp.FilterFunc {
	if cfg == nil || !cfg.GetEnabled() {
		return nil
	}

	// Get CORS options from configuration
	options := corsOptionsFromConfig(cfg)

	// Create and return the CORS filter function
	return handlers.CORS(options...)
}

// adaptServerConfig converts service configuration to Kratos HTTP server options.
func adaptServerConfig(cfg *transportv1.Server) ([]transhttp.ServerOption, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("server config is required for creation")
	}

	httpCfg := cfg.GetHttp()
	if httpCfg == nil {
		return nil, tkerrors.Errorf("http server config is required for creation")
	}

	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))

	// Start with default middlewares and options
	var opts []transhttp.ServerOption

	// Add CORS configuration if needed
	if httpCfg.GetCors() != nil && httpCfg.GetCors().GetEnabled() {
		corsHandler := NewCorsHandler(httpCfg.GetCors())
		opts = append(opts, transhttp.Filter(corsHandler))
	}

	// Add TLS configuration if needed
	if httpCfg.GetTlsConfig() != nil && httpCfg.GetTlsConfig().GetEnabled() {
		tlsConfig, err := tls.NewServerTLSConfig(httpCfg.GetTlsConfig())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation")
		}
		opts = append(opts, transhttp.TLSConfig(tlsConfig))
	}

	// Add network and address configurations
	if httpCfg.Network != "" {
		opts = append(opts, transhttp.Network(httpCfg.Network))
	}

	if httpCfg.Addr != "" {
		opts = append(opts, transhttp.Address(httpCfg.Addr))
	}

	// Configure timeout
	timeout := 5 * time.Second
	if httpCfg.Timeout != nil {
		timeout = httpCfg.Timeout.AsDuration()
	}
	opts = append(opts, transhttp.Timeout(timeout))

	// Handle endpoint configuration
	ll.Debugw("msg", "HTTP", "address", httpCfg.Addr)
	if httpCfg.Addr != "" {
		parsedAddr, err := url.Parse("http://" + httpCfg.Addr)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "failed to parse address for server creation")
		}
		opts = append(opts, transhttp.Endpoint(parsedAddr))
	}

	return opts, nil
}

// adaptClientConfig converts service configuration to Kratos HTTP client options.
func adaptClientConfig(cfg *transportv1.Client) ([]transhttp.ClientOption, error) {
	// 1. Validate the configuration
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation")
	}

	httpCfg := cfg.GetHttp()
	if httpCfg == nil {
		return nil, tkerrors.Errorf("http client config is required for creation")
	}

	var opts []transhttp.ClientOption

	// 2. Processing endpoints
	if httpCfg.Endpoint != "" {
		opts = append(opts, transhttp.WithEndpoint(httpCfg.Endpoint))
	}

	// 3. Processing timeout
	if httpCfg.Timeout != nil {
		opts = append(opts, transhttp.WithTimeout(httpCfg.Timeout.AsDuration()))
	}

	// 4. Handle TLS
	if httpCfg.GetTlsConfig() != nil && httpCfg.GetTlsConfig().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(httpCfg.GetTlsConfig())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		opts = append(opts, transhttp.WithTLSConfig(tlsConfig))
	}

	// 5. Process selectors
	if selectorCfg := httpCfg.GetSelector(); selectorCfg != nil {
		filter, err := selector.NewFilter(selectorCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation")
		}
		opts = append(opts, transhttp.WithNodeFilter(filter))
	}

	return opts, nil
}

// registerPprofHandlers registers pprof handlers to the HTTP server.
func registerPprofHandlers(srv *transhttp.Server) {
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
