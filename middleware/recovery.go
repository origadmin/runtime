/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type recoveryFactory struct {
}

func (r recoveryFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] Recovery client middleware enabled")
	if checkEnabled(middleware, "recovery") {
		return recovery.Recovery(), true
	}
	return nil, false
}

func (r recoveryFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] Recovery server middleware enabled")
	if checkEnabled(middleware, "recovery") {
		return recovery.Recovery(), true
	}
	return nil, false
}

func Recovery(ms []KMiddleware) []KMiddleware {
	log.Infof("[Middleware] Recovery middleware enabled")
	return append(ms, recovery.Recovery())
}
