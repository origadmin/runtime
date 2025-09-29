/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/goexts/generic/configure"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/log"
)

// Middleware names.
const (
	Jwt            = "jwt"
	CircuitBreaker = "circuit_breaker"
	Logging        = "logging"
	Metadata       = "metadata"
	RateLimit      = "rate_limit"
	Tracing        = "tracing"
	Validator      = "validator"
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

func buildClientMiddlewares(cfg *middlewarev1.Middlewares, ss ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware
	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := configure.Apply(&Options{
		Logger: log.DefaultLogger,
	}, ss)
	for _, middlewareConfig := range cfg.GetMiddlewares() {
		if !middlewareConfig.GetEnabled() {
			continue
		}
		switch middlewareConfig.GetType() {
		case Jwt:
			m, ok := JwtClient(middlewareConfig.GetJwt())
			if ok && middlewareConfig.GetSelector().GetEnabled() {
				m = SelectorClient(middlewareConfig.GetSelector(), option.MatchFunc, m)
			}
			middlewares = append(middlewares, m)
		case CircuitBreaker:
			middlewares = CircuitBreakerClient(middlewares)
		case Logging:
			middlewares = LoggingClient(middlewares, option.Logger)
		case Metadata:
			middlewares = MetadataClient(middlewares, middlewareConfig.GetMetadata())
		case RateLimit:
		//middlewares = RateLimitClient(middlewares, cfg.GetRateLimiter())
		case Tracing:
			middlewares = TracingClient(middlewares)
		case Validator:
			//middlewares = ValidateClient(middlewares, cfg.GetValidator())
		}
	}
	return middlewares
}

// NewServer creates a new server with the given configuration
func buildServerMiddlewares(cfg *middlewarev1.Middlewares, ss ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := configure.Apply(&Options{
		Logger: log.DefaultLogger,
	}, ss)
	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		switch ms.GetType() {
		case Jwt:
			m, ok := JwtServer(ms.GetJwt())
			if ok && ms.GetSelector().GetEnabled() {
				m = SelectorServer(ms.GetSelector(), option.MatchFunc, m)
			}
			middlewares = append(middlewares, m)
		case CircuitBreaker:
			//middlewares = CircuitBreakerServer(middlewares)
		case Logging:
			middlewares = LoggingServer(middlewares, option.Logger)
		case Metadata:
			middlewares = MetadataServer(middlewares, ms.GetMetadata())
		case RateLimit:
			middlewares = RateLimitServer(middlewares, ms.GetRateLimiter())
		case Tracing:
			middlewares = TracingServer(middlewares)
		case Validator:
			middlewares = ValidateServer(middlewares, ms.GetValidator())
		}
	}
	return middlewares
}
