/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

type Option func(l *loader)

func WithIgnores(ignores ...string) Option {
	return func(l *loader) {
		l.ignores = ignores
	}
}
