/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware/selector"
)

func TracingClient(selector selector.Selector) selector.Selector {
	log.Debug("[KMiddleware] Tracing client middleware enabled")
	return selector.Append("Metadata", tracing.Client())
}

func TracingServer(selector selector.Selector) selector.Selector {
	log.Debug("[KMiddleware] Tracing server middleware enabled")
	return selector.Append("Metadata", tracing.Server())
}
