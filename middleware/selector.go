/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/selector/v1"
	"github.com/origadmin/runtime/log"
)

type selectorFactory struct {
}

func (s selectorFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	selectorConfig := cfg.GetSelector()
	if selectorConfig == nil {
		return nil, false
	}

	helper.Infof("enabling client selector middleware")

	var mws []KMiddleware
	// Determine target middleware names by includes/excludes before fetching instances
	var names []string
	includes := selectorConfig.GetIncludes()
	excludes := selectorConfig.GetExcludes()
	names = append(names, includes...)
	// apply excludes filter
	if len(excludes) > 0 {
		ex := make(map[string]struct{}, len(excludes))
		for _, n := range excludes {
			ex[n] = struct{}{}
		}
		var filtered []string
		for _, n := range names {
			if _, skip := ex[n]; !skip {
				filtered = append(filtered, n)
			}
		}
		names = filtered
	}
	// fetch middlewares by final names
	for _, name := range names {
		helper.Infof("enabling client selector middleware: %s", name)
		middleware, ok := mwOpts.Carrier.Clients[name]
		if !ok {
			helper.Warnf("unknown client selector middleware: %s", name)
			continue
		}
		mws = append(mws, middleware)
	}

	// Create a selector builder that wraps no initial middlewares.
	// The actual middlewares to be selected will be determined by the selector logic
	// when it's applied in the chain.
	builder := selector.Client(mws...)

	return selectorBuilder(selectorConfig, builder, mwOpts.MatchFunc), true
}

func (s selectorFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	selectorConfig := cfg.GetSelector()
	if selectorConfig == nil {
		return nil, false
	}

	helper.Infof("enabling server selector middleware")
	var mws []KMiddleware
	// Determine target middleware names by includes/excludes before fetching instances
	var names []string
	includes := selectorConfig.GetIncludes()
	excludes := selectorConfig.GetExcludes()
	names = append(names, includes...)
	// apply excludes filter
	if len(excludes) > 0 {
		ex := make(map[string]struct{}, len(excludes))
		for _, n := range excludes {
			ex[n] = struct{}{}
		}
		var filtered []string
		for _, n := range names {
			if _, skip := ex[n]; !skip {
				filtered = append(filtered, n)
			}
		}
		names = filtered
	}
	// fetch middlewares by final names
	for _, name := range names {
		helper.Infof("enabling server selector middleware: %s", name)
		middleware, ok := mwOpts.Carrier.Servers[name]
		if !ok {
			helper.Warnf("unknown server selector middleware: %s", name)
			continue
		}
		mws = append(mws, middleware)
	}

	// Create a selector builder that wraps no initial middlewares.
	// The actual middlewares to be selected will be determined by the selector logic
	// when it's applied in the chain.
	builder := selector.Server(mws...)

	return selectorBuilder(selectorConfig, builder, mwOpts.MatchFunc), true
}

// SelectorServer creates a selector middleware for server-side.
// This helper function is still available for direct use if needed to wrap specific middlewares.
func SelectorServer(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middlewares ...KMiddleware) KMiddleware {
	return selectorBuilder(cfg, selector.Server(middlewares...), matchFunc)
}

// SelectorClient creates a selector middleware for client-side.
// This helper function is still available for direct use if needed to wrap specific middlewares.
func SelectorClient(cfg *selectorv1.Selector, matchFunc selector.MatchFunc, middlewares ...KMiddleware) KMiddleware {
	return selectorBuilder(cfg, selector.Client(middlewares...), matchFunc)
}

// selectorBuilder configures and builds a Kratos selector middleware.
// Enhance selectorBuilder to support exclude_middlewares
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
