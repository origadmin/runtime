/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
)

type Validator interface {
	Validate(ctx context.Context, req interface{}) error
}

// Server is a validator middleware.
func Server(ss ...Option) (middleware.Middleware, error) {
	cfg := settings.Apply(&Options{
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

func buildValidator(cfg *Options) (Validator, error) {
	switch cfg.version {
	case V1:
		return NewValidateV1(cfg.failFast, cfg.callback), nil
	case V2:
		if cfg.failFast {
			cfg.validatorOptions = append(cfg.validatorOptions, protovalidate.WithFailFast())
		}
		return NewValidateV2(cfg.validatorOptions...)
	default:
		return nil, fmt.Errorf("unsupported version: %d", cfg.version)
	}
}
