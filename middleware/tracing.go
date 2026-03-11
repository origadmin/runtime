/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and contracts for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
)

type tracingFactory struct {
}

func (t tracingFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	mwOpts.GetLogger("middleware.tracing").Debug("enabling tracing client middleware")
	return tracing.Client(), true
}

func (t tracingFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	mwOpts.GetLogger("middleware.tracing").Debug("enabling tracing server middleware")
	return tracing.Server(), true
}
