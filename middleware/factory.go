/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/goexts/generic/settings"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/log"
)

// DefaultBuilder is the default instance of the middlewareBuilder .
var DefaultBuilder = NewBuilder()

func init() {
	DefaultBuilder.Register("jwt", &jwtFactory{})
	DefaultBuilder.Register("circuit_breaker", &circuitBreakerFactory{})
	DefaultBuilder.Register("logging", &loggingFactory{})
	DefaultBuilder.Register("rate_limit", &rateLimitFactory{})
	DefaultBuilder.Register("metadata", &metadataFactory{})
	DefaultBuilder.Register("selector", &selectorFactory{})
	DefaultBuilder.Register("tracing", &tracingFactory{})
	DefaultBuilder.Register("validator", &validatorFactory{})
}

type middlewareBuilder struct {
	factory.Registry[Factory]
}

func (b *middlewareBuilder) BuildClient(cfg *middlewarev1.Middleware, options ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Options{
		Logger: log.DefaultLogger,
	}, options)
	log.Infof("build middleware client")
	for _, em := range cfg.EnabledMiddlewares {
		f, ok := b.Get(em)
		if !ok {
			continue
		}
		log.Infof("middleware: %s", em)
		m, ok := f.NewMiddlewareClient(cfg, option)
		if ok {
			middlewares = append(middlewares, m)
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
	//	jwt, exist := b.Get("jwt")
	//	if !exist {
	//		return middlewares
	//	}
	//	m, ok := jwt.NewMiddlewareClient(cfg, option)
	//	if ok && cfg.GetSelector().GetEnabled() {
	//		m = SelectorServer(cfg.GetSelector(), option.MatchFunc, m)
	//	}
	//	middlewares = append(middlewares, m)
	//}
	return middlewares
}

func (b *middlewareBuilder) BuildServer(cfg *middlewarev1.Middleware, options ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := settings.Apply(&Options{
		Logger: log.DefaultLogger,
	}, options)
	log.Infof("build middleware server")
	for _, em := range cfg.EnabledMiddlewares {
		f, ok := b.Get(em)
		if !ok {
			continue
		}
		log.Infof("middleware: %s", em)
		m, ok := f.NewMiddlewareServer(cfg, option)
		if ok {
			middlewares = append(middlewares, m)
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
	//	jwt, exist := b.Get("jwt")
	//	if !exist {
	//		return middlewares
	//	}
	//	m, ok := jwt.NewMiddlewareServer(cfg, option)
	//	if ok && cfg.GetSelector().GetEnabled() {
	//		m = SelectorServer(cfg.GetSelector(), option.MatchFunc, m)
	//	}
	//	middlewares = append(middlewares, m)
	//}
	return middlewares
}

func (b *middlewareBuilder) RegisterMiddlewareBuilder(name string, factory Factory) {
	b.Register(name, factory)
}

func NewBuilder() Builder {
	return &middlewareBuilder{
		Registry: factory.New[Factory](),
	}
}
