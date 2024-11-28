/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	middlewareValidate "github.com/go-kratos/kratos/v2/middleware/validate"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/validate"
)

func Validate(ms []Middleware, ok bool, validator *configv1.Middleware_Validator) []Middleware {
	if !ok {
		return ms
	}
	switch validate.Version(validator.Version) {
	case validate.V1:
		return append(ms, validateMiddlewareV1(validator))
	case validate.V2:
		return append(ms, validateMiddlewareV2(validator))
	}
	return ms
}

func validateMiddlewareV1(_ *configv1.Middleware_Validator) middleware.Middleware {
	return middlewareValidate.Validator()
}

func validateMiddlewareV2(validator *configv1.Middleware_Validator) middleware.Middleware {
	opts := []validate.ProtoValidatorOption{
		validate.WithFailFast(validator.GetFailFast()),
	}

	val, err := validate.NewValidate(opts...)
	if err != nil {
		return Empty()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(proto.Message); ok {
				if err := val.Validate(v); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}
