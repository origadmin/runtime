/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware/selector"
)

func Recovery(selector selector.Selector) selector.Selector {
	log.Infof("[KMiddleware] Recovery middleware enabled")
	return selector.Append("Recovery", recovery.Recovery())
}
