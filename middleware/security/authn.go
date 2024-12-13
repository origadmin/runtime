/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/security"
)

const (
	MetadataAuthZ            = "x-metadata-security-authz"
	MetadataAuthN            = "x-metadata-security-authn"
	MetadataSecurityTokenKey = "x-metadata-security-token-key"
	MetadataSecuritySkipKey  = "x-metadata-security-skip-key"
)

const (
	ErrorCreateOptionNil = errors.String("authenticator middleware create failed: option is nil")
)

// NewAuthNClient is a client authenticator middleware.
func NewAuthNClient(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authenticator == nil {
		return nil, ErrorCreateOptionNil
	}

	paths := append(cfg.GetPublicPaths(), cfg.GetAuthn().GetPublicPaths()...)
	if option.Skipper == nil {
		option.Skipper = defaultSkipper(paths...)
	}
	tokenParser := defaultTokenParser(FromTransportClient(security.HeaderAuthorize, security.SchemeBearer.String()))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if option.Skipper != nil {
				if tr, ok := transport.FromClientContext(ctx); ok {
					if option.Skipper(tr.Operation()) {
						ctx := WithSkipContextClient(NewSkipContext(ctx), option.SkipKey)
						return handler(ctx, req)
					}
				}
			}

			token := tokenParser(ctx)
			if token == "" {
				return nil, ErrMissingToken
			}
			ctx = metadata.AppendToClientContext(ctx, option.TokenKey, token)
			//claims, err := option.Authenticator.AuthenticateToken(token)
			//if err != nil {
			//	log.Errorf("authenticator middleware create token failed: %s", err.Error())
			//}

			//ctx = NewClaimsContext(ctx, claims)
			return handler(ctx, req)
		}
	}, nil
}

// NewAuthNServer is a server authenticator middleware.
func NewAuthNServer(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authenticator == nil {
		return nil, ErrorCreateOptionNil
	}

	tokenParser := defaultTokenParser(FromMetaData(option.TokenKey), FromTransportServer(security.HeaderAuthorize, security.SchemeBearer.String()))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if option.Skipper != nil {
				if tr, ok := transport.FromClientContext(ctx); ok {
					if option.Skipper(tr.Operation()) {
						ctx := WithSkipContextClient(NewSkipContext(ctx), option.SkipKey)
						return handler(ctx, req)
					}
				}
			}

			token := tokenParser(ctx)
			if token == "" {
				return nil, ErrMissingToken
			}
			var err error
			claims, err := option.Authenticator.AuthenticateToken(ctx, token)
			if err != nil {
				return nil, err
			}

			// set claims to context, so that the next middleware can get it
			ctx = NewClaimsContext(ctx, claims)
			return handler(ctx, req)
		}
	}, nil
}

// NewAuthN is a server authenticator middleware.
func NewAuthN(option *Option) (middleware.Middleware, error) {
	if option == nil || option.Authenticator == nil {
		return nil, ErrorCreateOptionNil
	}
	if option.TokenKey == "" {
		option.TokenKey = MetadataSecurityTokenKey
	}
	if option.SkipKey == "" {
		option.SkipKey = MetadataSecuritySkipKey
	}

	tokenParser := defaultTokenParser(FromMetaData(option.TokenKey), FromTransportServer(security.HeaderAuthorize, security.SchemeBearer.String()))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if IsSkipped(ctx, option.SkipKey) {
				return handler(NewSkipContext(ctx), req)
			}
			var err error
			token := tokenParser(ctx)
			if token == "" {
				return nil, ErrMissingToken
			}

			claims, err := option.Authenticator.AuthenticateToken(ctx, token)
			if err != nil {
				return nil, err
			}

			// set claims to context, so that the next middleware can get it
			ctx = NewClaimsContext(ctx, claims)
			return handler(ctx, req)
		}
	}, nil
}
