/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/selector"
)

type Options struct {
	Logger    log.Logger
	MatchFunc selector.MatchFunc
}

type Option = func(*Options)

func WithLogger(logger log.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithMatchFunc(matchFunc selector.MatchFunc) Option {
	return func(o *Options) {
		o.MatchFunc = matchFunc
	}
}
