package config

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/optionutil"
)

// KSource is an alias for config.Source, assuming its definition.
type KSource = config.Source

// configOptions holds the configuration for the config module.
type configOptions struct {
	ConfigOptions []config.Option
	EnvPrefixes   []string
	Sources       []KSource
}

// Key is the unique key for configOptions in the options.Option.
var Key = optionutil.Key[*configOptions]{}

// Option is a function that configures configOptions.
type Option options.OptionFunc

// WithConfigOption appends Kratos config.Option to the configOptions.
func WithConfigOption(opts ...config.Option) Option {
	return optionutil.Update(Key, func(c *configOptions) {
		c.ConfigOptions = append(c.ConfigOptions, opts...)
	})
}

// WithEnvPrefixes appends environment variable prefixes to the configOptions.
func WithEnvPrefixes(prefixes ...string) Option {
	return optionutil.Update(Key, func(c *configOptions) {
		c.EnvPrefixes = append(c.EnvPrefixes, prefixes...)
	})
}

// WithSource appends config.Source to the configOptions.
func WithSource(s ...config.Source) Option {
	return optionutil.Update(Key, func(c *configOptions) {
		c.Sources = append(c.Sources, s...)
	})
}

// FromOption retrieves configOptions pointer from the provided options.Option.
// It returns nil if the options are not found or opt is nil.
func FromOption(opt options.Option) *configOptions {
	if opt == nil {
		return nil
	}
	val, ok := optionutil.Value(opt, Key)
	// If not found or nil, return nil
	if !ok || val == nil {
		return nil
	}
	return val
}
