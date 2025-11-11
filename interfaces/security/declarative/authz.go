/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// Authorizer is responsible for checking if an authenticated Principal has permission to perform the current operation.
// It typically uses the Principal's roles, permissions, and an access control model (e.g., RBAC, ABAC) to make the decision.
type Authorizer interface {
	// Authorize checks if the Principal has permission to access the resource.
	// resourceIdentifier is a string representing the resource being accessed.
	// For gRPC, it's the full method name (e.g., "/api.Greeter/SayHello").
	// For HTTP, it's the method:path (e.g., "GET:/helloworld/{name}").
	// Returns true if authorized, false otherwise.
	Authorize(ctx context.Context, p Principal, resourceIdentifier string) (bool, error)
}
