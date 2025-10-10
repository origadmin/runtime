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
	Register(RateLimit, &rateLimitFactory{})
	Register(Metadata, &metadataFactory{})
	Register(Selector, &selectorFactory{})
	Register(Tracing, &tracingFactory{})
	Register(Validator, &validatorFactory{})
}

// middlewareBuilder is a builder for creating middleware chains.
type middlewareBuilder struct {
	factory.Registry[Factory]
}

// BuildClient 构建客户端中间件链
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

	// 创建中间件载体用于上下文传递
	carrier := &Carrier{
		Clients: make(map[string]KMiddleware),
		Servers: make(map[string]KMiddleware),
	}
	// 将载体添加到上下文中
	opt.Context = WithMiddlewaresToContext(opt.Context, carrier)

	helper := log.NewHelper(logger)
	helper.Info("building client middlewares")

	// 第一次遍历：分离普通中间件和selector配置
	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		middlewareName := ms.GetType()
		// 分离selector中间件配置，后面单独处理
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

		// 创建中间件
		m, ok := f.NewMiddlewareClient(ms, options.WithContext(opt.Context), withOptions(opt))
		if ok {
			middlewares = append(middlewares, m)
			// 将创建的中间件添加到载体中
			carrier.Clients[middlewareName] = m
		}
	}

	// 第二次遍历：处理selector中间件配置
	for _, ms := range selectorConfigs {
		middlewareName := ms.GetType()
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown client middleware: %s", middlewareName)
			continue
		}

		helper.Infof("enabling client middleware: %s", middlewareName)

		// 创建selector中间件（此时可以访问已创建的中间件）
		m, ok := f.NewMiddlewareClient(ms, options.WithContext(opt.Context), withOptions(opt))
		if ok {
			middlewares = append(middlewares, m)
		}
	}

	return middlewares
}

// BuildServer 构建服务端中间件链（类似BuildClient的修改）
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

	// 创建中间件载体用于上下文传递
	carrier := &Carrier{
		Clients: make(map[string]KMiddleware),
		Servers: make(map[string]KMiddleware),
	}
	// 将载体添加到上下文中
	opt.Context = WithMiddlewaresToContext(opt.Context, carrier)

	helper := log.NewHelper(logger)
	helper.Info("building server middlewares")

	// 第一次遍历：分离普通中间件和selector配置
	for _, ms := range cfg.GetMiddlewares() {
		if !ms.GetEnabled() {
			continue
		}
		middlewareName := ms.GetType()
		// 分离selector中间件配置，后面单独处理
		if middlewareName == string(Selector) {
			selectorConfigs = append(selectorConfigs, ms)
			continue
		}
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown server middleware: %s", middlewareName)
			continue
		}

		helper.Infof("enabling server middleware: %s", middlewareName)

		// 创建中间件
		m, ok := f.NewMiddlewareServer(ms, options.WithContext(opt.Context), withOptions(opt))
		if ok {
			middlewares = append(middlewares, m)
			// 将创建的中间件添加到载体中
			carrier.Servers[middlewareName] = m
		}
	}

	// 第二次遍历：处理selector中间件配置
	for _, ms := range selectorConfigs {
		middlewareName := ms.GetType()
		f, ok := b.Get(middlewareName)
		if !ok {
			helper.Warnf("unknown server middleware: %s", middlewareName)
			continue
		}

		helper.Infof("enabling server middleware: %s", middlewareName)

		// 创建selector中间件（此时可以访问已创建的中间件）
		m, ok := f.NewMiddlewareServer(ms, options.WithContext(opt.Context), withOptions(opt))
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