/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package selector implements the functions, types, and interfaces for the module.
package selector

import (
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"

	selectorv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/selector"
)

type Selector interface {
	Append(key string, middleware middleware.Middleware) Selector
	Middlewares(filtered bool) []middleware.Middleware
	Build(cfg *selectorv1.Selector, serv bool) []middleware.Middleware
	Total() int
}

type selectorFilter struct {
	matchFunc selector.MatchFunc
	keys      []string
	filtered  []middleware.Middleware
	all       []middleware.Middleware
}

type (
	builderFunc func(...middleware.Middleware) *selector.Builder
)

func (f *selectorFilter) Build(cfg *selectorv1.Selector, serv bool) []middleware.Middleware {
	if cfg == nil || !cfg.Enabled {
		return f.all
	}
	if serv {
		return f.build(cfg, selector.Server)
	}
	return f.build(cfg, selector.Client)
}

type unfilteredSelector struct {
	middlewares []middleware.Middleware
}

func (u unfilteredSelector) Append(key string, middleware middleware.Middleware) Selector {
	u.middlewares = append(u.middlewares, middleware)
	return u
}

func (u unfilteredSelector) Middlewares(bool) []middleware.Middleware {
	return u.middlewares
}

func (u unfilteredSelector) Build(*selectorv1.Selector, bool) []middleware.Middleware {
	return u.middlewares
}

func (u unfilteredSelector) Total() int {
	return len(u.middlewares)
}

func (f *selectorFilter) build(cfg *selectorv1.Selector, fn builderFunc) []middleware.Middleware {
	sc := fn(f.filtered...)
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
	return append([]middleware.Middleware{sc.Build()}, f.all...)
}

func (f *selectorFilter) Total() int {
	return len(f.all)
}

func (f *selectorFilter) Middlewares(filtered bool) []middleware.Middleware {
	if filtered {
		return f.filtered
	}
	return f.all
}

func (f *selectorFilter) Append(key string, m middleware.Middleware) Selector {
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

func NewSelectorFilter(keys []string, matchFunc selector.MatchFunc) Selector {
	return &selectorFilter{
		keys:      keys,
		matchFunc: matchFunc,
	}
}

func New(cfg *selectorv1.Selector, matchFunc selector.MatchFunc) Selector {
	if cfg == nil || !cfg.Enabled {
		return Unfiltered()
	}
	return NewSelectorFilter(cfg.GetNames(), matchFunc)
}

func Unfiltered() Selector {
	return &unfilteredSelector{}
}
