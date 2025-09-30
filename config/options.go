package config

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

// ConfigOptions holds the configuration for the config module.
type ConfigOptions struct {
	ConfigOptions []config.Option
	EnvPrefixes   []string
	Sources       []KSource
}

// Key is the unique key for ConfigOptions in the options.Option.
var Key = optionutil.Key[*ConfigOptions]{}

// Option is a function that configures ConfigOptions.
type Option = options.Option

// WithConfigOption appends Kratos config.Option to the ConfigOptions.
func WithConfigOption(opts ...config.Option) Option {
	return optionutil.Update(func(c *ConfigOptions) {
		c.ConfigOptions = append(c.ConfigOptions, opts...)
	})
}

// WithEnvPrefixes appends environment variable prefixes to the ConfigOptions.
func WithEnvPrefixes(prefixes ...string) Option {
	return optionutil.Update(func(c *ConfigOptions) {
		c.EnvPrefixes = append(c.EnvPrefixes, prefixes...)
	})
}

// WithSource appends config.Source to the ConfigOptions.
func WithSource(s ...config.Source) Option {
	return optionutil.Update(func(c *ConfigOptions) {
		c.Sources = append(c.Sources, s...)
	})
}

// FromOptions retrieves ConfigOptions pointer from the provided options.Option.
// It returns nil if the options are not found or opt is nil.
func FromOptions(opts ...Option) *ConfigOptions {
	var configOpt ConfigOptions
	optionutil.Apply(&configOpt, opts...)
	return &configOpt
}
