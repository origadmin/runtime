/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type circuitBreakerFactory struct {
}

func (c circuitBreakerFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debug("[Middleware] CircuitBreaker client middleware enabled")
	//if !cfg.GetEnabled() || cfg.GetType() != "circuit_breaker" {
	//	return nil, false
	//}

	return circuitbreaker.Client(), true
}

func (c circuitBreakerFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debug("[Middleware] CircuitBreaker server middleware enabled, not supported yet")
	return nil, false
}
