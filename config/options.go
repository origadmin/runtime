package config

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/optionutil"
)

// Options contains the options for creating service components.
// It embeds interfaces.Option for common context handling.
type configOptions struct {
	ConfigOptions []config.Option
	EnvPrefixes   []string
	Sources       []KSource
}

type Options = optionutil.Options[configOptions]

// Option is a function that configures service.Options.
type Option func(*Options)

func WithConfigOption(opts ...config.Option) Option {
	return func(o *Options) {
		o.Update(func(v *configOptions) {
			v.ConfigOptions = append(v.ConfigOptions, opts...)
		})
	}
}

func WithEnvPrefixes(prefixes ...string) Option {
	return func(o *Options) {
		o.Update(func(v *configOptions) {
			v.EnvPrefixes = append(v.EnvPrefixes, prefixes...)
		})
	}
}

func WithSource(s ...config.Source) Option {
	return func(o *Options) {
		o.Update(func(v *configOptions) {
			v.Sources = append(v.Sources, s...)
		})
	}
}
