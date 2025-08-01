/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"
)

const Type = "middleware"

type (
	KHandler    = middleware.Handler
	KMiddleware = middleware.Middleware
)

// Chain returns a middleware that executes a chain of middleware.
func Chain(m ...KMiddleware) KMiddleware {
	return middleware.Chain(m...)
}
