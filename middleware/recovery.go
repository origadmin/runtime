/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"github.com/origadmin/runtime/log"
)

func Recovery(ms []Middleware, ok bool) []Middleware {
	if ok {
		log.Infof("[Middleware] Recovery: %v", ok)
		ms = append(ms, recovery.Recovery())
	}
	return ms
}
