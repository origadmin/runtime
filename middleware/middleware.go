/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
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
	RateLimit      Name = "rate_limit"
	Tracing        Name = "tracing"
	Validator      Name = "validator"
	Optimize       Name = "optimize"
	Selector       Name = "selector"
)

type (
	// Builder is an interface that defines a method for registering a buildImpl.
	Builder interface {
		factory.Registry[Factory]
		BuildClient(*middlewarev1.Middlewares, ...Option) []KMiddleware
		BuildServer(*middlewarev1.Middlewares, ...Option) []KMiddleware
	}

	// Factory is an interface that defines a method for creating a new buildImpl.
	Factory interface {
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(*middlewarev1.MiddlewareConfig, *Options) (KMiddleware, bool)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(*middlewarev1.MiddlewareConfig, *Options) (KMiddleware, bool)
	}
)

type Middleware struct {
}

// NewClient creates a new client with the given configuration
func NewClient(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	return defaultBuilder.BuildClient(cfg, options...)
}

func NewServer(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	return defaultBuilder.BuildServer(cfg, options...)
}
