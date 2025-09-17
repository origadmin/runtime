/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package file implements the functions, types, and interfaces for the module.
package file

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/goexts/generic/configure"

	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/optionutil"
)

var optionKey optionutil.Key[[]Option]

type Option func(*file)

type Formatter func(key string, value []byte) (*config.KeyValue, error)

func WithIgnores(ignores ...string) runtimeconfig.Option {
	return func(options *runtimeconfig.Options) {
		opt := func(f *file) {
			f.ignores = append(f.ignores, ignores...)
		}
		optionutil.Append(options, optionKey, opt)
	}
}

func WithFormatter(formatter Formatter) runtimeconfig.Option {
	return func(options *runtimeconfig.Options) {
		// Create a new option that sets the formatter
		opt := func(f *file) {
			f.formatter = formatter
		}
		// Append the option to the options slice
		optionutil.Append(options, optionKey, opt)
	}
}

// FromOptions extracts file options from the provided runtime options.
// If options is nil or no file options are found, it returns the original file.
func applyFileOptions(f *file, options *runtimeconfig.Options) *file {
	return configure.Apply(f, FromOptions(options))
}

func FromOptions(options *runtimeconfig.Options) []Option {
	if options == nil {
		return nil
	}
	return optionutil.Slice(options, optionKey)
}
