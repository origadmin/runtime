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

// Get retrieves a component by name from a locator and asserts its type.
// Deprecated: Use comp.Get instead to avoid unnecessary dependencies.
func Get[T any](ctx context.Context, l component.Locator, name ...string) (T, error) {
	return comp.Get[T](ctx, l, name...)
}

// Iter returns a type-safe iterator for components in a locator.
func Iter[T any](ctx context.Context, l component.Locator) iter.Seq2[string, T] {
	return comp.Iter[T](ctx, l)
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
