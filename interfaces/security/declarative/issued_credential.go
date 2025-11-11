/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package declarative

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

// IssuedCredential represents a newly created, serializable credential.
// It provides a standard way to access the credential's data for various
// serialization formats like JSON and Protobuf.
type IssuedCredential interface {
	// Type returns the specific type of the credential, e.g., "jwt", "apikey".
	Type() string

	// Marshaler implements json.Marshaler for direct JSON serialization.
	// This is used for responding to HTTP/JSON APIs.
	json.Marshaler

	// ToProto converts the credential into its corresponding Protobuf message
	// for inter-service communication.
	ToProto() (proto.Message, error)
}
