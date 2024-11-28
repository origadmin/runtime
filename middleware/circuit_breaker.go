/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
)

func CircuitBreakerClient(ms []Middleware, ok bool) []Middleware {
	if !ok {
		return ms
	}
	return append(ms, circuitbreaker.Client())
}
