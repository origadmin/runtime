/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security provides declarative security interfaces for authentication and authorization.
package security

import (
	"github.com/origadmin/runtime/context"
)

// PolicyProvider defines the interface for providing security policies.
// It is responsible for mapping a request (e.g., gRPC method) to a specific security policy name,
// and then retrieving the actual SecurityPolicy implementation by that name.
type PolicyProvider interface {
	// GetPolicyNameForMethod retrieves the policy name for a given gRPC method.
	// The middleware uses this to determine which policy to apply.
	GetPolicyNameForMethod(ctx context.Context, fullMethodName string) (string, error)

	// GetPolicy retrieves the SecurityPolicy implementation by its name.
	// It returns an error if the policy is not found or not configured.
	GetPolicy(ctx context.Context, policyName string) (SecurityPolicy, error)
}

// SecurityPolicy combines an Authenticator and an Authorizer.
// It represents a complete security processing unit.
type SecurityPolicy interface {
	Authenticator
	Authorizer
}

// SecurityFactory is a factory function for creating SecurityPolicy instances.
// It receives the specific configuration for the policy from the config file.
type SecurityFactory func(config []byte) (SecurityPolicy, error)
