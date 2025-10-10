/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

type tracingFactory struct {
}

func (t tracingFactory) NewMiddlewareClient(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KratosMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	if !cfg.GetEnabled() || cfg.GetType() != "tracing" {
		return nil, false
	}

	helper.Debug("[Middleware] Tracing client middleware enabled")
	return tracing.Client(), true
}

func (t tracingFactory) NewMiddlewareServer(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KratosMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	if !cfg.GetEnabled() || cfg.GetType() != "tracing" {
		return nil, false
	}

	helper.Debug("[Middleware] Tracing server middleware enabled")
	return tracing.Server(), true
}
