// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package cors implements CORS middleware for the framework.
package cors

import (
	"net/http"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	corsv1 "github.com/origadmin/runtime/api/gen/go/config/middleware/cors/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware"
)

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
	var corsOptions []handlers.CORSOption

	// Handle origin configuration
	if cfg.GetAllowAnyOrigin() {
		corsOptions = append(corsOptions, handlers.AllowedOrigins([]string{"*"}))
	} else if len(cfg.GetAllowedOrigins()) > 0 {
		corsOptions = append(corsOptions, handlers.AllowedOrigins(cfg.GetAllowedOrigins()))
	} else {
		// Default to allow any origin if no specific origin configuration
		corsOptions = append(corsOptions, handlers.AllowedOrigins([]string{"*"}))
	}

	// Handle method configuration
	if cfg.GetAllowAnyMethod() {
		corsOptions = append(corsOptions, handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}))
	} else if len(cfg.GetAllowedMethods()) > 0 {
		corsOptions = append(corsOptions, handlers.AllowedMethods(cfg.GetAllowedMethods()))
	} else {
		// Default methods if none specified
		corsOptions = append(corsOptions, handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}))
	}

	// Handle header configuration
	if cfg.GetAllowAnyHeader() {
		corsOptions = append(corsOptions, handlers.AllowedHeaders([]string{"*"}))
	} else if len(cfg.GetAllowedHeaders()) > 0 {
		corsOptions = append(corsOptions, handlers.AllowedHeaders(cfg.GetAllowedHeaders()))
	} else {
		// Default headers if none specified
		corsOptions = append(corsOptions, handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}))
	}

	// Add exposed headers if specified
	if len(cfg.GetExposedHeaders()) > 0 {
		corsOptions = append(corsOptions, handlers.ExposedHeaders(cfg.GetExposedHeaders()))
	}

	// Add credentials option if specified
	if cfg.GetAllowCredentials() {
		corsOptions = append(corsOptions, handlers.AllowCredentials())
	}

	// Handle max age configuration
	if cfg.GetMaxAge() > 0 {
		corsOptions = append(corsOptions, handlers.MaxAge(int(cfg.GetMaxAge())))
	} else {
		// Default max age if not specified
		corsOptions = append(corsOptions, handlers.MaxAge(86400)) // 24 hours
	}

	return corsOptions
}

// NewGorillaCors creates a new CORS handler for Gorilla/Mux or standard HTTP servers
func NewGorillaCors(cfg *corsv1.Cors) func(http.Handler) http.Handler {
	if cfg == nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// Get CORS opts from configuration
	opts := corsOptionsFromConfig(cfg)

	return handlers.CORS(opts...)
}

// NewKratosCors creates a new CORS handler for Kratos HTTP servers
func NewKratosCors(cfg *corsv1.Cors) transhttp.FilterFunc {
	if cfg == nil {
		return nil
	}

	// Get CORS opts from configuration
	opts := corsOptionsFromConfig(cfg)

	return handlers.CORS(opts...)
}

// corsFactory implements middleware.Factory interface for CORS middleware
// This is a middleware pattern implementation that can be used with the framework's middleware registry
// Note: CORS is primarily for HTTP servers, so this implementation is designed for HTTP transport
// but follows the standard middleware interface for consistency

type corsFactory struct{}

// NewCorsFactory creates a new CORS middleware factory
func NewCorsFactory() middleware.Factory {
	return &corsFactory{}
}

// NewMiddlewareClient implements middleware.Factory interface
// Since CORS is primarily for servers, this returns nil for clients
func (f *corsFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...options.Option) (middleware.KMiddleware, bool) {
	// CORS is not typically used for client-side middleware
	return func(handler middleware.KHandler) middleware.KHandler {
		// This is a no-op middleware for the standard middleware chain
		// The actual CORS handling is done in the HTTP transport adapter
		return handler
	}, true
}

// NewMiddlewareServer implements middleware.Factory interface
// This creates a server-side CORS middleware handler
func (f *corsFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...options.Option) (middleware.KMiddleware, bool) {
	// Resolve common options
	mwOpts := middleware.FromOptions(opts...)
	logger := log.NewHelper(mwOpts.Logger)

	// Check if CORS is enabled and type is correct
	if !cfg.GetEnabled() || cfg.GetType() != "cors" {
		return nil, false
	}

	// Get the CORS configuration from the middleware config
	corsConfig := cfg.GetCors()
	if corsConfig == nil {
		logger.Errorf("CORS configuration is nil")
		return nil, false
	}
	logger.Debug("[Middleware] CORS server middleware enabled")

	// For HTTP servers, we use the standard middleware chain
	// The actual CORS handling is done in the HTTP transport layer
	// This middleware is a placeholder that follows the framework's pattern
	return func(handler middleware.KHandler) middleware.KHandler {
		// This is a no-op middleware for the standard middleware chain
		// The actual CORS handling is done in the HTTP transport adapter
		return handler
	}, true
}
