/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middlewares implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/builder"
	"github.com/origadmin/runtime/log"
)

type (
	// Builder is an interface that defines a method for registering a buildImpl.
	Builder interface {
		builder.Builder[Factory]
		Factory
	}

	Provider interface {
		BuildClient(cfg *middlewarev1.Middleware) KMiddleware
		BuildServer(cfg *middlewarev1.Middleware) KMiddleware
	}

	// Factory is an interface that defines a method for creating a new buildImpl.
	Factory interface {
		// NewMiddlewaresClient build middleware
		NewMiddlewaresClient([]KMiddleware, *configv1.Customize, ...Option) []KMiddleware
		// NewMiddlewaresServer build middleware
		NewMiddlewaresServer([]KMiddleware, *configv1.Customize, ...Option) []KMiddleware
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(string, *configv1.Customize_Config, ...Option) (KMiddleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(string, *configv1.Customize_Config, ...Option) (KMiddleware, error)
	}
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
	return middlewares
}
