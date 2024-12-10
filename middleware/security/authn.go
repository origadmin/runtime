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
)

const (
	MetadataAuthZ            = "x-metadata-security-authz"
	MetadataAuthN            = "x-metadata-security-authn"
	MetadataSecurityTokenKey = "x-metadata-security-token-key"
	MetadataSecuritySkipKey  = "x-metadata-security-skip-key"
)

// NewAuthNClient is a client authenticator middleware.
func NewAuthNClient(ss ...ConfigOptionSetting) middleware.Middleware {
	option := settings.Apply(&ConfigOption{}, ss)
	if option == nil || option.Authorizer == nil {
		return nil
	}
	if option.SecurityTokenKey == "" {
		option.SecurityTokenKey = MetadataSecurityTokenKey
	}
	if option.SecuritySkipKey == "" {
		option.SecuritySkipKey = MetadataSecuritySkipKey
	}

	if option.Skipper == nil {
		option.Skipper = defaultSkipper(option.PublicPaths...)
	}
	tokenParser := defaultTokenParser(FromTransportClient("Authorization", "Bearer"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if option.Skipper != nil {
				if tr, ok := transport.FromClientContext(ctx); ok {
					if option.Skipper(tr.Operation()) {
						ctx := WithSkipContextClient(NewSkipContext(ctx), option.SecuritySkipKey)
						return handler(ctx, req)
					}
				}
			}

			token := tokenParser(ctx)
			if token == "" {
				return nil, ErrMissingToken
			}
			ctx = metadata.AppendToClientContext(ctx, option.SecurityTokenKey, token)
			//claims, err := option.Authenticator.AuthenticateToken(token)
			//if err != nil {
			//	log.Errorf("authenticator middleware create token failed: %s", err.Error())
			//}

			//ctx = NewClaimsContext(ctx, claims)
			return handler(ctx, req)
		}
	}
}

// NewAuthNServer is a server authenticator middleware.
func NewAuthNServer(ss ...ConfigOptionSetting) middleware.Middleware {
	option := settings.Apply(&ConfigOption{}, ss)
	if option == nil || option.Authenticator == nil {
		return nil
	}
	if option.SecurityTokenKey == "" {
		option.SecurityTokenKey = MetadataSecurityTokenKey
	}
	if option.SecuritySkipKey == "" {
		option.SecuritySkipKey = MetadataSecuritySkipKey
	}

	tokenParser := defaultTokenParser(FromMetaData(option.SecurityTokenKey), FromTransportServer("Authorization", "Bearer"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if IsSkipped(ctx, option.SecuritySkipKey) {
				return handler(NewSkipContext(ctx), req)
			}
			var err error
			token := tokenParser(ctx)
			if token == "" {
				return nil, ErrMissingToken
			}

			claims, err := option.Authenticator.AuthenticateToken(token)
			if err != nil {
				return nil, err
			}

			// set claims to context, so that the next middleware can get it
			ctx = NewClaimsContext(ctx, claims)
			return handler(ctx, req)
		}
	}
}
