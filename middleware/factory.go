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
)

// defaultBuilder is the default instance of the middlewareBuilder .
var defaultBuilder = NewBuilder()

func init() {
	// The factories will be registered here once they are updated to the new interface.
	// optimizeFactory is removed from here as it's not a formal feature and should be registered by the user.
	// All other factories will be uncommented as they are updated.
	Register(Jwt, &jwtFactory{})
	Register(CircuitBreaker, &circuitBreakerFactory{})
	Register(Logging, &loggingFactory{})
	Register(RateLimiter, &rateLimitFactory{})
	Register(Metadata, &metadataFactory{})
	Register(Selector, &selectorFactory{})
	Register(Tracing, &tracingFactory{})
	Register(Validator, &validatorFactory{})
}

// middlewareBuilder is a builder for creating middleware chains.
type middlewareBuilder struct {
	factory.Registry[Factory]
}

// BuildClient builds the client middleware chain
func (b *middlewareBuilder) BuildClient(cfg *middlewarev1.Middlewares, opts ...options.Option) []KMiddleware {
	var middlewares []KMiddleware
	var selectorConfigs []*middlewarev1.MiddlewareConfig
	if cfg == nil {
		return middlewares
	}

	opt := FromOptions(opts...)
	if opt == nil {
		return middlewares
	}

	logger := opt.Logger
	if logger == nil {
		logger = log.DefaultLogger
	}

	// Create a middleware carrier for context propagation
	carrier := &Carrier{
		Clients: make(map[string]KMiddleware),
		Servers: make(map[string]KMiddleware),
	}

	helper := log.NewHelper(logger)
	helper.Info("building client middlewares")

	// First pass: separate regular middlewares and selector configs
	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		middlewareName := ms.GetType()
		// Defer selector middleware configs for later processing
		if middlewareName == string(Selector) {
			selectorConfigs = append(selectorConfigs, ms)
			continue
		}
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown client middleware: %s", middlewareName)
			continue
		}

		helper.Infof("enabling client middleware: %s", middlewareName)

		// Create middleware
		m, ok := f.NewMiddlewareClient(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
			// Add the created middleware into the carrier
			carrier.Clients[middlewareName] = m
		}
	}
	// Attach the carrier into options context
	opts = append(opts, WithCarrier(carrier))

	// Second pass: process selector middleware configs
	for _, ms := range selectorConfigs {
		middlewareName := ms.GetType()
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown client middleware: %s", middlewareName)
			continue
		}

		helper.Infof("enabling client middleware: %s", middlewareName)

		// Create selector middleware (can access previously created middlewares now)
		m, ok := f.NewMiddlewareClient(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
		}
	}

	return middlewares
}

// BuildServer builds the server middleware chain (similar to BuildClient)
func (b *middlewareBuilder) BuildServer(cfg *middlewarev1.Middlewares, opts ...options.Option) []KMiddleware {
	var middlewares []KMiddleware
	var selectorConfigs []*middlewarev1.MiddlewareConfig
	if cfg == nil {
		return middlewares
	}

	opt := FromOptions(opts...)
	if opt == nil {
		return middlewares
	}

	logger := opt.Logger
	if logger == nil {
		logger = log.DefaultLogger
	}

	// Create a middleware carrier for context propagation
	carrier := &Carrier{
		Clients: make(map[string]KMiddleware),
		Servers: make(map[string]KMiddleware),
	}

	helper := log.NewHelper(logger)
	helper.Info("building server middlewares")

	// First pass: separate regular middlewares and selector configs
	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		middlewareName := ms.GetType()
		// Defer selector middleware configs for later processing
		if middlewareName == string(Selector) {
			selectorConfigs = append(selectorConfigs, ms)
			continue
		}
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown server middleware: %s", middlewareName)
			continue
		}

		// Create middleware
		m, ok := f.NewMiddlewareServer(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
			// Add the created middleware into the carrier
			carrier.Servers[middlewareName] = m
		}
	}
	// Attach the carrier into options context
	opts = append(opts, WithCarrier(carrier))

	// Second pass: process selector middleware configs
	for _, ms := range selectorConfigs {
		middlewareName := ms.GetType()
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown server middleware: %s", middlewareName)
			continue
		}

		// Create selector middleware (can access previously created middlewares now)
		m, ok := f.NewMiddlewareServer(ms, opts...)
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
