package config

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces"
)

// Options contains the options for creating service components.
// It embeds interfaces.ContextOptions for common context handling.
type Options struct {
	interfaces.OptionValue
	ConfigOptions []config.Option
}

// Option is a function that configures service.Options.
type Option func(*Options)

func DefaultServerOptions() *Options {
	return &Options{
		OptionValue: interfaces.DefaultOptions(),
	}
}

func WithConfigOptions(opts ...config.Option) Option {
	return func(o *Options) {
		o.ConfigOptions = append(o.ConfigOptions, opts...)
	}
}
