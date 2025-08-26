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

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/middleware
//go:adapter:package:type *
//go:adapter:package:type:prefix Kratos
//go:adapter:package:func *
//go:adapter:package:func:prefix Kratos
