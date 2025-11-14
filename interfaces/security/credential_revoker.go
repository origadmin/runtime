/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security provides interfaces for declarative security policies.
package security

import (
	"context"
)

// CredentialRevoker is responsible for invalidating or revoking previously issued credentials.
// This is typically used for logout, forced sign-out, or security-related credential invalidation.
type CredentialRevoker interface {
	// Revoke invalidates the given raw credential string, making it unusable for authentication.
	// Implementations might add the credential to a blacklist, update a session store,
	// or perform other actions to ensure the credential is no longer valid.
	Revoke(ctx context.Context, rawCredential string) error
}
