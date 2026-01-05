/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type recoveryFactory struct {
}

func (r recoveryFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(log.With(mwOpts.Logger, "module", "middleware.recovery"))
	helper.Debug("enabling recovery client middleware")

	// The default Kratos recovery middleware is sufficient and includes logging.
	return recovery.Recovery(), true
}

func (r recoveryFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(log.With(mwOpts.Logger, "module", "middleware.recovery"))
	helper.Debug("enabling recovery server middleware")

	// The default Kratos recovery middleware is sufficient and includes logging.
	return recovery.Recovery(), true
}
