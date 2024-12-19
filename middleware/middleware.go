/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middlewares implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goexts/generic/settings"

	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
)

const Type = "middleware"

type (
	Handler    = middleware.Handler
	Middleware = middleware.Middleware
)

//type Middleware struct {
//
//}

// Chain returns a middleware that executes a chain of middleware.
func Chain(m ...Middleware) Middleware {
	return middleware.Chain(m...)
}

// NewClient creates a new client with the given configuration
func NewClient(cfg *middlewarev1.Middleware, ss ...OptionSetting) []Middleware {
	// Create an empty slice of Middleware
	var middlewares []Middleware
	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Option{
		Logger: log.DefaultLogger,
	}, ss)

	if cfg.Logging {
		// Add the LoggingClient middleware to the slice
		middlewares = LoggingClient(middlewares, option.Logger)
	}
	if cfg.Recovery {
		// Add the Recovery middleware to the slice
		middlewares = Recovery(middlewares)
	}
	if cfg.GetMetadata().GetEnabled() {
		// Add the MetadataClient middleware to the slice
		middlewares = MetadataClient(middlewares, cfg.GetMetadata())
	}
	if cfg.Tracing {
		// Add the TracingClient middleware to the slice
		middlewares = TracingClient(middlewares)
	}
	if cfg.CircuitBreaker {
		// Add the CircuitBreakerClient middleware to the slice
		middlewares = CircuitBreakerClient(middlewares)
	}
	if cfg.GetJwt().GetEnabled() {
		middlewares = JwtClient(middlewares, cfg.GetJwt())
	}
	if cfg.GetSelector().GetEnabled() {
		return SelectorClient(middlewares, cfg.GetSelector(), option.MatchFunc)
	}
	// Add the Security middleware to the slice
	return middlewares
}

// NewServer creates a new server with the given configuration
func NewServer(cfg *middlewarev1.Middleware, ss ...OptionSetting) []Middleware {
	// Create an empty slice of Middleware
	var middlewares []Middleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Option{
		Logger: log.DefaultLogger,
	}, ss)
	if cfg.Logging {
		middlewares = LoggingServer(middlewares, option.Logger)
	}
	if cfg.Recovery {
		// Add the Recovery middleware to the slice
		middlewares = Recovery(middlewares)
	}
	if cfg.GetValidator().GetEnabled() {
		// Add the ValidateServer middleware to the slice
		middlewares = ValidateServer(middlewares, cfg.Validator)
	}
	if cfg.Tracing {
		// Add the TracingServer middleware to the slice
		middlewares = TracingServer(middlewares)
	}
	if cfg.GetMetadata().GetEnabled() {
		// Add the MetadataServer middleware to the slice
		middlewares = MetadataServer(middlewares, cfg.Metadata)
	}
	if cfg.GetRateLimiter().GetEnabled() {
		// Add the RateLimitServer middleware to the slice
		middlewares = RateLimitServer(middlewares, cfg.RateLimiter)
	}
	if cfg.GetJwt().GetEnabled() {
		middlewares = JwtServer(middlewares, cfg.Jwt)
	}
	if cfg.GetSelector().GetEnabled() {
		return SelectorServer(middlewares, cfg.GetSelector(), option.MatchFunc)
	}
	// Return the slice of middlewares
	return middlewares
}
