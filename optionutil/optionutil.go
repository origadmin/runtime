// Package optionutil provides utility functions for working with emptyContext.
package optionutil

import (
	"github.com/origadmin/runtime/interfaces/options"
)

// Key is a generic struct used as a type-safe identifier for option keys
type Key[T any] struct{}

// emptyContext implement the optionvalue interface
type emptyContext struct {
	parent options.Context
	key    any
	value  any
}

// Value retrieves a value from the ctx based on the given key
// It first checks the current ctx, and if not found, recursively checks parent ctxs
// Parameters:
//   - key: The key to look up
// Returns:
//   - any: The value associated with the key, or nil if not found
func (o *emptyContext) Value(key any) any {
	if o.key != nil && o.key == key {
		return o.value
	}
	if o.parent != nil {
		return o.parent.Value(key)
	}
	return nil
}

// with creates a new ctx with the specified key-value pair
// Parameters:
//   - key: The key, cannot be nil
//   - value: The value
// Returns:
//   - *emptyContext: The newly created ctx
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

// With implements the Option interface's With method to create a ctx with a new key-value pair
// Parameters:
//   - key: The key
//   - value: The value
// Returns:
//   - options.Context: The Option instance containing the new key-value pair
func (o *emptyContext) With(key any, value any) options.Context {
	return o.with(key, value)
}

// With sets a key-value pair in the specified Option
// If the passed opt is nil, a new empty ctx is created
// Parameters:
//   - ctx: The original Option instance
//   - key: The generic key
//   - value: The value to set
// Returns:
//   - options.Context: The Option instance with the new value set
func With[T any](ctx options.Context, key Key[T], value T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}
	return ctx.With(key, value)
}

// Value gets a value from the Empty
// Get a value of the specified type from the Option
// Parameters:
//   - ctx: The Option instance
//   - key: The generic key
// Returns:
//   - T: The retrieved value
//   - bool: Whether the value was successfully retrieved
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

// ValueOr returns the value associated with the given key in the options.Context, or a default value if the key is not found.
// Parameters:
//   - ctx: The options.Context instance to retrieve the value from
//   - key: The key to look up in the context
//   - defaultValue: The default value to return if the key is not found
// Returns:
//   - T: The retrieved value or the default value if the key is not found
func ValueOr[T any](ctx options.Context, key Key[T], defaultValue T) T {
	if v, ok := Value(ctx, key); ok {
		return v
	}
	return defaultValue
}

// Append appends values to a slice in the emptyContext
// Append values to a slice in the Option
// Parameters:
//   - ctx: The Option instance
//   - key: The slice-type key
//   - values: The values to append
// Returns:
//   - options.Context: The updated Option instance
func Append[T any](ctx options.Context, key Key[[]T], values ...T) options.Context {
	if ctx == nil {
		ctx = &emptyContext{}
	}

	// Get existing slice
	existing := Slice(ctx, key)
	// Append new values to the existing slice
	newSlice := append(existing, values...)
	// Store the combined slice
	return With(ctx, key, newSlice)
}

// Slice gets a slice from the Empty
// Get a copy of a slice from the Option
// Parameters:
//   - ctx: The Option instance
//   - key: The slice-type key
// Returns:
//   - []T: A copy of the slice, or nil if it doesn't exist
func Slice[T any](ctx options.Context, key Key[[]T]) []T {
	if ctx == nil {
		return nil
	}

	val, ok := Value(ctx, key)
	if !ok {
		return nil
	}

	// Return a copy to prevent external modifications
	result := make([]T, len(val))
	copy(result, val)
	return result
}

// Empty creates an empty Option instance
// Returns:
//   - options.Context: An empty Option instance
func Empty() options.Context {
	return &emptyContext{}
}

// Default creates a default Option instance
// Returns:
//   - options.Context: A default Option instance
func Default() options.Context {
	return &emptyContext{}
}

