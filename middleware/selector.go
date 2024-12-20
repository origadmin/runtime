/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	selectorv1 "github.com/origadmin/runtime/gen/go/middleware/selector/v1"
)

func SelectorClient(f Filter, cfg *selectorv1.Selector, matchFunc selector.MatchFunc) []KMiddleware {
	s := selector.Client(f.Filtered()...)
	if matchFunc != nil {
		s.Match(matchFunc)
	}
	if path := cfg.GetPaths(); path != nil {
		s.Path(path...)
	}
	if prefixes := cfg.GetPrefixes(); prefixes != nil {
		s.Prefix(prefixes...)
	}
	if regex := cfg.GetRegex(); regex != "" {
		s.Regex(regex)
	}
	return append([]KMiddleware{s.Build()}, f.All()...)
}

func SelectorServer(f Filter, cfg *selectorv1.Selector, matchFunc selector.MatchFunc) []KMiddleware {
	s := selector.Server(f.Filtered()...)
	if matchFunc != nil {
		s.Match(matchFunc)
	}
	if path := cfg.GetPaths(); path != nil {
		s.Path(path...)
	}
	if prefixes := cfg.GetPrefixes(); prefixes != nil {
		s.Prefix(prefixes...)
	}
	if regex := cfg.GetRegex(); regex != "" {
		s.Regex(regex)
	}
	return append([]KMiddleware{s.Build()}, f.All()...)
}
