/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package declarative

import (
	"fmt"
	"strings"

	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/security/declarative"
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
	authHeaders := provider.Get(AuthorizationHeader)
	if len(authHeaders) == 0 {
		return nil, errors.New(401, "AUTHORIZATION_HEADER_MISSING", "authorization header is missing")
	}

	// In case of multiple Authorization headers, RFC 7235 states they should not be combined.
	// We will process the first one found, as this is the most common use case.
	authHeader := authHeaders[0]

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
	cred, err := NewCredential(scheme, rawCredential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential object: %w", err)
	}

	return cred, nil
}

