/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/selector"
)

type Option struct {
	Logger    log.Logger
	MatchFunc selector.MatchFunc
}

type OptionSetting = func(*Option)

func WithLogger(logger log.Logger) OptionSetting {
	return func(o *Option) {
		o.Logger = logger
	}
}

func WithMatchFunc(matchFunc selector.MatchFunc) OptionSetting {
	return func(o *Option) {
		o.MatchFunc = matchFunc
	}
}
