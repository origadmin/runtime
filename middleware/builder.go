/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/builder"
)

// DefaultServiceBuilder is the default instance of the buildImpl.
var DefaultServiceBuilder Factory = &factory{}

// ServiceBuilder is a struct that implements the buildImpl interface.
// It provides methods for creating new gRPC and HTTP servers and clients.
type factory struct{}

func (f factory) NewMiddlewaresClient(middlewares []KMiddleware, customize *configv1.Customize, option ...Option) []KMiddleware {
	//TODO implement me
	panic("implement me")
}

func (f factory) NewMiddlewaresServer(middlewares []KMiddleware, customize *configv1.Customize, option ...Option) []KMiddleware {
	//TODO implement me
	panic("implement me")
}

func (f factory) NewMiddlewareClient(s string, config *configv1.Customize_Config, option ...Option) (KMiddleware, error) {
	//TODO implement me
	panic("implement me")
}

func (f factory) NewMiddlewareServer(s string, config *configv1.Customize_Config, option ...Option) (KMiddleware, error) {
	//TODO implement me
	panic("implement me")
}

type buildImpl struct {
	builder.Builder[Factory]
}

func (b *buildImpl) NewMiddlewaresServer(middlewares []KMiddleware, customize *configv1.Customize, option ...Option) []KMiddleware {
	//TODO implement me
	panic("implement me")
}

func (b *buildImpl) NewMiddlewareServer(s string, config *configv1.Customize_Config, option ...Option) (KMiddleware, error) {
	//TODO implement me
	panic("implement me")
}

func (b *buildImpl) RegisterMiddlewareBuilder(name string, factory Factory) {

}

func (b *buildImpl) NewMiddlewareClient(name string, config *configv1.Customize_Config, option ...Option) (KMiddleware, error) {
	//TODO implement me
	panic("implement me")
}

func (b *buildImpl) NewMiddlewaresClient(middlewares []KMiddleware, customize *configv1.Customize, option ...Option) []KMiddleware {
	//TODO implement me
	panic("implement me")
}

func NewBuilder() Builder {
	return &buildImpl{
		Builder: builder.New[Factory](),
	}
}
