/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
)

type Resolver interface {
	Resolve(config KConfig) (Resolved, error)
}

type Resolved interface {
	WithDecode(name string, v any, decode func([]byte, any) error) error
	Value(name string) (any, error)
	Middleware() *middlewarev1.Middleware
	Service() *configv1.Service
	Logger() *configv1.Logger
}

type ResolveObserver interface {
	Observer(string, KValue)
}

type ResolveFunc func(config KConfig) (Resolved, error)

func (r ResolveFunc) Resolve(config KConfig) (Resolved, error) {
	return r(config)
}
