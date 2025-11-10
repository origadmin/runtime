/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

// Name is the name of a middleware.
type Name string

// Middleware names.
const (
	Jwt                 Name = "jwt"
	CircuitBreaker      Name = "circuit_breaker"
	Logging             Name = "logging"
	Metadata            Name = "metadata"
	RateLimiter         Name = "rate_limiter"
	Tracing             Name = "tracing"
	Validator           Name = "validator"
	Optimize            Name = "optimize"
	Recovery            Name = "recovery"
	Selector            Name = "selector"
	Security            Name = "security"
	DeclarativeSecurity Name = "declarative_security"
	Cors                Name = "cors"
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
		BuildClientMiddlewares(*middlewarev1.Middlewares, ...Option) []KMiddleware
		BuildServerMiddlewares(*middlewarev1.Middlewares, ...Option) []KMiddleware
	}

	// Factory is an interface that defines a method for creating a new buildImpl.
	// It receives the middleware-specific Protobuf configuration and the generic Option slice.
	// Each factory is responsible for parsing the options it cares about (e.g., by using log.FromOptions).
	Factory interface {
		// NewMiddlewareClient builds a client-side middleware.
		NewMiddlewareClient(*middlewarev1.Middleware, ...Option) (KMiddleware, bool)
		// NewMiddlewareServer builds a server-side middleware.
		NewMiddlewareServer(*middlewarev1.Middleware, ...Option) (KMiddleware, bool)
	}
)

// NewClient creates a new client with the given configuration.
// This function is a convenience wrapper around the default builder.
func NewClient(mc *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
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
func NewServer(mc *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
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

// BuildClientMiddlewares build a client middleware chain
func BuildClientMiddlewares(middlewares *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildClientMiddlewares(middlewares, opts...)
}

// BuildServerMiddlewares build a server middleware chain
func BuildServerMiddlewares(middlewares *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildServerMiddlewares(middlewares, opts...)
}
