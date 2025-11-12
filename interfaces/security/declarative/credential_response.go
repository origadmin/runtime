/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides declarative security interfaces for authentication and authorization.
package declarative

import (
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/metadata"
)

// CredentialResponse represents a credential structure intended for
// transmission to clients (e.g., frontend applications). It contains
// processed and safe-to-expose credential information.
type CredentialResponse interface {
	// GetType returns the type of the credential.
	GetType() string

	// Payload returns the payload of the credential.
	Payload() *securityv1.Payload

	// GetMeta returns the metadata associated with the credential as a map.
	// The values are converted from google.protobuf.Value to Go's any type.
	GetMeta() metadata.Meta

	// ToProto converts the CredentialResponse to its protobuf representation.
	// This allows direct access to the underlying protobuf message, including its oneof fields.
	ToProto() *securityv1.CredentialResponse

	// ToJSON serializes the entire CredentialResponse to a JSON string.
	ToJSON() (string, error)

	// MarshalPayload serializes only the value within the Payload field to JSON.
	// This provides direct access to the core payload's JSON representation.
	MarshalPayload() ([]byte, error)
}
