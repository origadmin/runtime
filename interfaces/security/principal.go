/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security provides interfaces for declarative security policies.
package security

import (
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
)

// Principal defines the identity of an authenticated user or system.
// It contains information such as user ID, roles, and claims, and is injected into context.Context for downstream business logic.
type Principal interface {
	// GetID returns the unique identifier of the principal.
	GetID() string

	// GetRoles returns the roles assigned to the principal.
	GetRoles() []string

	// GetPermissions returns the permissions assigned to the principal.
	GetPermissions() []string

	// GetScopes returns the scopes assigned to the principal.
	GetScopes() map[string]bool

	// GetClaims returns a standard, type-safe accessor for the principal's claims.
	// This is the single, recommended way to access all custom claims data.
	GetClaims() Claims

	// Export returns the serializable, transportable representation of the principal.
	// This method is guaranteed to succeed for a valid Principal.
	Export() *securityv1.Principal
}
