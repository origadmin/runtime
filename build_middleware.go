/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/middleware"
)

type (
	// MiddlewareBuilders middleware builders for runtime
	MiddlewareBuilders interface {
		// NewMiddlewaresClient build middleware
		NewMiddlewaresClient(*middlewarev1.Middleware, ...middleware.Option) []middleware.KMiddleware
		// NewMiddlewaresServer build middleware
		NewMiddlewaresServer(*middlewarev1.Middleware, ...middleware.Option) []middleware.KMiddleware
		// NewMiddlewareClient build middleware
		NewMiddlewareClient(string, *middlewarev1.Middleware, ...middleware.Option) (middleware.KMiddleware, error)
		// NewMiddlewareServer build middleware
		NewMiddlewareServer(string, *middlewarev1.Middleware, ...middleware.Option) (middleware.KMiddleware, error)
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

func (b *builder) NewMiddlewareClient(name string, config *middlewarev1.Middleware,
	ss ...middleware.Option) (middleware.KMiddleware, error) {
	factory, ok := b.MiddlewareBuilder.Get(name)
	if !ok {
		return nil, ErrNotFound
	}
	options := settings.ApplyZero(ss)
	m, ok := factory.NewMiddlewareClient(config, options)
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (b *builder) NewMiddlewareServer(name string, config *middlewarev1.Middleware, ss ...middleware.Option) (middleware.KMiddleware, error) {
	factory, ok := b.MiddlewareBuilder.Get(name)
	if !ok {
		return nil, ErrNotFound
	}
	options := settings.ApplyZero(ss)
	m, ok := factory.NewMiddlewareServer(config, options)
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (b *builder) NewMiddlewaresClient(config *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware {
	return b.MiddlewareBuilder.BuildClient(config, ss...)
}

func (b *builder) NewMiddlewaresServer(config *middlewarev1.Middleware, ss ...middleware.Option) []middleware.KMiddleware {
	return b.MiddlewareBuilder.BuildServer(config, ss...)
}

// RegisterMiddlewareBuilder registers a new MiddlewareBuilder with the given name.
func (b *builder) RegisterMiddlewareBuilder(name string, builder middleware.Factory) {
	b.MiddlewareBuilder.Register(name, builder)
}
