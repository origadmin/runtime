/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

func Empty() Middleware {
	return func(handler Handler) Handler {
		return handler
	}
}
