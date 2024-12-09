/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware/security"
)

type ContextType int

const (
	ContextTypeGrpc = iota
	ContextTypeMetaData
)

func SecurityClient(middlewares []Middleware, cfg *configv1.Security, option *config.MiddlewareOption) []Middleware {
	security.NewClient()
	return middlewares
}

func SecurityServer(middlewares []Middleware, cfg *configv1.Security, option *config.MiddlewareOption) []Middleware {

	return middlewares
}
