/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

	var options []authjwt.Option
	if cfg.GetTokenHeader() != nil {
		// TokenHeader sets additional header fields for the JWT token itself (e.g., kid, typ)
		options = append(options, authjwt.WithTokenHeader(cfg.GetTokenHeader().AsMap()))
	}
	// Use the provided SigningMethod from options if available, otherwise get from config
	var signMethod jwt.SigningMethod
	if opts != nil && opts.SigningMethod != nil {
		signMethod = opts.SigningMethod
	} else {
		signMethod = getSigningMethod(config.SigningMethod)
	}
	options = append(options, authjwt.WithSigningMethod(signMethod))
	// Check for insecure signing method
	if signMethod == jwt.SigningMethodNone {
		log.Warn("Using insecure signing method 'none'. This should only be used for testing.")
	}

	kf := getKeyFunc(config.SigningKey, signMethod.Alg())

	// For server middleware, we need to validate tokens, not generate them
	// Use a claims factory that creates appropriate claims for parsing
	var claimsFactory func() jwt.Claims
	if opts != nil && opts.ClaimsFactory != nil {
		claimsFactory = opts.ClaimsFactory
	} else {
		// For server validation, we need a factory that creates the correct claims type
		// This ensures the token can be parsed into the expected claims structure
		claimType := "registered" // default claim type
		if cfg.GetClaimType() != "" {
			claimType = cfg.GetClaimType()
		}
		claimsFactory = getClaimsFactory(claimType)
	}
	options = append(options, authjwt.WithClaims(claimsFactory))

	return authjwt.Server(kf, options...), true
}

// JwtClient creates a Kratos client middleware for JWT token generation and injection.
// This middleware dynamically creates a new `authjwt.Client` instance for each request,
// ensuring that the JWT is generated based on the claims present in the request's context.
func JwtClient(cfg *jwtv1.JWT, opts *Options) (KMiddleware, bool) {
	config := cfg.GetConfig()
	if config == nil {
		return nil, false
	}

	var options []authjwt.Option
	if cfg.GetTokenHeader() != nil {
		// TokenHeader sets additional header fields for the JWT token itself (e.g., kid, typ)
		options = append(options, authjwt.WithTokenHeader(cfg.GetTokenHeader().AsMap()))
	}
	// Use the provided SigningMethod from options if available, otherwise get from config
	var signMethod jwt.SigningMethod
	if opts != nil && opts.SigningMethod != nil {
		signMethod = opts.SigningMethod
	} else {
		signMethod = getSigningMethod(config.SigningMethod)
	}
	options = append(options, authjwt.WithSigningMethod(signMethod))
	// Check for insecure signing method
	if signMethod == jwt.SigningMethodNone {
		log.Warn("Using insecure signing method 'none'. This should only be used for testing.")
	}

	kf := getKeyFunc(config.SigningKey, signMethod.Alg())

	// For client middleware, we need to generate tokens, not validate them
	// Use the claims factory to create claims for token generation
	var claimsFactory func() jwt.Claims
	if opts != nil && opts.ClaimsFactory != nil {
		claimsFactory = opts.ClaimsFactory
	} else {
		// NewClaimsFactory internally reads cfg.GetClaimType() to ensure
		// the generated claims match the expected type for consistency
		claimsFactory = NewClaimsFactory(cfg, opts)
	}
	options = append(options, authjwt.WithClaims(claimsFactory))
	return authjwt.Client(kf, options...), true
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
// Exported for testing purposes.
func NewClaimsFactory(cfg *jwtv1.JWT, opts *Options) func() jwt.Claims {
	// Pre-compute static configuration
	config := cfg.GetConfig()
	if config == nil {
		// Return a no-op factory if config is missing
		return func() jwt.Claims { return &jwt.RegisteredClaims{} }
	}

	// Determine claim type from configuration
	claimType := "registered" // default claim type
	if cfg.GetClaimType() != "" {
		claimType = cfg.GetClaimType()
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

	// Pre-create the appropriate claims factory based on claim type
	// This avoids repeated switch statements on each call
	switch claimType {
	case "map":
		// Return a factory that creates MapClaims
		return func() jwt.Claims {
			now := time.Now()
			claims := jwt.MapClaims{
				"iss": cc.issuer,
				"sub": cc.getSubject(),
				"iat": now.Unix(),
				"nbf": now.Unix(),
				"exp": now.Add(cc.lifetime).Unix(),
			}
			
			// Add audience if configured
			if len(cc.audience) > 0 {
				claims["aud"] = cc.audience
			}
			
			return claims
		}
	default:
		// Return a factory that creates RegisteredClaims
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
}

// GetSigningMethod returns the jwt.SigningMethod based on the provided string.
// Exported for testing purposes.
func GetSigningMethod(alg string) jwt.SigningMethod {
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

// getSigningMethod returns the jwt.SigningMethod based on the provided string.
func getSigningMethod(alg string) jwt.SigningMethod {
	return GetSigningMethod(alg)
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

		// Handle different signing methods appropriately
		switch method {
		case "HS256", "HS384", "HS512":
			// HMAC methods use the same key for signing and verification
			return []byte(key), nil
		case "RS256", "RS384", "RS512":
			// RSA methods require proper key parsing
			pubKey, err := parseRSAPublicKey(key)
			if err != nil {
				log.Errorf("Failed to parse RSA public key: %v", err)
				return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
			}
			return pubKey, nil
		case "ES256", "ES384", "ES512":
			// ECDSA methods require proper key parsing
			pubKey, err := parseECDSAPublicKey(key)
			if err != nil {
				log.Errorf("Failed to parse ECDSA public key: %v", err)
				return nil, fmt.Errorf("failed to parse ECDSA public key: %w", err)
			}
			return pubKey, nil
		case "none":
			// No signature verification
			return nil, nil
		default:
			return nil, authjwt.ErrUnSupportSigningMethod
		}
	}
}

// parseRSAPublicKey parses a PEM-encoded RSA public key.
func parseRSAPublicKey(keyData string) (interface{}, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing RSA public key")
	}

	// Try PKIX format first
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// Fallback to PKCS1 format
		if key, errPkcs1 := x509.ParsePKCS1PublicKey(block.Bytes); errPkcs1 == nil {
			return key, nil
		}
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA public key")
	}
	return rsaKey, nil
}

// parseECDSAPublicKey parses a PEM-encoded ECDSA public key.
func parseECDSAPublicKey(keyData string) (interface{}, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing ECDSA public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ECDSA public key: %w", err)
	}

	// Verify it's a valid public key (ECDSA or other types)
	// The actual type checking will be done by the JWT library during verification
	return pub, nil
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
