// Package optionutil provides utility functions for working with emptyContext.
package optionutil

import (
	"github.com/origadmin/runtime/interfaces/options"
)

// =============================
// 1) Core Types & Implementation
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
// 2) Constructors
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
// 3) Get/Set & Collection Utilities
// =============================

// WithValue sets a key-value pair in the specified options.Context.
// If the context is nil, a new empty context is created.
func WithValue[T any](ctx options.Context, key Key[T], value T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}
	return ctx.With(key, value)
}

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

// ValueOr returns the value associated with the given key in the options.Context,
// or a default value if the key is not found.
func ValueOr[T any](ctx options.Context, key Key[T], defaultValue T) T {
	if v, ok := Value(ctx, key); ok {
		return v
	}
	return defaultValue
}

// WithUpdate returns an options.Option that modifies a configuration struct *T in the options.Context chain.
func WithUpdate[T any](updater func(T)) options.Option {
	return func(ctx options.Context) options.Context {
		key := Key[T]{}
		if v, ok := Value(ctx, key); ok {
			updater(v)
		}
		return ctx
	}
}

// WithAppend appends values to a slice in the options.Context.
func WithAppend[T any](ctx options.Context, key Key[[]T], values ...T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}

	existing := Slice(ctx, key)
	newSlice := append(existing, values...)
	return WithValue(ctx, key, newSlice)
}

// Slice retrieves a copy of a slice from the options.Context.
func Slice[T any](ctx options.Context, key Key[[]T]) []T {
	if ctx == nil {
		return nil
	}

	val, ok := Value(ctx, key)
	if !ok {
		return nil
	}

	result := make([]T, len(val))
	copy(result, val)
	return result
}

// =============================
// 4) Apply/Update & Context Utilities
// =============================

// New creates a new instance of T, applies options to it, and returns a pointer to the configured instance.
func New[T any](opts ...options.Option) (options.Context, *T) {
	cfg := new(T)
	ctx := WithValue(Empty(), Key[*T]{}, cfg)
	for _, option := range opts {
		ctx = option(ctx)
	}
	return ctx, cfg
}

// Apply applies a series of options.Option to a given configuration object.
func Apply[T any](cfg T, opts ...options.Option) options.Context {
	ctx := WithValue(Empty(), Key[T]{}, cfg)
	for _, option := range opts {
		ctx = option(ctx)
	}
	return ctx
}

// ApplyContext applies a series of options.Option to a given options.Context.
func ApplyContext(ctx options.Context, opts ...options.Option) options.Context {
	if ctx == nil {
		ctx = Empty()
	}
	for _, option := range opts {
		ctx = option(ctx)
	}
	return ctx
}

// WithContext returns an options.Option that sets the given options.Context as the current context.
func WithContext(ctx options.Context) options.Option {
	return func(o options.Context) options.Context {
		return ctx
	}
}

// Extract retrieves a value of type T from the given options.Context.
func Extract[T any](ctx options.Context) (T, bool) {
	if v, ok := Value(ctx, Key[T]{}); ok {
		return v, ok
	}
	var zero T
	return zero, false
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
		return ApplyContext(ctx, opts...)
	}
}
