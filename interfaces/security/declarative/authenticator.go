/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"

	"google.golang.org/protobuf/types/known/anypb"
)

// principalKey is the context key for storing the Principal object.
type principalKey struct{}

// Credential abstracts the source from which credential data can be extracted.
// It defines a contract for accessing credential information without coupling
// to the underlying transport (e.g., HTTP headers, gRPC metadata).
type Credential interface {
	// Type returns the type of the credential (e.g., "jwt", "apikey").
	Type() string

	// Raw returns the original, unparsed credential string.
	Raw() string

	// Payload returns the parsed credential payload.
	// The type of the message inside Any should correspond to the 'type' field.
	Payload() anypb.Any

	// String returns a string representation of the credential.
	String() string
}

// Authenticator is responsible for validating the identity of the request initiator.
// It parses credential data from a Credential and returns a Principal object.
type Authenticator interface {
	// Authenticate extracts credential data from the provided credential contract
	// and validates them, returning a Principal object if successful.
	Authenticate(ctx context.Context, cred Credential) (Principal, error)
}
