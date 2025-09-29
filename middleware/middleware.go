/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middlewares implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/goexts/generic/configure"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/log"
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
	//if defaultBuilder != nil {
	//	return defaultBuilder.BuildClient(cfg, options...)
	//}
	//return buildClientMiddlewares(cfg, options...)
}

func NewServer(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	return defaultBuilder.BuildServer(cfg, options...)
	//if defaultBuilder != nil {
	//	return defaultBuilder.BuildServer(cfg, options...)
	//}
	//return buildServerMiddlewares(cfg, options...)
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
		case "jwt":
			m, ok := JwtClient(middlewareConfig.GetJwt())
			if ok && middlewareConfig.GetSelector().GetEnabled() {
				m = SelectorClient(middlewareConfig.GetSelector(), option.MatchFunc, m)
			}
			middlewares = append(middlewares, m)
		case "circuit_breaker":
			middlewares = CircuitBreakerClient(middlewares)
		case "logging":
			middlewares = LoggingClient(middlewares, option.Logger)
		case "metadata":
			middlewares = MetadataClient(middlewares, middlewareConfig.GetMetadata())
		case "rate_limit":
		//middlewares = RateLimitClient(middlewares, cfg.GetRateLimiter())
		case "tracing":
			middlewares = TracingClient(middlewares)
		case "validator":
			//middlewares = ValidateClient(middlewares, cfg.GetValidator())
		}
	}
	//if cfg.Logging {
	//	// Add the LoggingClient middleware to the slice
	//	middlewares = LoggingClient(middlewares, option.Logger)
	//}
	//if cfg.Recovery {
	//	// Add the Recovery middleware to the slice
	//	middlewares = Recovery(middlewares)
	//}
	//if cfg.GetMetadata().GetEnabled() {
	//	// Add the MetadataClient middleware to the slice
	//	middlewares = MetadataClient(middlewares, cfg.GetMetadata())
	//}
	//if cfg.Tracing {
	//	// Add the TracingClient middleware to the slice
	//	middlewares = TracingClient(middlewares)
	//}
	//if cfg.CircuitBreaker {
	//	// Add the CircuitBreakerClient middleware to the slice
	//	middlewares = CircuitBreakerClient(middlewares)
	//}
	//if cfg.GetJwt().GetEnabled() {
	//	m, ok := JwtClient(cfg.GetJwt())
	//	if ok && cfg.GetSelector().GetEnabled() {
	//		m = SelectorClient(cfg.GetSelector(), option.MatchFunc, m)
	//	}
	//	middlewares = append(middlewares, m)
	//}
	// Add the Security middleware to the slice
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
		case "jwt":
			m, ok := JwtServer(ms.GetJwt())
			if ok && ms.GetSelector().GetEnabled() {
				m = SelectorServer(ms.GetSelector(), option.MatchFunc, m)
			}
			middlewares = append(middlewares, m)
		case "circuit_breaker":
			//middlewares = CircuitBreakerServer(middlewares)
		case "logging":
			middlewares = LoggingServer(middlewares, option.Logger)
		case "metadata":
			middlewares = MetadataServer(middlewares, ms.GetMetadata())
		case "rate_limit":
			middlewares = RateLimitServer(middlewares, ms.GetRateLimiter())
		case "tracing":
			middlewares = TracingServer(middlewares)
		case "validator":
			middlewares = ValidateServer(middlewares, ms.GetValidator())
		}
	}
	//if cfg.Logging {
	//	middlewares = LoggingServer(middlewares, option.Logger)
	//}
	//if cfg.Recovery {
	//	// Add the Recovery middleware to the slice
	//	middlewares = Recovery(middlewares)
	//}
	//if cfg.GetValidator().GetEnabled() {
	//	// Add the ValidateServer middleware to the slice
	//	middlewares = ValidateServer(middlewares, cfg.Validator)
	//}
	//if cfg.Tracing {
	//	// Add the TracingServer middleware to the slice
	//	middlewares = TracingServer(middlewares)
	//}
	//if cfg.GetMetadata().GetEnabled() {
	//	// Add the MetadataServer middleware to the slice
	//	middlewares = MetadataServer(middlewares, cfg.Metadata)
	//}
	//if cfg.GetRateLimiter().GetEnabled() {
	//	// Add the RateLimitServer middleware to the slice
	//	middlewares = RateLimitServer(middlewares, cfg.RateLimiter)
	//}
	//if cfg.GetJwt().GetEnabled() {
	//	m, ok := JwtServer(cfg.Jwt)
	//	if ok && cfg.GetSelector().GetEnabled() {
	//		m = SelectorServer(cfg.GetSelector(), option.MatchFunc, m)
	//	}
	//	middlewares = append(middlewares, m)
	//}
	return middlewares
}
