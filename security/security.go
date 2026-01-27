// Copyright 2024 The OrigAdmin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package security

import "sync"

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
	mu              sync.RWMutex

	// policyNameCache stores pre-filtered service methods by policy name for quick lookup.
	// This cache is built once after all init() functions have run.
	policyNameCache     map[string]map[string]struct{}
	policyNameCacheOnce sync.Once
)

// buildPolicyNameCache initializes policyNameCache from unifiedPolicies.
// This function is designed to be called exactly once.
func buildPolicyNameCache() {
	policyNameCache = make(map[string]map[string]struct{})
	for _, p := range unifiedPolicies {
		if _, ok := policyNameCache[p.Name]; !ok {
			policyNameCache[p.Name] = make(map[string]struct{})
		}
		policyNameCache[p.Name][p.ServiceMethod] = struct{}{}
	}
}

// RegisterPolicies is a public function called by generated code in init() functions.
// It appends a slice of policies to the global unifiedPolicies registry.
func RegisterPolicies(policies []Policy) {
	mu.Lock()
	defer mu.Unlock()
	unifiedPolicies = append(unifiedPolicies, policies...)
}

// RegisteredPolicies returns a copy of all policy registrations.
// This is called once at application startup to sync policies to the database (the "Resource").
func RegisteredPolicies() []Policy {
	mu.RLock()
	defer mu.RUnlock()

	// Return a copy to prevent external modifications to the original slice.
	clone := make([]Policy, len(unifiedPolicies))
	copy(clone, unifiedPolicies)
	return clone
}

// GetServiceMethodsByPolicyNames returns a map of service methods that have any of the given policy names.
// This function provides a convenient way to access pre-filtered lists of code-defined policies.
// It's intended for static, compile-time defined policies, not for dynamic, user-editable policies.
func GetServiceMethodsByPolicyNames(names ...string) map[string]struct{} {
	policyNameCacheOnce.Do(func() {
		mu.RLock() // Lock for reading unifiedPolicies during cache build
		defer mu.RUnlock()
		buildPolicyNameCache()
	})

	mu.RLock() // Lock for reading policyNameCache
	defer mu.RUnlock()

	result := make(map[string]struct{})
	for _, name := range names {
		if methods, ok := policyNameCache[name]; ok {
			// Copy methods to the result map
			for k, v := range methods {
				result[k] = v
			}
		}
	}
	return result
}
