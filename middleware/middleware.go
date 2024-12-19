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

	f := MakeFilter(cfg.GetSelector().GetNames())
	if cfg.Logging {
		// Add the LoggingClient middleware to the slice
		f = LoggingClient(f, option.Logger)
	}
	if cfg.Recovery {
		// Add the Recovery middleware to the slice
		f = Recovery(f)
	}
	if cfg.GetMetadata().GetEnabled() {
		// Add the MetadataClient middleware to the slice
		f = MetadataClient(f, cfg.GetMetadata())
	}
	if cfg.Tracing {
		// Add the TracingClient middleware to the slice
		f = TracingClient(f)
	}
	if cfg.CircuitBreaker {
		// Add the CircuitBreakerClient middleware to the slice
		f = CircuitBreakerClient(f)
	}
	if cfg.GetJwt().GetEnabled() {
		f = JwtClient(f, cfg.GetJwt())
	}
	if cfg.GetSelector().GetEnabled() {
		return SelectorClient(f, cfg.GetSelector(), option.MatchFunc)
	}
	// Add the Security middleware to the slice
	return f.All()
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
	f := MakeFilter(cfg.GetSelector().GetNames())
	if cfg.Logging {
		f = LoggingServer(f, option.Logger)
	}
	if cfg.Recovery {
		// Add the Recovery middleware to the slice
		f = Recovery(f)
	}
	if cfg.GetValidator().GetEnabled() {
		// Add the ValidateServer middleware to the slice
		f = ValidateServer(f, cfg.Validator)
	}
	if cfg.Tracing {
		// Add the TracingServer middleware to the slice
		f = TracingServer(f)
	}
	if cfg.GetMetadata().GetEnabled() {
		// Add the MetadataServer middleware to the slice
		f = MetadataServer(f, cfg.Metadata)
	}
	if cfg.GetRateLimiter().GetEnabled() {
		// Add the RateLimitServer middleware to the slice
		f = RateLimitServer(f, cfg.RateLimiter)
	}
	if cfg.GetJwt().GetEnabled() {
		f = JwtServer(f, cfg.Jwt)
	}
	if cfg.GetSelector().GetEnabled() {
		return SelectorServer(f, cfg.GetSelector(), option.MatchFunc)
	}
	// Return the slice of middlewares
	return f.All()
}
