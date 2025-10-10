package config

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// Options holds the configuration for the config module.
type Options struct {
	ConfigOptions []KOption
	//EnvPrefixes   []string
	Sources []KSource
}

// WithConfigOption appends Kratos config.Option to the Options.
func WithConfigOption(opts ...config.Option) options.Option {
	return optionutil.WithUpdate(func(c *Options) {
		c.ConfigOptions = append(c.ConfigOptions, opts...)
	})
}

// WithEnvPrefixes appends environment variable prefixes to the Options.
func WithEnvPrefixes(prefixes ...string) options.Option {
	return optionutil.WithUpdate(func(c *Options) {
		//c.EnvPrefixes = append(c.EnvPrefixes, prefixes...)
	})
}

// WithSource appends config.Source to the Options.
func WithSource(s ...config.Source) options.Option {
	return optionutil.WithUpdate(func(c *Options) {
		c.Sources = append(c.Sources, s...)
	})
}

// FromOptions retrieves Options pointer from the provided options.Option.
// It returns nil if the options are not found or opt is nil.
func FromOptions(opts ...options.Option) *Options {
	var configOpt Options
	optionutil.Apply(&configOpt, opts...)
	return &configOpt
}
