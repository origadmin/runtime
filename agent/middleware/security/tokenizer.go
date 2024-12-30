/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/security"
)

type tokenizer struct {
	Scheme        security.Scheme
	CacheStorage  security.CacheStorage
	Authenticator security.Authenticator
}

func (obj *tokenizer) DestroyToken(ctx context.Context, tokenStr string) error {
	if obj.CacheStorage == nil {
		return nil
	}
	return obj.CacheStorage.Remove(ctx, tokenStr)
}

func (obj *tokenizer) DestroyTokenContext(ctx context.Context, tokenType security.TokenType) error {
	tokenStr, err := TokenFromTypeContext(ctx, tokenType, obj.schemeString())
	if err != nil || tokenStr == "" {
		return ErrInvalidToken
	}
	// Destroy the token string.
	return obj.DestroyToken(ctx, tokenStr)
}

func (obj *tokenizer) Close(ctx context.Context) error {
	if obj.CacheStorage == nil {
		return nil
	}
	return obj.CacheStorage.Close(ctx)
}

// CreateToken creates a token string from the claims.
func (obj *tokenizer) CreateToken(ctx context.Context, claims security.Claims) (string, error) {
	return obj.Authenticator.CreateToken(ctx, claims)
}

// CreateTokenContext creates a token string from the claims and adds it to the context.
func (obj *tokenizer) CreateTokenContext(ctx context.Context, tokenType security.TokenType, claims security.Claims) (context.Context, error) {
	// Create the token string.
	tokenStr, err := obj.Authenticator.CreateToken(ctx, claims)
	if err != nil {
		return ctx, err
	}
	// Add the token string to the context.
	ctx = TokenToTypeContext(ctx, tokenType, obj.schemeString(), tokenStr)
	return ctx, nil
}

func (obj *tokenizer) CreateIdentityClaims(ctx context.Context, id string, refresh bool) (security.Claims, error) {
	claims, err := obj.Authenticator.CreateIdentityClaims(ctx, id, refresh)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (obj *tokenizer) CreateIdentityClaimsContext(ctx context.Context, _ security.TokenType, id string) (context.Context, error) {
	// Create the claims.
	claims, err := obj.CreateIdentityClaims(ctx, id, false)
	if err != nil {
		return ctx, err
	}

	// Add the token to the context.
	ctx = security.NewClaimsContext(ctx, claims)
	return ctx, nil
}

func (obj *tokenizer) Authenticate(ctx context.Context, tokenStr string) (security.Claims, error) {
	log.Debugf("Authenticating token string: %s", tokenStr)
	// If the token cache service is not nil, validate the token.
	if obj.CacheStorage != nil {
		log.Debugf("Validating token using cache service")
		ok, err := obj.CacheStorage.Exist(ctx, tokenStr)
		switch {
		case err != nil:
			log.Errorf("Error validating token: %v", err)
			return nil, ErrInvalidToken
		case !ok:
			log.Debugf("Token not found in cache")
			return nil, ErrTokenNotFound
		}
	}

	return obj.Authenticator.Authenticate(ctx, tokenStr)
}

func (obj *tokenizer) AuthenticateContext(ctx context.Context, tokenType security.TokenType) (security.Claims, error) {
	log.Debugf("Entering AuthenticateContext with tokenType: %s", tokenType)
	// Get the token string from the context.
	tokenStr, err := TokenFromTypeContext(ctx, tokenType, obj.schemeString())
	if err != nil {
		log.Errorf("Error getting token from context: %v", err)
	} else if tokenStr == "" {
		log.Debugf("Token string is empty")
	}
	if err != nil || tokenStr == "" {
		log.Errorf("Invalid token or token string is empty")
		return nil, ErrInvalidToken
	}
	log.Debugf("Token string retrieved from context: %s", tokenStr)
	// Authenticate the token string.
	log.Debugf("Authenticating token string")
	claims, err := obj.Authenticate(ctx, tokenStr)
	if err != nil {
		log.Errorf("Error authenticating token: %v", err)
	}
	log.Debugf("Authentication result: %+v", claims)
	return claims, err
}

func (obj *tokenizer) Verify(ctx context.Context, tokenStr string) (bool, error) {
	return obj.Authenticator.Verify(ctx, tokenStr)
}

func (obj *tokenizer) VerifyContext(ctx context.Context, tokenType security.TokenType) (bool, error) {
	// Get the token string from the context.
	tokenStr, err := TokenFromTypeContext(ctx, tokenType, obj.schemeString())
	if err != nil || tokenStr == "" {
		return false, ErrInvalidToken
	}
	// Authenticate the token string.
	return obj.Verify(ctx, tokenStr)
}

func (obj *tokenizer) schemeString() string {
	return obj.Scheme.String()
}

func NewTokenizer(authenticator security.Authenticator, ss ...TokenizerSetting) security.Tokenizer {
	t := settings.Apply(&tokenizer{
		Authenticator: authenticator,
		CacheStorage:  security.NewCacheStorage(),
		Scheme:        security.SchemeBearer,
	}, ss)
	return t
}

var _ security.Tokenizer = (*tokenizer)(nil)
