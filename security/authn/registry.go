package authn

import (
	"cmp"
	"fmt"
	"sync"

	authnv1 "github.com/origadmin/runtime/api/gen/go/config/security/authn/v1"
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/security"
)

// Provider is an interface for a security component that can provide various security capabilities.
// This allows a single configured component to expose multiple, separate interfaces (like Authenticator, CredentialCreator, etc.)
// in a type-safe manner.
type Provider interface {
	GetAuthenticator() (security.Authenticator, bool)
	GetCredentialCreator() (security.CredentialCreator, bool)
	GetCredentialRevoker() (security.CredentialRevoker, bool)
}

// Factory is a function type that creates a Provider instance.
type Factory func(config *authnv1.Authenticator, opts ...options.Option) (Provider, error)

// registry is the central store for authenticator factories and provider instances.
type registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
	providers map[string]Provider
}

var defaultRegistry = &registry{
	factories: make(map[string]Factory),
	providers: make(map[string]Provider),
}

// RegisterAuthenticatorFactory registers a Provider Factory with a given name.
func RegisterAuthenticatorFactory(name string, factory Factory) {
	defaultRegistry.mu.Lock()
	defer defaultRegistry.mu.Unlock()
	if _, exists := defaultRegistry.factories[name]; exists {
		panic(fmt.Sprintf("Authenticator factory with name '%s' already registered", name))
	}
	defaultRegistry.factories[name] = factory
}

// BuildProviders creates and stores provider instances based on the provided configuration.
func BuildProviders(authnConfigs *securityv1.AuthenticatorConfigs, opts ...options.Option) error {
	if authnConfigs == nil {
		return nil // No config provided.
	}

	defaultRegistry.mu.Lock()
	defer defaultRegistry.mu.Unlock()

	for _, config := range authnConfigs.GetConfigs() {
		factory, ok := defaultRegistry.factories[config.GetType()]
		if !ok {
			return fmt.Errorf("authenticator factory '%s' not registered", config.GetType())
		}
		provider, err := factory(config, opts...)
		if err != nil {
			return fmt.Errorf("failed to create provider for '%s': %w", config.GetType(), err)
		}
		key := cmp.Or(config.GetName(), config.GetType())
		defaultRegistry.providers[key] = provider
	}
	return nil
}

// GetAuthenticator returns a previously created authenticator instance by name.
func GetAuthenticator(name string) (security.Authenticator, bool) {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()
	provider, ok := defaultRegistry.providers[name]
	if !ok {
		return nil, false
	}
	return provider.GetAuthenticator()
}

// GetCredentialCreator returns the CredentialCreator for the given authenticator name.
func GetCredentialCreator(name string) (security.CredentialCreator, error) {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()
	provider, ok := defaultRegistry.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	creator, ok := provider.GetCredentialCreator()
	if !ok {
		return nil, fmt.Errorf("provider '%s' does not support credential creation", name)
	}
	return creator, nil
}

// GetCredentialRevoker returns the CredentialRevoker for the given authenticator name.
func GetCredentialRevoker(name string) (security.CredentialRevoker, error) {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()
	provider, ok := defaultRegistry.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	revoker, ok := provider.GetCredentialRevoker()
	if !ok {
		return nil, fmt.Errorf("provider '%s' does not support credential revocation", name)
	}
	return revoker, nil
}
