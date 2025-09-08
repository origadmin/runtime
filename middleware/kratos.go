/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

const Type = "middleware"

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/middleware
//go:adapter:package:type *
//go:adapter:package:type:prefix K
//go:adapter:package:func *
//go:adapter:package:func:regex New([A-Z])=NewK$1
//go:adapter:package:func *
//go:adapter:package:func:prefix K
