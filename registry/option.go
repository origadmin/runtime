/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"time"

	"github.com/origadmin/runtime/optionutil"
)

type registryOption struct {
	Addrs   []string
	Timeout time.Duration
	Secure  bool
}

// Options contains the options for creating registry components.
// It embeds interfaces.ContextOptions for common context handling.
type Options = optionutil.Options[registryOption]

// Option is a function that configures registry.Options.
type Option func(*Options)

// WithAddrs sets the addresses for the registry.
func WithAddrs(addrs ...string) Option {
	return func(o *Options) {
		o.Update(func(v *registryOption) {
			v.Addrs = append(v.Addrs, addrs...)
		})
	}
}

// WithTimeout sets the timeout for registry operations.
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Update(func(v *registryOption) {
			v.Timeout = d
		})
	}
}

// WithSecure enables or disables secure connection.
func WithSecure(secure bool) Option {
	return func(o *Options) {
		o.Update(func(v *registryOption) {
			v.Secure = secure
		})
	}
}
