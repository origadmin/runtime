/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type recoveryFactory struct {
}

func (r recoveryFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debug("[Middleware] Recovery client middleware enabled")

	//recoveryConfig := cfg.GetRecovery()
	//if recoveryConfig == nil {
	//	return nil, false
	//}
	return recovery.Recovery(), true
}

func (r recoveryFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debug("[Middleware] Recovery server middleware enabled")

	//recoveryConfig := cfg.GetRecovery()
	//if recoveryConfig == nil {
	//	return nil, false
	//}
	return recovery.Recovery(), true
}
