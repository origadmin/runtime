/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// Authenticator is responsible for validating the identity of the request initiator.
// It receives credential data and returns a Principal object.
type Authenticator interface {
	// Authenticate validates the provided credential and returns a Principal object if successful.
	Authenticate(ctx context.Context, cred Credential) (Principal, error)
}
