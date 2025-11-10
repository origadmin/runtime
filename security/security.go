// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"sync"
)

// --- Global Policy Registry (for init() function registration) ---

var (
	// servicePolicies stores policies generated from .proto annotations for gRPC methods.
	// The key is the full gRPC method name (e.g., /helloworld.v1.Greeter/SayHello),
	// and the value is the policy name (e.g., "jwt-auth").
	servicePolicies = make(map[string]string)
	// gatewayPolicies stores policies generated from .proto for HTTP paths.
	gatewayPolicies = make(map[string]string)
	mu              sync.RWMutex
)

// RegisterPolicies is a public function, intended to be called by generated code in init().
// It merges the incoming policies into the global default maps.
func RegisterPolicies(serviceMap, gatewayMap map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	// Merge policies for gRPC methods. The full RPC method name is the canonical key.
	for method, policy := range serviceMap {
		servicePolicies[method] = policy
	}

	// Merge policies for HTTP gateway paths. The key is typically "METHOD:/path/template".
	for path, policy := range gatewayMap {
		gatewayPolicies[path] = policy
	}
}

// GetDefaultServicePolicies returns a copy of the currently registered default policies for gRPC services.
// This is primarily for the PolicyProvider to use during its initialization or sync process.
func GetDefaultServicePolicies() map[string]string {
	mu.RLock()
	defer mu.RUnlock()
	// Return a copy to prevent external modification
	copiedMap := make(map[string]string, len(servicePolicies))
	for k, v := range servicePolicies {
		copiedMap[k] = v
	}
	return copiedMap
}

// GetDefaultGatewayPolicies returns a copy of the currently registered default policies for HTTP gateway paths.
func GetDefaultGatewayPolicies() map[string]string {
	mu.RLock()
	defer mu.RUnlock()
	copiedMap := make(map[string]string, len(gatewayPolicies))
	for k, v := range gatewayPolicies {
		copiedMap[k] = v
	}
	return copiedMap
}
