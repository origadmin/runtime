/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package file implements the functions, types, and interfaces for the module.
package file

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

type Formatter func(key string, value []byte) (*config.KeyValue, error)

type Option = options.Option

func WithIgnores(ignores ...string) options.Option {
	return optionutil.WithUpdate(func(o *file) {
		o.ignores = append(o.ignores, ignores...)
	})
}

func WithFormatter(formatter Formatter) options.Option {
	return optionutil.WithUpdate(func(o *file) {
		o.formatter = formatter
	})
}

// FromOptions extracts file options from the provided runtime options.
// WithCond options is nil or no file options are found, it returns the original file.
func applyFileOptions(f *file, opts ...options.Option) *file {
	optionutil.Apply(f, opts...)
	return f
}
