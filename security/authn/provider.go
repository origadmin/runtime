/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package authn

import (
	authnv1 "github.com/origadmin/runtime/api/gen/go/config/security/authn/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/security"
)

// Provider is an interface for a security component that can provide various authentication-related capabilities.
type Provider interface {
	// Authenticator returns the Authenticator capability, if supported.
	Authenticator() (security.Authenticator, bool)
	// CredentialCreator returns the CredentialCreator capability, if supported.
	CredentialCreator() (security.CredentialCreator, bool)
	// CredentialRevoker returns the CredentialRevoker capability, if supported.
	CredentialRevoker() (security.CredentialRevoker, bool)
}

// Factory is a function type that creates a Provider instance.
type Factory func(config *authnv1.Authenticator, opts ...options.Option) (Provider, error)
