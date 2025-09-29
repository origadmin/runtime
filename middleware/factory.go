/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/optionutil"
)

// defaultBuilder is the default instance of the middlewareBuilder .
var defaultBuilder = NewBuilder()

func init() {
	// The factories will be registered here once they are updated to the new interface.
	// optimizeFactory is removed from here as it's not a formal feature and should be registered by the user.
	// All other factories will be uncommented as they are updated.
	defaultBuilder.Register(Jwt, &jwtFactory{})
	defaultBuilder.Register(CircuitBreaker, &circuitBreakerFactory{})
	defaultBuilder.Register(Logging, &loggingFactory{})
	defaultBuilder.Register(RateLimit, &rateLimitFactory{})
	defaultBuilder.Register(Metadata, &metadataFactory{})
	defaultBuilder.Register(Selector, &selectorFactory{})
	defaultBuilder.Register(Tracing, &tracingFactory{})
	defaultBuilder.Register(Validator, &validatorFactory{})
}

// middlewareBuilder is a builder for creating middleware chains.
type middlewareBuilder struct {
	factory.Registry[Factory]
}

// BuildClient builds a client-side middleware chain from the given configuration.
func (b *middlewareBuilder) BuildClient(cfg *middlewarev1.Middlewares, opts ...options.Option) []KMiddleware {
	var middlewares []KMiddleware
	if cfg == nil {
		return middlewares
	}

	ctx, opt := FromOptions(opts...)
	var logger log.Logger
	if opt != nil && opt.Logger != nil {
		logger = opt.Logger
	} else {
		logger = log.FromContext(ctx)
	}

	// This logger is for the factory's own internal logging, not for the middlewares themselves.
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

		// Pass the raw options slice directly to the factory.
		// The factory is responsible for parsing the options it needs.
		m, ok := f.NewMiddlewareClient(ms, optionutil.WithContext(ctx), withOptions(opt))
		if ok {
			middlewares = append(middlewares, m)
		}
	}
	return middlewares
}

// BuildServer builds a server-side middleware chain from the given configuration.
func (b *middlewareBuilder) BuildServer(cfg *middlewarev1.Middlewares, opts ...options.Option) []KMiddleware {
	var middlewares []KMiddleware
	if cfg == nil {
		return middlewares
	}

	ctx, opt := FromOptions(opts...)
	var logger log.Logger
	if opt != nil && opt.Logger != nil {
		logger = opt.Logger
	} else {
		logger = log.FromContext(ctx)
	}

	// This logger is for the factory's own internal logging.
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

		// Pass the raw options slice directly to the factory.
		// The factory is responsible for parsing the options it needs.
		m, ok := f.NewMiddlewareServer(ms, optionutil.WithContext(ctx), withOptions(opt))
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

// NewBuilder creates a new middleware builder.
func NewBuilder() Builder {
	return &middlewareBuilder{
		Registry: factory.New[Factory](),
	}
}
