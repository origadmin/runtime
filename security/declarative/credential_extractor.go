/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package declarative

import (
	"context"
	"strings"

	"google.golang.org/protobuf/proto"

	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
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

// Extract is responsible for all extraction and parsing logic. It prepares all
// necessary components and then calls the pure NewCredential constructor.
func (e *HeaderCredentialExtractor) Extract(ctx context.Context, provider declarative.SecurityRequest) (declarative.Credential, error) {
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

	// Prepare all components for the constructor
	var credentialType string
	var payload proto.Message
	switch strings.ToLower(scheme) {
	case "bearer":
		credentialType = "jwt"
		payload = &securityv1.BearerCredential{
			Token: rawCredential,
		}
	default:
		credentialType = scheme
	}

	// Directly get Go-idiomatic metadata from the provider.
	goMeta := provider.GetAll()

	// Call the pure constructor with the final, prepared components.
	// NewCredential is now defined in credential.go
	return NewCredential(credentialType, authHeader, payload, goMeta)
}
