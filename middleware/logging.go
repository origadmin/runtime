/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"

	"github.com/origadmin/runtime/middleware/selector"
)

func LoggingServer(selector selector.Selector, logger log.Logger) selector.Selector {
	log.Debug("[KMiddleware] Logging server middleware enabled")
	return selector.Append("Logging", logging.Server(logger))
}

func LoggingClient(selector selector.Selector, logger log.Logger) selector.Selector {
	log.Debug("[KMiddleware] Logging client middleware enabled")
	return selector.Append("Logging", logging.Client(logger))
}
