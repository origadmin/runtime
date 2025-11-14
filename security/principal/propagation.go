package principal

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"

	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/interfaces/security"
)

const (
	// MetadataKey is the key used to store the Principal in gRPC metadata or HTTP headers.
	MetadataKey = "x-principal-proto"
)

// FromContext extracts the Principal from the given context.
// It returns the Principal and a boolean indicating if it was found.
func FromContext(ctx context.Context) (security.Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(security.Principal)
	return p, ok
}

// WithContext returns a new context with the given Principal attached.
// It is used to inject the Principal into the context for downstream business logic.
func WithContext(ctx context.Context, p security.Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}

// EncodePrincipal encodes a security.Principal into a base64-encoded Protobuf string.
func EncodePrincipal(p security.Principal) (string, error) {
	if p == nil {
		return "", nil
	}
	protoP, err := ToProto(p)
	if err != nil {
		return "", fmt.Errorf("failed to convert security.Principal to proto: %w", err)
	}
	data, err := proto.Marshal(protoP)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proto.Principal: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodePrincipal decodes a base64-encoded Protobuf string into a security.Principal.
func DecodePrincipal(encoded string) (security.Principal, error) {
	if encoded == "" {
		return nil, nil
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %w", err)
	}
	protoP := &securityv1.Principal{}
	if err := proto.Unmarshal(data, protoP); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto.Principal: %w", err)
	}
	return FromProto(protoP)
}
