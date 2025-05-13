/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package file implements the functions, types, and interfaces for the module.
package file

type Option func(*file)

func WithIgnores(ignores ...string) Option {
	return func(o *file) {
		o.ignores = ignores
	}
}
