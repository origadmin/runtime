/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package declarative

import (
	"github.com/go-kratos/kratos/v2/transport"

	iface "github.com/origadmin/runtime/interfaces/security/declarative"
)

// KratosTransportHeaderAdapter implements CredentialSource for Kratos transport.Header.
// This adapter is used by the declarative security middleware to bridge Kratos's transport.Header
// with the generic CredentialSource interface.
type KratosTransportHeaderAdapter struct {
	Header transport.Header
}

// NewKratosTransportHeaderAdapter creates a new adapter for the given Kratos header.
func NewKratosTransportHeaderAdapter(header transport.Header) iface.CredentialSource {
	return &KratosTransportHeaderAdapter{
		Header: header,
	}
}

// GetAuthorization returns the value of the Authorization header, if present.
func (a *KratosTransportHeaderAdapter) GetAuthorization() (string, bool) {
	auth := a.Header.Get("Authorization")
	return auth, auth != ""
}

// Get returns the value of a specific header/metadata key.
func (a *KratosTransportHeaderAdapter) Get(key string) (string, bool) {
	val := a.Header.Get(key)
	return val, val != ""
}

// GetAll returns all available headers/metadata as a map.
func (a *KratosTransportHeaderAdapter) GetAll() map[string][]string {
	keys := a.Header.Keys()
	m := make(map[string][]string)
	for _, key := range keys {
		m[key] = a.Header.Values(key)
	}
	return m
}
