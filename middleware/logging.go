/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
)

func LoggingServer(f Filter, logger log.Logger) Filter {
	log.Debug("[KMiddleware] Logging server middleware enabled")
	return f.Filter("Logging", logging.Server(logger))
}

func LoggingClient(f Filter, logger log.Logger) Filter {
	log.Debug("[KMiddleware] Logging client middleware enabled")
	return f.Filter("Logging", logging.Client(logger))
}
