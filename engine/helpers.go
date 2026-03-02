package engine

import (
	"context"
	"fmt"
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

// BindConfig is a generic helper to populate a target pointer using type assertion.
func BindConfig[T any](h Handle, target *T) error {
	return h.BindConfig(target)
}
