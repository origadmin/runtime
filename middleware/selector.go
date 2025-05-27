/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/selector"
)

type selectorFactory struct {
}

func (s selectorFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	if !middleware.GetSelector().GetEnabled() {
		return nil, false
	}
	return SelectorClient(middleware.Selector, options.MatchFunc, options.Middlewares...), true
}

func (s selectorFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	if !middleware.GetSelector().GetEnabled() {
		return nil, false
	}
	return SelectorServer(middleware.Selector, options.MatchFunc, options.Middlewares...), true
}

func SelectorServer(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middlewares ...KMiddleware) KMiddleware {
	return selectorBuilder(cfg, selector.Server(middlewares...), matchFunc)
}

func SelectorClient(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middlewares ...KMiddleware) KMiddleware {
	return selectorBuilder(cfg, selector.Client(middlewares...), matchFunc)
}

func selectorBuilder(cfg *selectorv1.Selector, builder *selector.Builder, matchFunc selector.MatchFunc) KMiddleware {
	if matchFunc != nil {
		builder.Match(matchFunc)
	}
	if cfg == nil {
		return builder.Build()
	}
	if cfg.Paths != nil {
		builder.Path(cfg.Paths...)
	}
	if cfg.Prefixes != nil {
		builder.Prefix(cfg.Prefixes...)
	}
	if cfg.Regex != "" {
		builder.Regex(cfg.Regex)
	}
	return builder.Build()
}
