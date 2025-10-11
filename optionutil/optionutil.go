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
// Data Access Functions
// =============================

// value retrieves a value with an explicit key.
func value[T any](ctx options.Context, key Key[T]) (T, bool) {
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

// withValue sets a key-value pair with an explicit key.
func withValue[T any](ctx options.Context, key Key[T], value T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}
	return ctx.With(key, value)
}

// slice retrieves a slice with an explicit key.
func slice[T any](ctx options.Context, key Key[[]T]) []T {
	if ctx == nil {
		return nil
	}

	val, ok := value(ctx, key)
	if !ok {
		return nil
	}

	return val
}

// appendValues appends values to a slice with an explicit key.
func appendValues[T any](ctx options.Context, key Key[[]T], values ...T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}

	existing := slice(ctx, key)
	newSlice := append(existing, values...)
	return withValue(ctx, key, newSlice)
}

// --- Implicit Key --- //

// Value retrieves a value using an implicit key.
func Value[T any](ctx options.Context) (T, bool) {
	return value(ctx, Key[T]{})
}

// ValueOr retrieves a value, or a default if not found.
func ValueOr[T any](ctx options.Context, defaultValue T) T {
	if v, ok := Value[T](ctx); ok {
		return v
	}
	return defaultValue
}

// ValueCond retrieves a value if the condition is true, or a default otherwise.
func ValueCond[T any](ctx options.Context, condition func(T) bool, defaultValue T) T {
	if v, ok := Value[T](ctx); ok && condition(v) {
		return v
	}
	return defaultValue
}

// WithValue sets a value using an implicit key.
func WithValue[T any](ctx options.Context, value T) options.Context {
	return withValue(ctx, Key[T]{}, value)
}

// Slice retrieves a slice using an implicit key.
func Slice[T any](ctx options.Context) []T {
	return slice[T](ctx, Key[[]T]{})
}

// SliceOr retrieves a slice, or a default if not found.
func SliceOr[T any](ctx options.Context, defaultValue []T) []T {
	if v := Slice[T](ctx); v != nil {
		return v
	}
	return defaultValue
}

func SliceCond[T any](ctx options.Context, condition func([]T) bool) []T {
	if v := Slice[T](ctx); v != nil && condition(v) {
		return v
	}
	return nil
}

// Append appends values to a slice using an implicit key.
func Append[T any](ctx options.Context, values ...T) options.Context {
	return appendValues[T](ctx, Key[[]T]{}, values...)
}

// Update returns an options.Option that modifies a configuration struct *T using an implicit key.
func Update[T any](updaters ...func(*T)) options.Option {
	return func(ctx options.Context) options.Context {
		if v, ok := value(ctx, Key[*T]{}); ok {
			for _, updater := range updaters {
				updater(v)
			}
		}
		return ctx
	}
}

// =============================
// Context Application Core
// =============================

// apply is the internal core for applying options.
func apply[T any](cfg *T, opts ...options.Option) (options.Context, *T) {
	key := Key[*T]{}
	ctx := withValue(Empty(), key, cfg)

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
func NewT[T any](opts ...options.Option) *T {
	_, cfg := apply[T](new(T), opts...)
	return cfg
}

// NewContext creates a new instance of T, applies options, and returns the resulting context.
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
