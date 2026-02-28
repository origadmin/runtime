// Package builderutil implements the functions, types, and contracts for the module.
package builderutil

import (
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/contracts/builder"
	"github.com/origadmin/runtime/contracts/options"
)

// FuncBuilder is a generic builder implementation that uses a Func
type FuncBuilder[T any, C any] struct {
	config C
	opts   []options.Option
	fn     builder.Func[T, C]
}

// NewFuncBuilder creates a new FuncBuilder with the given Func
func NewFuncBuilder[T any, C any](fn builder.Func[T, C]) *FuncBuilder[T, C] {
	return &FuncBuilder[T, C]{fn: fn}
}

// WithConfig implements Builder.WithConfig
func (b *FuncBuilder[T, C]) WithConfig(config C) builder.Builder[T, C] {
	b.config = config
	return b
}

// WithOptions implements Builder.WithOptions
func (b *FuncBuilder[T, C]) WithOptions(opts ...options.Option) builder.Builder[T, C] {
	b.opts = append(b.opts, opts...)
	return b
}

// Build implements Builder.Build
func (b *FuncBuilder[T, C]) Build() (T, error) {
	return b.fn(b.config, b.opts...)
}

// ContextFuncBuilder is a generic builder implementation that uses a ContextFunc
type ContextFuncBuilder[T any, C any] struct {
	ctx    context.Context
	config C
	opts   []options.Option
	fn     builder.ContextFunc[T, C]
}

// NewContextFuncBuilder creates a new ContextFuncBuilder with the given ContextFunc
func NewContextFuncBuilder[T any, C any](fn builder.ContextFunc[T, C]) *ContextFuncBuilder[T, C] {
	return &ContextFuncBuilder[T, C]{
		ctx: context.Background(),
		fn:  fn,
	}
}

// WithConfig implements Builder.WithConfig for ContextFuncBuilder
func (b *ContextFuncBuilder[T, C]) WithConfig(config C) builder.Builder[T, C] {
	b.config = config
	return b
}

// WithContext sets the context for the ContextFuncBuilder
func (b *ContextFuncBuilder[T, C]) WithContext(ctx context.Context) *ContextFuncBuilder[T, C] {
	b.ctx = ctx
	return b
}

// WithOptions implements Builder.WithOptions for ContextFuncBuilder
func (b *ContextFuncBuilder[T, C]) WithOptions(opts ...options.Option) builder.Builder[T, C] {
	b.opts = append(b.opts, opts...)
	return b
}

// Build implements Builder.Build for ContextFuncBuilder
func (b *ContextFuncBuilder[T, C]) Build() (T, error) {
	return b.fn(b.ctx, b.config, b.opts...)
}
