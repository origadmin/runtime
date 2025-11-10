/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

// SecurityPolicy combines an Authenticator and an Authorizer.
// It represents a complete security processing unit.
type SecurityPolicy interface {
	Authenticator
	Authorizer
}

// SecurityFactory is a factory function for creating SecurityPolicy instances.
// It receives the specific configuration for the policy from the config file.
type SecurityFactory func(config []byte) (SecurityPolicy, error)
