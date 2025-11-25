package runtime

import (
	"errors"

	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/container"
	"github.com/origadmin/runtime/interfaces/options"
)

// Builder is the interface for building a runtime.App.
type Builder interface {
	// WithAppInfo creates and configures the application's metadata in a single call.
	// It uses the internal newAppInfo to create a concrete *appInfo struct.
	WithAppInfo(name, version string, opts ...AppInfoOption) Builder

	// WithBootstrapOptions applies options to the underlying bootstrap process.
	WithBootstrapOptions(opts ...options.Option) Builder
	// WithContainerOptions applies options to the underlying container.
	WithContainerOptions(opts ...options.Option) Builder
	// Build creates a new App instance.
	Build(bootstrapPath string) (*App, error)
}

type builder struct {
	opts appOptions
}

// NewBuilder creates a new runtime builder.
func NewBuilder() Builder {
	return &builder{}
}

// Build creates a new App instance from the builder's configuration.
func (b *builder) Build(bootstrapPath string) (*App, error) {
	// 1. Validate that AppInfo has been configured.
	if b.opts.appInfo == nil {
		return nil, errors.New("application info is not configured, use WithAppInfo()")
	}

	// 2. Call bootstrap.New with the collected bootstrap options.
	bootstrapResult, err := bootstrap.New(bootstrapPath, b.opts.bootstrapOpts...)
	if err != nil {
		return nil, err
	}

	// 3. Prepare container options, ensuring the pre-built concrete *appInfo is passed.
	// The container will accept it as an interfaces.AppInfo.
	ctnOpts := []options.Option{container.WithAppInfo(b.opts.appInfo)}
	if len(b.opts.containerOpts) > 0 {
		ctnOpts = append(ctnOpts, b.opts.containerOpts...)
	}

	// 4. Create the App instance.
	rt := New(bootstrapResult, ctnOpts...)
	return rt, nil
}

// WithAppInfo creates and stores the concrete *appInfo instance immediately.
func (b *builder) WithAppInfo(name, version string, opts ...AppInfoOption) Builder {
	b.opts.appInfo = newAppInfo(name, version, opts...)
	return b
}

// WithBootstrapOptions applies options to the underlying bootstrap process.
func (b *builder) WithBootstrapOptions(opts ...options.Option) Builder {
	b.opts.bootstrapOpts = append(b.opts.bootstrapOpts, opts...)
	return b
}

// WithContainerOptions applies options to the underlying container.
func (b *builder) WithContainerOptions(opts ...options.Option) Builder {
	b.opts.containerOpts = append(b.opts.containerOpts, opts...)
	return b
}
