/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type tracingFactory struct {
}

func (t tracingFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] Tracing client middleware enabled")
	if checkEnabled(middleware, "tracing") {
		return tracing.Client(), true
	}
	return nil, false
}

func (t tracingFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] Tracing server middleware enabled")
	if checkEnabled(middleware, "tracing") {
		return tracing.Server(), true
	}
	return nil, false
}

func TracingClient(ms []KMiddleware) []KMiddleware {
	log.Debug("[Middleware] Tracing client middleware enabled")
	return append(ms, tracing.Client())
}

func TracingServer(ms []KMiddleware) []KMiddleware {
	log.Debug("[Middleware] Tracing server middleware enabled")
	return append(ms, tracing.Server())
}
