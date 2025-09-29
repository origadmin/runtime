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
	_, mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	selectorConfig := cfg.GetSelector()
	if selectorConfig == nil || !selectorConfig.GetEnabled() {
		return nil, false
	}

	helper.Infof("enabling client selector middleware")

	// Create a selector builder that wraps no initial middlewares.
	// The actual middlewares to be selected will be determined by the selector logic
	// when it's applied in the chain.
	builder := selector.Client()

	return selectorBuilder(selectorConfig, builder, mwOpts.MatchFunc), true
}

func (s selectorFactory) NewMiddlewareServer(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	_, mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	selectorConfig := cfg.GetSelector()
	if selectorConfig == nil || !selectorConfig.GetEnabled() {
		return nil, false
	}

	helper.Infof("enabling server selector middleware")

	// Create a selector builder that wraps no initial middlewares.
	// The actual middlewares to be selected will be determined by the selector logic
	// when it's applied in the chain.
	builder := selector.Server()

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
