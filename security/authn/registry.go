package security

import (
	"cmp"
	"fmt"

	authnv1 "github.com/origadmin/runtime/api/gen/go/config/security/authn/v1"
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// AuthenticatorFactory is a function type that creates an Authenticator instance
// from a specific Protobuf configuration message.
type AuthenticatorFactory func(config *authnv1.Authenticator, opts ...options.Option) (declarative.Authenticator, error)

// authenticatorFactories stores all registered Authenticator factories.
// The key is a string identifier for the authenticator type (e.g., "jwt", "apikey").
var authenticatorFactories = make(map[string]AuthenticatorFactory)

// RegisterAuthenticatorFactory registers an AuthenticatorFactory with a given name.
// This should be called during application initialization (e.g., in an init() function).
func RegisterAuthenticatorFactory(name string, factory AuthenticatorFactory) {
	if _, exists := authenticatorFactories[name]; exists {
		panic(fmt.Sprintf("Authenticator factory with name '%s' already registered", name))
	}
	authenticatorFactories[name] = factory
}

// NewAuthenticators creates an Authenticator instance based on the provided AuthnConfig Protobuf message.
// It inspects the oneof field in AuthnConfig to determine which specific authenticator to create.
func NewAuthenticators(authnConfigs *securityv1.AuthenticatorConfigs, opts ...options.Option) (map[string]declarative.Authenticator, error) {
	if authnConfigs == nil {
		return nil, nil // No config provided, so no authenticator is needed
	}

	configs := authnConfigs.GetConfigs()
	authenticators := make(map[string]declarative.Authenticator)
	for _, config := range configs {
		authenticator, err := NewAuthenticator(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticator '%s': %w", config.GetType(), err)
		}
		key := cmp.Or(config.GetName(), config.GetType())
		authenticators[key] = authenticator
	}
	return authenticators, nil
}

func NewAuthenticator(authnConfig *authnv1.Authenticator, opts ...options.Option) (declarative.Authenticator, error) {
	if authnConfig == nil {
		return nil, nil // No config provided, so no authenticator is needed
	}
	factory, ok := authenticatorFactories[authnConfig.GetType()]
	if !ok {
		return nil, fmt.Errorf("authenticator factory '%s' not registered", authnConfig.GetType())
	}
	return factory(authnConfig)
}
