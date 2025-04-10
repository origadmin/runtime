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
		RegisterMiddlewareBuilder(string, MiddlewareBuilder)
	}

	// MiddlewareBuilders middleware builders for runtime
	MiddlewareBuilders interface {
		// NewMiddlewaresClient build middleware
		NewMiddlewaresClient([]middleware.KMiddleware, *configv1.Customize, ...middleware.Option) []middleware.KMiddleware
		// NewMiddlewaresServer build middleware
		NewMiddlewaresServer([]middleware.KMiddleware, *configv1.Customize, ...middleware.Option) []middleware.KMiddleware
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(string, *configv1.Customize_Config, ...middleware.Option) (middleware.KMiddleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(string, *configv1.Customize_Config, ...middleware.Option) (middleware.KMiddleware, error)
	}

	// MiddlewareBuilder middleware builder interface
	MiddlewareBuilder interface {
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(*configv1.Customize_Config, ...middleware.Option) (middleware.KMiddleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(*configv1.Customize_Config, ...middleware.Option) (middleware.KMiddleware, error)
	}

	// MiddlewareBuildFunc is an interface that defines methods for creating middleware.
	MiddlewareBuildFunc = func(*configv1.Customize_Config, ...middleware.Option) (middleware.KMiddleware, error)
)

func (b *builder) NewMiddlewareClient(name string, config *configv1.Customize_Config, ss ...middleware.Option) (middleware.KMiddleware, error) {
	b.middlewareMux.RLock()
	defer b.middlewareMux.RUnlock()
	if builder, ok := b.middlewares[name]; ok {
		return builder.NewMiddlewareClient(config, ss...)
	}
	return nil, ErrNotFound
}

func (b *builder) NewMiddlewareServer(name string, config *configv1.Customize_Config, ss ...middleware.Option) (middleware.KMiddleware, error) {
	b.middlewareMux.RLock()
	defer b.middlewareMux.RUnlock()
	if builder, ok := b.middlewares[name]; ok {
		return builder.NewMiddlewareServer(config, ss...)
	}
	return nil, ErrNotFound
}

func (b *builder) NewMiddlewaresClient(mms []middleware.KMiddleware, cc *configv1.Customize, ss ...middleware.Option) []middleware.KMiddleware {
	configs := customize.ConfigsFromType(cc, middleware.Type)
	var mbs []*middlewareWrap
	b.middlewareMux.RLock()
	for name := range configs {
		if mb, ok := b.middlewares[name]; ok {
			mbs = append(mbs, &middlewareWrap{
				Name:    name,
				Config:  configs[name],
				Builder: mb,
			})
		}
	}
	b.middlewareMux.RUnlock()
	for _, mb := range mbs {
		if m, err := mb.NewClient(ss...); err == nil {
			mms = append(mms, m)
		}
	}
	return mms
}

func (b *builder) NewMiddlewaresServer(mms []middleware.KMiddleware, cc *configv1.Customize, ss ...middleware.Option) []middleware.KMiddleware {
	configs := customize.ConfigsFromType(cc, middleware.Type)
	var mbs []*middlewareWrap
	b.middlewareMux.RLock()
	for name := range configs {
		if mb, ok := b.middlewares[name]; ok {
			mbs = append(mbs, &middlewareWrap{
				Name:    name,
				Config:  configs[name],
				Builder: mb,
			})
		}
	}
	b.middlewareMux.RUnlock()
	for _, mb := range mbs {
		if m, err := mb.NewServer(ss...); err == nil {
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

type middlewareWrap struct {
	Name    string
	Config  *configv1.Customize_Config
	Builder MiddlewareBuilder
}

func (m middlewareWrap) NewClient(ss ...middleware.Option) (middleware.KMiddleware, error) {
	return m.Builder.NewMiddlewareClient(m.Config, ss...)
}

func (m middlewareWrap) NewServer(ss ...middleware.Option) (middleware.KMiddleware, error) {
	return m.Builder.NewMiddlewareServer(m.Config, ss...)
}
