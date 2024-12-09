/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package helper implements the functions, types, and interfaces for the module.
package helper

import (
	"github.com/go-kratos/kratos/v2/metadata"

	"github.com/origadmin/runtime/context"
)

func FromMD(ctx context.Context, key string) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(key)
	}
	return ""
}
