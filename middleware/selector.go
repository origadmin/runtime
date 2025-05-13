/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	selectorv1 "github.com/origadmin/runtime/gen/go/middleware/selector/v1"
	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
)

type selectorFactory struct {
}

func (s selectorFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	if middleware.Selector == nil || !middleware.Selector.Enabled {
		return nil, false
	}
	return SelectorClient(middleware.Selector, options.MatchFunc, options.Middlewares[0]), true
}

func (s selectorFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	if middleware.Selector == nil || !middleware.Selector.Enabled {
		return nil, false
	}
	return SelectorServer(middleware.Selector, options.MatchFunc, options.Middlewares[0]), true
}

func SelectorServer(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middleware KMiddleware) KMiddleware {
	if cfg == nil || !cfg.Enabled {
		return middleware
	}
	return selectorBuilder(cfg, selector.Server(middleware), matchFunc)
}

func SelectorClient(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middleware KMiddleware) KMiddleware {
	if cfg == nil || !cfg.Enabled {
		return middleware
	}
	return selectorBuilder(cfg, selector.Client(middleware), matchFunc)
}

func selectorBuilder(cfg *selectorv1.Selector, builder *selector.Builder, matchFunc selector.MatchFunc) KMiddleware {
	if matchFunc != nil {
		builder.Match(matchFunc)
	}
	if path := cfg.GetPaths(); path != nil {
		builder.Path(path...)
	}
	if prefixes := cfg.GetPrefixes(); prefixes != nil {
		builder.Prefix(prefixes...)
	}
	if regex := cfg.GetRegex(); regex != "" {
		builder.Regex(regex)
	}
	return builder.Match(matchFunc).Build()
}
