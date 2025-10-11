// Package optionutil provides utility functions for working with options.Context.
package optionutil

import (
	"github.com/origadmin/runtime/interfaces/options"
)

// =============================
// Core Types & Implementation
// =============================

// Key is a generic struct used as a type-safe identifier for option keys.
// It ensures type safety when storing and retrieving values from options.Context.
type Key[T any] struct{}

// emptyContext implements the options.Context interface.
// It forms a chain of contexts for storing key-value pairs.
type emptyContext struct {
	parent options.Context
	key    any
	value  any
}

// Value retrieves a value from the context based on the given key.
// It first checks the current context, and if not found, recursively checks parent contexts.
func (o *emptyContext) Value(key any) any {
	if o.key != nil && o.key == key {
		return o.value
	}
	if o.parent != nil {
		return o.parent.Value(key)
	}
	return nil
}

// with creates a new context with the specified key-value pair.
func (o *emptyContext) with(key any, value any) *emptyContext {
	if key == nil {
		panic("key cannot be nil")
	}

	return &emptyContext{
		parent: o,
		key:    key,
		value:  value,
	}
}

// With implements the options.Context interface to create a context with a new key-value pair.
func (o *emptyContext) With(key any, value any) options.Context {
	return o.with(key, value)
}

// =============================
// Constructors
// =============================

// Empty creates an empty options.Context instance.
func Empty() options.Context {
	return &emptyContext{}
}

// Default creates a default options.Context instance.
func Default() options.Context {
	return &emptyContext{}
}

// =============================
// Explicit Key Functions
// =============================

// --- Setters (Explicit Key) ---

// WithValue sets a key-value pair in the specified options.Context.
// If the context is nil, a new empty context is created.
func WithValue[T any](ctx options.Context, key Key[T], value T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}
	return ctx.With(key, value)
}

// Append appends values to a slice in the options.Context.
func Append[T any](ctx options.Context, key Key[[]T], values ...T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}

	existing := SliceValue(ctx, key)
	newSlice := append(existing, values...)
	return WithValue(ctx, key, newSlice)
}

// --- Getters (Explicit Key) ---

// Value retrieves a value of the specified type from the options.Context.
func Value[T any](ctx options.Context, key Key[T]) (T, bool) {
	var zero T
	if ctx == nil {
		return zero, false
	}

	val := ctx.Value(key)
	if val == nil {
		return zero, false
	}

	if v, ok := val.(T); ok {
		return v, true
	}

	return zero, false
}

// SliceValue retrieves a copy of a slice from the options.Context.
func SliceValue[T any](ctx options.Context, key Key[[]T]) []T {
	if ctx == nil {
		return nil
	}

	val, ok := Value(ctx, key)
	if !ok {
		return nil
	}

	return val
}

// =============================
// Implicit Key Functions
// =============================

// --- Setters (Implicit Key) ---

// AppendValues appends values to a slice in the options.Context using an implicit key.
func AppendValues[T any](ctx options.Context, values ...T) options.Context {
	key := Key[[]T]{}
	return Append[T](ctx, key, values...)
}

// Update returns an options.Option that modifies a configuration struct *T in the options.Context chain.
func Update[T any](updaters ...func(*T)) options.Option {
	return func(ctx options.Context) options.Context {
		if v, ok := Value(ctx, Key[*T]{}); ok {
			for _, updater := range updaters {
				updater(v)
			}
		}
		return ctx
	}
}

// --- Getters (Implicit Key) ---

// Extract retrieves a value of type T from the given options.Context.
func Extract[T any](ctx options.Context) (T, bool) {
	if v, ok := Value(ctx, Key[T]{}); ok {
		return v, ok
	}
	var zero T
	return zero, false
}

// ExtractOr retrieves a value of type T from the given options.Context, or a default value if not found.
func ExtractOr[T any](ctx options.Context, defaultValue T) T {
	if v, ok := Extract[T](ctx); ok {
		return v
	}
	return defaultValue
}

// ExtractSlice retrieves a slice of type []T from the given options.Context.
func ExtractSlice[T any](ctx options.Context) []T {
	val, ok := Extract[[]T](ctx)
	if !ok {
		return nil
	}
	return val
}

// ExtractSliceOr retrieves a slice of type []T from the given options.Context, or a default value if not found.
func ExtractSliceOr[T any](ctx options.Context, defaultValue []T) []T {
	if v := ExtractSlice[T](ctx); v != nil {
		return v
	}
	return defaultValue
}

// =============================
// Context Application Core
// =============================

// apply is the internal core for applying options.
// It handles instance creation, context initialization, and option application.
func apply[T any](cfg *T, opts ...options.Option) (options.Context, *T) {
	// Centralized setup logic
	key := Key[*T]{}
	ctx := WithValue(Empty(), key, cfg)

	// Apply all options
	for _, option := range opts {
		ctx = option(ctx)
	}

	return ctx, cfg
}

// New creates a new instance of T, applies options, and returns both the final
// context and the configured instance.
func New[T any](opts ...options.Option) (options.Context, *T) {
	return apply[T](new(T), opts...)
}

// NewT creates a new instance of T, applies options, and returns the configured instance.
// This is a convenience wrapper around New for when the context is not needed.
func NewT[T any](opts ...options.Option) *T {
	_, cfg := apply[T](new(T), opts...)
	return cfg
}

// NewContext creates a new instance of T, applies options, and returns the resulting context.
// It's useful for creating a configured context template.
func NewContext[T any](opts ...options.Option) options.Context {
	ctx, _ := apply[T](new(T), opts...)
	return ctx
}

// Apply applies a series of options.Option to a given configuration object.
func Apply[T any](cfg *T, opts ...options.Option) options.Context {
	ctx, _ := apply[T](cfg, opts...)
	return ctx
}

// =============================
// Context Utilities
// =============================

// WithContext returns an options.Option that sets the given options.Context as the current context.
func WithContext(ctx options.Context) options.Option {
	return func(o options.Context) options.Context {
		return ctx
	}
}

// WithCond returns an options.Option that applies the given options.Option if the condition is true.
func WithCond(condition bool, opt options.Option) options.Option {
	if condition {
		return opt
	}
	return func(ctx options.Context) options.Context {
		return ctx
	}
}

// WithGroup returns an options.Option that applies a group of options.Option to the context.
func WithGroup(opts ...options.Option) options.Option {
	return func(ctx options.Context) options.Context {
		if ctx == nil {
			ctx = Empty()
		}
		for _, option := range opts {
			ctx = option(ctx)
		}
		return ctx
	}
}
