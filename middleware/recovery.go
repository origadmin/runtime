/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"github.com/origadmin/runtime/log"
)

func Recovery(f Filter) Filter {
	log.Infof("[Middleware] Recovery middleware enabled")
	return f.Filter("Recovery", recovery.Recovery())
}
