/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// Authorizer is responsible for checking if an authenticated Principal has permission
// to perform a specific action on a given resource.
type Authorizer interface {
	// Authorize checks if the given Principal is authorized to perform the specified
	// action on the target resource.
	//
	// The resourceIdentifier typically uniquely identifies the resource being accessed,
	// e.g., "/users/123", "order:create".
	// The action specifies the operation being attempted, e.g., "read", "write", "delete".
	Authorize(ctx context.Context, p Principal, resourceIdentifier string, action string) (bool, error)
}
