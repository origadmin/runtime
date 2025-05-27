/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type circuitBreakerFactory struct {
}

func (c circuitBreakerFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] CircuitBreaker client middleware enabled")
	if checkEnabled(middleware, "circuit_breaker") {
		return circuitbreaker.Client(), true
	}
	return nil, false
}

func (c circuitBreakerFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] CircuitBreaker server middleware enabled, not supported yet")
	return nil, false
}

func CircuitBreakerClient(ms []KMiddleware) []KMiddleware {
	log.Debug("[Middleware] CircuitBreaker client middleware enabled")
	return append(ms, circuitbreaker.Client())
}
