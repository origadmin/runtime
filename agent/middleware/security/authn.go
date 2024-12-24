/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/security"
)

const (
	MetadataAuthZ = "x-md-global-security-authz"
	MetadataAuthN = "x-md-global-security-authn"
)

const (
	ErrorCreateOptionNil = errors.String("authenticator middleware create failed: option is nil")
)

// NewAuthNClient is a client authenticator middleware.
func NewAuthNClient(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	log.Debugf("NewAuthNClient: creating client authenticator middleware with config: %+v", cfg)
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authenticator == nil {
		log.Errorf("NewAuthNClient: authenticator is nil, returning error")
		return nil, ErrorCreateOptionNil
	}

	paths := append(cfg.GetPublicPaths(), cfg.GetAuthn().GetPublicPaths()...)
	log.Debugf("NewAuthNClient: public paths: %+v", paths)
	if option.Skipper == nil {
		log.Debugf("NewAuthNClient: skipper is nil, setting default skipper")
		option.Skipper = defaultSkipper(paths...)
	}
	//tokenParser := aggregateTokenParsers(FromTransportClient(security.HeaderAuthorize, security.SchemeBearer.String()))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Debugf("NewAuthNClient: handling request: %+v", req)
			if option.Skipper != nil {
				if tr, ok := transport.FromClientContext(ctx); ok {
					log.Debugf("NewAuthNClient: checking skipper for operation: %+v", tr.Operation())
					if option.Skipper(tr.Operation()) {
						log.Debugf("NewAuthNClient: skipping request")
						ctx := WithSkipContextClient(NewSkipContext(ctx), option.SkipKey)
						return handler(ctx, req)
					}
				}
			}

			log.Debugf("NewAuthNClient: authenticating context")
			claims, err := option.Authenticator.AuthenticateContext(ctx, security.ContextTypeHeader)
			if err != nil {
				log.Errorf("NewAuthNClient: authentication failed: %s", err.Error())
				return nil, err
			}

			log.Debugf("NewAuthNClient: creating token context")
			ctx, err = option.Authenticator.CreateTokenContext(ctx, security.ContextTypeMetadata, claims)
			if err != nil {
				log.Errorf("NewAuthNClient: creating token context failed: %s", err.Error())
				return nil, err
			}
			log.Debugf("NewAuthNClient: calling next handler")
			return handler(ctx, req)
		}
	}, nil
}

// NewAuthNServer is a server authenticator middleware.
func NewAuthNServer(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	log.Debugf("NewAuthNServer: creating server authenticator middleware with config: %+v", cfg)
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Authenticator == nil {
		log.Errorf("NewAuthNServer: authenticator is nil, returning error")
		return nil, ErrorCreateOptionNil
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Debugf("NewAuthNServer: handling request: %+v", req)
			if option.Skipper != nil {
				if tr, ok := transport.FromClientContext(ctx); ok {
					log.Debugf("NewAuthNServer: checking skipper for operation: %+v", tr.Operation())
					if option.Skipper(tr.Operation()) {
						log.Debugf("NewAuthNServer: skipping request")
						ctx := WithSkipContextClient(NewSkipContext(ctx), option.SkipKey)
						return handler(ctx, req)
					}
				} else {
					log.Debugf("NewAuthNServer: unable to get transport from client context")
				}
			} else {
				log.Debugf("NewAuthNServer: skipper is nil")
			}

			log.Debugf("NewAuthNServer: authenticating context")
			var err error
			claims, err := option.Authenticator.AuthenticateContext(ctx, security.ContextTypeMetadata)
			if err != nil {
				log.Errorf("NewAuthNServer: authentication failed: %s", err.Error())
				return nil, err
			}

			log.Debugf("NewAuthNServer: setting claims to context")
			ctx = NewClaimsContext(ctx, claims)
			log.Debugf("NewAuthNServer: calling next handler")
			return handler(ctx, req)
		}
	}, nil
}

// NewAuthN is a server authenticator middleware.
func NewAuthN(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, error) {
	log.Debugf("NewAuthN: creating server authenticator middleware with config: %+v", cfg)
	option := settings.ApplyDefaultsOrZero(ss...)
	if option == nil || option.Authenticator == nil {
		log.Errorf("NewAuthN: option or authenticator is nil, returning error")
		return nil, ErrorCreateOptionNil
	}

	log.Debugf("NewAuthN: applying defaults and creating token parser")
	//tokenParsers := aggregateTokenParsers(
	//	option.TokenParser,
	//	FromTransportClient(option.HeaderAuthorize, option.Scheme),
	//	FromTransportServer(option.HeaderAuthorize, option.Scheme))

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Debugf("NewAuthN: handling request: %+v", req)
			if IsSkipped(ctx, option.SkipKey) {
				log.Debugf("NewAuthN: skipping request due to skip key")
				return handler(ctx, req)
			}

			log.Debugf("NewAuthN: parsing token from context")
			var err error
			token := option.TokenParser(ctx)
			if token == "" {
				log.Errorf("NewAuthN: missing token, returning error")
				return nil, ErrMissingToken
			}

			log.Debugf("NewAuthN: authenticating token")
			claims, err := option.Authenticator.Authenticate(ctx, token)
			if err != nil {
				log.Errorf("NewAuthN: authentication failed: %s", err.Error())
				return nil, err
			}
			if option.IsRoot(ctx, claims) {
				ctx = security.WithRootContext(ctx)
			}

			log.Debugf("NewAuthN: setting claims to context")
			ctx = NewClaimsContext(ctx, claims)
			log.Debugf("NewAuthN: calling next handler")
			return handler(ctx, req)
		}
	}, nil
}
