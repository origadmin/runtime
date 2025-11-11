/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// TokenParser defines the contract for parsing and validating a raw token string
// and converting it into a Principal. It acts as a reusable engine for handling
// token-specific technical details.
type TokenParser interface {
	// Parse converts a raw token string into a validated Principal object.
	// It is responsible for all technical validation of the token itself
	// (e.g., signature, expiration, format).
	Parse(ctx context.Context, rawToken string) (Principal, error)
}
