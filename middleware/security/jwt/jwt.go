/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package jwt implements the functions, types, and interfaces for the module.
package jwt

import (
	"context"
	"errors"

	"github.com/goexts/generic/settings"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/origadmin/toolkits/security"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware/security/helper"
)

var _ security.Authenticator = (*Authenticator)(nil)

type Option struct {
	signingMethod jwtv5.SigningMethod
	keyFunc       func(*jwtv5.Token) (any, error)
}

type Setting = func(*Option)

type Authenticator struct {
	option   *Option
	ExtraKey string
}

func (jwt *Authenticator) AuthenticateToken(tokenString string) (security.Claims, error) {
	jwtToken, err := jwt.parseToken(tokenString)

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

	if jwtToken.Method != jwt.option.signingMethod {
		return nil, ErrUnsupportedSigningMethod
	}

	if jwtToken.Claims == nil {
		return nil, ErrInvalidClaims
	}

	authClaims, err := ToClaims(jwtToken.Claims, jwt.ExtraKey)
	if err != nil {
		return nil, err
	}

	return authClaims, nil
}

func (jwt *Authenticator) AuthenticateTokenContext(ctx context.Context, contextType security.ContextType) (security.Claims, error) {
	tokenString := helper.FromMD(ctx, "bearer")
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	return jwt.AuthenticateToken(tokenString)
}

func (jwt *Authenticator) Authenticate(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.AuthenticateToken(tokenString)
	if err != nil {
		return false, err
	}
	// TODO: check token
	_ = token
	return true, nil
}

func (jwt *Authenticator) AuthenticateContext(ctx context.Context, contextType security.ContextType, tokenString string) (bool, error) {
	return jwt.Authenticate(ctx, tokenString)
}

func (jwt *Authenticator) CreateToken(claims security.Claims) (string, error) {
	jwtToken := jwtv5.NewWithClaims(jwt.option.signingMethod, ClaimsToJwtClaims(claims))

	strToken, err := jwt.generateToken(jwtToken)
	if err != nil {
		return "", err
	}

	return strToken, nil
}

func (jwt *Authenticator) CreateTokenContext(ctx context.Context, contextType security.ContextType, claims security.Claims) (context.Context, error) {
	strToken, err := jwt.CreateToken(claims)
	if err != nil {
		return ctx, err
	}

	//ctx = utils.MDWithAuth(ctx, utils.BearerWord, strToken, contextType)
	_ = strToken
	return ctx, nil
}

func (jwt *Authenticator) DestroyToken(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (jwt *Authenticator) DestroyTokenContext(ctx context.Context, contextType security.ContextType, s string) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthenticator(security *configv1.Security, ss ...Setting) (security.Authenticator, error) {
	config := security.GetJwt()
	if config == nil {
		return nil, errors.New("security config is empty")
	}
	option := settings.Apply(&Option{}, ss)
	if option.signingMethod == nil {
		option.signingMethod = GetSigningMethodFromAlg(config.GetAlgorithm())
	}
	if config.SigningKey == "" {
		return nil, errors.New("signing key is empty")
	}
	if option.keyFunc == nil {
		option.keyFunc = GetKeyFuncWithAlg(config.Algorithm, config.SigningKey)
	}

	auth := &Authenticator{
		option: option,
		//Claims: &SecurityClaims{},
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

//// Authenticate authenticates the token string and returns the claims.
//func (jwt *Authenticator) Authenticate(ctx context.Context, contextType security.ContextType) (security.Claims, error) {
//	tokenString := helper.FromMD(ctx, utils.BearerWord, contextType)
//	if err != nil {
//		return nil, ErrMissingBearerToken
//	}
//
//	return jwt.AuthenticateToken(tokenString)
//}
//
//// AuthenticateToken authenticates the token string and returns the claims.
//func (jwt *Authenticator) AuthenticateToken(tokenString string) (security.Claims, error) {
//	jwtToken, err := jwt.parseToken(tokenString)
//
//	if jwtToken == nil {
//		return nil, ErrInvalidToken
//	}
//
//	if err != nil {
//		switch {
//		case errors.Is(err, jwtv5.ErrTokenMalformed):
//			return nil, ErrInvalidToken
//		case errors.Is(err, jwtv5.ErrTokenSignatureInvalid):
//			return nil, ErrSignTokenFailed
//		case errors.Is(err, jwtv5.ErrTokenExpired) || errors.Is(err, jwtv5.ErrTokenNotValidYet):
//			return nil, ErrTokenExpired
//		default:
//			return nil, ErrInvalidToken
//		}
//	}
//
//	if !jwtToken.Valid {
//		return nil, ErrInvalidToken
//	}
//	if jwtToken.Method != jwt.option.signingMethod {
//		return nil, ErrUnsupportedSigningMethod
//	}
//	if jwtToken.Claims == nil {
//		return nil, ErrInvalidClaims
//	}
//
//	claims, ok := jwtToken.Claims.(jwtv5.MapClaims)
//	if !ok {
//		return nil, ErrInvalidClaims
//	}
//
//	authClaims, err := utils.MapClaimsToAuthClaims(claims)
//	if err != nil {
//		return nil, err
//	}
//
//	return authClaims, nil
//}
//
//// CreateIdentityWithContext creates a signed token string from the claims and sets it to the context.
//func (jwt *Authenticator) CreateIdentityWithContext(ctx context.Context, contextType security.ContextType, claims security.Claims) (context.Context, error) {
//	strToken, err := jwt.CreateIdentity(claims)
//	if err != nil {
//		return ctx, err
//	}
//
//	ctx = utils.MDWithAuth(ctx, utils.BearerWord, strToken, contextType)
//
//	return ctx, nil
//}
//
//// CreateIdentity creates a signed token string from the claims.
//func (jwt *Authenticator) CreateIdentity(claims security.Claims) (string, error) {
//	jwtToken := jwtv5.NewWithClaims(jwt.option.signingMethod, utils.AuthClaimsToJwtClaims(claims))
//
//	strToken, err := jwt.generateToken(jwtToken)
//	if err != nil {
//		return "", err
//	}
//
//	return strToken, nil
//}

func (jwt *Authenticator) Close() {}

// parseToken parses the token string and returns the token.
func (jwt *Authenticator) parseToken(token string) (*jwtv5.Token, error) {
	if jwt.option.keyFunc == nil {
		return nil, ErrMissingKeyFunc
	}
	if jwt.ExtraKey != "" {
		return jwtv5.ParseWithClaims(token, &jwtv5.RegisteredClaims{}, jwt.option.keyFunc)
	}

	return jwtv5.Parse(token, jwt.option.keyFunc)
}

// generateToken generates a signed token string from the token.
func (jwt *Authenticator) generateToken(jwtToken *jwtv5.Token) (string, error) {
	if jwt.option.keyFunc == nil {
		return "", ErrMissingKeyFunc
	}

	key, err := jwt.option.keyFunc(jwtToken)
	if err != nil {
		return "", ErrGetKeyFailed
	}

	strToken, err := jwtToken.SignedString(key)
	if err != nil {
		return "", ErrSignTokenFailed
	}

	return strToken, nil
}
