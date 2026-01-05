/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	internalfactory "github.com/origadmin/runtime/internal/factory"
	"github.com/origadmin/runtime/log"
)

// defaultBuilder is the default instance of the Builder.
var defaultBuilder = NewBuilder()

func init() {
	// The factories will be registered here once they are updated to the new interface.
	// optimizeFactory is removed from here as it's not a formal feature and should be registered by the user.
	// All other factories will be uncommented as they are updated.
	RegisterFactory(Recovery, &recoveryFactory{})
	RegisterFactory(Jwt, &jwtFactory{})
	RegisterFactory(CircuitBreaker, &circuitBreakerFactory{})
	RegisterFactory(Logging, &loggingFactory{})
	RegisterFactory(RateLimiter, &rateLimitFactory{})
	RegisterFactory(Metadata, &metadataFactory{})
	RegisterFactory(Selector, &selectorFactory{})
	RegisterFactory(Tracing, &tracingFactory{})
	RegisterFactory(Validator, &validatorFactory{})
}

// Builder is a builder for creating middleware chains.
type Builder struct {
	factory.Registry[Factory]
}

// BuildClientMiddlewares builds the client middleware chain
func (b *Builder) BuildClientMiddlewares(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
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

	// Create a middleware carrier for context propagation
	clients := make(map[string]KMiddleware)

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
			clients[middlewareName] = m
		}
	}
	// Attach the carrier into options context
	opts = append(opts, WithClientCarrier(clients))
	helper.Debugf("carrier: %+v", clients)
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
			//clients[middlewareName] = m
		}
	}

	return middlewares
}

// BuildServerMiddlewares builds the server middleware chain (similar to BuildClient)
func (b *Builder) BuildServerMiddlewares(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
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

	helper := log.NewHelper(logger)
	helper.Debug("building server middlewares")

	// Create a middleware carrier for context propagation
	servers := make(map[string]KMiddleware)

	// First pass: process non-selector middlewares
	for _, mwCfg := range cfg.GetConfigs() {
		if !mwCfg.GetEnabled() {
			continue
		}

		// Use Name if provided, otherwise fall back to Type
		middlewareType := mwCfg.GetType()
		middlewareName := mwCfg.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		// Defer selector middleware configs for later processing
		if middlewareType == string(Selector) {
			selectorConfigs = append(selectorConfigs, mwCfg)
			continue
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			helper.Warnf("unknown server middleware type: %s", middlewareType)
			continue
		}

		helper.Debugf("enabling server middleware: %s (type: %s)", middlewareName, middlewareType)

		// Create middleware
		ms, ok := f.NewMiddlewareServer(mwCfg, opts...)
		if ok {
			middlewares = append(middlewares, ms)
			// Add the created middleware into the carrier
			servers[middlewareName] = ms
		}
	}

	// Attach the carrier into options context for selector middlewares
	opts = append(opts, WithServerCarrier(servers))
	helper.Debugf("carrier: %+v", servers)
	// Second pass: process selector middlewares

	for _, mwCfg := range selectorConfigs {
		middlewareType := mwCfg.GetType()
		middlewareName := mwCfg.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			helper.Warnf("unknown server middleware type: %s", middlewareType)
			continue
		}

		helper.Debugf("enabling server selector middleware: %s (type: %s)", middlewareName, middlewareType)

		// Create selector middleware (can access previously created middlewares)
		ms, ok := f.NewMiddlewareServer(mwCfg, opts...)
		if ok {
			middlewares = append(middlewares, ms)
			// Add the created middleware into the carrier
			//servers[middlewareName] = ms
		}
	}

	return middlewares
}

// NewClient creates a single client-side middleware instance.
// This is the public API for creating individual client-side middlewares.
func NewClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	middlewareType := cfg.GetType()
	f, ok := defaultBuilder.Get(middlewareType)
	if !ok {
		log.Warnf("unknown client middleware type: %s", middlewareType)
		return nil, false
	}
	return f.NewMiddlewareClient(cfg, opts...)
}

// NewServer creates a single server-side middleware instance.
// This is the public API for creating individual server-side middlewares.
func NewServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	middlewareType := cfg.GetType()
	f, ok := defaultBuilder.Get(middlewareType)
	if !ok {
		log.Warnf("unknown server middleware type: %s", middlewareType)
		return nil, false
	}
	return f.NewMiddlewareServer(cfg, opts...)
}

// BuildClients creates a new client middleware chain using the default builder.
// This is the public API for building client-side middlewares.
func BuildClients(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildClientMiddlewares(cfg, opts...)
}

// BuildServers creates a new server middleware chain using the default builder.
// This is the public API for building server-side middlewares.
func BuildServers(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildServerMiddlewares(cfg, opts...)
}

// RegisterFactory registers a middleware factory with the given name.
func RegisterFactory(name Name, factory Factory) {
	defaultBuilder.Register(string(name), factory)
}

// NewBuilder creates a new middleware builder.
func NewBuilder() *Builder {
	return &Builder{
		Registry: internalfactory.New[Factory](),
	}
}
