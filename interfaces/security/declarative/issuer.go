/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// TokenIssuer defines the contract for issuing new credentials.
type TokenIssuer interface {
	// Issue issues a new credential for the given principal and returns
	// a standard, serializable IssuedCredential.
	Issue(ctx context.Context, p Principal) (IssuedCredential, error)
}
