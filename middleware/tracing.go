/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/origadmin/runtime/log"
)

func TracingClient(ms []Middleware) []Middleware {
	log.Debug("[Middleware] Tracing client middleware enabled")
	return append(ms, tracing.Client())
}

func TracingServer(ms []Middleware) []Middleware {
	log.Debug("[Middleware] Tracing server middleware enabled")
	return append(ms, tracing.Server())
}
