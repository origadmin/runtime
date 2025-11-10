/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	"context"
)

// Principal defines the identity of an authenticated user or system.
// It contains information such as user ID, roles, and claims, and is injected into context.Context for downstream business logic.
type Principal interface {
	// GetID returns the unique identifier of the user/service.
	GetID() string
	// GetRoles returns the list of roles associated with the principal.
	GetRoles() []string
	// GetClaims returns all claims associated with the principal as a map.
	GetClaims() map[string]interface{}
}

// PrincipalFromContext extracts the Principal from the given context.
// It returns the Principal and a boolean indicating if it was found.
func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(contextKeyPrincipal).(Principal)
	return p, ok
}

// NewContextWithPrincipal injects the given Principal into the context.
// It returns a new context with the Principal value.
func NewContextWithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, contextKeyPrincipal, p)
}

type contextKey string

const (
	contextKeyPrincipal contextKey = "security-principal"
)
