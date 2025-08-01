/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"context"
)

// Authenticator interface
type Authenticator interface {
	// Authenticate returns a nil error and the AuthClaims info (if available).
	Authenticate(context.Context, string) (Claims, error)
	// AuthenticateContext returns a nil error and the AuthClaims info (if available).
	// if the subject is authenticated or a non-nil error with an appropriate error cause otherwise.
	AuthenticateContext(context.Context, TokenSource) (Claims, error)
	// DestroyToken invalidate a token by removing it from the token store.
	DestroyToken(context.Context, string) error
	// DestroyRefreshToken by removing from the token store to invalidate a refresh token
	DestroyRefreshToken(context.Context, string) error
}
