/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"
	middlewareValidate "github.com/go-kratos/kratos/v2/middleware/validate"

	validatorv1 "github.com/origadmin/runtime/gen/go/middleware/validator/v1"
	"github.com/origadmin/runtime/middleware/validate"
)

// Validate is a middleware validator.
// Deprecated: use ValidateServer
func Validate(f Filter, validator *validatorv1.Validator) Filter {
	switch validate.Version(validator.Version) {
	case validate.V1:
		return f.Filter("Validate", validateMiddlewareV1(validator))
	case validate.V2:
		return ValidateServer(f, validator)
	}
	return f
}

func ValidateServer(f Filter, validator *validatorv1.Validator) Filter {
	opts := []validate.OptionSetting{
		validate.WithFailFast(validator.GetFailFast()),
	}
	if validate.Version(validator.Version) == validate.V2 {
		opts = append(opts, validate.WithV2ProtoValidatorOptions())
	}
	if m, err := validate.Server(opts...); err == nil {
		f.Filter("Validate", m)
	}
	return f
}

func validateMiddlewareV1(_ *validatorv1.Validator) middleware.Middleware {
	return middlewareValidate.Validator()
}
