package principal

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/security"
)

// principalKey is an unexported type for context.Context keys.
type principalKey struct{}

// concretePrincipal is a concrete implementation of the security.Principal interface.
type concretePrincipal struct {
	id          string
	roles       []string
	permissions []string        // New field
	scopes      map[string]bool // New field
	extraClaims map[string]any  // Renamed from 'claims'
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

// New creates a new security.Principal instance.
func New(id string, roles []string, permissions []string, scopes map[string]bool, extraClaims map[string]any) security.Principal {
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

// ToProto converts a security.Principal to a *securityv1.Principal Protobuf message.
// It attempts to pack various Go types into anypb.Any.
func ToProto(p security.Principal) (*securityv1.Principal, error) {
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

// FromProto converts a *securityv1.Principal Protobuf message to a security.Principal.
// It stores anypb.Any directly in the extraClaims map; consumers will need to unpack them.
func FromProto(protoP *securityv1.Principal) (security.Principal, error) {
	if protoP == nil {
		return nil, nil
	}

	extraClaims := make(map[string]any)
	for key, anyValue := range protoP.GetExtraClaims() { // Renamed field
		extraClaims[key] = anyValue
	}

	return New(protoP.GetId(), protoP.GetRoles(), protoP.GetPermissions(), protoP.GetScopes(), extraClaims), nil
}
