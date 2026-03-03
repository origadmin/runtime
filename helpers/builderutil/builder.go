package builderutil

import (
	"context"
	"fmt"
	"sync"
)

type (
	Registry[T any] interface {
		Register(name string, f T)
		Get(name string) (T, bool)
		Names() []string
		RegisteredFactories() map[string]T
		Reset()
	}

	registry[T any] struct {
		mu        sync.RWMutex
		factories map[string]T
	}
)

func New[T any]() Registry[T] {
	return &registry[T]{
		factories: make(map[string]T),
	}
}

func (r *registry[T]) Register(name string, f T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = f
}

func (r *registry[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[name]
	return f, ok
}

func (r *registry[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

func (r *registry[T]) RegisteredFactories() map[string]T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make(map[string]T, len(r.factories))
	for k, v := range r.factories {
		res[k] = v
	}
	return res
}

func (r *registry[T]) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories = make(map[string]T)
}

func Get[T any](ctx context.Context, r Registry[T], name string) (T, error) {
	f, ok := r.Get(name)
	if !ok {
		var zero T
		return zero, fmt.Errorf("factory %s not found", name)
	}
	return f, nil
}

func MustGet[T any](ctx context.Context, r Registry[T], name string) T {
	f, err := Get(ctx, r, name)
	if err != nil {
		panic(err)
	}
	return f
}

func Names[T any](ctx context.Context, r Registry[T]) []string {
	return r.Names()
}
