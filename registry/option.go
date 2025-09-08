/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"context"
	"time"
)

// Options contains the options for creating registry components.
// It uses a context to carry implementation-specific values.
type Options struct {
	Context context.Context
	Addrs   []string
	Timeout time.Duration
	Secure  bool
}

// Option is a function that configures Options.
type Option func(o *Options)

// WithContext sets the context for the options.
func WithContext(ctx context.Context) Option {
	return func(o *Options) { o.Context = ctx }
}

// WithAddrs sets the addresses for the registry.
func WithAddrs(addrs ...string) Option {
	return func(o *Options) { o.Addrs = addrs }
}

// WithTimeout sets the timeout for registry operations.
func WithTimeout(d time.Duration) Option {
	return func(o *Options) { o.Timeout = d }
}

// WithSecure enables or disables secure connection.
func WithSecure(secure bool) Option {
	return func(o *Options) { o.Secure = secure }
}
