/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"github.com/origadmin/runtime/middleware/security"
)

type Option struct {
	security []security.OptionSetting
}

type OptionSetting = func(*Option)

func (o Option) SecurityOptions() []security.OptionSetting {
	return o.security
}

func WithSecurityOptions(ss ...security.OptionSetting) OptionSetting {
	return func(option *Option) {
		if option.security == nil {
			option.security = ss
		}
		option.security = append(option.security, ss...)
	}
}
