/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"context"
	"time"
)

// CacheStorage is the interface for cache the Authenticator token.
type CacheStorage interface {
	// Store stores the token with a specific expiration time to TokenService
	Store(context.Context, string, time.Duration) error
	// Exist checks if the token exists in the TokenService
	Exist(context.Context, string) (bool, error)
	// Remove deletes the token from the TokenService
	Remove(context.Context, string) error
	// Close closes the TokenService
	Close(context.Context) error
}
