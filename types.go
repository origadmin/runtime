/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/origadmin/runtime/engine"
)

func WithScope(s engine.Scope) engine.RegisterOption {
	return engine.WithScope(s)
}

func WithPriority(p int) engine.RegisterOption {
	return engine.WithPriority(p)
}
