/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"
	middlewareValidate "github.com/go-kratos/kratos/v2/middleware/validate"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware/validate"
)

// Validate is a middleware validator.
// Deprecated: use ValidateServer
func Validate(ms []Middleware, ok bool, validator *configv1.Middleware_Validator) []Middleware {
	if !ok {
		return ms
	}
	switch validate.Version(validator.Version) {
	case validate.V1:
		return append(ms, validateMiddlewareV1(validator))
	case validate.V2:
		return ValidateServer(ms, ok, validator)
	}
	return ms
}

func ValidateServer(ms []Middleware, ok bool, validator *configv1.Middleware_Validator) []Middleware {
	if !ok {
		return ms
	}
	opts := []validate.OptionSetting{
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

func validateMiddlewareV1(_ *configv1.Middleware_Validator) middleware.Middleware {
	return middlewareValidate.Validator()
}
