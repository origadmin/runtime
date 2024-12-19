/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/middleware/security"
)

type Option struct {
	securities []security.OptionSetting
	Logger     log.Logger
}

type OptionSetting = func(*Option)

func (o Option) Securities() []security.OptionSetting {
	return o.securities
}

func WithSecurityOptions(ss ...security.OptionSetting) OptionSetting {
	return func(option *Option) {
		if option.securities == nil {
			option.securities = ss
		}
		option.securities = append(option.securities, ss...)
	}
}
