/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package security

import (
	"strings"
)

// Scheme represents the type of authorization.
//
//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer@latest -type=Scheme -trimprefix=Scheme
type Scheme int

const (
	// SchemeAnonymous represents an anonymous authorization.
	SchemeAnonymous Scheme = iota
	// SchemeBasic represents a basic authorization.
	SchemeBasic
	// SchemeBearer represents a bearer authorization.
	SchemeBearer
	// SchemeDigest represents a digest authorization.
	SchemeDigest
	// SchemeHOBA represents a HTTP Origin-Bound Authentication (HOBA) authorization.
	SchemeHOBA
	// SchemeMutual represents a mutual authentication.
	SchemeMutual
	// SchemeNegotiate represents a negotiate authorization.
	SchemeNegotiate
	// SchemeVapid represents a VAPID authorization.
	SchemeVapid
	// SchemeSCRAM represents a SCRAM authorization.
	SchemeSCRAM
	// SchemeAWS4HMAC256 represents an AWS4-HMAC-SHA256 authorization.
	SchemeAWS4HMAC256
	// SchemeDPoP represents a DPoP authorization.
	SchemeDPoP
	// SchemeGNAP represents a GNAP authorization.
	SchemeGNAP
	// SchemePrivate represents a private authorization.
	SchemePrivate
	// SchemeOAuth represents an OAuth authorization.
	SchemeOAuth
	// SchemeUnknown represents an unknown authorization.
	SchemeUnknown
	SchemeMax
)

const (
	// SchemeNTLM represents an NTLM authorization.
	SchemeNTLM = SchemeNegotiate
)

// Lower returns the lowercase string representation of the Type.
func (t Scheme) Lower() string {
	return strings.ToLower(t.String())
}

func (t Scheme) Equal(other string) bool {
	schemeType := t.String() + " "
	if len(other) < len(schemeType) {
		return false
	}

	return strings.EqualFold(other[:len(schemeType)], schemeType)
}
