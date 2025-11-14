/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security provides interfaces for declarative security policies.
package security

import (
	"google.golang.org/protobuf/types/known/structpb" // Import structpb
)

// Claims defines a standard, type-safe interface for accessing principal's claims.
// It is designed to be created from a raw map[string]any by a factory,
// ensuring that all data is validated and normalized at the source.
// The Claims object is the single source of truth for custom claims data,
// and it does not provide any backdoor to access the raw underlying map.
type Claims interface {
	// Get looks up a key and returns the raw, validated value.
	// This method is discouraged for direct use in business logic.
	// Prefer using the type-safe accessors like GetString, GetInt64, etc.
	Get(key string) (any, bool)

	// GetString retrieves a value as a string.
	GetString(key string) (string, bool)

	// GetInt64 retrieves a value as an int64, handling conversions from float64/string.
	GetInt64(key string) (int64, bool)

	// GetFloat64 retrieves a value as a float64.
	GetFloat64(key string) (float64, bool)

	// GetBool retrieves a value as a bool.
	GetBool(key string) (bool, bool)

	// GetStringSlice retrieves a value as a slice of strings.
	GetStringSlice(key string) ([]string, bool)

	// Export converts the claims into a map of standard google.protobuf.Value messages.
	// This method is for internal use by the Principal's Export method.
	// It is guaranteed to succeed for a valid Claims object.
	Export() map[string]*structpb.Value
}
