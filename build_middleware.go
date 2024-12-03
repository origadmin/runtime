/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
)

type (
	// registryBuildRegistry is an interface that defines a method for registering a RegistryBuilder.
	middlewareBuildRegistry interface {
		RegisterMiddlewareBuilder(string, MiddlewareBuilder)
	}

	// MiddlewareBuilders middleware builders for runtime
	MiddlewareBuilders interface {
		// NewMiddlewaresClient build middleware
		NewMiddlewaresClient([]middleware.Middleware, *configv1.Customize, *config.RuntimeConfig) []middleware.Middleware
		// NewMiddlewaresServer build middleware
		NewMiddlewaresServer([]middleware.Middleware, *configv1.Customize, *config.RuntimeConfig) []middleware.Middleware
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(string, *configv1.Customize_Config, *config.RuntimeConfig) (middleware.Middleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(string, *configv1.Customize_Config, *config.RuntimeConfig) (middleware.Middleware, error)
	}

	// MiddlewareBuilder middleware builder interface
	MiddlewareBuilder interface {
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(*configv1.Customize_Config, *config.RuntimeConfig) (middleware.Middleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(*configv1.Customize_Config, *config.RuntimeConfig) (middleware.Middleware, error)
	}

	// MiddlewareBuildFunc is an interface that defines methods for creating middleware.
	MiddlewareBuildFunc = func(*configv1.Customize_Config, *config.RuntimeConfig) (middleware.Middleware, error)
)

func (b *builder) NewMiddlewareClient(name string, config *configv1.Customize_Config, runtimeConfig *config.RuntimeConfig) (middleware.Middleware, error) {
	b.middlewareMux.RLock()
	defer b.middlewareMux.RUnlock()
	if builder, ok := b.middlewares[name]; ok {
		return builder.NewMiddlewareClient(config, runtimeConfig)
	}
	return nil, ErrNotFound
}

func (b *builder) NewMiddlewareServer(name string, config *configv1.Customize_Config, runtimeConfig *config.RuntimeConfig) (middleware.Middleware, error) {
	b.middlewareMux.RLock()
	defer b.middlewareMux.RUnlock()
	if builder, ok := b.middlewares[name]; ok {
		return builder.NewMiddlewareServer(config, runtimeConfig)
	}
	return nil, ErrNotFound
}

func (b *builder) NewMiddlewaresClient(mms []middleware.Middleware, cc *configv1.Customize, rc *config.RuntimeConfig) []middleware.Middleware {
	configs := config.GetTypeConfigs(cc, middleware.Type)
	var mbs []*middlewareBuilderWrap
	b.middlewareMux.RLock()
	for name := range configs {
		if mb, ok := b.middlewares[name]; ok {
			mbs = append(mbs, &middlewareBuilderWrap{
				Name:    name,
				Config:  configs[name],
				Builder: mb,
			})
		}
	}
	b.middlewareMux.RUnlock()
	for _, mb := range mbs {
		if m, err := mb.NewClient(rc); err == nil {
			mms = append(mms, m)
		}
	}
	return mms
}

func (b *builder) NewMiddlewaresServer(mms []middleware.Middleware, cc *configv1.Customize, rc *config.RuntimeConfig) []middleware.Middleware {
	configs := config.GetTypeConfigs(cc, middleware.Type)
	var mbs []*middlewareBuilderWrap
	b.middlewareMux.RLock()
	for name := range configs {
		if mb, ok := b.middlewares[name]; ok {
			mbs = append(mbs, &middlewareBuilderWrap{
				Name:    name,
				Config:  configs[name],
				Builder: mb,
			})
		}
	}
	b.middlewareMux.RUnlock()
	for _, mb := range mbs {
		if m, err := mb.NewServer(rc); err == nil {
			mms = append(mms, m)
		}
	}
	return mms
}

// RegisterMiddlewareBuilder registers a new MiddlewareBuilder with the given name.
func (b *builder) RegisterMiddlewareBuilder(name string, builder MiddlewareBuilder) {
	b.middlewareMux.Lock()
	defer b.middlewareMux.Unlock()
	b.middlewares[name] = builder
}

type middlewareBuilderWrap struct {
	Name    string
	Config  *configv1.Customize_Config
	Builder MiddlewareBuilder
}

func (m middlewareBuilderWrap) NewClient(runtimeConfig *config.RuntimeConfig) (middleware.Middleware, error) {
	return m.Builder.NewMiddlewareClient(m.Config, runtimeConfig)
}

func (m middlewareBuilderWrap) NewServer(runtimeConfig *config.RuntimeConfig) (middleware.Middleware, error) {
	return m.Builder.NewMiddlewareServer(m.Config, runtimeConfig)
}
