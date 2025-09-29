// Package optionutil provides utility functions for working with emptyContext.
package optionutil

import (
	"context" // Import context package

	"github.com/origadmin/runtime/interfaces"
)

// Key is a generic struct used as a type-safe identifier for option keys
type Key[T any] struct{}

// emptyContext implement the optionvalue interface
type emptyContext struct {
	parent interfaces.Context
	key    any
	value  any
}

// contextKey is an unexported type for context keys.
type contextKey int

const (
	optionContextKey contextKey = iota
)

// ToContext stores the interfaces.Context into a standard library context.Context.
func ToContext(ctx context.Context, opt interfaces.Context) context.Context {
	return context.WithValue(ctx, optionContextKey, opt)
}

// FromContext retrieves the interfaces.Context from a standard library context.Context.
// It returns nil if no interfaces.Context is found.
func FromContext(ctx context.Context) interfaces.Context {
	if ctx == nil {
		return nil
	}
	if opt, ok := ctx.Value(optionContextKey).(interfaces.Context); ok {
		return opt
	}
	return nil
}

// Value retrieves a value from the context based on the given key
// It first checks the current context, and if not found, recursively checks parent contexts
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

// with creates a new context with the specified key-value pair
// Parameters:
//   - key: The key, cannot be nil
//   - value: The value
// Returns:
//   - *emptyContext: The newly created context
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

// With implements the Option interface's With method to create a context with a new key-value pair
// Parameters:
//   - key: The key
//   - value: The value
// Returns:
//   - interfaces.Context: The Option instance containing the new key-value pair
func (o *emptyContext) With(key any, value any) interfaces.Context {
	return o.with(key, value)
}

// With sets a key-value pair in the specified Option
// If the passed opt is nil, a new empty context is created
// Parameters:
//   - opt: The original Option instance
//   - key: The generic key
//   - value: The value to set
// Returns:
//   - interfaces.Context: The Option instance with the new value set
func With[T any](opt interfaces.Context, key Key[T], value T) interfaces.Context {
	if opt == nil {
		opt = &emptyContext{}
	}
	return opt.With(key, value)
}

// Value gets a value from the Empty
// Get a value of the specified type from the Option
// Parameters:
//   - opt: The Option instance
//   - key: The generic key
// Returns:
//   - T: The retrieved value
//   - bool: Whether the value was successfully retrieved
func Value[T any](opt interfaces.Context, key Key[T]) (T, bool) {
	var zero T
	if opt == nil {
		return zero, false
	}

	val := opt.Value(key)
	if val == nil {
		return zero, false
	}

	if v, ok := val.(T); ok {
		return v, true
	}

	return zero, false
}

// Append appends values to a slice in the emptyContext
// Append values to a slice in the Option
// Parameters:
//   - opt: The Option instance
//   - key: The slice-type key
//   - values: The values to append
// Returns:
//   - interfaces.Context: The updated Option instance
func Append[T any](opt interfaces.Context, key Key[[]T], values ...T) interfaces.Context {
	if opt == nil {
		opt = &emptyContext{}
	}

	// Get existing slice
	existing := Slice(opt, key)
	// Append new values to the existing slice
	newSlice := append(existing, values...)
	// Store the combined slice
	return With(opt, key, newSlice)
}

// Slice gets a slice from the Empty
// Get a copy of a slice from the Option
// Parameters:
//   - opt: The Option instance
//   - key: The slice-type key
// Returns:
//   - []T: A copy of the slice, or nil if it doesn't exist
func Slice[T any](opt interfaces.Context, key Key[[]T]) []T {
	if opt == nil {
		return nil
	}

	val, ok := Value(opt, key)
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
//   - interfaces.Context: An empty Option instance
func Empty() interfaces.Context {
	return &emptyContext{}
}

// Default creates a default Option instance
// Returns:
//   - interfaces.Context: A default Option instance
func Default() interfaces.Context {
	return &emptyContext{}
}

// Update returns an interfaces.ContextFunc that modifies a configuration struct *T in the interfaces.Context chain.
// It finds the *T instance using the provided key and applies the updater function if it exists.
// If the value for the key does not exist, it does nothing.
// Create an OptionFunc that updates a configuration struct in the Option chain
// Parameters:
//   - updater: The update function that takes a T parameter
// Returns:
//   - interfaces.ContextFunc: A function that can be applied to an Option
func Update[T any](updater func(T)) interfaces.Option {
	return func(opt interfaces.Context) {
		// Use a zero-value Key[T] as the type-specific key to find the value.
		key := Key[T]{}
		if v, ok := Value(opt, key); ok {
			// If the value is found, apply the updater function.
			// Since T is expected to be a pointer type, the update modifies the original object.
			updater(v)
		}
	}
}

// Apply applies a series of OptionFuncs to a given configuration object and returns the resulting interfaces.Context.
// It iterates through the provided option functions and applies each one to the object.
// Parameters:
//   - cfg: The configuration object
//   - opts: The list of OptionFuncs to apply
// Returns:
//   - interfaces.Context: The interfaces.Context instance with the applied options.
func Apply[T any](cfg T, opts ...interfaces.Option) interfaces.Context {
	// Start with an option chain that contains the configuration object,
	// keyed by its type via Key[T].
	options := With(Empty(), Key[T]{}, cfg)
	for _, o := range opts {
		o(options)
	}
	return options
}

// ApplyNew applies a series of OptionFuncs to a new instance of a configuration object and returns a pointer to it.
// It iterates through the provided option functions and applies each one to the object.
// Parameters:
//   - opts: The list of OptionFuncs to apply
// Returns:
//   - *T: A pointer to the newly created configuration object with the applied options.
// Note: This function is useful when you want to create a new configuration object with specific options applied.
func ApplyNew[T any](opts ...interfaces.Option) T {
	var cfg T
	// Start with an option chain that contains the configuration object,
	// keyed by its type via Key[T].
	options := With(Empty(), Key[T]{}, cfg)
	for _, o := range opts {
		o(options)
	}
	return cfg
}
