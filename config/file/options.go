/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package file implements the functions, types, and interfaces for the module.
package file

import (
	"github.com/go-kratos/kratos/v2/config"

	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/service/optionutil"
)

var fileOptionKey optionutil.OptionKey[[]Option]

type Option func(*file)

type Formatter func(key string, value []byte) (*config.KeyValue, error)

func WithIgnores(ignores ...string) runtimeconfig.Option {
	return func(options *runtimeconfig.Options) {
		optionutil.WithSliceOption(options.OptionValue, fileOptionKey, func(f *file) {
			f.ignores = append(f.ignores, ignores...)
		})
	}
}

func WithFormatter(formatter Formatter) runtimeconfig.Option {
	return func(options *runtimeconfig.Options) {
		optionutil.WithSliceOption(options.OptionValue, fileOptionKey, func(f *file) {
			f.formatter = formatter
		})
	}
}

func FromOptions(options *runtimeconfig.Options) []Option {
	return optionutil.GetSliceOption(options.OptionValue, fileOptionKey)
}
