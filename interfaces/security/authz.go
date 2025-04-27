/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"context"
)

type PolicyMap map[string]any
type RoleMap map[string]any

// PolicyRegistry is an interface that defines the methods for a policy manager
type PolicyRegistry interface {
	// AddPolicy adds a policy for a given subject, object, action, domain
	AddPolicy(sec string, pt string, rule []string) error
	// RemovePolicy removes a policy for a given subject, object, action, domain
	RemovePolicy(sec string, pt string, rule []string) error
	// SetRoles Set the role for a given context
	SetRoles(ctx context.Context, roles RoleMap) error
	// SetPolicies sets the policies for a given context
	SetPolicies(ctx context.Context, policies PolicyMap) error
	// SetPolicyRoles sets the policies for a given context.
	// It takes a context, a map of policies, and a map of roles as input.
	// It returns an error if the policies cannot be set.
	SetPolicyRoles(ctx context.Context, policies PolicyMap, roles RoleMap) error
}

// Authorizer is an interface that defines the methods for an authorizer.
// It is used to manage policies and check authorization.
type Authorizer interface {
	// Authorized checks if a user is authorized to perform an action.
	// It takes a context and a UserClaims object as input.
	// It returns a boolean indicating whether the user is authorized and an error if the check fails.
	Authorized(ctx context.Context, policy Policy, object string, action string) (bool, error)
	// AuthorizedWithDomain checks if a user is authorized to perform an action within a specific domain.
	// It takes a context, a UserClaims object, a domain, an object, and an action as input.
	// It returns a boolean indicating whether the user is authorized and an error if the check fails.
	AuthorizedWithDomain(ctx context.Context, policy Policy, object string, action string, domain string) (bool, error)
	// AuthorizedWithExtra checks if a user is authorized to perform an action within a specific domain.
	// It takes a context, a UserClaims object, a domain, an object, and an action as input.
	// It returns a boolean indicating whether the user is authorized and an error if the check fails.
	AuthorizedWithExtra(ctx context.Context, data ExtraData) (bool, error)
}
