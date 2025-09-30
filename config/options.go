package config

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// configOptions holds the configuration for the config module.
type configOptions struct {
	ConfigOptions []config.Option
	EnvPrefixes   []string
	Sources       []KSource
}

// Key is the unique key for configOptions in the options.Option.
var Key = optionutil.Key[*configOptions]{}

// Option is a function that configures configOptions.
type Option = options.Option

// WithConfigOption appends Kratos config.Option to the configOptions.
func WithConfigOption(opts ...config.Option) Option {
	return optionutil.Update(func(c *configOptions) {
		c.ConfigOptions = append(c.ConfigOptions, opts...)
	})
}

// WithEnvPrefixes appends environment variable prefixes to the configOptions.
func WithEnvPrefixes(prefixes ...string) Option {
	return optionutil.Update(func(c *configOptions) {
		c.EnvPrefixes = append(c.EnvPrefixes, prefixes...)
	})
}

// WithSource appends config.Source to the configOptions.
func WithSource(s ...config.Source) Option {
	return optionutil.Update(func(c *configOptions) {
		c.Sources = append(c.Sources, s...)
	})
}

// FromOptions retrieves configOptions pointer from the provided options.Option.
// It returns nil if the options are not found or opt is nil.
func FromOptions(opts ...Option) *configOptions {
	var configOpt configOptions
	optionutil.Apply(&configOpt, opts...)
	return &configOpt
}
