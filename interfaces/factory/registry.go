package factory

import (
	"context"

	"github.com/origadmin/runtime/interfaces/options"
)

// Registry defines the interface for managing factory functions.
// It allows for registering and retrieving factory functions by name.
type Registry[F any] interface {
	// Get retrieves a factory by name
	Get(name string) (F, bool)
	// Register adds or updates a factory with the given name
	Register(name string, factory F)
	// RegisteredFactories returns all registered factories
	RegisteredFactories() map[string]F
	// Reset removes all registered factories
	Reset()
}

// Builder defines the interface for building instances
// with optional configuration and dependencies.
type Builder[T any, C any] interface {
	// WithConfig sets the configuration for the builder
	WithConfig(config C) Builder[T, C]
	// WithOptions sets the dependencies for the builder
	WithOptions(opts ...options.Option) Builder[T, C]
	// Build creates a new instance of T
	Build() (T, error)
}

// Func is a function that creates an instance of T
// using the provided configuration and options.
type Func[T any, C any] func(config C, opts ...options.Option) (T, error)

// ContextFunc is a function that creates an instance of T with context support
// using the provided context, configuration and options.
type ContextFunc[T any, C any] func(ctx context.Context, config C, opts ...options.Option) (T, error)

// FuncRegistry is a specialized registry for factory functions
// that create instances of type T using configuration of type C.
type FuncRegistry[T any, C any] struct {
	Registry[Func[T, C]]
}

// CtxFuncRegistry is a specialized registry for context-aware factory functions
type CtxFuncRegistry[T any, C any] struct {
	Registry[ContextFunc[T, C]]
}

// NewFuncRegistry creates a new FuncRegistry
func NewFuncRegistry[T any, C any]() *FuncRegistry[T, C] {
	return &FuncRegistry[T, C]{}
}

// NewCtxFuncRegistry creates a new CtxFuncRegistry
func NewCtxFuncRegistry[T any, C any]() *CtxFuncRegistry[T, C] {
	return &CtxFuncRegistry[T, C]{}
}

// Create creates a new instance using the factory registered with the given name
func (r *FuncRegistry[T, C]) Create(name string, config C, opts ...options.Option) (T, error) {
	factory, exists := r.Get(name)
	if !exists {
		var zero T
		return zero, ErrFactoryNotFound{Name: name}
	}
	return factory(config, opts...)
}

// Create creates a new instance using the context-aware factory registered with the given name
func (r *CtxFuncRegistry[T, C]) Create(ctx context.Context, name string, config C, opts ...options.Option) (T, error) {
	factory, exists := r.Get(name)
	if !exists {
		var zero T
		return zero, ErrFactoryNotFound{Name: name}
	}
	return factory(ctx, config, opts...)
}

// RegistryBuilder is a specialized registry for builders
// that create instances of type T using configuration of type C.
type RegistryBuilder[T any, C any] struct {
	Registry[func() Builder[T, C]]
}

// NewRegistryBuilder creates a new RegistryBuilder
func NewRegistryBuilder[T any, C any]() *RegistryBuilder[T, C] {
	return &RegistryBuilder[T, C]{}
}

// CreateBuilder creates a new builder using the factory registered with the given name
func (r *RegistryBuilder[T, C]) CreateBuilder(name string) (Builder[T, C], error) {
	factory, exists := r.Get(name)
	if !exists {
		return nil, ErrFactoryNotFound{Name: name}
	}
	return factory(), nil
}

// Create creates a new instance using the builder registered with the given name
func (r *RegistryBuilder[T, C]) Create(name string, config C, opts ...options.Option) (T, error) {
	builder, err := r.CreateBuilder(name)
	if err != nil {
		var zero T
		return zero, err
	}
	return builder.WithConfig(config).WithOptions(opts...).Build()
}

// ErrFactoryNotFound is returned when a requested factory is not found
type ErrFactoryNotFound struct {
	Name string
}

// Error implements error interface
func (e ErrFactoryNotFound) Error() string {
	return "factory not found: " + e.Name
}

// Is implements the interface used by errors.Is
func (e ErrFactoryNotFound) Is(target error) bool {
	_, ok := target.(ErrFactoryNotFound)
	return ok
}
