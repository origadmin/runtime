/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and contracts for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
)

type rateLimitFactory struct {
}

// NewMiddlewareClient creates a new client-side rate limit middleware.
func (r rateLimitFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	return nil, false
}

// NewMiddlewareServer creates a new server-side rate limit middleware.
func (r rateLimitFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	logger := mwOpts.GetLogger("middleware.ratelimit")
	logger.Debug("enabling rate limit server middleware")

	ratelimitConfig := cfg.GetRateLimiter()
	if ratelimitConfig == nil {
		logger.Debugf("using default BBR rate limiter")
		return ratelimit.Server(), true
	}

	var rlOpts []ratelimit.Option
	switch ratelimitConfig.GetName() {
	case "redis":
		// TODO: Implement Redis rate limiter options
		logger.Warnf("Redis rate limiter not yet implemented")
		return nil, false
	case "memory":
		// TODO: Implement Memory rate limiter options
		logger.Warnf("Memory rate limiter not yet implemented")
		return nil, false
	//case "bbr":
	// default is bbr
	// rlOpts = append(rlOpts, middlewareRateLimit.WithLimiter(bbr.NewLimiter()))
	default:
		// Default to BBR if no specific name is provided or recognized
		logger.Debugf("using default BBR rate limiter")
	}

	return ratelimit.Server(rlOpts...), true
}
