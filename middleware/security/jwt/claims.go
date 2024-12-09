/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package jwt implements the functions, types, and interfaces for the module.
package jwt

import (
	"bytes"
	"maps"
	"strings"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/origadmin/toolkits/security"
	"google.golang.org/protobuf/types/known/timestamppb"

	securityv1 "github.com/origadmin/runtime/gen/go/security/v1"
)

var (
	ErrInvalidToken             = securityv1.ErrorAuthErrorReasonBearerTokenMissing("invalid bearer token")
	ErrTokenMalformed           = securityv1.ErrorAuthErrorReasonBearerTokenMissing("token malformed")
	ErrTokenSignatureInvalid    = securityv1.ErrorAuthErrorReasonSignTokenFailed("token signature invalid")
	ErrTokenExpired             = securityv1.ErrorAuthErrorReasonTokenExpired("token expired")
	ErrTokenNotValidYet         = securityv1.ErrorAuthErrorReasonTokenExpired("token not valid yet")
	ErrUnsupportedSigningMethod = securityv1.ErrorAuthErrorReasonUnsupportedSigningMethod("unsupported signing method")
	ErrInvalidClaims            = securityv1.ErrorAuthErrorReasonInvalidClaims("invalid Claims")
	ErrBearerTokenMissing       = securityv1.ErrorAuthErrorReasonBearerTokenMissing("bearer token missing")
	ErrSignTokenFailed          = securityv1.ErrorAuthErrorReasonSignTokenFailed("sign token failed")
	ErrMissingKeyFunc           = securityv1.ErrorAuthErrorReasonMissingKeyFunc("missing key function")
	ErrGetKeyFailed             = securityv1.ErrorAuthErrorReasonGetKeyFailed("get key failed")
	ErrInvalidSubject           = securityv1.ErrorAuthErrorReasonInvalidSubject("invalid subject")
	ErrInvalidIssuer            = securityv1.ErrorAuthErrorReasonInvalidIssuer("invalid issuer")
	ErrInvalidAudience          = securityv1.ErrorAuthErrorReasonInvalidAudience("invalid audience")
	ErrInvalidExpiration        = securityv1.ErrorAuthErrorReasonInvalidExpiration("invalid expiration")
	//ErrInvalidNotBefore         = securityv1.ErrorAuthErrorReasonInvalidNotBefore("invalid not before")
	//ErrInvalidIssuedAt          = securityv1.ErrorAuthErrorReasonInvalidIssuedAt("invalid issued at")
)

type SecurityClaims struct {
	*securityv1.Claims
	Extra map[string]string
}

func (s *SecurityClaims) GetSubject() string {
	return s.Claims.Sub
}

func (s *SecurityClaims) GetIssuer() string {
	return s.Claims.Iss
}

func (s *SecurityClaims) GetAudience() []string {
	return s.Claims.Aud
}

func (s *SecurityClaims) GetExpiration() time.Time {
	return s.Claims.Exp.AsTime()
}

func (s *SecurityClaims) GetNotBefore() time.Time {
	return s.Claims.Nbf.AsTime()
}

func (s *SecurityClaims) GetIssuedAt() time.Time {
	return s.Claims.Iat.AsTime()
}

func (s *SecurityClaims) GetJWTID() string {
	return s.Claims.Jti
}

func (s *SecurityClaims) GetExtra() map[string]string {
	return s.Extra
}

func (s *SecurityClaims) GetScopes() map[string]bool {
	return s.Claims.Scopes
}

func ClaimsToJwtClaims(raw security.Claims) jwtv5.Claims {
	Claims := jwtv5.MapClaims{
		"sub": raw.GetSubject(),
	}

	if iss := raw.GetIssuer(); iss != "" {
		Claims["iss"] = raw.GetIssuer()
	}
	if aud := raw.GetAudience(); len(aud) > 0 {
		Claims["aud"] = aud
	}
	if exp := raw.GetExpiration(); !exp.IsZero() {
		Claims["exp"] = exp
	}

	var buffer bytes.Buffer
	count := len(raw.GetScopes())
	idx := 0
	for scope := range raw.GetScopes() {
		buffer.WriteString(scope)
		if idx != count-1 {
			buffer.WriteString(" ")
		}
		idx++
	}
	str := buffer.String()
	if len(str) > 0 {
		Claims["scope"] = buffer.String()
	}

	return Claims
}

