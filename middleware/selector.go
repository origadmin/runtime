/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	selectorv1 "github.com/origadmin/runtime/gen/go/middleware/selector/v1"
)

func SelectorServer(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middleware KMiddleware) KMiddleware {
	if cfg == nil || !cfg.Enabled {
		return middleware
	}
	return selectorBuilder(selector.Server(middleware), cfg, matchFunc)
}

func SelectorClient(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middleware KMiddleware) KMiddleware {
	if cfg == nil || !cfg.Enabled {
		return middleware
	}
	return selectorBuilder(selector.Client(middleware), cfg, matchFunc)
}

func selectorBuilder(builder *selector.Builder, cfg *selectorv1.Selector, matchFunc selector.MatchFunc) KMiddleware {
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
