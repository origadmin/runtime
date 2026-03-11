/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and contracts for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
)

type circuitBreakerFactory struct {
}

func (c circuitBreakerFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	mwOpts := FromOptions(opts...)
	mwOpts.GetLogger("middleware.circuit_breaker").Debug("enabling circuit_breaker client middleware")
	return circuitbreaker.Client(), true
}

func (c circuitBreakerFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	mwOpts := FromOptions(opts...)
	mwOpts.GetLogger("middleware.circuit_breaker").Debug("enabling circuit_breaker server middleware, not supported yet")
	return nil, false
}
