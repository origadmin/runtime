package bootstrap

import (
	runtimeconfig "github.com/origadmin/runtime/config"
)

// Options holds the configuration options for bootstrap
type Options struct {
	BootstrapOptions []runtimeconfig.KOption
	ConfigOptions    []runtimeconfig.Option
}

// Option is a function that configures bootstrap.Options
type Option func(*Options)

// WithBootstrapOptions sets the options for loading bootstrap config
func WithBootstrapOptions(opts ...runtimeconfig.KOption) Option {
	return func(o *Options) {
		o.BootstrapOptions = append(o.BootstrapOptions, opts...)
	}
}

// WithConfigOptions sets the options for creating runtime config
func WithConfigOptions(opts ...runtimeconfig.Option) Option {
	return func(o *Options) {
		o.ConfigOptions = append(o.ConfigOptions, opts...)
	}
}
