/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security provides interfaces for declarative security policies.
package security

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

	// GetClaims returns all claims associated with the principal as a map of string to *anypb.Any.
	// This allows for type-safe unpacking of claims into specific protobuf messages.
	GetClaims() map[string]any
}
