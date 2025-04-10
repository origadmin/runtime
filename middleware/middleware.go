/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middlewares implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/goexts/generic/settings"

	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type Middleware struct {
}

// NewClient creates a new client with the given configuration
func NewClient(cfg *middlewarev1.Middleware, ss ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware
	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Options{
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
		m, ok := JwtClient(cfg.GetJwt())
		if ok && cfg.GetSelector().GetEnabled() {
			m = SelectorClient(cfg.GetSelector(), option.MatchFunc, m)
		}
		middlewares = append(middlewares, m)
	}
	// Add the Security middleware to the slice
	return middlewares
}

// NewServer creates a new server with the given configuration
func NewServer(cfg *middlewarev1.Middleware, ss ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Options{
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
		m, ok := JwtServer(cfg.Jwt)
		if ok && cfg.GetSelector().GetEnabled() {
			m = SelectorServer(cfg.GetSelector(), option.MatchFunc, m)
		}
		middlewares = append(middlewares, m)
	}
	//if cfg.GetSelector().GetEnabled() {
	//	return SelectorServer(filter, cfg.GetSelector(), option.MatchFunc)
	//}
	// Return the slice of middlewares
	return middlewares
}
