// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// credential is the internal implementation of the Credential interface.
type credential struct {
	c *securityv1.CredentialSource
}

// ParsedPayload unmarshals the credential's payload into the provided protobuf message.
// It returns an error if the unmarshalling fails.
func (c *credential) ParsedPayload(message proto.Message) error {
	if c.c == nil || c.c.Payload == nil {
		return fmt.Errorf("credential source or payload is nil")
	}
	return c.c.Payload.UnmarshalTo(message)
}

// Get returns the string value of a specific metadata key.
// It returns the value and true if the key exists and its value is a string, otherwise it returns an empty string and false.
func (c *credential) Get(key string) (string, bool) {
	meta := c.c.GetMeta()
	if meta == nil {
		return "", false
	}
	if v, ok := meta[key]; ok {
		stringValue := new(structpb.Value)
		if v.MessageIs(stringValue) {
			err := v.UnmarshalTo(stringValue)
			if err != nil {
				return "", false
			}
			return stringValue.GetStringValue(), true
		}
	}
	return "", false
}

// GetAll returns all available metadata as a map.
// The values in the map are unmarshalled from their protobuf Any type to their corresponding Go types.
func (c *credential) GetAll() map[string]any {
	meta := c.c.GetMeta()
	if meta == nil {
		return nil
	}
	result := make(map[string]any)
	for k, v := range meta {
		var val structpb.Value
		if err := v.UnmarshalTo(&val); err != nil {
			log.Printf("failed to unmarshal metadata key %s: %v", k, err) // Log the error
			continue
		}
		switch kind := val.Kind.(type) {
		case *structpb.Value_StringValue:
			result[k] = kind.StringValue
		case *structpb.Value_NumberValue:
			result[k] = kind.NumberValue
		case *structpb.Value_BoolValue:
			result[k] = kind.BoolValue
		case *structpb.Value_ListValue:
			// Handle list values if necessary, potentially recursively
			var list []any
			for _, item := range kind.ListValue.Values {
				// This is a simplified conversion, might need a more robust one for nested structures
				list = append(list, item.AsInterface())
			}
			result[k] = list
		case *structpb.Value_StructValue:
			// Handle struct values if necessary, potentially recursively
			result[k] = kind.StructValue.AsMap()
		default:
			// For other types or unknown types, we can just store the interface representation
			result[k] = val.AsInterface()
		}
	}
	return result
}

// NewCredential creates a new Credential instance from a securityv1.Credential protobuf message.
func NewCredential(typo, raw string, src proto.Message) (declarative.Credential, error) {
	payload, err := anypb.New(src)
	if err != nil {
		return nil, err
	}
	return &credential{c: &securityv1.CredentialSource{
		Type:    typo,
		Raw:     raw,
		Payload: payload,
	}}, nil
}

// Type returns the type of the credential.
func (c *credential) Type() string {
	return c.c.GetType()
}

// Raw returns the original, unparsed credential string.
func (c *credential) Raw() string {
	return c.c.GetRaw()
}

// credentialResponse is the internal implementation of the CredentialResponse interface.
type credentialResponse struct {
	cr *securityv1.CredentialResponse
}

func (c *credentialResponse) Payload() *securityv1.Payload {
	return c.cr.GetPayload()
}

// GetType returns the type of the credential.
func (c *credentialResponse) GetType() string {
	if c.cr == nil {
		return ""
	}
	return c.cr.GetType()
}

// GetMeta returns the metadata associated with the credential as a map.
// The values are converted from google.protobuf.Value to Go's any type.
func (c *credentialResponse) GetMeta() map[string]any {
	if c.cr == nil || c.cr.GetMeta() == nil {
		return nil
	}

	result := make(map[string]any)
	for k, v := range c.cr.GetMeta() {
		result[k] = v.AsInterface()
	}
	return result
}

// ToProto converts the CredentialResponse to its protobuf representation.
// This allows direct access to the underlying protobuf message, including its oneof fields.
func (c *credentialResponse) ToProto() *securityv1.CredentialResponse {
	return c.cr
}

// ToJSON serializes the entire CredentialResponse to a JSON string.
func (c *credentialResponse) ToJSON() (string, error) {
	if c.cr == nil {
		return "{}", nil
	}
	b, err := protojson.Marshal(c.cr)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CredentialResponse to JSON: %w", err)
	}
	return string(b), nil
}

// MarshalPayload serializes only the value within the Payload field to JSON.
func (c *credentialResponse) MarshalPayload() ([]byte, error) {
	if c.cr == nil || c.cr.GetPayload() == nil {
		return []byte("null"), nil
	}

	payload := c.cr.GetPayload()
	var msg proto.Message

	// Check each optional field in Payload and assign to msg if present
	if payload.Basic != nil {
		msg = payload.Basic
	} else if payload.Key != nil {
		msg = payload.Key
	} else if payload.Oidc != nil {
		msg = payload.Oidc
	} else if payload.Token != nil {
		msg = payload.Token
	} else if payload.RawData != nil {
		// If RawData is present, marshal it directly as a string
		return []byte(payload.GetRawData()), nil
	}

	if msg != nil {
		return protojson.Marshal(msg)
	}

	return []byte("null"), nil
}

// NewCredentialResponse creates a CredentialResponse from a protobuf message.
func NewCredentialResponse(typo string, pb *securityv1.Payload) declarative.CredentialResponse {
	return &credentialResponse{
		cr: &securityv1.CredentialResponse{
			Type:    typo,
			Payload: pb,
		},
	}
}
