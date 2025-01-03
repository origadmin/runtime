/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/security"
)

type BridgeSetting = func(*Bridge)

type Data interface {
	QueryRoles(ctx context.Context, subject string) ([]string, error)
	QueryPermissions(ctx context.Context, subject string) ([]string, error)
}

type Bridge struct {
	// TokenSource is the source of the token.
	TokenSource security.TokenSource
	// Scheme is the scheme used for the authorization header.
	Scheme security.Scheme
	// AuthenticationHeader is the header used for the authorization header.
	AuthenticationHeader string
	// Authenticator is the authenticator used for the authorization header.
	Authenticator security.Authenticator
	// Authorizer is the authorizer used for the authorization header.
	Authorizer security.Authorizer
	// SkipKey is the key used to skip authentication.
	SkipKey string
	// PublicPaths are the public paths that do not require authentication.
	PublicPaths []string
	// Skipper is the function used to skip authentication.
	Skipper func(string) bool
	// IsRoot is the function used to check if the request is root.
	IsRoot func(ctx context.Context, claims security.Claims) bool
	// Data is the permission data from the database.
	Data Data
}

func (obj Bridge) SkipFromContext(ctx context.Context) (context.Context, bool) {
	if IsSkipped(ctx, obj.SkipKey) {
		log.Debugf("NewAuthN: skipping request due to skip key")
		return ctx, true
	}
	if obj.Skipper == nil {
		log.Debugf("NewAuthNServer: skipper is nil")
		return ctx, false
	}
	if tr, ok := transport.FromClientContext(ctx); ok {
		log.Debugf("NewAuthNServer: checking skipper for operation: %+v", tr.Operation())
		if obj.Skipper(tr.Operation()) {
			log.Debugf("NewAuthNServer: skipping request")
			ctx := WithSkipContextClient(NewSkipContext(ctx), obj.SkipKey)
			return ctx, true
		}
	} else {
		log.Debugf("NewAuthNServer: unable to get transport from client context")
	}
	return ctx, false
}

func (obj Bridge) schemeString() string {
	return obj.Scheme.String()
}

func (obj Bridge) TokenParser(ctx context.Context) string {
	return obj.TokenParser(ctx)
}

func (obj Bridge) Build() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			log.Debugf("NewAuthN: handling request: %+v", req)
			if ctx, ok := obj.SkipFromContext(ctx); ok {
				return handler(ctx, req)
			}
			log.Debugf("NewAuthN: parsing token from context")
			token := obj.TokenParser(ctx)
			if token == "" {
				log.Errorf("NewAuthN: missing token, returning error")
				return nil, ErrMissingToken
			}

			log.Debugf("NewAuthN: authenticating token")
			claims, err := obj.Authenticator.Authenticate(ctx, token)
			if err != nil {
				log.Errorf("NewAuthN: authentication failed: %s", err.Error())
				return nil, err
			}
			if obj.IsRoot(ctx, claims) {
				ctx = security.WithRootContext(ctx)
			}

			policy, err := obj.PolicyParser(ctx, claims)
			if err != nil {
				return nil, err
			}
			if ok, err := obj.Authorizer.Authorized(ctx, policy, policy.GetAction(), policy.GetObject()); err != nil {
				log.Errorf("NewAuthN: authorization failed")
				return nil, ErrInvalidAuthorization
			} else if !ok {
				log.Errorf("NewAuthN: authorization check failed")
				return nil, ErrInvalidAuthorization
			} else {
				log.Debugf("NewAuthN: authorization successful, proceeding with request")
			}

			log.Debugf("NewAuthN: setting claims to context")
			ctx = obj.WithContext(ctx, token)
			log.Debugf("NewAuthN: calling next handler")
			return handler(ctx, req)
		}
	}
}

func (obj Bridge) WithContext(ctx context.Context, token string) context.Context {
	return TokenToContext(ctx, obj.TokenSource, obj.schemeString(), token)
}

func (obj Bridge) PolicyParser(ctx context.Context, claims security.Claims) (security.Policy, error) {
	roles, err := obj.Data.QueryRoles(ctx, claims.GetSubject())
	if err != nil {
		return nil, err
	}
	permissions, err := obj.Data.QueryPermissions(ctx, claims.GetSubject())
	if err != nil {
		return nil, err
	}
	policy := security.RegisteredPolicy{
		Subject: claims.GetSubject(),
		//Object:     claims.GetObject(),
		//Action:     claims.GetAction(),
		Domain:     claims.GetIssuer(),
		Roles:      roles,
		Permission: permissions,
	}
	return &policy, nil
}
func BridgeMiddleware(authenticator security.Authenticator, authorizer security.Authorizer, bss ...BridgeSetting) middleware.Middleware {
	bridge := &Bridge{
		Authenticator: authenticator,
		Authorizer:    authorizer,
	}
	bridge = settings.Apply(bridge, bss)
	return bridge.Build()
}
