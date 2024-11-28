/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/origadmin/runtime/customize"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware"
)

type (
	// registryBuildRegistry is an interface that defines a method for registering a RegistryBuilder.
	middlewareBuildRegistry interface {
		RegisterMiddlewareBuilder(name string, registryBuilder MiddlewareBuilder)
	}

	// MiddlewareBuilders middleware builders for runtime
	MiddlewareBuilders interface {
		// NewMiddlewaresClient build middleware
		NewMiddlewaresClient([]middleware.Middleware, *configv1.Customize) []middleware.Middleware
		// NewMiddlewaresServer build middleware
		NewMiddlewaresServer([]middleware.Middleware, *configv1.Customize) []middleware.Middleware
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(name string, config *configv1.Customize_Config) (middleware.Middleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(name string, config *configv1.Customize_Config) (middleware.Middleware, error)
	}

	// MiddlewareBuilder middleware builder interface
	MiddlewareBuilder interface {
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(config *configv1.Customize_Config) (middleware.Middleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(config *configv1.Customize_Config) (middleware.Middleware, error)
	}

	// MiddlewareBuildFunc is an interface that defines methods for creating middleware.
	MiddlewareBuildFunc = func(*configv1.Customize_Config) (middleware.Middleware, error)
)

func (b *builder) NewMiddlewareClient(name string, config *configv1.Customize_Config) (middleware.Middleware, error) {
	b.middlewareMux.RLock()
	defer b.middlewareMux.RUnlock()
	if builder, ok := b.middlewares[name]; ok {
		return builder.NewMiddlewareClient(config)
	}
	return nil, ErrNotFound
}

func (b *builder) NewMiddlewareServer(name string, config *configv1.Customize_Config) (middleware.Middleware, error) {
	b.middlewareMux.RLock()
	defer b.middlewareMux.RUnlock()
	if builder, ok := b.middlewares[name]; ok {
		return builder.NewMiddlewareServer(config)
	}
	return nil, ErrNotFound
}

func (b *builder) NewMiddlewaresClient(mms []middleware.Middleware, cc *configv1.Customize) []middleware.Middleware {
	configs := customize.GetTypeConfigs(cc, middleware.Type)
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
		if m, err := mb.NewClient(); err == nil {
			mms = append(mms, m)
		}
	}
	return mms
}

func (b *builder) NewMiddlewaresServer(mms []middleware.Middleware, cc *configv1.Customize) []middleware.Middleware {
	configs := customize.GetTypeConfigs(cc, middleware.Type)
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
		if m, err := mb.NewServer(); err == nil {
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

func (m middlewareBuilderWrap) NewClient() (middleware.Middleware, error) {
	return m.Builder.NewMiddlewareClient(m.Config)
}

func (m middlewareBuilderWrap) NewServer() (middleware.Middleware, error) {
	return m.Builder.NewMiddlewareServer(m.Config)
}
