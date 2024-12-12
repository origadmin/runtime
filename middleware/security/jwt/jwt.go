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

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware/security/internal/helper"
	"github.com/origadmin/toolkits/security"
	"github.com/origadmin/toolkits/storage/cache"
)

// Setting is a function type for setting the Authenticator.
type Setting = func(*Authenticator)

// Authenticator is a struct that implements the Authenticator interface.
type Authenticator struct {
	// signingMethod is the signing method for the token.
	signingMethod jwtv5.SigningMethod
	// keyFunc is a function that returns the key for signing the token.
	keyFunc func(*jwtv5.Token) (any, error)
	// schemeType is the scheme type for the token.
	schemeType security.Scheme
	// cache is the token cache service.
	cache security.TokenCacheService
	// extraKeys are the extra keys for the token.
	extraKeys []string
	// scoped enabled for the token.
	scoped bool
	//// user parser is the user parser for the token. It is optional.    If it is not set, the user will not be parsed.
	//userParser security.UserClaimsParser
}

// schemeString returns the scheme type as a string.
func (obj *Authenticator) schemeString() string {
	return obj.schemeType.String()
}

// AuthenticateToken authenticates the token string.
func (obj *Authenticator) AuthenticateToken(ctx context.Context, tokenStr string) (security.Claims, error) {
	// If the token cache service is not nil, validate the token.
	if obj.cache != nil {
		ok, err := obj.cache.Validate(ctx, tokenStr)
		switch {
		case err != nil:
			return nil, ErrInvalidToken
		case !ok:
			return nil, ErrTokenNotFound
		}
	}
	// Parse the token string.
	jwtToken, err := obj.parseToken(tokenStr)

	// If the token is nil, return an error.
	if jwtToken == nil {
		return nil, ErrInvalidToken
	}

	// If there is an error, return the appropriate error.
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

	// If the token is not valid, return an error.
	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	// If the signing method is not supported, return an error.
	if jwtToken.Method != obj.signingMethod {
		return nil, ErrUnsupportedSigningMethod
	}

	// If the claims are nil, return an error.
	if jwtToken.Claims == nil {
		return nil, ErrInvalidClaims
	}

	// Convert the claims to security.Claims.
	securityClaims, err := ToClaims(jwtToken.Claims, obj.extraKeys...)
	if err != nil {
		return nil, err
	}
	return securityClaims, nil
}

// AuthenticateTokenContext authenticates the token string from the context.
func (obj *Authenticator) AuthenticateTokenContext(ctx context.Context, tokenType security.TokenType) (security.Claims, error) {
	// Get the token string from the context.
	tokenStr, err := helper.FromTokenTypeContext(ctx, tokenType, obj.schemeString())
	if err != nil || tokenStr == "" {
		return nil, ErrInvalidToken
	}
	// Authenticate the token string.
	return obj.AuthenticateToken(ctx, tokenStr)
}

// Authenticate authenticates the token string.
func (obj *Authenticator) Authenticate(ctx context.Context, tokenStr string) (bool, error) {
	// Authenticate the token string.
	_, err := obj.AuthenticateToken(ctx, tokenStr)
	// If there is an error, return false and the error.
	if err != nil {
		return false, err
	}
	// Otherwise, return true.
	return true, nil
}

// AuthenticateContext authenticates the token string from the context.
func (obj *Authenticator) AuthenticateContext(ctx context.Context, tokenType security.TokenType) (bool, error) {
	// Get the token string from the context.
	tokenStr, err := helper.FromTokenTypeContext(ctx, tokenType, obj.schemeString())
	if err != nil || tokenStr == "" {
		return false, ErrInvalidToken
	}
	// Authenticate the token string.
	return obj.Authenticate(ctx, tokenStr)
}

