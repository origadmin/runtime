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

// --- Retrieval Helpers ---

// Get retrieves a component by name from a locator and asserts its type.
func Get[T any](ctx context.Context, l component.Locator, name string) (T, error) {
	var zero T
	inst, err := l.Get(ctx, name)
	if err != nil {
		return zero, err
	}
	if t, ok := inst.(T); ok {
		return t, nil
	}
	return zero, fmt.Errorf("engine: component '%s' in category '%s' is not of type %T", name, l.Category(), zero)
}

// GetWithTag retrieves a component by tag from a locator and asserts its type.
func GetWithTag[T any](ctx context.Context, l component.Locator, tag string) (T, error) {
	// Fluent API: directly use the new interface method
	return Get[T](ctx, l.WithInTags(tag), component.DefaultName)
}

// GetWithFallback retrieves a component by name, falling back to the default name if not found.
func GetWithFallback[T any](ctx context.Context, l component.Locator, name string) (T, error) {
	inst, err := Get[T](ctx, l, name)
	if err == nil {
		return inst, nil
	}
	return GetDefault[T](ctx, l)
}

// GetDefault retrieves the default component from a locator and asserts its type.
func GetDefault[T any](ctx context.Context, l component.Locator) (T, error) {
	return Get[T](ctx, l, component.DefaultName)
}

// --- Iteration Helpers ---

// Iter returns a type-safe iterator for components in a locator.
func Iter[T any](ctx context.Context, l component.Locator) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		for name, inst := range l.Iter(ctx) {
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
	for name, t := range Iter[T](ctx, l) {
		m[name] = t
	}
	return m, nil
}
