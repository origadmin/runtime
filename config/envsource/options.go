// Package envsrouce implements the functions, types, and interfaces for the module.
package envsource

import (
	options "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/optionutil"
)

// Option defines a function type that is used to configure the source
type envOptions struct {
	prefixes []string
}

// WithPrefixes creates an option to set the environment variable prefix
// Parameter prefixes: One or more string prefixes by which environment variables will be filtered
// Return value: Returns a options. Option function,
// which applies the prefix configuration to the configuration options
func WithPrefixes(prefixes ...string) options.Option {
	return optionutil.WithUpdate(func(o *envOptions) {
		o.prefixes = append(o.prefixes, prefixes...)
	})
}

// FromOptions extracts the environment variable prefix from the configuration options
// Parameter options: Point to options. Options, which contains configuration options
// Return value: String slice containing the environment variable prefix set in the configuration
func FromOptions(opts ...options.Option) []string {
	var envOpts envOptions
	optionutil.Apply(&envOpts, opts...)
	return envOpts.prefixes
}
