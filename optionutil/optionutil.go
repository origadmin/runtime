// Package optionutil provides utility functions for working with Context.
package optionutil

import (
	"github.com/origadmin/runtime/interfaces"
)

type Key[T any]  struct{}

// Context implement the optionvalue interface
type Context struct {
	parent interfaces.Option
	key    any
	value  any
}

func (o *Context) Value(key any) any {
	if o.key != nil && o.key == key {
		return o.value
	}
	if o.parent != nil {
		return o.parent.Value(key)
	}
	return nil
}

func (o *Context) with(key any, value any) *Context {
	if key == nil {
		panic("key cannot be nil")
	}

	return &Context{
		parent: o,
		key:    key,
		value:  value,
	}
}

func (o *Context) With(key any, value any) interfaces.Option {
	return o.with(key, value)
}

// Options 是一个通用的选项包装器
type Options[T any] struct {
	option  interfaces.Option
	wrapped *T
}

func (o *Options[T]) Wrap(value *T) *Options[T] {
	o.wrapped = value
	return o
}
func (o *Options[T]) Unwrap() *T {
	return o.wrapped
}

func (o *Options[T]) Update(updater func(*T) *T) *Options[T] {
	updated := updater(o.wrapped)
	return &Options[T]{
		option:  o.option,
		wrapped: updated,
	}
}

func (o *Options[T]) Value(key any) any {
	if o == nil || o.option == nil {
		return nil
	}
	return o.option.Value(key)
}

func (o *Options[T]) With(key any, value any) interfaces.Option {
	if o.option == nil {
		o.option = &Context{}
	}
	o.option = o.option.With(key, value)
	return o
}
func With[T any](opt interfaces.Option, key Key[T], value T) interfaces.Option {
	if opt == nil {
		opt = &Context{}
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

// Append appends values to a slice in the Context
func Append[T any](opt interfaces.Option, key Key[[]T], values ...T) interfaces.Option {
	if opt == nil {
		opt = &Context{}
	}

	// Get existing slice
	existing := Slice(opt, key)
	// Append new values to the existing slice
	newSlice := append(existing, values...)
	// Store the combined slice
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

func Option() interfaces.Option {
	return &Context{}
}
