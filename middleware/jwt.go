/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"time"

	authjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

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
	// Use the provided SigningMethod from options if available, otherwise get from config
	var signMethod jwt.SigningMethod
	if opts != nil && opts.SigningMethod != nil {
		signMethod = opts.SigningMethod
	} else {
		signMethod = getSigningMethod(config.SigningMethod)
	}
	kf := getKeyFunc(config.SigningKey, signMethod.Alg())

	// Prioritize user-provided ClaimsFactory.
	// If not provided, use the default factory which is driven by config.
	claimsFactory := opts.ClaimsFactory
	if claimsFactory == nil {
		// Create a default claims factory using the provided configuration.
		claimsFactory = NewClaimsFactory(cfg, opts)
	}

	return authjwt.Server(kf, authjwt.WithClaims(claimsFactory), authjwt.WithSigningMethod(signMethod)), true
}

// JwtClient creates a Kratos client middleware for JWT token generation and injection.
// This middleware dynamically creates a new `authjwt.Client` instance for each request,
// ensuring that the JWT is generated based on the claims present in the request's context.
func JwtClient(cfg *jwtv1.JWT, opts *Options) (KMiddleware, bool) {
	config := cfg.GetConfig()
	if config == nil {
		return nil, false
	}
	// Use the provided SigningMethod from options if available, otherwise get from config
	var signMethod jwt.SigningMethod
	if opts != nil && opts.SigningMethod != nil {
		signMethod = opts.SigningMethod
	} else {
		signMethod = getSigningMethod(config.SigningMethod)
	}
	kf := getKeyFunc(config.SigningKey, signMethod.Alg())

	// Prioritize user-provided ClaimsFactory.
	// If not provided, use the default factory which is driven by config.
	claimsFactory := opts.ClaimsFactory
	if claimsFactory == nil {
		// Create a default claims factory using the provided configuration.
		claimsFactory = NewClaimsFactory(cfg, opts)
	}

	return authjwt.Client(kf, authjwt.WithClaims(claimsFactory), authjwt.WithSigningMethod(signMethod)), true
}

// claimsConfig holds the pre-computed configuration for JWT claims
type claimsConfig struct {
	lifetime   time.Duration
	issuer     string
	audience   []string
	getSubject func() string
}

// NewClaimsFactory is a higher-order function that takes JWT configuration
// and returns a `func() jwt.Claims`. This returned function, when called,
// generates a new set of claims based on the provided configuration.
// It automatically handles issuer, lifetime, and generates a random UUID for the subject.
// The implementation is optimized to avoid recreating static configuration on each call.
func NewClaimsFactory(cfg *jwtv1.JWT, opts *Options) func() jwt.Claims {
	// Pre-compute static configuration
	config := cfg.GetConfig()
	if config == nil {
		// Return a no-op factory if config is missing
		return func() jwt.Claims { return &jwt.RegisteredClaims{} }
	}

	// Initialize claims configuration
	cc := claimsConfig{
		lifetime: 2 * time.Hour, // default lifetime
		issuer:   "origadmin",   // default issuer
	}

	// Set up lifetime if configured
	if config.GetAccessTokenLifetime() != 0 {
		cc.lifetime = time.Duration(config.GetAccessTokenLifetime()) * time.Second
	}

	// Set up issuer if configured
	if config.GetIssuer() != "" {
		cc.issuer = config.GetIssuer()
	}

	// Set up audience if configured
	if len(config.GetAudience()) > 0 {
		cc.audience = append([]string{}, config.GetAudience()...)
	}

	// Set up subject factory
	if opts != nil && opts.SubjectFactory != nil {
		cc.getSubject = opts.SubjectFactory
	} else {
		log.Warn("JWT 'subject' is being generated as a random UUID. For production use, provide a meaningful user identifier via middleware.WithSubjectFactory().")
		cc.getSubject = func() string { return uuid.New().String() }
	}

	// Return a closure that only updates dynamic values
	return func() jwt.Claims {
		now := time.Now()
		return &jwt.RegisteredClaims{
			Issuer:    cc.issuer,
			Subject:   cc.getSubject(),
			Audience:  cc.audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(cc.lifetime)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		}
	}
}

// getSigningMethod returns the jwt.SigningMethod based on the provided string.
func getSigningMethod(alg string) jwt.SigningMethod {
	switch alg {
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
