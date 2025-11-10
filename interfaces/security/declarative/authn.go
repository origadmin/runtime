/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// principalKey is the context key for storing the Principal object.
type principalKey struct{}

// CredentialSource abstracts the source from which credentials can be extracted.
type CredentialSource interface {
	// GetAuthorization returns the value of the Authorization header, if present.
	GetAuthorization() (string, bool)
	// Get returns the value of a specific header/metadata key.
	Get(key string) (string, bool)
	// GetAll returns all available headers/metadata as a map.
	GetAll() map[string][]string
}

// Authenticator is responsible for validating the identity of the request initiator.
// It parses credentials from a CredentialSource and returns a Principal object.
type Authenticator interface {
	// Authenticate extracts credentials from the provided source and validates them,
	// returning a Principal object if successful.
	Authenticate(ctx context.Context, source CredentialSource) (Principal, error)
}
