// Package declarative implements the functions, types, and interfaces for the module.
package declarative

import (
	"google.golang.org/protobuf/proto" // Import proto for proto.Message
)

// Credential represents a credential, either received from a request or newly issued.
// It provides a unified interface to access credential data and supports various
// serialization formats like JSON and Protobuf.
type Credential interface {
	// Type returns the type of the credential (e.g., "jwt", "apikey").
	Type() string

	// Raw returns the original, unparsed credential string.
	// For example, the full "Bearer eyJ..." JWT string, or the API key string.
	Raw() string

	// ParsedPayload returns the parsed credential payload as a google.protobuf.Any message.
	// This allows for type-safe unmarshalling into specific protobuf messages.
	ParsedPayload(message proto.Message) error

	// Get returns the value of a specific header/metadata key.
	Get(key string) (string, bool)

	// GetAll returns all available headers/metadata as a map.
	GetAll() map[string]any
}
