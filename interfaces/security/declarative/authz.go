/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security_declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// Authorizer is responsible for checking if an authenticated Principal has permission to perform the current operation.
// It typically uses the Principal's roles, permissions, and an access control model (e.g., RBAC, ABAC) to make the decision.
type Authorizer interface {
	// Authorize checks if the Principal has permission to access the resource.
	// fullMethodName is the full name of the method being called (e.g., "/api.Greeter/SayHello").
	// Returns true if authorized, false otherwise.
	Authorize(ctx context.Context, p Principal, fullMethodName string) (bool, error)
}
