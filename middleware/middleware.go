/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middlewares implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

const Type = "middleware"

type (
	Handler    = middleware.Handler
	Middleware = middleware.Middleware
)

// Chain returns a middleware that executes a chain of middleware.
func Chain(m ...Middleware) Middleware {
	return middleware.Chain(m...)
}

func NewClient(cfg *configv1.Middleware) []Middleware {
	var middlewares []Middleware

	if cfg == nil {
		return middlewares
	}
	middlewares = Recovery(middlewares, cfg.EnableRecovery)
	//middlewares = Validate(middlewares, cfg.EnableValidate, cfg.Validator)
	middlewares = SecurityClient(middlewares, cfg.Security) // 新增 JWT 中间件
	//middlewares = CasbinClient(middlewares, cfg.EnableCasbin, cfg.CasbinConfig) // 新增 Casbin 中间件
	middlewares = MetadataClient(middlewares, cfg.EnableMetadata, cfg.Metadata)
	middlewares = TracingClient(middlewares, cfg.EnableTracing)
	middlewares = CircuitBreakerClient(middlewares, cfg.EnableCircuitBreaker)
	return middlewares
}

func NewServer(cfg *configv1.Middleware) []Middleware {
	var middlewares []Middleware

	if cfg == nil {
		return middlewares
	}
	middlewares = Recovery(middlewares, cfg.EnableRecovery)
	middlewares = ValidateServer(middlewares, cfg.EnableValidate, cfg.Validator)
	//middlewares = JWTServer(middlewares, cfg.EnableJWT, cfg.JWTConfig)         added jwt middleware
	middlewares = SecurityServer(middlewares, cfg.Security) // added casbin middleware
	middlewares = TracingServer(middlewares, cfg.EnableTracing)
	middlewares = MetadataServer(middlewares, cfg.EnableMetadata, cfg.Metadata)
	middlewares = RateLimitServer(middlewares, cfg.RateLimiter)
	return middlewares
}
