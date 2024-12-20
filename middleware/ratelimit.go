/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"

	ratelimitv1 "github.com/origadmin/runtime/gen/go/middleware/ratelimit/v1"
	"github.com/origadmin/runtime/log"
)

func RateLimitServer(f Filter, cfg *ratelimitv1.RateLimiter) Filter {
	if cfg == nil {
		return f
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
	log.Debugf("[KMiddleware] Rate limit server middleware enabled with %v", cfg.GetName())
	return f.Filter("RateLimit", ratelimit.Server(options...))
}
