/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package empty implements the functions, types, and interfaces for the module.
package empty

import (
	"github.com/go-kratos/kratos/v2/middleware"
)

func Empty() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return handler
	}
}