func MapToClaims(rawClaims jwtv5.MapClaims, extraKey string) (security.Claims, error) {
	//Claims := security.Claims{
	//	Scopes: make(ScopeSet),
	//}
	Claims := &securityv1.Claims{
		Scopes: make(map[string]bool),
	}

	// optional Subject
	if subjectClaim, err := rawClaims.GetSubject(); err == nil {
		Claims.Sub = subjectClaim
	} else {
		return nil, ErrInvalidSubject
	}
	// optional Issuer
	if issuerClaim, err := rawClaims.GetIssuer(); err == nil {
		Claims.Iss = issuerClaim
	} else {
		return nil, ErrInvalidIssuer
	}
	// optional Audience
	if audienceClaim, err := rawClaims.GetAudience(); err == nil {
		Claims.Aud = append(Claims.Aud, audienceClaim...)
	} else {
		return nil, ErrInvalidAudience
	}
	// optional Expiration
	if expClaim, err := rawClaims.GetExpirationTime(); err == nil {
		if expClaim != nil {
			Claims.Exp = timestamppb.New(expClaim.Time)
		}
	} else {
		return nil, ErrInvalidExpiration
	}
	// optional scopes
	if scopeKey, ok := rawClaims["scope"]; ok {
		if scope, ok := scopeKey.(string); ok {
			scopes := strings.Split(scope, " ")
			for _, s := range scopes {
				Claims.Scopes[s] = true
			}
		}
	}

	extra := make(map[string]string)
	if extraVal, ok := rawClaims[extraKey]; ok {
		if extras, ok := extraVal.(map[string]string); ok {
			extra = maps.Clone(extras)
		}
	}

	return &SecurityClaims{
		Claims: Claims,
		Extra:  extra,
	}, nil
}

func RegisteredToClaims(rawClaims *jwtv5.RegisteredClaims) (security.Claims, error) {
	Claims := &securityv1.Claims{
		Scopes: make(map[string]bool),
	}

	// optional Subject
	if subjectClaim, err := rawClaims.GetSubject(); err == nil {
		Claims.Sub = subjectClaim
	} else {
		return nil, ErrInvalidSubject
	}
	// optional Issuer
	if issuerClaim, err := rawClaims.GetIssuer(); err == nil {
		Claims.Iss = issuerClaim
	} else {
		return nil, ErrInvalidIssuer
	}
	// optional Audience
	if audienceClaim, err := rawClaims.GetAudience(); err == nil {
		Claims.Aud = append(Claims.Aud, audienceClaim...)
	} else {
		return nil, ErrInvalidAudience
	}
	// optional Expiration
	if expClaim, err := rawClaims.GetExpirationTime(); err == nil {
		if expClaim != nil {
			Claims.Exp = timestamppb.New(expClaim.Time)
		}
	} else {
		return nil, ErrInvalidExpiration
	}
	// optional scopes
	//if scopeKey, ok := rawClaims.Scope["scope"]; ok {
	//	if scope, ok := scopeKey.(string); ok {
	//		scopes := strings.Split(scope, " ")
	//		for _, s := range scopes {
	//			Claims.Scopes[s] = true
	//		}
	//	}
	//}

	return &SecurityClaims{
		Claims: Claims,
	}, nil
}

func ToClaims(rawClaims jwtv5.Claims, extraKey string) (security.Claims, error) {
	if Claims, ok := rawClaims.(*jwtv5.RegisteredClaims); ok {
		return RegisteredToClaims(Claims)
	}
	if Claims, ok := rawClaims.(jwtv5.MapClaims); ok {
		return MapToClaims(Claims, extraKey)
	}
	return nil, ErrInvalidClaims
}
