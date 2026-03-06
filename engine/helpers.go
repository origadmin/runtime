/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package engine

import (
	"context"
	"iter"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
)

// Get retrieves a component by name from a handle and asserts its type.
// Deprecated: Use comp.Get instead to avoid unnecessary dependencies.
func Get[T any](ctx context.Context, h component.Handle, name string) (T, error) {
	return comp.Get[T](ctx, h, name)
}

// GetDefault retrieves the default component from a handle and asserts its type.
// Deprecated: Use comp.GetDefault instead.
func GetDefault[T any](ctx context.Context, h component.Handle) (T, error) {
	return comp.GetDefault[T](ctx, h)
}

// Iter returns a type-safe iterator for components in a handle.
func Iter[T any](ctx context.Context, h component.Handle) iter.Seq2[string, T] {
	return comp.Iter[T](ctx, h)
}

// AsConfig extracts and asserts the configuration from a handle.
func AsConfig[T any](h component.Handle) (*T, error) {
	return comp.AsConfig[T](h)
}

// BindConfig is a convenience function that applies the configuration from a handle to a target pointer.
// Discouraged: Use comp.AsConfig[T](h) instead for more idiomatic Go code.
func BindConfig[T any](h component.Handle, target *T) error {
	cfg, err := comp.AsConfig[T](h)
	if err != nil {
		return err
	}
	*target = *cfg
	return nil
}
