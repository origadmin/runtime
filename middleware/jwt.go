/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	authjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"

	jwtv1 "github.com/origadmin/runtime/api/gen/go/config/middleware/jwt/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/log"
)

type jwtFactory struct{}

func (f jwtFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debugf("enabling jwt client middleware")

	jwtConfig := cfg.GetJwt()
	if jwtConfig == nil {
		return nil, false
	}

	return JwtClient(jwtConfig, mwOpts)
}

func (f jwtFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)
	helper.Debugf("enabling jwt server middleware")

	jwtConfig := cfg.GetJwt()
	if jwtConfig == nil {
		return nil, false
	}

	return JwtServer(jwtConfig, mwOpts)
}

// JwtServer creates a Kratos server middleware for JWT authentication.
// It uses the provided JWT configuration to validate incoming tokens.
func JwtServer(cfg *jwtv1.JWT, opts *Options) (KMiddleware, bool) {
	config := cfg.GetConfig()
	if config == nil {
		return nil, false
	}
	// The key function can be created once and reused for all requests, as the key and method are static.
	kf := getKeyFunc(config.SigningKey, config.SigningMethod)
	claimsFactory := getClaimsFactory(cfg.GetClaimType())
	return authjwt.Server(kf, authjwt.WithClaims(claimsFactory)), true
}

// JwtClient creates a Kratos client middleware for JWT token generation and injection.
// This middleware dynamically creates a new `authjwt.Client` instance for each request,
// ensuring that the JWT is generated based on the claims present in the request's context.
func JwtClient(cfg *jwtv1.JWT, opts *Options) (KMiddleware, bool) {
	config := cfg.GetConfig()
	if config == nil {
		return nil, false
	}
	// The key function can be created once and reused for all requests, as the key and method are static.
	kf := getKeyFunc(config.SigningKey, config.SigningMethod)
	claimsFactory := getClaimsFactory(cfg.GetClaimType())
	return authjwt.Client(kf, authjwt.WithClaims(claimsFactory)), true
}

// getSigningMethod returns the jwt.SigningMethod based on the provided string.
func getSigningMethod(sm string) jwt.SigningMethod {
	switch sm {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS384":
		return jwt.SigningMethodRS384
	case "RS512":
		return jwt.SigningMethodRS512
	case "ES256":
		return jwt.SigningMethodES256
	case "ES384":
		return jwt.SigningMethodES384
	case "ES512":
		return jwt.SigningMethodES512
	default:
		return jwt.SigningMethodNone
	}
}

// getKeyFunc returns a jwt.Keyfunc for token validation.
func getKeyFunc(key string, method string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if key == "" {
			return nil, authjwt.ErrMissingKeyFunc
		}
		if token.Method.Alg() != method {
			return nil, authjwt.ErrUnSupportSigningMethod
		}
		return []byte(key), nil
	}
}

// getClaimsFactory returns a function that creates an empty jwt.Claims object
// based on the claimType. This is used by the server to parse incoming tokens.
func getClaimsFactory(claimType string) func() jwt.Claims {
	switch claimType {
	case "map":
		return func() jwt.Claims {
			return jwt.MapClaims{}
		}
	default:
		return func() jwt.Claims {
			return &jwt.RegisteredClaims{}
		}
	}
}
