/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/origadmin/runtime/log"
)

func TracingClient(f Filter) Filter {
	log.Debug("[KMiddleware] Tracing client middleware enabled")
	return f.Filter("Metadata", tracing.Client())
}

func TracingServer(f Filter) Filter {
	log.Debug("[KMiddleware] Tracing server middleware enabled")
	return f.Filter("Metadata", tracing.Server())
}
