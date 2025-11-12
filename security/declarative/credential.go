package declarative

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/metadata"
	"github.com/origadmin/runtime/interfaces/security/declarative"
	"github.com/origadmin/runtime/security/meta"
)

// credential is the internal implementation of the Credential interface.
type credential struct {
	c *securityv1.CredentialSource
	m meta.Meta
}

func (c *credential) ToJSON() (string, error) {
	if c.c == nil {
		return "", fmt.Errorf("credential source is nil")
	}
	if c.c.Metadata == nil {
		c.c.Metadata = c.m.ToProto()
	}
	jsonStr, err := protojson.Marshal(c.c)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

func (c *credential) ToProto() ([]byte, error) {
	if c.c == nil {
		return nil, fmt.Errorf("credential source is nil")
	}
	if c.c.Metadata == nil {
		c.c.Metadata = c.m.ToProto()
	}
	return protojson.Marshal(c.c)
}

func (c *credential) FromMeta(meta metadata.Meta) {
	if c == nil || meta == nil {
		return
	}
	c.m = meta.GetAll()
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
	metaValue := c.m.Get(key)
	if metaValue == "" {
		return "", false
	}
	return "", false
}

// GetAll returns all available metadata as a map.
// The values in the map are unmarshalled from their protobuf Any type to their corresponding Go types.
func (c *credential) GetAll() map[string][]string {
	if c == nil || c.m == nil {
		return nil
	}
	return c.m
}

// Type returns the type of the credential.
func (c *credential) Type() string {
	return c.c.GetType()
}

// Raw returns the original, unparsed credential string.
func (c *credential) Raw() string {
	return c.c.GetRaw()
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

// credentialResponse is the internal implementation of the CredentialResponse interface.
type credentialResponse struct {
	cr *securityv1.CredentialResponse
	m  meta.Meta
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
func (c *credentialResponse) GetMeta() metadata.Meta {
	if c == nil || c.m == nil {
		return nil
	}
	return c.m
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
