/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package customize

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// Option is a struct that holds a value.
type Option struct {
	Customize *configv1.Customize
}

// OptionSetting is a function that sets a value on a Setting.
type OptionSetting = func(config *Option)

// WithCustomize sets the customize config.
func WithCustomize(customize *configv1.Customize) OptionSetting {
	return func(option *Option) {
		option.Customize = customize
	}
}
