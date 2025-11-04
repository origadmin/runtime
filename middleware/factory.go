/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/log"
)

// defaultBuilder is the default instance of the middlewareBuilder .
var defaultBuilder = NewBuilder()

func init() {
	// The factories will be registered here once they are updated to the new interface.
	// optimizeFactory is removed from here as it's not a formal feature and should be registered by the user.
	// All other factories will be uncommented as they are updated.
	Register(Recovery, &recoveryFactory{})
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

// BuildClientMiddlewares builds the client middleware chain
func (b *middlewareBuilder) BuildClientMiddlewares(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	var middlewares []KMiddleware
	var selectorConfigs []*middlewarev1.Middleware
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
	helper.Debug("building client middlewares")

	// First pass: separate regular middlewares and selector configs
	for _, ms := range cfg.GetConfigs() {
		if !ms.GetEnabled() {
			continue
		}
		// Use Name if provided, otherwise fall back to Type
		middlewareType := ms.GetType()
		middlewareName := ms.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		// Defer selector middleware configs for later processing
		if middlewareType == string(Selector) {
			selectorConfigs = append(selectorConfigs, ms)
			continue
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			helper.Warnf("unknown client middleware type: %s", middlewareType)
			continue
		}

		helper.Infof("enabling client middleware: %s (type: %s)", middlewareName, middlewareType)

		// Create middleware
		m, ok := f.NewMiddlewareClient(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
			// Add the created middleware into the carrier using the middlewareName
			carrier.Clients[middlewareName] = m
		}
	}
	// Attach the carrier into options context
	opts = append(opts, WithCarrier(carrier))

	// Second pass: process selector middleware configs
	for _, ms := range selectorConfigs {
		middlewareType := ms.GetType()
		middlewareName := ms.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			helper.Warnf("unknown client middleware type: %s", middlewareType)
			continue
		}

		helper.Infof("enabling client middleware: %s (type: %s)", middlewareName, middlewareType)

		// Create selector middleware (can access previously created middlewares now)
		m, ok := f.NewMiddlewareClient(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
			// Add the created middleware into the carrier using the middlewareName
			carrier.Clients[middlewareName] = m
		}
	}

	return middlewares
}

// BuildServerMiddlewares builds the server middleware chain (similar to BuildClient)
func (b *middlewareBuilder) BuildServerMiddlewares(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	var middlewares []KMiddleware
	var selectorConfigs []*middlewarev1.Middleware
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

	helper := log.NewHelper(logger)
	helper.Debug("building server middlewares")

	// Create a middleware carrier for context propagation
	carrier := &Carrier{
		Clients: make(map[string]KMiddleware),
		Servers: make(map[string]KMiddleware),
	}

	// First pass: separate regular middlewares and selector configs
	for _, ms := range cfg.GetConfigs() {
		if !ms.GetEnabled() {
			continue
		}

		// Use Name if provided, otherwise fall back to Type
		middlewareType := ms.GetType()
		middlewareName := ms.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		// Defer selector middleware configs for later processing
		if middlewareType == string(Selector) {
			selectorConfigs = append(selectorConfigs, ms)
			continue
		}

		f, ok := b.Get(middlewareName)
		if !ok {
			f, ok = b.Get(middlewareType)
			if !ok {
				helper.Warnf("unknown server middleware type: %s", middlewareType)
				continue
			}
		}

		helper.Debugf("enabling server middleware: %s (type: %s)", middlewareName, middlewareType)

		// Create middleware
		m, ok := f.NewMiddlewareServer(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
			// Add the created middleware into the carrier
			carrier.Servers[middlewareName] = m
		}
	}

	// Second pass: process selector middleware configs
	for _, ms := range selectorConfigs {
		middlewareType := ms.GetType()
		middlewareName := ms.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			helper.Warnf("unknown server middleware type: %s", middlewareType)
			continue
		}

		helper.Debugf("enabling server selector middleware: %s (type: %s)", middlewareName, middlewareType)

		// Create selector middleware (can access previously created middlewares now)
		m, ok := f.NewMiddlewareServer(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
			// Add the created middleware into the carrier
			carrier.Servers[middlewareName] = m
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
