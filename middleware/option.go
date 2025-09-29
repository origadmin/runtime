/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/selector"
)

type Options struct {
	Logger      log.Logger
	MatchFunc   selector.MatchFunc
	Middlewares []KMiddleware
}

type Option = func(*Options)

func WithMatchFunc(matchFunc selector.MatchFunc) Option {
	return func(o *Options) {
		o.MatchFunc = matchFunc
	}
}

func WithMiddlewares(middlewares ...KMiddleware) Option {
	return func(o *Options) {
		o.Middlewares = middlewares
	}

}
