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

// defaultBuilder is the default instance of the middlewareBuilder .
var defaultBuilder = NewBuilder()

func init() {
	defaultBuilder.Register("jwt", &jwtFactory{})
	defaultBuilder.Register("circuit_breaker", &circuitBreakerFactory{})
	defaultBuilder.Register("logging", &loggingFactory{})
	defaultBuilder.Register("rate_limit", &rateLimitFactory{})
	defaultBuilder.Register("metadata", &metadataFactory{})
	defaultBuilder.Register("selector", &selectorFactory{})
	defaultBuilder.Register("tracing", &tracingFactory{})
	defaultBuilder.Register("validator", &validatorFactory{})
}

type middlewareBuilder struct {
	factory.Registry[Factory]
}

func (b *middlewareBuilder) BuildClient(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := configure.Apply(&Options{
		Logger: log.DefaultLogger,
	}, options)
	log.Infof("build middleware client")
	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		f, ok := b.Get(ms.GetType())
		if !ok {
			continue
		}
		log.Infof("middleware: %s", ms.GetType())
		m, ok := f.NewMiddlewareClient(ms, option)
		if ok {
			middlewares = append(middlewares, m)
		}

	}
	return middlewares
}

func (b *middlewareBuilder) BuildServer(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	// Create an empty slice of KMiddleware
	var middlewares []KMiddleware

	// If the configuration is nil, return the empty slice
	if cfg == nil {
		return middlewares
	}
	option := configure.Apply(&Options{
		Logger: log.DefaultLogger,
	}, options)
	log.Infof("build middleware server")
	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		f, ok := b.Get(ms.GetType())
		if !ok {
			continue
		}
		log.Infof("middleware: %s", ms.GetType())
		m, ok := f.NewMiddlewareServer(ms, option)
		if ok {
			middlewares = append(middlewares, m)
		}
	}
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

func DefaultBuilder() Builder {
	return defaultBuilder
}
