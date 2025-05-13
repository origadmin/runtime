/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"

	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
)

func LoggingServer(ms []KMiddleware, logger log.Logger) []KMiddleware {
	log.Debug("[Middleware] Logging server middleware enabled")
	return append(ms, logging.Server(logger))
}

func LoggingClient(ms []KMiddleware, logger log.Logger) []KMiddleware {
	log.Debug("[Middleware] Logging client middleware enabled")
	return append(ms, logging.Client(logger))
}

type loggingFactory struct {
}

func (l loggingFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	if middleware.GetLogging() {
		return logging.Client(options.Logger), true
	}
	return nil, false
}

func (l loggingFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	if middleware.GetLogging() {
		return logging.Server(options.Logger), true
	}
	return nil, false
}
