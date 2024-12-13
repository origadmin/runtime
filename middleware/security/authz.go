/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// NewAuthZServer returns a new server middleware.
func NewAuthZServer(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authorizer == nil {
		return nil, ErrorCreateOptionNil
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if IsSkipped(ctx, option.SkipKey) {
				return handler(NewSkipContext(ctx), req)
			}
			var (
				allowed bool
				err     error
			)

			claims := UserClaimsFromContext(ctx)
			if claims == nil {
				return nil, ErrMissingToken
			}

			if claims.GetSubject() == "" || claims.GetAction() == "" || claims.GetObject() == "" {
				return nil, ErrInvalidClaims
			}

			//var project []string
			//if domains := claims.GetDomain(); domains != nil {
			//	project = domains
			//}
			// todo add domain project

			allowed, err = option.Authorizer.Authorized(ctx, claims)
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, ErrInvalidAuth
			}

			return handler(ctx, req)
		}
	}, nil
}
