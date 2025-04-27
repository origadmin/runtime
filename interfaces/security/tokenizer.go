/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"context"
)

type Tokenizer interface {
	// CreateClaims creates a new identity claims.
	CreateClaims(context.Context, string) (Claims, error)
	// CreateToken inject user claims into token string.
	CreateToken(context.Context, Claims) (string, error)
	// ParseClaims parses a token string and returns the Claims.
	ParseClaims(context.Context, string) (Claims, error)
	// Validate validates if a token is valid.
	Validate(context.Context, string) (bool, error)
}

type RefreshTokenizer interface {
	Tokenizer
	// CreateRefreshClaims creates a new identity claims specifically for a refresh token.
	CreateRefreshClaims(context.Context, string) (Claims, error)
}
