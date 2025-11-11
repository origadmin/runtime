// Package declarative implements the functions, types, and interfaces for the module.
package declarative

import (
	"github.com/origadmin/runtime/context"
	"google.golang.org/protobuf/proto" // Import proto for proto.Message
)

// CredentialParser defines the contract for parsing and validating a credential
// and converting it into a Principal. It acts as a reusable engine for handling
// token-specific technical details.
type CredentialParser interface {
	// ParseCredential converts a credential into a validated Principal object.
	// It is responsible for all technical validation of the credential itself
	// (e.g., signature, expiration, format).
	ParseCredential(ctx context.Context, rawToken string) (Credential, error)

	// ParseCredentialFrom parses a credential from a source of type sourceType.
	// The source can be a string, map, or any other type that the parser supports.
	// sourceType specifies the type of the source (e.g., "header", "query", "body").
	// source is the actual data to be parsed.
	ParseCredentialFrom(ctx context.Context, sourceType string, source any) (Credential, error)
}
