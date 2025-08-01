/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"time"
)

type Options struct {
	Timeout time.Duration
	Retries int
}

type Option = func(o *Options)

func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

func WithRetries(n int) Option {
	return func(o *Options) {
		o.Retries = n
	}
}
