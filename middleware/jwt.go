/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"time"

	authjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/goexts/generic/maps"
	"github.com/golang-jwt/jwt/v5"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	jwtv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/jwt"
	secjwtv1 "github.com/origadmin/runtime/api/gen/go/security/jwt/v1"
	"github.com/origadmin/runtime/log"
)

type jwtFactory struct{}

func (f jwtFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] Jwt client middleware enabled")
	if middleware.GetJwt().GetEnabled() {
		m, ok := JwtClient(middleware.GetJwt())
		if ok && middleware.GetSelector().GetEnabled() {
			m = SelectorServer(middleware.GetSelector(), options.MatchFunc, m)
		}
	}
	return nil, false
}

func (f jwtFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	log.Debug("[Middleware] Jwt server middleware enabled")
	if middleware.GetJwt().GetEnabled() {
		m, ok := JwtServer(middleware.GetJwt())
		if ok && middleware.GetSelector().GetEnabled() {
			m = SelectorServer(middleware.GetSelector(), options.MatchFunc, m)
		}
		return m, ok
	}
	return nil, false
}

func JwtServer(cfg *jwtv1.JWT) (KMiddleware, bool) {
	config := cfg.GetConfig()
	if config == nil {
		return nil, false
	}
	kf := getKeyFunc(config.Key, config.SigningMethod)
	opts := fromJwtConfig(config, cfg.GetSubject(), cfg.GetClaimType(), cfg.GetTokenHeader())
	return authjwt.Server(kf, opts...), true
}
func JwtClient(cfg *jwtv1.JWT) (KMiddleware, bool) {
	config := cfg.GetConfig()
	if config == nil {
		return nil, false
	}
	kf := getKeyFunc(config.Key, config.SigningMethod)
	opts := fromJwtConfig(config, cfg.GetSubject(), cfg.GetClaimType(), cfg.GetTokenHeader())
	return authjwt.Client(kf, opts...), true
}

func fromJwtConfig(cfg *secjwtv1.Config, subject string, ctp string, header map[string]string) []authjwt.Option {
	sm := getSigningMethod(cfg.SigningMethod)
	jcf := getClaimsFunc(subject, ctp, cfg)
	tkh := getTokenHeader(header)
	return []authjwt.Option{
		authjwt.WithSigningMethod(sm),
		authjwt.WithClaims(jcf),
		authjwt.WithTokenHeader(tkh),
	}
}

func getTokenHeader(header map[string]string) map[string]any {
	if header == nil {
		return map[string]any{}
	}
	return maps.Transform(header, func(k, v string) (string, any, bool) {
		return k, v, true
	})
}
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

func getKeyFunc(key string, method string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if key == "" {
			return nil, authjwt.ErrMissingKeyFunc
		}
		if token.Method.Alg() != method {
			return nil, authjwt.ErrUnSupportSigningMethod
		}
		return key, nil
	}
}

func getClaimsFunc(subject string, claimType string, cfg *secjwtv1.Config) func() jwt.Claims {
	if subject == "" {
		subject = "anonymous"
	}
	exp := time.Duration(cfg.GetAccessTokenLifetime())
	if exp == 0 {
		exp = time.Hour
	}
	switch claimType {
	case "map":
		return func() jwt.Claims {
			now := time.Now()
			return jwt.MapClaims{
				"iss": cfg.Issuer,
				"sub": subject,
				"aud": cfg.Audience,
				"exp": now.Add(exp).Unix(),
				"nbf": now.Unix(),
				"iat": now.Unix(),
			}
		}
	default:
		return func() jwt.Claims {
			now := time.Now()
			return &jwt.RegisteredClaims{
				Issuer:    cfg.Issuer,
				Subject:   subject,
				Audience:  cfg.Audience,
				ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			}
		}
	}
}
