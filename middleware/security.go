/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

func SecurityClient(middlewares []Middleware, cfg *configv1.Security) []Middleware {
	return middlewares
}

func SecurityServer(middlewares []Middleware, cfg *configv1.Security) []Middleware {
	return middlewares
}
