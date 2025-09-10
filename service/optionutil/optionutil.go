// Package optionutil provides utility functions for working with options.
package optionutil

import (
	"github.com/origadmin/runtime/interfaces"
)

// OptionKey is a generic key type for option values
type OptionKey[T any] struct{}

// WithOption adds a value to the OptionSet
func WithOption[T any](opt interfaces.OptionValue, key OptionKey[T], value T) interfaces.OptionValue {
	if opt == nil {
		opt = interfaces.DefaultOptions()
	}
	return opt.WithOption(key, value)
}

// GetOption gets a value from the Option
func GetOption[T any](opt interfaces.OptionValue, key OptionKey[T]) (T, bool) {
	if opt == nil {
		var zero T
		return zero, false
	}

	val := opt.Option(key)
	if val == nil {
		var zero T
		return zero, false
	}

	if v, ok := val.(T); ok {
		return v, true
	}

	var zero T
	return zero, false
}

// WithSliceOption appends values to a slice in the OptionSet
func WithSliceOption[T any](opt interfaces.OptionValue, key OptionKey[[]T], values ...T) interfaces.OptionValue {
	if opt == nil {
		opt = interfaces.DefaultOptions()
	}

	// Get existing slice if any
	existing, _ := GetOption[[]T](opt, key)
	if existing == nil {
		existing = []T{}
	}

	// Create a new slice with appended values
	newSlice := make([]T, len(existing), len(existing)+len(values))
	copy(newSlice, existing)
	newSlice = append(newSlice, values...)

	return WithOption(opt, key, newSlice)
}

// GetSliceOption gets a slice from the Option
func GetSliceOption[T any](opt interfaces.OptionValue, key OptionKey[[]T]) []T {
	if opt == nil {
		return nil
	}

	val, ok := GetOption[[]T](opt, key)
	if !ok {
		return nil
	}

	// Return a copy to prevent external modifications
	result := make([]T, len(val))
	copy(result, val)
	return result
}
