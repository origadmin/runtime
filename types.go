/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/metadata"
)

type RegisterOption = component.RegisterOption
type InOption = component.InOption

// WithScope specifies the exact perspective during In().
func WithScope(s metadata.Scope) InOption {
	return engine.WithScope(s)
}

// WithScopes specifies multiple visibilities during registration.
func WithScopes(ss ...metadata.Scope) RegisterOption {
	return engine.WithScopes(ss...)
}

// WithPriority is a functional option to specify the initialization priority.
func WithPriority(p int) RegisterOption {
	return engine.WithPriority(p)
}
