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

// Credential abstracts the source from which credential data can be extracted.
// It defines a contract for accessing credential information without coupling
// to the underlying transport (e.g., HTTP headers, gRPC metadata).
type Credential interface {
	// GetAuthorization returns the value of the Authorization header, if present.
	GetAuthorization() (string, bool)
	// Get returns the value of a specific header/metadata key.
	Get(key string) (string, bool)
	// GetAll returns all available headers/metadata as a map.
	GetAll() map[string][]string
}

// Authenticator is responsible for validating the identity of the request initiator.
// It parses credential data from a Credential and returns a Principal object.
type Authenticator interface {
	// Authenticate extracts credential data from the provided credential contract
	// and validates them, returning a Principal object if successful.
	Authenticate(ctx context.Context, cred Credential) (Principal, error)
}
