// Package source implements the functions, types, and interfaces for the module.
package envsource

import (
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/service/optionutil"
)

var envOptionKey optionutil.OptionKey[[]string]

type Option func(*source)

func WithPrefixes(prefixes ...string) runtimeconfig.Option {
	return func(options *runtimeconfig.Options) {
		optionutil.WithSliceOption(options.OptionValue, envOptionKey, prefixes...)
	}
}

func FromOptions(options *runtimeconfig.Options) []string {
	return optionutil.GetSliceOption(options.OptionValue, envOptionKey)
}
