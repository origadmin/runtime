/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
)

func RateLimitServer(ms []Middleware, cfg *configv1.Middleware_RateLimiter) []Middleware {
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
