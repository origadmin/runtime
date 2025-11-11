/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides interfaces for declarative security policies.
package declarative

import (
	"github.com/origadmin/runtime/context"
)

// Principal defines the identity of an authenticated user or system.
// It contains information such as user ID, roles, and claims, and is injected into context.Context for downstream business logic.
type Principal interface {
	GetID() string
	GetRoles() []string
	// GetClaims returns all claims associated with the principal as a map.
	GetClaims() map[string]interface{}
}

type principalKey struct{}

// PrincipalFromContext extracts the Principal from the given context.
// It returns the Principal and a boolean indicating if it was found.
func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}

// PrincipalWithPrincipal returns a new context with the given Principal attached.
// It is used to inject the Principal into the context for downstream business logic.
func PrincipalWithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}
