package config

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces"
)

// Options contains the options for creating service components.
// It embeds optionvalue.OptionValue for common context handling.
type Options struct {
	interfaces.OptionValue // Updated type
	ConfigOptions          []config.Option
	EnvPrefixes            []string
}

// Option is a function that configures service.Options.
type Option func(*Options)

func DefaultServerOptions() *Options {
	return &Options{
		OptionValue: interfaces.DefaultOptions(), // Updated function call
	}
}

func WithConfigOptions(opts ...config.Option) Option {
	return func(o *Options) {
		o.ConfigOptions = append(o.ConfigOptions, opts...)
	}
}

func WithEnvPrefixes(prefixes ...string) Option {
	return func(o *Options) {
		o.EnvPrefixes = append(o.EnvPrefixes, prefixes...)
	}
}
