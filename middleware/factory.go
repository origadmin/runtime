/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and contracts for the module.
package middleware

import (
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/contracts/builder"
	internalfactory "github.com/origadmin/runtime/helpers/builderutil"
	"github.com/origadmin/runtime/log"
)

// defaultBuilder is the default instance of the Builder.
var defaultBuilder = NewBuilder()

func init() {
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
	builder.Registry[Factory]
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

	// Create a middleware carrier for context propagation
	clients := make(map[string]KMiddleware)
	logger := log.NewHelper(opt.Logger)

	// First pass: separate regular middlewares and selector configs
	for _, ms := range cfg.GetConfigs() {
		if !ms.GetEnabled() {
			continue
		}
		middlewareType := ms.GetType()
		if middlewareType == "" {
			continue
		}
		middlewareName := ms.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		if middlewareType == string(Selector) {
			selectorConfigs = append(selectorConfigs, ms)
			continue
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			logger.Warnf("unknown client middleware type: %s", middlewareType)
			continue
		}

		// Create middleware
		m, ok := f.NewMiddlewareClient(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
			clients[middlewareName] = m
		}
	}

	// Attach carrier and process selectors
	opts = append(opts, WithClientCarrier(clients))
	for _, ms := range selectorConfigs {
		middlewareType := ms.GetType()
		middlewareName := ms.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			continue
		}

		m, ok := f.NewMiddlewareClient(ms, opts...)
		if ok {
			middlewares = append(middlewares, m)
		}
	}

	return middlewares
}

// BuildServerMiddlewares builds the server middleware chain
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

	servers := make(map[string]KMiddleware)
	logger := log.NewHelper(opt.Logger)

	// First pass
	for _, mwCfg := range cfg.GetConfigs() {
		if !mwCfg.GetEnabled() {
			continue
		}

		middlewareType := mwCfg.GetType()
		if middlewareType == "" {
			continue
		}
		middlewareName := mwCfg.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		if middlewareType == string(Selector) {
			selectorConfigs = append(selectorConfigs, mwCfg)
			continue
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			logger.Warnf("unknown server middleware type: %s", middlewareType)
			continue
		}

		ms, ok := f.NewMiddlewareServer(mwCfg, opts...)
		if ok {
			middlewares = append(middlewares, ms)
			servers[middlewareName] = ms
		}
	}

	// Attach carrier and process selectors
	opts = append(opts, WithServerCarrier(servers))
	for _, mwCfg := range selectorConfigs {
		middlewareType := mwCfg.GetType()
		middlewareName := mwCfg.GetName()
		if middlewareName == "" {
			middlewareName = middlewareType
		}

		f, ok := b.Get(middlewareType)
		if !ok {
			continue
		}

		ms, ok := f.NewMiddlewareServer(mwCfg, opts...)
		if ok {
			middlewares = append(middlewares, ms)
		}
	}

	return middlewares
}

// NewClient creates a single client-side middleware instance.
func NewClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	if cfg == nil || !cfg.GetEnabled() {
		return nil, false
	}
	middlewareType := cfg.GetType()
	f, ok := defaultBuilder.Get(middlewareType)
	if !ok {
		return nil, false
	}
	return f.NewMiddlewareClient(cfg, opts...)
}

// NewServer creates a single server-side middleware instance.
func NewServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	if cfg == nil || !cfg.GetEnabled() {
		return nil, false
	}
	middlewareType := cfg.GetType()
	f, ok := defaultBuilder.Get(middlewareType)
	if !ok {
		return nil, false
	}
	return f.NewMiddlewareServer(cfg, opts...)
}

func BuildClients(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildClientMiddlewares(cfg, opts...)
}

func BuildServers(cfg *middlewarev1.Middlewares, opts ...Option) []KMiddleware {
	return defaultBuilder.BuildServerMiddlewares(cfg, opts...)
}

func RegisterFactory(name Name, factory Factory) {
	defaultBuilder.Register(string(name), factory)
}

func NewBuilder() *Builder {
	return &Builder{
		Registry: internalfactory.New[Factory](),
	}
}
