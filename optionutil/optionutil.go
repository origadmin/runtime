// Package optionutil provides utility functions for working with options.
package optionutil

import (
	"github.com/origadmin/runtime/interfaces"
)

type Key[T any]  struct{}

// options implement the optionvalue interface
type options struct {
	parent *options
	key    any
	value  any
}

func (o *options) replaceInChain(key interface{}) *options {
	if o == nil {
		return nil
	}

	if o.key == key {
		return o.parent
	}

	return &options{
		parent: o.parent.replaceInChain(key),
		key:    o.key,
		value:  o.value,
	}
}

func (o *options) Value(key any) any {
	if o == nil {
		return nil
	}
	if o.key != nil && o.key == key {
		return o.value
	}
	if o.parent != nil {
		return o.parent.Value(key)
	}
	return nil
}

func (o *options) With(key any, value any) interfaces.Option {
	if key == nil {
		panic("key cannot be nil")
	}

	return &options{
		parent: o,
		key:    key,
		value:  value,
	}
}

func Options() interfaces.Option {
	return &options{}
}

// With adds a value to the options
func With[T any](opt interfaces.Option, key Key[T], value T) interfaces.Option {
	if opt == nil {
		opt = Options()
	}
	return opt.With(key, value)
}

// Value gets a value from the Option
func Value[T any](opt interfaces.Option, key Key[T]) (T, bool) {
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

// Append appends values to a slice in the options
func Append[T any](opt interfaces.Option, key Key[[]T], values ...T) interfaces.Option {
	if opt == nil {
		opt = Options()
	}

	// Get existing slice if any
	existing, _ := Value(opt, key)
	if existing == nil {
		existing = []T{}
	}

	// Create a new slice with appended values
	newSlice := make([]T, len(existing), len(existing)+len(values))
	copy(newSlice, existing)
	newSlice = append(newSlice, values...)

	return With(opt, key, newSlice)
}

// Slice gets a slice from the Option
func Slice[T any](opt interfaces.Option, key Key[[]T]) []T {
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
