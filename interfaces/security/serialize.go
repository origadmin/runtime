/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security provides interfaces and types for security-related operations
package security

import (
	"context"
)

// Serializer is an interface that defines the methods for a serializer
type Serializer interface {
	// Serialize serializes the given data into a byte slice
	Serialize(ctx context.Context, data Claims) ([]byte, error)
	// Deserialize deserializes the given byte slice into the given data
	Deserialize(ctx context.Context, data []byte) (Claims, error)
}
