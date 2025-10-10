/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/selector"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

type selectorFactory struct {
}

func (s selectorFactory) NewMiddlewareClient(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	selectorConfig := cfg.GetSelector()
	if selectorConfig == nil {
		return nil, false
	}

	helper.Infof("enabling client selector middleware")

	var mws []KMiddleware
	for _, name := range selectorConfig.Names {
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

func (s selectorFactory) NewMiddlewareServer(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	selectorConfig := cfg.GetSelector()
	if selectorConfig == nil {
		return nil, false
	}

	helper.Infof("enabling server selector middleware")
	var mws []KMiddleware
	for _, name := range cfg.Names {
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
// 增强selectorBuilder函数，支持exclude_middlewares
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
	// 实现excludePaths逻辑（如果需要）
	if cfg.ExcludePaths != nil {
		builder.ExcludePath(cfg.ExcludePaths...)
	}
	if cfg.ExcludePrefixes != nil {
		builder.ExcludePrefix(cfg.ExcludePrefixes...)
	}
	if cfg.ExcludeRegex != "" {
		builder.ExcludeRegex(cfg.ExcludeRegex)
	}
	return builder.Build()
}