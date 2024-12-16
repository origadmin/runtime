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
	"github.com/origadmin/runtime/log"
)

// NewAuthZClient returns a new server middleware.
func NewAuthZClient(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	log.Debugf("NewAuthZClient: creating new server middleware with config %+v and options %+v", cfg, ss)
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authorizer == nil {
		log.Errorf("NewAuthZClient: authorizer is nil")
		return nil, ErrorCreateOptionNil
	}
	log.Debugf("NewAuthZClient: authorizer is not nil, proceeding with middleware creation")
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Debugf("NewAuthZClient: handling request %+v with context %+v", req, ctx)
			if IsSkipped(ctx, option.SkipKey) {
				log.Debugf("NewAuthZClient: skipping request due to skip key")
				return handler(NewSkipContext(ctx), req)
			}

			claims := ClaimsFromContext(ctx)
			if claims == nil {
				log.Errorf("NewAuthZClient: claims are nil")
				return nil, ErrMissingToken
			}
			log.Debugf("NewAuthZClient: claims are not nil, proceeding with user claims parsing")
			userClaims := option.ParserUserClaims(ctx, claims)

			if userClaims.GetSubject() == "" || userClaims.GetAction() == "" || userClaims.GetObject() == "" {
				log.Errorf("NewAuthZClient: invalid user claims")
				return nil, ErrInvalidClaims
			}

			log.Debugf("NewAuthZClient: user claims are valid, proceeding with authorization")
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

			log.Debugf("NewAuthZClient: returning handler with context %+v and request %+v", ctx, req)
			return handler(ctx, req)
		}
	}, nil
}

// NewAuthZServer returns a new server middleware.
func NewAuthZServer(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	log.Debugf("NewAuthZServer: creating new server middleware with config %+v and options %+v", cfg, ss)
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authorizer == nil {
		log.Errorf("NewAuthZServer: authorizer is nil")
		return nil, ErrorCreateOptionNil
	}
	log.Debugf("NewAuthZServer: authorizer is not nil, proceeding with middleware creation")
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Debugf("NewAuthZServer: handling request %+v with context %+v", req, ctx)
			if IsSkipped(ctx, option.SkipKey) {
				log.Debugf("NewAuthZServer: skipping request due to skip key")
				return handler(NewSkipContext(ctx), req)
			}
			var (
				allowed bool
				err     error
			)

			claims := UserClaimsFromContext(ctx)
			if claims == nil {
				log.Errorf("NewAuthZServer: claims are nil")
				return nil, ErrMissingToken
			}
			log.Debugf("NewAuthZServer: claims are not nil, proceeding with authorization")
			if claims.GetSubject() == "" || claims.GetAction() == "" || claims.GetObject() == "" {
				log.Errorf("NewAuthZServer: invalid claims")
				return nil, ErrInvalidClaims
			}

			log.Debugf("NewAuthZServer: claims are valid, proceeding with authorization")
			//var project []string
			//if domains := claims.GetDomain(); domains != nil {
			//	project = domains
			//}
			// todo add domain project

			allowed, err = option.Authorizer.Authorized(ctx, claims)
			if err != nil {
				log.Errorf("NewAuthZServer: authorization error %+v", err)
				return nil, err
			}
			if !allowed {
				log.Errorf("NewAuthZServer: authorization denied")
				return nil, ErrInvalidAuth
			}

			log.Debugf("NewAuthZServer: returning handler with context %+v and request %+v", ctx, req)
			return handler(ctx, req)
		}
	}, nil
}

// NewAuthZ returns a new server middleware.
func NewAuthZ(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	log.Debugf("NewAuthZ: creating new server middleware with config %+v and options %+v", cfg, ss)
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authorizer == nil {
		log.Errorf("NewAuthZ: authorizer is nil")
		return nil, ErrorCreateOptionNil
	}
	log.Debugf("NewAuthZ: authorizer is not nil, proceeding with middleware creation")
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Debugf("NewAuthZ: handling request %+v with context %+v", req, ctx)
			if IsSkipped(ctx, option.SkipKey) {
				log.Debugf("NewAuthZ: skipping request due to skip key")
				return handler(ctx, req)
			}

			var (
				allowed bool
				err     error
			)

			claims := ClaimsFromContext(ctx)
			if claims == nil {
				log.Errorf("NewAuthZ: claims are nil")
				return nil, ErrMissingToken
			}
			log.Debugf("NewAuthZ: claims are not nil, subject: %s", claims.GetSubject())
			if option.Parser == nil {
				log.Errorf("NewAuthZ: parser is nil")
				return nil, ErrMissingClaims
			}
			userClaims, err := option.Parser(ctx, claims.GetSubject())
			if err != nil {
				log.Errorf("NewAuthZ: error parsing user claims: %v", err)
				return nil, err
			}
			log.Debugf("NewAuthZ: user claims: %+v", userClaims)

			if userClaims.GetSubject() == "" || userClaims.GetAction() == "" || userClaims.GetObject() == "" {
				log.Errorf("NewAuthZ: invalid user claims")
				return nil, ErrInvalidClaims
			}

			//var project []string
			//if domains := claims.GetDomain(); domains != nil {
			//	project = domains
			//}
			// todo add domain project
			log.Debugf("NewAuthZ: checking authorization for user claims %+v", userClaims)
			allowed, err = option.Authorizer.Authorized(ctx, userClaims)
			if err != nil {
				log.Errorf("NewAuthZ: authorization error: %v", err)
				return nil, err
			}
			if !allowed {
				log.Errorf("NewAuthZ: authorization denied")
				return nil, ErrInvalidAuth
			}

			log.Debugf("NewAuthZ: authorization successful, proceeding with request")
			return handler(ctx, req)
		}
	}, nil
}
