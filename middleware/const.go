/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/goexts/generic/settings"

	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
	"github.com/origadmin/runtime/log"
)

const Type = "middleware"

type (
	KHandler    = middleware.Handler
	KMiddleware = middleware.Middleware
)

// Chain returns a middleware that executes a chain of middleware.
func Chain(m ...KMiddleware) KMiddleware {
	return middleware.Chain(m...)
}

// NewClient creates a new client with the given configuration
func NewClient(cfg *middlewarev1.Middleware, ss ...OptionSetting) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware
	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Option{
		Logger: log.DefaultLogger,
	}, ss)

	filter := WithSelector(cfg.GetSelector(), option.MatchFunc)
	if cfg.Logging {
		// Add the LoggingClient middleware to the slice
		filter = LoggingClient(filter, option.Logger)
	}
	if cfg.Recovery {
		// Add the Recovery middleware to the slice
		filter = Recovery(filter)
	}
	if cfg.GetMetadata().GetEnabled() {
		// Add the MetadataClient middleware to the slice
		filter = MetadataClient(filter, cfg.GetMetadata())
	}
	if cfg.Tracing {
		// Add the TracingClient middleware to the slice
		filter = TracingClient(filter)
	}
	if cfg.CircuitBreaker {
		// Add the CircuitBreakerClient middleware to the slice
		filter = CircuitBreakerClient(filter)
	}
	if cfg.GetJwt().GetEnabled() {
		filter = JwtClient(filter, cfg.GetJwt())
	}
	//if cfg.GetSelector().GetEnabled() {
	//	return SelectorClient(filter, cfg.GetSelector(), option.MatchFunc)
	//}
	// Add the Security middleware to the slice
	return filter.Build(cfg.Selector, selector.Client)
}

// NewServer creates a new server with the given configuration
func NewServer(cfg *middlewarev1.Middleware, ss ...OptionSetting) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Option{
		Logger: log.DefaultLogger,
	}, ss)
	filter := WithSelector(cfg.GetSelector(), option.MatchFunc)
	if cfg.Logging {
		filter = LoggingServer(filter, option.Logger)
	}
	if cfg.Recovery {
		// Add the Recovery middleware to the slice
		filter = Recovery(filter)
	}
	if cfg.GetValidator().GetEnabled() {
		// Add the ValidateServer middleware to the slice
		filter = ValidateServer(filter, cfg.Validator)
	}
	if cfg.Tracing {
		// Add the TracingServer middleware to the slice
		filter = TracingServer(filter)
	}
	if cfg.GetMetadata().GetEnabled() {
		// Add the MetadataServer middleware to the slice
		filter = MetadataServer(filter, cfg.Metadata)
	}
	if cfg.GetRateLimiter().GetEnabled() {
		// Add the RateLimitServer middleware to the slice
		filter = RateLimitServer(filter, cfg.RateLimiter)
	}
	if cfg.GetJwt().GetEnabled() {
		filter = JwtServer(filter, cfg.Jwt)
	}
	//if cfg.GetSelector().GetEnabled() {
	//	return SelectorServer(filter, cfg.GetSelector(), option.MatchFunc)
	//}
	// Return the slice of middlewares
	return filter.Build(cfg.Selector, selector.Server)
}
