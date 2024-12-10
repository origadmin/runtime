/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
)

// NewAuthZServer returns a new server middleware.
func NewAuthZServer(ss ...ConfigOptionSetting) middleware.Middleware {
	option := settings.Apply(&ConfigOption{}, ss)
	if option == nil || option.Authorizer == nil {
		return nil
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if IsSkipped(ctx, option.SecuritySkipKey) {
				return handler(NewSkipContext(ctx), req)
			}
			var (
				allowed bool
				err     error
			)

			policy := PolicyFromContext(ctx)
			if policy == nil {
				return nil, ErrMissingToken
			}

			if policy.GetSubject() == "" || policy.GetAction() == "" || policy.GetObject() == "" {
				return nil, ErrInvalidClaims
			}

			//var project []string
			//if domains := policy.GetDomain(); domains != nil {
			//	project = domains
			//}
			// todo add domain project

			allowed, err = option.Authorizer.Authorized(ctx, policy)
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, ErrInvalidAuth
			}

			return handler(ctx, req)
		}
	}
}
