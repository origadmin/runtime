// Package envsrouce implements the functions, types, and interfaces for the module.
package envsrouce

import (
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/optionutil"
)

// optionKey A key used to store the environment variable prefix option in the configuration
var optionKey optionutil.Key[[]string]

// Option defines a function type that is used to configure the source
type Option func(*source)

// WithPrefixes creates an option to set the environment variable prefix
// Parameter prefixes: One or more string prefixes by which environment variables will be filtered
// Return value: Returns a runtimeconfig. Option function,
// which applies the prefix configuration to the configuration options
func WithPrefixes(prefixes ...string) runtimeconfig.Option {
	return func(options *runtimeconfig.Options) {
		optionutil.Append(options, optionKey, prefixes...)
	}
}

// FromOptions extracts the environment variable prefix from the configuration options
// Parameter options: Point to runtimeconfig. Options, which contains configuration options
// Return value: String slice containing the environment variable prefix set in the configuration
func FromOptions(options *runtimeconfig.Options) []string {
	return optionutil.Slice(options, optionKey)
}
