/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/origadmin/runtime/log"
)

func TracingClient(selector Selector) Selector {
	log.Debug("[KMiddleware] Tracing client middleware enabled")
	return selector.Append("Metadata", tracing.Client())
}

func TracingServer(selector Selector) Selector {
	log.Debug("[KMiddleware] Tracing server middleware enabled")
	return selector.Append("Metadata", tracing.Server())
}
