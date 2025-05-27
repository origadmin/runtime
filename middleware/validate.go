/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	validatorv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/validator"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware/validate"
)

type validatorFactory struct {
}

func (f validatorFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	return nil, false
}

func (f validatorFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] Validator server middleware enabled")
	if !checkEnabled(middleware, "validator") {
		return nil, false
	}
	cfg := middleware.GetValidator()
	switch validate.Version(cfg.GetVersion()) {
	case validate.V2:
		opts := []validate.Option{
			validate.WithFailFast(cfg.GetFailFast()),
		}
		if validate.Version(cfg.Version) == validate.V2 {
			opts = append(opts, validate.WithV2ProtoValidatorOptions())
		}
		if m, err := validate.Server(opts...); err == nil {
			return m, true
		}
	default:
		return validateMiddlewareV1(cfg), true
	}
	return nil, false
}

// Validate is a middleware validator.
// Deprecated: use ValidateServer
func Validate(ms []KMiddleware, validator *validatorv1.Validator) []KMiddleware {
	switch validate.Version(validator.Version) {
	case validate.V1:
		return append(ms, validateMiddlewareV1(validator))
	case validate.V2:
		return ValidateServer(ms, validator)
	}
	return ms
}

func ValidateServer(ms []KMiddleware, validator *validatorv1.Validator) []KMiddleware {
	opts := []validate.Option{
		validate.WithFailFast(validator.GetFailFast()),
	}
	if validate.Version(validator.Version) == validate.V2 {
		opts = append(opts, validate.WithV2ProtoValidatorOptions())
	}
	if m, err := validate.Server(opts...); err == nil {
		ms = append(ms, m)
	}
	return ms
}

func validateMiddlewareV1(_ *validatorv1.Validator) middleware.Middleware {
	return Validator()
}

type validator interface {
	Validate() error
}

func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}
