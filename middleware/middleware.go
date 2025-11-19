/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
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

// Factory is an interface that defines a method for creating a new buildImpl.
// It receives the middleware-specific Protobuf configuration and the generic Option slice.
// Each factory is responsible for parsing the options it cares about (e.g., by using log.FromOptions).
type Factory interface {
	// NewMiddlewareClient builds a client-side middleware.
	NewMiddlewareClient(*middlewarev1.Middleware, ...Option) (KMiddleware, bool)
	// NewMiddlewareServer builds a server-side middleware.
	NewMiddlewareServer(*middlewarev1.Middleware, ...Option) (KMiddleware, bool)
}


