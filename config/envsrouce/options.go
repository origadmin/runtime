// Package source implements the functions, types, and interfaces for the module.
package envsrouce

import (
	"github.com/origadmin/runtime/bootstrap"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/service/optionutil"
)

var envOptionKey optionutil.OptionKey[[]string]

type Option func(*source)

func WithPrefixes(prefixes ...string) bootstrap.Option {
	return func(options *runtimeconfig.Options) {
		optionutil.WithSliceOption(options.OptionValue, envOptionKey, prefixes...)
	}
}

func FromOptions(options *runtimeconfig.Options) []string {
	return optionutil.GetSliceOption(options.OptionValue, envOptionKey)
}
