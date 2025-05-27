/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package file implements the functions, types, and interfaces for the module.
package file

import (
	"github.com/go-kratos/kratos/v2/config"
)

type Option func(*file)

type Formatter func(key string, value []byte) (*config.KeyValue, error)

func WithIgnores(ignores ...string) Option {
	return func(o *file) {
		o.ignores = ignores
	}
}

func WithFormatter(f Formatter) Option {
	return func(o *file) {
		o.formatter = f
	}
}
