/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

type circuitBreakerFactory struct {
}

func (c circuitBreakerFactory) NewMiddlewareClient(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	_, mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	if !cfg.GetEnabled() || cfg.GetType() != "circuit_breaker" {
		return nil, false
	}

	helper.Debug("[Middleware] CircuitBreaker client middleware enabled")
	return circuitbreaker.Client(), true
}

func (c circuitBreakerFactory) NewMiddlewareServer(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	_, mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debug("[Middleware] CircuitBreaker server middleware enabled, not supported yet")
	return nil, false
}
