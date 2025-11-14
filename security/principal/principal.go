package principal

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/origadmin/runtime/interfaces/security/declarative"
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// MetadataKey is the key used to store the Principal in gRPC metadata or HTTP headers.
	MetadataKey = "x-principal-proto"
)

// principalKey is an unexported type for context.Context keys.
type principalKey struct{}

// concretePrincipal is a concrete implementation of the declarative.Principal interface.
type concretePrincipal struct {
	id          string
	roles       []string
	permissions []string // New field
	scopes      map[string]bool // New field
	extraClaims map[string]any // Renamed from 'claims'
}

// GetID returns the unique identifier of the principal.
func (p *concretePrincipal) GetID() string {
	return p.id
}

// GetRoles returns a list of roles assigned to the principal.
func (p *concretePrincipal) GetRoles() []string {
	return p.roles
}

// GetPermissions returns a list of permissions assigned to the principal.
func (p *concretePrincipal) GetPermissions() []string {
	return p.permissions
}

// GetScopes returns a map of scopes assigned to the principal.
func (p *concretePrincipal) GetScopes() map[string]bool {
	return p.scopes
}

// GetClaims returns a map of all extra claims associated with the principal.
func (p *concretePrincipal) GetClaims() map[string]any {
	return p.extraClaims
}

// New creates a new declarative.Principal instance.
func New(id string, roles []string, permissions []string, scopes map[string]bool, extraClaims map[string]any) declarative.Principal {
	if extraClaims == nil {
		extraClaims = make(map[string]any)
	}
	if scopes == nil {
		scopes = make(map[string]bool)
	}
	return &concretePrincipal{
		id:          id,
		roles:       roles,
		permissions: permissions,
		scopes:      scopes,
		extraClaims: extraClaims,
	}
}

// ToProto converts a declarative.Principal to a *securityv1.Principal Protobuf message.
// It attempts to pack various Go types into anypb.Any.
func ToProto(p declarative.Principal) (*securityv1.Principal, error) {
	if p == nil {
		return nil, nil
	}

	protoExtraClaims := make(map[string]*anypb.Any)
	for key, value := range p.GetClaims() { // GetClaims now returns extraClaims
		var anyValue *anypb.Any
		var err error

		switch v := value.(type) {
		case proto.Message: // Already a protobuf message
			anyValue, err = anypb.New(v)
		case string:
			anyValue, err = anypb.New(wrapperspb.String(v))
		case int32:
			anyValue, err = anypb.New(wrapperspb.Int32(v))
		case int64:
			anyValue, err = anypb.New(wrapperspb.Int64(v))
		case bool:
			anyValue, err = anypb.New(wrapperspb.Bool(v))
		case float32:
			anyValue, err = anypb.New(wrapperspb.Float(v))
		case float64:
			anyValue, err = anypb.New(wrapperspb.Double(v))
		default:
			return nil, fmt.Errorf("unsupported claim type for key '%s': %T, cannot pack into anypb.Any", key, v)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to pack claim '%s' into anypb.Any: %w", key, err)
		}
		protoExtraClaims[key] = anyValue
	}

	return &securityv1.Principal{
		Id:          p.GetID(),
		Roles:       p.GetRoles(),
		Permissions: p.GetPermissions(), // New field
		Scopes:      p.GetScopes(),      // New field
		ExtraClaims: protoExtraClaims,   // Renamed field
	}, nil
}

// FromProto converts a *securityv1.Principal Protobuf message to a declarative.Principal.
// It stores anypb.Any directly in the extraClaims map; consumers will need to unpack them.
func FromProto(protoP *securityv1.Principal) (declarative.Principal, error) {
	if protoP == nil {
		return nil, nil
	}

	extraClaims := make(map[string]any)
	for key, anyValue := range protoP.GetExtraClaims() { // Renamed field
		extraClaims[key] = anyValue
	}

	return New(protoP.GetId(), protoP.GetRoles(), protoP.GetPermissions(), protoP.GetScopes(), extraClaims), nil
}

// PrincipalFromContext extracts the Principal from the given context.
// It returns the Principal and a boolean indicating if it was found.
func PrincipalFromContext(ctx context.Context) (declarative.Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(declarative.Principal)
	return p, ok
}

// PrincipalWithContext returns a new context with the given Principal attached.
// It is used to inject the Principal into the context for downstream business logic.
func PrincipalWithContext(ctx context.Context, p declarative.Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}

// EncodePrincipal encodes a declarative.Principal into a base64-encoded Protobuf string.
func EncodePrincipal(p declarative.Principal) (string, error) {
	if p == nil {
		return "", nil
	}
	protoP, err := ToProto(p)
	if err != nil {
		return "", fmt.Errorf("failed to convert declarative.Principal to proto: %w", err)
	}
	data, err := proto.Marshal(protoP)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proto.Principal: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodePrincipal decodes a base64-encoded Protobuf string into a declarative.Principal.
func DecodePrincipal(encoded string) (declarative.Principal, error) {
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
