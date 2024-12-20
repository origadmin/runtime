/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"

	"github.com/origadmin/runtime/log"
)

func CircuitBreakerClient(f Filter) Filter {
	log.Debug("[KMiddleware] CircuitBreaker client middleware enabled")
	return f.Filter("CircuitBreaker", circuitbreaker.Client())
}
