/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package engine

import (
	"context"
	"fmt"
	"iter"

	"github.com/origadmin/runtime/contracts/component"
)

// Get retrieves a component by name from the handle and asserts its type.
func Get[T any](ctx context.Context, h Handle, name string) (T, error) {
	var zero T
	inst, err := h.Get(ctx, name)
	if err != nil {
		return zero, err
	}
	typed, ok := inst.(T)
	if !ok {
		return zero, fmt.Errorf("engine: component %s/%s is not of type %T", h.Category(), name, zero)
	}
	return typed, nil
}

// GetDefault retrieves the active/default instance for the current handle.
func GetDefault[T any](ctx context.Context, h Handle) (T, error) {
	return Get[T](ctx, h, "")
}

// GetOr retrieves a component by name, or falls back to the default instance if not found.
func GetOr[T any](ctx context.Context, h Handle, name string) (T, error) {
	if name != "" && name != component.DefaultName {
		inst, err := Get[T](ctx, h, name)
		if err == nil {
			return inst, nil
		}
	}
	return Get[T](ctx, h, "")
}

// ToMap collects all typed instances from the current handle into a map.
func ToMap[T any](ctx context.Context, h Handle) (map[string]T, error) {
	res := make(map[string]T)
	for name, inst := range h.Iter(ctx) {
		if typed, ok := inst.(T); ok {
			res[name] = typed
		}
	}
	return res, nil
}

// Iter returns an iterator that yields all instances of the desired type from the current handle.
func Iter[T any](ctx context.Context, h Handle) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		for name, inst := range h.Iter(ctx) {
			if typed, ok := inst.(T); ok {
				if !yield(name, typed) {
					return
				}
			}
		}
	}
}

// AsConfig asserts the configuration associated with the handle to *T.
// This is the recommended way to consume configuration in providers.
func AsConfig[T any](h Handle) (*T, error) {
	val := h.Config()
	if val == nil {
		return nil, fmt.Errorf("engine: configuration is nil for %s/%s", h.Category(), h.Scope())
	}
	// Priority 1: Direct pointer match
	if typed, ok := val.(*T); ok {
		return typed, nil
	}
	// Priority 2: Value match (wrap in pointer)
	if typed, ok := val.(T); ok {
		return &typed, nil
	}
	return nil, fmt.Errorf("engine: cannot cast config %T to %T", val, (*T)(nil))
}

// BindConfig is a generic helper to populate a target pointer using type assertion on Config().
// Discouraged: Use AsConfig[T](h) instead for more idiomatic Go code.
func BindConfig[T any](h Handle, target *T) error {
	cfg, err := AsConfig[T](h)
	if err != nil {
		return err
	}
	*target = *cfg
	return nil
}