// Update returns an options.ContextFunc that modifies a configuration struct *T in the options.Context chain.
// It finds the *T instance using the provided key and applies the updater function if it exists.
// If the value for the key does not exist, it does nothing.
// Create an OptionFunc that updates a configuration struct in the Option chain
// Parameters:
//   - updater: The update function that takes a T parameter
// Returns:
//   - options.Option: A function that can be applied to an Option
func Update[T any](updater func(T)) options.Option {
	return func(ctx options.Context) options.Context {
		// Use a zero-value Key[T] as the type-specific key to find the value.
		key := Key[T]{}
		if v, ok := Value(ctx, key); ok {
			// If the value is found, apply the updater function.
			// Since T is expected to be a pointer type, the update modifies the original object.
			updater(v)
		}
		return ctx
	}
}

// Apply applies a series of OptionFuncs to a given configuration object and returns the resulting options.Context.
// It iterates through the provided option functions and applies each one to the object.
// Parameters:
//   - cfg: The configuration object
//   - opts: The list of OptionFuncs to apply
// Returns:
//   - options.Context: The options.Context instance with the applied options.
func Apply[T any](cfg T, opts ...options.Option) options.Context {
	// Start with an option chain that contains the configuration object,
	// keyed by its type via Key[T].
	ctx := With(Empty(), Key[T]{}, cfg)
	for _, option := range opts {
		ctx = option(ctx)
	}
	return ctx
}

// ApplyNew creates a new instance of T, applies options to it, and returns a pointer to the configured instance.
// The generic type T should be the struct type itself, not a pointer to it (e.g., use serverConfig, not *serverConfig).
// Parameters:
//   - opts: The list of OptionFuncs to apply to the new instance
// Returns:
//   - options.Context: The options.Context instance with the applied options.
func ApplyNew[T any](opts ...options.Option) (options.Context, *T) {
	// Create a zero-value instance of the struct T.
	cfg := new(T)

	// Store this valid pointer in the context, keyed by its type *T.
	ctx := With(Empty(), Key[*T]{}, cfg)

	// Apply all opts. The Update function will now receive a valid pointer.
	for _, option := range opts {
		ctx = option(ctx)
	}

	// Return the pointer to the configured struct.
	return ctx, cfg
}

// ApplyContext applies a series of OptionFuncs to a given options.Context and returns the resulting options.Context.
// It iterates through the provided option functions and applies each one to the context.
// Parameters:
//   - ctx: The options.Context instance to apply the options to
//   - opts: The list of OptionFuncs to apply
// Returns:
//   - options.Context: The options.Context instance with the applied options.
func ApplyContext(ctx options.Context, opts ...options.Option) options.Context {
	// Start with the provided context.
	if ctx == nil {
		ctx = Empty()
	}

	for _, option := range opts {
		ctx = option(ctx)
	}
	return ctx
}

// WithContext returns an OptionFunc that sets the given options.Context as the current context.
// Parameters:
//   - ctx: The options.Context instance to set as the current context
// Returns:
//   - options.Option: An OptionFunc that sets the context
func WithContext(ctx options.Context) options.Option {
	return func(o options.Context) options.Context {
		return ctx
	}
}

// FromContext retrieves a value of type T from the given options.Context.
// Parameters:
//   - ctx: The options.Context instance to retrieve the value from
// Returns:
//   - T: The retrieved value
//   - bool: Whether the value was successfully retrieved
func FromContext[T any](ctx options.Context) (T, bool) {
	if v, ok := Value(ctx, Key[T]{}); ok {
		return v, ok
	}
	var zero T
	return zero, false
}

// If returns an OptionFunc that applies the given OptionFunc if the condition is true.
// Otherwise, it returns an OptionFunc that does nothing.
// Parameters:
//   - condition: The condition to check
//   - opt: The OptionFunc to apply if the condition is true
// Returns:
//   - options.Option: The OptionFunc that applies the option if the condition is true
//   - options.Option: The OptionFunc that does nothing if the condition is false
func If(condition bool, opt options.Option) options.Option {
	if condition {
		return opt
	}
	return func(ctx options.Context) options.Context {
		return ctx
	}
}

// Group returns an OptionFunc that applies a group of OptionFuncs to the context.
// Parameters:
//   - opts: The list of OptionFuncs to apply
// Returns:
//   - options.Option: The OptionFunc that applies the group of options
func Group(opts ...options.Option) options.Option {
	return func(ctx options.Context) options.Context {
		return ApplyContext(ctx, opts...)
	}
}
