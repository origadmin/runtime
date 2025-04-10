/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package customize

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// Options is a struct that holds a value.
type Options struct {
	Customize *configv1.Customize
}

// Option is a function that sets a value on a Setting.
type Option = func(config *Options)

// WithCustomize sets the customize config.
func WithCustomize(customize *configv1.Customize) Option {
	return func(option *Options) {
		option.Customize = customize
	}
}
