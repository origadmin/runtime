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

// NewAuthZClient returns a new server middleware.
func NewAuthZClient(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authorizer == nil {
		return nil, ErrorCreateOptionNil
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if IsSkipped(ctx, option.SkipKey) {
				return handler(NewSkipContext(ctx), req)
			}

			claims := ClaimsFromContext(ctx)
			if claims == nil {
				return nil, ErrMissingToken
			}
			userClaims := option.ParserUserClaims(ctx, claims)

			if userClaims.GetSubject() == "" || userClaims.GetAction() == "" || userClaims.GetObject() == "" {
				return nil, ErrInvalidClaims
			}

			//var project []string
			//if domains := claims.GetDomain(); domains != nil {
			//	project = domains
			//}
			// todo add domain project

			//allowed, err = option.Authorizer.Authorized(ctx, userClaims)
			//if err != nil {
			//	return nil, err
			//}
			//if !allowed {
			//	return nil, ErrInvalidAuth
			//}

			return handler(ctx, req)
		}
	}, nil
}

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

// NewAuthZ returns a new server middleware.
func NewAuthZ(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authorizer == nil {
		return nil, ErrorCreateOptionNil
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if IsSkipped(ctx, option.SkipKey) {
				return handler(ctx, req)
			}
			var (
				allowed bool
				err     error
			)

			if option.Parser == nil {
				return nil, ErrMissingClaims
			}
			claims := ClaimsFromContext(ctx)
			if claims == nil {
				return nil, ErrMissingToken
			}
			userClaims, err := option.Parser.Parse(ctx, claims.GetSubject())
			if claims == nil {
				return nil, ErrMissingToken
			}

			if userClaims.GetSubject() == "" || userClaims.GetAction() == "" || userClaims.GetObject() == "" {
				return nil, ErrInvalidClaims
			}

			//var project []string
			//if domains := claims.GetDomain(); domains != nil {
			//	project = domains
			//}
			// todo add domain project
			allowed, err = option.Authorizer.Authorized(ctx, userClaims)
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
