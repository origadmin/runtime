/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
)

func LoggingServer(ms []Middleware, logger log.Logger) []Middleware {
	log.Debug("[Middleware] Logging server middleware enabled")
	return append(ms, logging.Server(logger))
}

func LoggingClient(ms []Middleware, logger log.Logger) []Middleware {
	log.Debug("[Middleware] Logging client middleware enabled")
	return append(ms, logging.Client(logger))
}
