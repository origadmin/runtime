/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package jwt implements the functions, types, and interfaces for the module.
package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/goexts/generic/settings"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/origadmin/toolkits/security"
	"github.com/origadmin/toolkits/storage/cache"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware/security/internal/helper"
)

type Setting = func(*Authenticator)

type Authenticator struct {
	signingMethod jwtv5.SigningMethod
	keyFunc       func(*jwtv5.Token) (any, error)
	schemeType    security.Scheme
	cache         security.TokenCacheService
	extraKeys     []string
}

func (obj *Authenticator) schemeString() string {
	return obj.schemeType.String()
}

func (obj *Authenticator) AuthenticateToken(ctx context.Context, tokenStr string) (security.Claims, error) {
	if obj.cache != nil {
		ok, err := obj.cache.Validate(ctx, tokenStr)
		switch {
		case err != nil:
			return nil, ErrInvalidToken
		case !ok:
			return nil, ErrTokenNotFound
		}
	}
	jwtToken, err := obj.parseToken(tokenStr)

	if jwtToken == nil {
		return nil, ErrInvalidToken
	}

	if err != nil {
		switch {
		case errors.Is(err, jwtv5.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		case errors.Is(err, jwtv5.ErrTokenSignatureInvalid):
			return nil, ErrTokenSignatureInvalid
		case errors.Is(err, jwtv5.ErrTokenExpired) || errors.Is(err, jwtv5.ErrTokenNotValidYet):
			return nil, ErrTokenExpired
		default:
			return nil, ErrInvalidToken
		}
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	if jwtToken.Method != obj.signingMethod {
		return nil, ErrUnsupportedSigningMethod
	}

	if jwtToken.Claims == nil {
		return nil, ErrInvalidClaims
	}

	securityClaims, err := ToClaims(jwtToken.Claims, obj.extraKeys...)
	if err != nil {
		return nil, err
	}
	return securityClaims, nil
}

func (obj *Authenticator) AuthenticateTokenContext(ctx context.Context, tokenType security.TokenType) (security.Claims, error) {
	tokenStr, err := helper.FromTokenTypeContext(ctx, tokenType, obj.schemeString())
	if err != nil || tokenStr == "" {
		return nil, ErrInvalidToken
	}
	return obj.AuthenticateToken(ctx, tokenStr)
}

func (obj *Authenticator) Authenticate(ctx context.Context, tokenStr string) (bool, error) {
	_, err := obj.AuthenticateToken(ctx, tokenStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (obj *Authenticator) AuthenticateContext(ctx context.Context, tokenType security.TokenType) (bool, error) {
	tokenStr, err := helper.FromTokenTypeContext(ctx, tokenType, obj.schemeString())
	if err != nil || tokenStr == "" {
		return false, ErrInvalidToken
	}
	return obj.Authenticate(ctx, tokenStr)
}

func (obj *Authenticator) CreateToken(ctx context.Context, claims security.Claims) (string, error) {
	jwtToken := jwtv5.NewWithClaims(obj.signingMethod, ClaimsToJwtClaims(claims))

	tokenStr, err := obj.generateToken(jwtToken)
	if err != nil || tokenStr == "" {
		return "", err
	}
	exp := time.Duration(claims.GetExpiration().UnixMilli())
	if obj.cache != nil {
		if err := obj.cache.Store(ctx, tokenStr, exp); err != nil {
			return tokenStr, err
		}
	}
	return tokenStr, nil
}

func (obj *Authenticator) CreateTokenContext(ctx context.Context, tokenType security.TokenType, claims security.Claims) (context.Context, error) {
	tokenStr, err := obj.CreateToken(ctx, claims)
	if err != nil {
		return ctx, err
	}
	ctx = helper.WithTokenTypeContext(ctx, tokenType, obj.schemeString(), tokenStr)
	return ctx, nil
}

func (obj *Authenticator) DestroyToken(ctx context.Context, tokenStr string) error {
	if obj.cache != nil {
		err := obj.cache.Remove(ctx, tokenStr)
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return err
		}
	}
	return nil
}

func (obj *Authenticator) DestroyTokenContext(ctx context.Context, token security.TokenType) error {
	tokenStr, err := helper.FromTokenTypeContext(ctx, token, obj.schemeString())
	if err != nil || tokenStr == "" {
		return ErrInvalidToken
	}
	return obj.DestroyToken(ctx, tokenStr)
}

func NewAuthenticator(cfg *configv1.Security, ss ...Setting) (security.Authenticator, error) {
	config := cfg.GetAuthn().GetJwt()
	if config == nil {
		return nil, errors.New("cfg config is empty")
	}
	auth := settings.Apply(&Authenticator{}, ss)
	if auth.signingMethod == nil {
		auth.signingMethod = GetSigningMethodFromAlg(config.Algorithm)
	}
	if config.SigningKey == "" {
		return nil, errors.New("signing key is empty")
	}
	if auth.keyFunc == nil {
		auth.keyFunc = GetKeyFuncWithAlg(config.Algorithm, config.SigningKey)
	}
	return auth, nil
}

func GetKeyFunc(key string) func(token *jwtv5.Token) (any, error) {
	return func(token *jwtv5.Token) (any, error) {
		if token.Method.Alg() == "" {
			return nil, ErrInvalidToken
		}
		return key, nil
	}
}

func GetKeyFuncWithAlg(alg, key string) func(token *jwtv5.Token) (any, error) {
	return func(token *jwtv5.Token) (any, error) {
		if token.Method.Alg() == "" || alg != token.Method.Alg() {
			return nil, ErrInvalidToken
		}
		// jwtv5 must be []byte
		return []byte(key), nil
	}
}

func GetSigningMethodFromAlg(algorithm string) jwtv5.SigningMethod {
	switch algorithm {
	case "HS256":
		return jwtv5.SigningMethodHS256
	case "HS384":
		return jwtv5.SigningMethodHS384
	case "HS512":
		return jwtv5.SigningMethodHS512
	case "RS256":
		return jwtv5.SigningMethodRS256
	case "RS384":
		return jwtv5.SigningMethodRS384
	case "RS512":
		return jwtv5.SigningMethodRS512
	case "ES256":
		return jwtv5.SigningMethodES256
	case "ES384":
		return jwtv5.SigningMethodES384
	case "ES512":
		return jwtv5.SigningMethodES512
	case "EdDSA":
		return jwtv5.SigningMethodEdDSA
	default:
		return nil
	}
}

func (obj *Authenticator) Close(ctx context.Context) error {
	if obj.cache != nil {
		return obj.cache.Close(ctx)
	}
	return nil
}

// parseToken parses the token string and returns the token.
func (obj *Authenticator) parseToken(token string) (*jwtv5.Token, error) {
	if obj.keyFunc == nil {
		return nil, ErrMissingKeyFunc
	}
	if obj.extraKeys == nil {
		return jwtv5.ParseWithClaims(token, &jwtv5.RegisteredClaims{}, obj.keyFunc)
	}

	return jwtv5.Parse(token, obj.keyFunc)
}

// generateToken generates a signed token string from the token.
func (obj *Authenticator) generateToken(jwtToken *jwtv5.Token) (string, error) {
	if obj.keyFunc == nil {
		return "", ErrMissingKeyFunc
	}

	key, err := obj.keyFunc(jwtToken)
	if err != nil {
		return "", ErrGetKeyFailed
	}

	strToken, err := jwtToken.SignedString(key)
	if err != nil {
		return "", ErrSignTokenFailed
	}

	return strToken, nil
}

var _ security.Authenticator = (*Authenticator)(nil)
