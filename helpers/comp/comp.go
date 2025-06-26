/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package comp provides concise generic utility functions for component handling.
package comp

import (
	"context"
	"fmt"
	"iter"

	"github.com/origadmin/runtime/contracts/component"
)

// --- Extraction Helpers ---

// RequireTyped retrieves a requirement by purpose from a handle and asserts its type.
func RequireTyped[T any](h component.Handle, purpose string) (T, error) {
	var zero T
	res, err := h.Require(purpose)
	if err != nil {
		return zero, err
	}
	if res == nil {
		return zero, fmt.Errorf("engine: component '%s/%s' required '%s' but got nil", h.Category(), h.Name(), purpose)
	}
	if t, ok := res.(T); ok {
		return t, nil
	}
	return zero, fmt.Errorf("engine: component '%s/%s' required '%s' with type %T, but got %T", h.Category(), h.Name(), purpose, (*T)(nil), res)
}

// AsConfig extracts and asserts the configuration from a handle.
func AsConfig[T any](h component.Handle) (*T, error) {
	cfg := h.Config()
	if cfg == nil {
		return nil, nil
	}
	if t, ok := cfg.(*T); ok {
		return t, nil
	}
	return nil, fmt.Errorf("engine: component '%s/%s' expected config type %T, but got %T", h.Category(), h.Name(), (*T)(nil), cfg)
}

// --- Retrieval Helpers ---
// Get retrieves a component by name from a locator and asserts its type.
// If no name is provided, the default component is retrieved.
func Get[T any](ctx context.Context, l component.Locator, name ...string) (T, error) {
	inst, err := l.Get(ctx, name...)

	if err != nil {
		return *new(T), err
	}
	if inst == nil {
		reqName := ""
		if len(name) > 0 {
			reqName = name[0]
		}
		return *new(T), fmt.Errorf("engine: locator for '%s' found component '%s' but instance is nil", l.Category(), reqName)
	}
	if t, ok := inst.(T); ok {
		return t, nil
	}
	reqName := ""
	if len(name) > 0 {
		reqName = name[0]
	}
	return *new(T), fmt.Errorf("engine: locator for '%s' found component '%s' with type %T, but expected %T", l.Category(), reqName, inst, (*T)(nil))
}

// GetWithTag retrieves a component by tag from a locator and asserts its type.
func GetWithTag[T any](ctx context.Context, l component.Locator, tag string) (T, error) {
	// Fluent API: directly use the new interface method
	return Get[T](ctx, l.WithInTags(tag))
}

// GetWithFallback retrieves a component by name, falling back to the default name if not found.
func GetWithFallback[T any](ctx context.Context, l component.Locator, name string) (T, error) {
	inst, err := Get[T](ctx, l, name)
	if err == nil {
		return inst, nil
	}
	return Get[T](ctx, l)
}

// --- Iteration Helpers ---

// Iter returns a type-safe iterator for components in a locator.
func Iter[T any](ctx context.Context, l component.Locator) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		var it = l.Iter(ctx)
		for it.Next() {
			name, inst := it.Value()
			if t, ok := inst.(T); ok {
				if !yield(name, t) {
					return
				}
			}
		}
	}
}

// GetMap collects all components from the given locator as a map and asserts their type.
func GetMap[T any](ctx context.Context, l component.Locator) (map[string]T, error) {
	m := make(map[string]T)
	it := l.Iter(ctx)
	for it.Next() {
		name, inst := it.Value()
		if t, ok := inst.(T); ok {
			m[name] = t
		}
	}
	if err := it.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
