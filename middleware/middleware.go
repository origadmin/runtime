/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middlewares implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

const Type = "middleware"

type (
	Handler    = middleware.Handler
	Middleware = middleware.Middleware
)

// Chain returns a middleware that executes a chain of middleware.
func Chain(m ...Middleware) Middleware {
	return middleware.Chain(m...)
}

// NewClient creates a new client with the given configuration
func NewClient(cfg *configv1.Middleware, option *config.MiddlewareOption) []Middleware {
	// Create an empty slice of Middleware
	var middlewares []Middleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	// Add the Recovery middleware to the slice
	middlewares = Recovery(middlewares, cfg.EnableRecovery)
	// Add the SecurityClient middleware to the slice
	middlewares = SecurityClient(middlewares, cfg.Security, option)
	// Add the MetadataClient middleware to the slice
	middlewares = MetadataClient(middlewares, cfg.EnableMetadata, cfg.Metadata)
	// Add the TracingClient middleware to the slice
	middlewares = TracingClient(middlewares, cfg.EnableTracing)
	// Add the CircuitBreakerClient middleware to the slice
	middlewares = CircuitBreakerClient(middlewares, cfg.EnableCircuitBreaker)
	// Return the slice of middlewares
	return middlewares
}

// NewServer creates a new server with the given configuration
func NewServer(cfg *configv1.Middleware, option *config.MiddlewareOption) []Middleware {
	// Create an empty slice of Middleware
	var middlewares []Middleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	// Add the Recovery middleware to the slice
	middlewares = Recovery(middlewares, cfg.EnableRecovery)
	// Add the ValidateServer middleware to the slice
	middlewares = ValidateServer(middlewares, cfg.EnableValidate, cfg.Validator)
	// Add the SecurityServer middleware to the slice
	middlewares = SecurityServer(middlewares, cfg.Security, option)
	// Add the TracingServer middleware to the slice
	middlewares = TracingServer(middlewares, cfg.EnableTracing)
	// Add the MetadataServer middleware to the slice
	middlewares = MetadataServer(middlewares, cfg.EnableMetadata, cfg.Metadata)
	// Add the RateLimitServer middleware to the slice
	middlewares = RateLimitServer(middlewares, cfg.RateLimiter)
	// Return the slice of middlewares
	return middlewares
}
