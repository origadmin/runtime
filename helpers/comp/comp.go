/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package comp provides concise generic utility functions for component handling.
package comp

import (
	"context"
	"iter"

	"github.com/origadmin/runtime/contracts/component"
)

// AsConfig extracts and asserts the configuration from a handle.
func AsConfig[T any](h component.Handle) (*T, error) {
	cfg := h.Config()
	if cfg == nil {
		return nil, nil
	}
	if t, ok := cfg.(*T); ok {
		return t, nil
	}
	return nil, nil
}

// Get retrieves a component by name from a handle and asserts its type.
func Get[T any](ctx context.Context, h component.Handle, name string) (T, error) {
	var zero T
	inst, err := h.Get(ctx, name)
	if err != nil {
		return zero, err
	}
	if t, ok := inst.(T); ok {
		return t, nil
	}
	return zero, nil
}

// GetDefault retrieves the default component from a handle and asserts its type.
func GetDefault[T any](ctx context.Context, h component.Handle) (T, error) {
	return Get[T](ctx, h, component.DefaultName)
}

// Iter returns a type-safe iterator for components in a handle.
func Iter[T any](ctx context.Context, h component.Handle) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		for name, inst := range h.Iter(ctx) {
			if t, ok := inst.(T); ok {
				if !yield(name, t) {
					return
				}
			}
		}
	}
}

// GetMap collects all components from the given handle as a map and asserts their type.
func GetMap[T any](ctx context.Context, h component.Handle) (map[string]T, error) {
	m := make(map[string]T)
	for name, inst := range h.Iter(ctx) {
		if t, ok := inst.(T); ok {
			m[name] = t
		}
	}
	return m, nil
}
