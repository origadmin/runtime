/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

type rateLimitFactory struct {
}

// NewMiddlewareClient creates a new client-side rate limit middleware.
func (r rateLimitFactory) NewMiddlewareClient(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	_, mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debug("[Middleware] Rate limit client middleware enabled, not supported yet")
	return nil, false
}

// NewMiddlewareServer creates a new server-side rate limit middleware.
func (r rateLimitFactory) NewMiddlewareServer(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	_, mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debug("[Middleware] Rate limit server middleware enabled")

	ratelimitConfig := cfg.GetRateLimiter()
	if ratelimitConfig == nil || !ratelimitConfig.GetEnabled() {
		return nil, false
	}

	var rlOpts []ratelimit.Option
	switch ratelimitConfig.GetName() {
	case "redis":
		// TODO: Implement Redis rate limiter options
		helper.Warnf("Redis rate limiter not yet implemented")
	case "memory":
		// TODO: Implement Memory rate limiter options
		helper.Warnf("Memory rate limiter not yet implemented")
	//case "bbr":
	// default is bbr
	// rlOpts = append(rlOpts, middlewareRateLimit.WithLimiter(bbr.NewLimiter()))
	default:
		// Default to BBR if no specific name is provided or recognized
		helper.Infof("Using default BBR rate limiter")
	}

	return ratelimit.Server(rlOpts...), true
}
