/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
)

type Validator interface {
	Validate(ctx context.Context, req interface{}) error
}

// Server is a validator middleware.
func Server(ss ...ConfigSetting) (middleware.Middleware, error) {
	cfg := settings.Apply(&Config{
		version:  V1,
		failFast: true,
	}, ss)
	v, err := buildValidator(cfg)
	if v == nil {
		return nil, err
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if err = v.Validate(ctx, req); err != nil {
				return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
			}
			return handler(ctx, req)
		}
	}, nil
}

func buildValidator(cfg *Config) (Validator, error) {
	switch cfg.version {
	case V1:
		return NewValidateV1(cfg.failFast, cfg.onValidationErrCallback), nil
	case V2:
		return NewValidateV2(protovalidate.WithFailFast(cfg.failFast))
	default:
		return nil, fmt.Errorf("unsupported version: %d", cfg.version)
	}
}
