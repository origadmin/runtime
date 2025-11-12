/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package declarative

import (
	"fmt"
	"strings"

	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/security/declarative"
	"github.com/origadmin/runtime/security/meta"
)

const (
	// AuthorizationHeader is the canonical header name for authorization.
	AuthorizationHeader = "Authorization"
)

// HeaderCredentialExtractor implements the declarative.CredentialExtractor interface.
// It extracts credentials from the "Authorization" HTTP header.
type HeaderCredentialExtractor struct{}

// NewHeaderCredentialExtractor creates a new instance of HeaderCredentialExtractor.
func NewHeaderCredentialExtractor() declarative.CredentialExtractor {
	return &HeaderCredentialExtractor{}
}

// Extract retrieves a credential from the "Authorization" header provided by a ValueProvider.
// It expects the header value to be in the format "Scheme CredentialString".
func (e *HeaderCredentialExtractor) Extract(provider declarative.ValueProvider) (declarative.Credential, error) {
	authHeader := provider.Get(AuthorizationHeader)
	if authHeader == "" {
		return nil, errors.New(401, "AUTHORIZATION_HEADER_NOT_FOUND", "authorization header not found")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return nil, errors.New(401, "INVALID_AUTHORIZATION_HEADER", "invalid authorization header format")
	}

	scheme := strings.TrimSpace(parts[0])
	rawCredential := strings.TrimSpace(parts[1])

	if scheme == "" || rawCredential == "" {
		return nil, errors.New(401, "INVALID_AUTHORIZATION_HEADER", "invalid authorization header format")
	}

	// Here, we use the NewCredential function from the same package.
	// Convert scheme to lowercase for case-insensitive comparison
	// Create a simple string value as the payload

	payload := &securityv1.BearerCredential{
		Token: rawCredential,
	}

	t := ""
	switch strings.ToLower(scheme) {
	case "bearer":
		t = "jwt"
	default:
		t = scheme
	}

	cred, err := NewCredential(t, authHeader, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential object: %w", err)
	}

	cred.FromMeta(meta.FromProvider(provider))

	return cred, nil
}
