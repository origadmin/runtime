package security

import (
	"context"
	"sync"

	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// Policy holds all information for a single resource's policy.
// This struct is created by generated code and registered via init().
type Policy struct {
	ServiceMethod string // gRPC full method name, e.g., "/user.v1.UserService/GetUser"
	GatewayPath   string // HTTP path and method, e.g., "GET:/api/v1/users/{id}"
	Name          string // The policy name/definition string from the proto annotation, e.g., "admin-only"
	VersionID     string // A hash representing the version of this policy definition
}

// --- Global Policy Registry ---

var (
	// unifiedPolicies stores all policy registrations from generated code.
	// It's populated by init() functions via the RegisterPolicies function.
	unifiedPolicies []Policy

	mu sync.RWMutex
)

// contextKey is an unexported type for context keys to avoid collisions.
type contextKey string

const (
	principalContextKey contextKey = "principal"
)

// ContextWithPrincipal returns a new context with the given Principal attached.
func ContextWithPrincipal(ctx context.Context, p declarative.Principal) context.Context {
	return context.WithValue(ctx, principalContextKey, p)
}

// PrincipalFromContext extracts the Principal from the context.
// It returns the Principal and a boolean indicating whether a Principal was found.
func PrincipalFromContext(ctx context.Context) (declarative.Principal, bool) {
	p, ok := ctx.Value(principalContextKey).(declarative.Principal)
	return p, ok
}

// RegisterPolicies is a public function called by generated code in init() functions.
// It appends a slice of policies to the global unifiedPolicies registry.
func RegisterPolicies(policies []Policy) {
	mu.Lock()
	defer mu.Unlock()
	unifiedPolicies = append(unifiedPolicies, policies...)
}

// RegisteredPolicies returns a copy of all policy registrations.
// This is called once at application startup to sync policies to the database.
func RegisteredPolicies() []Policy {
	mu.RLock()
	defer mu.RUnlock()
	
	clone := make([]Policy, len(unifiedPolicies))
	copy(clone, unifiedPolicies)
	return clone
}
