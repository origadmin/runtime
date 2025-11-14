/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security provides interfaces for declarative security policies.
package security

import (
	"context"
)

// Authorizer is responsible for checking if an authenticated Principal has permission
// to perform a specific action on a given resource.
type Authorizer interface {
	// Authorize checks if an authenticated Principal has permission to perform a specific action on a given resource.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - p: The authenticated Principal.
	//   - resourceIdentifier: The identifier of the resource being accessed (e.g., "data1", "/articles/123").
	//   - action: The action being performed on the resource (e.g., "read", "write", "delete").
	//
	// Returns:
	//   - bool: True if the Principal is authorized, false otherwise.
	//   - error: An error if the authorization check fails for any reason.
	Authorize(ctx context.Context, p Principal, resourceIdentifier string, action string) (bool, error)
}
