/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

// Name is the name of a middleware.
type Name string

// Middleware names.
const (
	Jwt            Name = "jwt"
	CircuitBreaker Name = "circuit_breaker"
	Logging        Name = "logging"
	Metadata       Name = "metadata"
	RateLimiter    Name = "rate_limiter"
	Tracing        Name = "tracing"
	Validator      Name = "validator"
	Optimize       Name = "optimize"
	Selector       Name = "selector"
)

// Carrier is a struct that holds the middlewares for client and server.
type Carrier struct {
	Clients map[string]KMiddleware
	Servers map[string]KMiddleware
}

type (
	// Builder is an interface that defines a method for registering a buildImpl.
	Builder interface {
		factory.Registry[Factory]
		BuildClient(*middlewarev1.Middlewares, ...Option) []KMiddleware
		BuildServer(*middlewarev1.Middlewares, ...Option) []KMiddleware
	}

	// Factory is an interface that defines a method for creating a new buildImpl.
	// It receives the middleware-specific Protobuf configuration and the generic Option slice.
	// Each factory is responsible for parsing the options it cares about (e.g., by using log.FromOptions).
	Factory interface {
		// NewMiddlewareClient builds a client-side middleware.
		NewMiddlewareClient(*middlewarev1.MiddlewareConfig, ...Option) (KMiddleware, bool)
		// NewMiddlewareServer builds a server-side middleware.
		NewMiddlewareServer(*middlewarev1.MiddlewareConfig, ...Option) (KMiddleware, bool)
	}
)

// NewClient creates a new client with the given configuration.
// This function is a convenience wrapper around the default builder.
func NewClient(mc *middlewarev1.MiddlewareConfig, opts ...Option) (KMiddleware, bool) {
	if mc == nil || !mc.GetEnabled() {
		return nil, false
	}
	// Get the middleware name.
	middlewareName := mc.GetType()
	f, ok := defaultBuilder.Get(middlewareName)
	if !ok {
		return nil, false
	}
	return f.NewMiddlewareClient(mc, opts...)
}

// NewServer creates a new server with the given configuration.
// This function is a convenience wrapper around the default builder.
func NewServer(mc *middlewarev1.MiddlewareConfig, opts ...Option) (KMiddleware, bool) {
	if mc == nil || !mc.GetEnabled() {
		return nil, false
	}
	// Get the middleware name.
	middlewareName := mc.GetType()
	f, ok := defaultBuilder.Get(middlewareName)
	if !ok {
		return nil, false
	}
	return f.NewMiddlewareServer(mc, opts...)
}

// BuildClient build a client middleware chain
func BuildClient(middlewares *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildClient(middlewares, opts...)
}

// BuildServer build a server middleware chain
func BuildServer(middlewares *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildServer(middlewares, opts...)
}
