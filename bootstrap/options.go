package bootstrap

import (
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces"
)

// Options holds the configuration options for bootstrap
// This struct is used to configure the bootstrap process itself.
type Options struct {
	BootstrapOptions []runtimeconfig.KOption
	ConfigOptions    []runtimeconfig.Option
	DecoderProvider  interfaces.ConfigDecoderProvider
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

// WithDecoderProvider sets the custom decoder provider.
// If not set, a default decoder provider will be used.
func WithDecoderProvider(p interfaces.ConfigDecoderProvider) Option {
	return func(o *Options) {
		o.DecoderProvider = p
	}
}
