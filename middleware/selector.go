/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"strings"

	"github.com/go-kratos/kratos/v2/middleware/selector"

	selectorv1 "github.com/origadmin/runtime/gen/go/middleware/selector/v1"
)

type Selector interface {
	Append(key string, middleware KMiddleware) Selector
	All() []KMiddleware
	Filtered() []KMiddleware
	Build(cfg *selectorv1.Selector, buildFn func(...KMiddleware) *selector.Builder) []KMiddleware
	Total() int
}

type selectorFilter struct {
	matchFunc selector.MatchFunc
	keys      []string
	filtered  []KMiddleware
	all       []KMiddleware
}

type unfilteredSelector struct {
	middlewares []KMiddleware
}

func (u unfilteredSelector) Append(key string, middleware KMiddleware) Selector {
	u.middlewares = append(u.middlewares, middleware)
	return u
}

func (u unfilteredSelector) All() []KMiddleware {
	return u.middlewares
}

func (u unfilteredSelector) Filtered() []KMiddleware {
	return u.middlewares
}

func (u unfilteredSelector) Build(*selectorv1.Selector, func(...KMiddleware) *selector.Builder) []KMiddleware {
	return u.middlewares
}

func (u unfilteredSelector) Total() int {
	return len(u.middlewares)
}

func (f *selectorFilter) Build(cfg *selectorv1.Selector, fn func(...KMiddleware) *selector.Builder) []KMiddleware {
	sc := selector.Client(f.Filtered()...)
	if f.matchFunc != nil {
		sc.Match(f.matchFunc)
	}
	if path := cfg.GetPaths(); path != nil {
		sc.Path(path...)
	}
	if prefixes := cfg.GetPrefixes(); prefixes != nil {
		sc.Prefix(prefixes...)
	}
	if regex := cfg.GetRegex(); regex != "" {
		sc.Regex(regex)
	}
	return append([]KMiddleware{sc.Build()}, f.All()...)
}

func (f *selectorFilter) Total() int {
	return len(f.all)
}

func (f *selectorFilter) All() []KMiddleware {
	return f.all
}

func (f *selectorFilter) Filtered() []KMiddleware {
	return f.filtered
}

func (f *selectorFilter) Append(key string, m KMiddleware) Selector {
	f.all = append(f.all, m)
	if len(f.keys) == 0 {
		return f
	}
	var kee string
	for _, kee = range f.keys {
		if strings.EqualFold(kee, key) {
			f.filtered = append(f.filtered, m)
			return f
		}
	}
	return f
}

func SelectorFilter(keys []string, matchFunc selector.MatchFunc) Selector {
	return &selectorFilter{
		keys:      keys,
		matchFunc: matchFunc,
	}
}

//func SelectorClient(s Selector, cfg *selectorv1.Selector, matchFunc selector.MatchFunc) []KMiddleware {
//	sc := selector.Client(s.Filtered()...)
//	if matchFunc != nil {
//		sc.Match(matchFunc)
//	}
//	if path := cfg.GetPaths(); path != nil {
//		sc.Path(path...)
//	}
//	if prefixes := cfg.GetPrefixes(); prefixes != nil {
//		sc.Prefix(prefixes...)
//	}
//	if regex := cfg.GetRegex(); regex != "" {
//		sc.Regex(regex)
//	}
//	return append([]KMiddleware{sc.Build()}, s.All()...)
//}
//
//func SelectorServer(s Selector, cfg *selectorv1.Selector, matchFunc selector.MatchFunc) []KMiddleware {
//	sc := selector.Server(s.Filtered()...)
//	if matchFunc != nil {
//		sc.Match(matchFunc)
//	}
//	if path := cfg.GetPaths(); path != nil {
//		sc.Path(path...)
//	}
//	if prefixes := cfg.GetPrefixes(); prefixes != nil {
//		sc.Prefix(prefixes...)
//	}
//	if regex := cfg.GetRegex(); regex != "" {
//		sc.Regex(regex)
//	}
//	return append([]KMiddleware{sc.Build()}, s.All()...)
//}

func WithSelector(cfg *selectorv1.Selector, matchFunc selector.MatchFunc) Selector {
	if cfg == nil || !cfg.Enabled {
		return Unfiltered()
	}
	return SelectorFilter(cfg.GetNames(), matchFunc)
}

func Unfiltered() Selector {
	return &unfilteredSelector{}
}
