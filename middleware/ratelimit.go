/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"

	ratelimitv1 "github.com/origadmin/runtime/gen/go/middleware/ratelimit/v1"
	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type rateLimitFactory struct {
}

func (r rateLimitFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	return nil, false
}

func (r rateLimitFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	if middleware.GetRateLimiter().GetEnabled() {
		options := make([]ratelimit.Option, 0)
		return ratelimit.Server(options...), true
	}
	return nil, false
}

func RateLimitServer(ms []KMiddleware, cfg *ratelimitv1.RateLimiter) []KMiddleware {
	if cfg == nil {
		return ms
	}
	var options []ratelimit.Option
	switch cfg.GetName() {
	case "redis":
		// TODO:
	case "memory":
		// TODO:
	//case "bbr":
	// default is bbr
	// options = append(options, middlewareRateLimit.WithLimiter(bbr.NewLimiter()))
	default:
		// do nothing
	}
	log.Debugf("[Middleware] Rate limit server middleware enabled with %v", cfg.GetName())
	return append(ms, ratelimit.Server(options...))
}
