/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	validatorv1 "github.com/origadmin/runtime/api/gen/go/config/middleware/validator/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware/validate"
)

type validatorFactory struct {
}

func (f validatorFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	return nil, false
}

func (f validatorFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(log.With(mwOpts.Logger, "module", "middleware.validator"))
	helper.Debug("enabling validator server middleware")

	if !middleware.GetEnabled() {
		return nil, false
	}
	cfg := middleware.GetValidator()
	switch middleware.GetType() {
	case string(Validator):
		switch validate.Version(cfg.GetVersion()) {
		case validate.V2:
			opts := []validate.Option{
				validate.WithFailFast(cfg.GetFailFast()),
			}
			if validate.Version(cfg.Version) == validate.V2 {
				opts = append(opts, validate.WithV2ProtoValidatorOption())
			}
			if m, err := validate.Server(opts...); err == nil {
				return m, true
			}
		default:
			return newValidatorV1(helper), true
		}
	}
	return nil, false
}

// Deprecated: use ValidateServer
func Validate(ms []KMiddleware, validator *validatorv1.Validator) []KMiddleware {
	switch validate.Version(validator.Version) {
	case validate.V1:
		helper := log.NewHelper(log.DefaultLogger) // Cannot inject logger here easily
		return append(ms, newValidatorV1(helper))
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
		opts = append(opts, validate.WithV2ProtoValidatorOption())
	}
	if m, err := validate.Server(opts...); err == nil {
		ms = append(ms, m)
	}
	return ms
}

type validator interface {
	Validate() error
}

// newValidatorV1 is the constructor for the v1 validation middleware.
func newValidatorV1(logger *log.Helper) KMiddleware {
	return func(handler KHandler) KHandler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					logger.WithContext(ctx).Debugf("[Validate] Validation failed for request: %v", err)
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}
