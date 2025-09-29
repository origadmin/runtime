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
	Register(Jwt, &jwtFactory{})
	Register(CircuitBreaker, &circuitBreakerFactory{})
	Register(Logging, &loggingFactory{})
	Register(RateLimit, &rateLimitFactory{})
	Register(Metadata, &metadataFactory{})
	Register(Selector, &selectorFactory{})
	Register(Tracing, &tracingFactory{})
	Register(Validator, &validatorFactory{})
	//Register(Optimize, &optimize.Factory{})
}

// middlewareBuilder is a builder for creating middleware chains.
type middlewareBuilder struct {
	factory.Registry[Factory]
}

// BuildClient builds a client-side middleware chain from the given configuration.
func (b *middlewareBuilder) BuildClient(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	var middlewares []KMiddleware
	if cfg == nil {
		return middlewares
	}

	// Apply options to get the logger and other settings.
	option := configure.Apply(&Options{}, options)

	var logger log.Logger
	if option.Context != nil {
		logger = log.FromContext(option.Context)
	} else {
		logger = log.DefaultLogger
	}

	helper := log.NewHelper(logger)
	helper.Info("building client middlewares")

	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		middlewareName := ms.GetType()
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown client middleware: %s", middlewareName)
			continue
		}

		helper.Infof("enabling client middleware: %s", middlewareName)
		m, ok := f.NewMiddlewareClient(ms, option)
		if ok {
			middlewares = append(middlewares, m)
		}
	}
	return middlewares
}

// BuildServer builds a server-side middleware chain from the given configuration.
func (b *middlewareBuilder) BuildServer(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	var middlewares []KMiddleware
	if cfg == nil {
		return middlewares
	}

	// Apply options to get the logger and other settings.
	option := configure.Apply(&Options{}, options)

	var logger log.Logger
	if option.Context != nil {
		logger = log.FromContext(option.Context)
	} else {
		logger = log.DefaultLogger
	}

	helper := log.NewHelper(logger)
	helper.Info("building server middlewares")

	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		middlewareName := ms.GetType()
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown server middleware: %s", middlewareName)
			continue
		}

		helper.Infof("enabling server middleware: %s", middlewareName)
		m, ok := f.NewMiddlewareServer(ms, option)
		if ok {
			middlewares = append(middlewares, m)
		}
	}
	return middlewares
}

// Register registers a middleware factory with the given name.
func Register(name Name, factory Factory) {
	defaultBuilder.Register(string(name), factory)
}

// BuildClient builds middlewares for clients.
func BuildClient(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	return defaultBuilder.BuildClient(cfg, options...)
}

// BuildServer builds middlewares for servers.
func BuildServer(cfg *middlewarev1.Middlewares, options ...Option) []KMiddleware {
	return defaultBuilder.BuildServer(cfg, options...)
}

// NewBuilder creates a new middleware builder.
func NewBuilder() Builder {
	return &middlewareBuilder{
		Registry: factory.New[Factory](),
	}
}

// DefaultBuilder returns the default middleware builder.
func DefaultBuilder() Builder {
	return defaultBuilder
}
