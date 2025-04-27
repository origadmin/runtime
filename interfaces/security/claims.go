/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security is a toolkit for security check and authorization
package security

// Claims is an interface that defines the methods that a security claims object should have
// It provides methods for getting the subject, issuer, audience, expiration, not before, issued at, JWT ID, and scopes of the claims
type Claims interface {
	// GetSubject returns the subject of the security
	GetSubject() string
	// GetIssuer returns the issuer of the security
	GetIssuer() string
	// GetAudience returns the audience of the security
	GetAudience() []string
	// GetExpiration returns the expiration time of the security
	GetExpiration() int64
	// GetNotBefore returns the time before which the security cannot be accepted
	GetNotBefore() int64
	// GetIssuedAt returns the time at which the security was issued
	GetIssuedAt() int64
	// GetID returns the unique identifier for the security
	GetID() string
	// GetScopes returns the scopes associated with the security
	GetScopes() map[string]bool
}

// RegisteredClaims is a struct that implements the Claims interface
// It provides fields for the subject, issuer, audience, expiration, not before, issued at, JWT ID, and scopes of the claims
// json example:
//
//	{
//	  "sub": "test_subject",
//	  "iss": "test_issuer",
//	  "aud": [
//	    "test_audience1",
//	    "test_audience2"
//	  ],
//	  "exp": 1735647621,
//	  "nbf": 1735644021,
//	  "iat": 1735644021,
//	  "jti": "test_jti",
//	  "scopes": {
//	    "scope1": true,
//	    "scope2": false
//	  }
//	}
type RegisteredClaims struct {
	ID         string          `json:"jti,omitempty"`
	Subject    string          `json:"sub,omitempty"`
	Issuer     string          `json:"iss,omitempty"`
	Audience   []string        `json:"aud,omitempty"`
	Expiration int64           `json:"exp,omitempty"`
	NotBefore  int64           `json:"nbf,omitempty"`
	IssuedAt   int64           `json:"iat,omitempty"`
	Scopes     map[string]bool `json:"scopes,omitempty"`
}

// GetSubject returns the subject of the claims
func (r RegisteredClaims) GetSubject() string {
	return r.Subject
}

// GetIssuer returns the issuer of the claims
func (r RegisteredClaims) GetIssuer() string {
	return r.Issuer
}

// GetAudience returns the audience of the claims
func (r RegisteredClaims) GetAudience() []string {
	return r.Audience
}

// GetExpiration returns the expiration time of the claims
func (r RegisteredClaims) GetExpiration() int64 {
	return r.Expiration
}

// GetNotBefore returns the time before which the claims cannot be accepted
func (r RegisteredClaims) GetNotBefore() int64 {
	return r.NotBefore
}

// GetIssuedAt returns the time at which the claims were issued
func (r RegisteredClaims) GetIssuedAt() int64 {
	return r.IssuedAt
}

// GetID returns the unique identifier for the claims
func (r RegisteredClaims) GetID() string {
	return r.ID
}

// GetScopes returns the scopes associated with the RegisteredClaims.
func (r RegisteredClaims) GetScopes() map[string]bool {
	// Return the scopes map directly, as it is already a map[string]bool.
	return r.Scopes
}

// ExtraClaims represents a claims object that contains both registered claims and extra claims.
type ExtraClaims struct {
	// Claims is the registered claims part of the ExtraClaims object.
	Claims Claims
	// Extra is the extra claims part of the ExtraClaims object.
	Extra Extra
}

// GetSubject returns the subject of the claims.
func (e ExtraClaims) GetSubject() string {
	// Delegate to the Claims field to get the subject.
	return e.Claims.GetSubject()
}

// GetIssuer returns the issuer of the claims.
func (e ExtraClaims) GetIssuer() string {
	// Delegate to the Claims field to get the issuer.
	return e.Claims.GetIssuer()
}

// GetAudience returns the audience of the claims.
func (e ExtraClaims) GetAudience() []string {
	// Delegate to the Claims field to get the audience.
	return e.Claims.GetAudience()
}

// GetExpiration returns the expiration time of the claims.
func (e ExtraClaims) GetExpiration() int64 {
	// Delegate to the Claims field to get the expiration time.
	return e.Claims.GetExpiration()
}

// GetNotBefore returns the time before which the claims cannot be accepted.
func (e ExtraClaims) GetNotBefore() int64 {
	// Delegate to the Claims field to get the not-before time.
	return e.Claims.GetNotBefore()
}

// GetIssuedAt returns the time at which the claims were issued.
func (e ExtraClaims) GetIssuedAt() int64 {
	// Delegate to the Claims field to get the issued-at time.
	return e.Claims.GetIssuedAt()
}

// GetID returns the unique identifier for the claims.
func (e ExtraClaims) GetID() string {
	// Delegate to the Claims field to get the ID.
	return e.Claims.GetID()
}

// GetScopes returns the scopes associated with the claims.
func (e ExtraClaims) GetScopes() map[string]bool {
	// Delegate to the Claims field to get the scopes.
	return e.Claims.GetScopes()
}

// GetExtra returns the extra claims as a map of strings.
func (e ExtraClaims) GetExtra() map[string]string {
	// Delegate to the Extra field to get the extra claims.
	return e.Extra.GetExtra()
}

// Get returns the value of the given key from the extra claims.
func (e ExtraClaims) Get(key string) (string, bool) {
	// Delegate to the Extra field to get the value of the given key.
	return e.Extra.Get(key)
}

// Set sets the value of the given key in the extra claims.
func (e ExtraClaims) Set(key string, value string) {
	// Delegate to the Extra field to set the value of the given key.
	e.Extra.Set(key, value)
}

var _ Claims = (*RegisteredClaims)(nil)
var _ Claims = (*ExtraClaims)(nil)
var _ Extra = (*ExtraClaims)(nil)