// CreateToken creates a token string from the claims.
func (obj *Authenticator) CreateToken(ctx context.Context, claims security.Claims) (string, error) {
	// Create a new token with the claims.
	jwtToken := jwtv5.NewWithClaims(obj.signingMethod, ClaimsToJwtClaims(claims))

	// Generate the token string.
	tokenStr, err := obj.generateToken(jwtToken)
	if err != nil || tokenStr == "" {
		return "", err
	}
	// Get the expiration time from the claims.
	exp := time.Duration(claims.GetExpiration().UnixMilli())
	// If the token cache service is not nil, store the token.
	if obj.cache != nil {
		if err := obj.cache.Store(ctx, tokenStr, exp); err != nil {
			return tokenStr, err
		}
	}
	return tokenStr, nil
}

// CreateTokenContext creates a token string from the claims and adds it to the context.
func (obj *Authenticator) CreateTokenContext(ctx context.Context, tokenType security.TokenType, claims security.Claims) (context.Context, error) {
	// Create the token string.
	tokenStr, err := obj.CreateToken(ctx, claims)
	if err != nil {
		return ctx, err
	}
	// Add the token string to the context.
	ctx = helper.WithTokenTypeContext(ctx, tokenType, obj.schemeString(), tokenStr)
	return ctx, nil
}

// DestroyToken destroys the token string.
func (obj *Authenticator) DestroyToken(ctx context.Context, tokenStr string) error {
	// If the token cache service is not nil, remove the token.
	if obj.cache != nil {
		err := obj.cache.Remove(ctx, tokenStr)
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return err
		}
	}
	return nil
}

// DestroyTokenContext destroys the token string from the context.
func (obj *Authenticator) DestroyTokenContext(ctx context.Context, token security.TokenType) error {
	// Get the token string from the context.
	tokenStr, err := helper.FromTokenTypeContext(ctx, token, obj.schemeString())
	if err != nil || tokenStr == "" {
		return ErrInvalidToken
	}
	// Destroy the token string.
	return obj.DestroyToken(ctx, tokenStr)
}

// Close closes the token cache service.
func (obj *Authenticator) Close(ctx context.Context) error {
	// If the token cache service is not nil, close it.
	if obj.cache != nil {
		return obj.cache.Close(ctx)
	}
	return nil
}

// parseToken parses the token string and returns the token.
func (obj *Authenticator) parseToken(token string) (*jwtv5.Token, error) {
	// If the key function is nil, return an error.
	if obj.keyFunc == nil {
		return nil, ErrMissingKeyFunc
	}
	// If the extra keys are nil, parse the token with the key function.
	if obj.extraKeys == nil {
		return jwtv5.ParseWithClaims(token, &jwtv5.RegisteredClaims{}, obj.keyFunc)
	}

	// Otherwise, parse the token with the key function and the extra keys.
	return jwtv5.Parse(token, obj.keyFunc)
}

// generateToken generates a signed token string from the token.
func (obj *Authenticator) generateToken(jwtToken *jwtv5.Token) (string, error) {
	// If the key function is nil, return an error.
	if obj.keyFunc == nil {
		return "", ErrMissingKeyFunc
	}

	// Get the key from the key function.
	key, err := obj.keyFunc(jwtToken)
	if err != nil {
		return "", ErrGetKeyFailed
	}

	// Generate the token string.
	strToken, err := jwtToken.SignedString(key)
	if err != nil {
		return "", ErrSignTokenFailed
	}

	return strToken, nil
}

// NewAuthenticator creates a new Authenticator.
func NewAuthenticator(cfg *configv1.Security, ss ...Setting) (security.Authenticator, error) {
	// Get the JWT config from the security config.
	config := cfg.GetAuthn().GetJwt()
	if config == nil {
		return nil, errors.New("authenticator jwt config is empty")
	}
	// If the signing key is empty, return an error.
	signingKey := config.SigningKey
	if signingKey == "" {
		return nil, errors.New("signing key is empty")
	}

	// Get the signing method and key function from the signing key.
	signingMethod, keyFunc, err := getSigningMethodAndKeyFunc(config.Algorithm, config.SigningKey)
	if err != nil {
		return nil, err
	}

	// Apply the settings to the Authenticator.
	auth := settings.Apply(&Authenticator{
		signingMethod: signingMethod,
		keyFunc:       keyFunc,
	}, ss)
	return auth, nil
}

var _ security.Authenticator = (*Authenticator)(nil)
