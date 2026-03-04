package engine

import (
	"context"
	"fmt"
	"iter"
)

// Cast gets a component from the current handle context and type-asserts it.
// Focuses strictly on conversion.
func Cast[T any](ctx context.Context, h Handle, name string) (T, error) {
	var zero T
	inst, err := h.Get(ctx, name)
	if err != nil {
		return zero, err
	}
	typed, ok := inst.(T)
	if !ok {
		return zero, fmt.Errorf("engine: component %s/%s is not of type %T", h.Category(), name, zero)
	}
	return typed, nil
}

// GetDefault retrieves the winner instance for the current handle context.
func GetDefault[T any](ctx context.Context, h Handle) (T, error) {
	return Cast[T](ctx, h, "")
}

// Shortcuts for common categories - focus strictly on current context.
func GetRegistry[T any](ctx context.Context, h Handle, name string) (T, error) {
	return Cast[T](ctx, h, name)
}

func GetMiddleware[T any](ctx context.Context, h Handle, name string) (T, error) {
	return Cast[T](ctx, h, name)
}

func GetDatabase[T any](ctx context.Context, h Handle, name string) (T, error) {
	return Cast[T](ctx, h, name)
}

// Get retrieves a component by name from the current handle context and type-asserts it.
func Get[T any](ctx context.Context, h Handle, name string) (T, error) {
	return Cast[T](ctx, h, name)
}

// All retrieves all instances from the current handle and converts them to the desired type.
func All[T any](ctx context.Context, h Handle) (map[string]T, error) {
	res := make(map[string]T)
	for name, inst := range h.Iter(ctx) {
		if typed, ok := inst.(T); ok {
			res[name] = typed
		}
	}
	return res, nil
}

// Iter returns an iterator that yields all instances of the desired type from the current handle.
func Iter[T any](ctx context.Context, h Handle) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		for name, inst := range h.Iter(ctx) {
			if typed, ok := inst.(T); ok {
				if !yield(name, typed) {
					return
				}
			}
		}
	}
}

// BindConfig is a generic helper to populate a target pointer using type assertion.
func BindConfig[T any](h Handle, target *T) error {
	return h.BindConfig(target)
}
