/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/selector"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
)

type Options struct {
	Logger      log.Logger
	MatchFunc   selector.MatchFunc
	Customize   *configv1.Customize
	Middlewares []KMiddleware
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

func WithCustomize(customize *configv1.Customize) Option {
	return func(o *Options) {
		o.Customize = customize
	}
}

func WithMiddlewares(middlewares ...KMiddleware) Option {
	return func(o *Options) {
		o.Middlewares = middlewares
	}

}
