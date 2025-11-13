package security

import (
	"cmp"
	"fmt"

	authnv1 "github.com/origadmin/runtime/api/gen/go/config/security/authn/v1"
	authzv1 "github.com/origadmin/runtime/api/gen/go/config/security/authz/v1"
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// AuthenticatorFactory is a function type that creates an Authenticator instance
// from a specific Protobuf configuration message.
type AuthenticatorFactory func(config *authnv1.Authenticator, opts ...options.Option) (declarative.Authenticator, error)

// AuthorizerFactory is a function type that creates an Authorizer instance
// from a specific Protobuf configuration message.
type AuthorizerFactory func(config *authzv1.Authorizer, opts ...options.Option) (declarative.Authorizer, error)

// authenticatorFactories stores all registered Authenticator factories.
// The key is a string identifier for the authenticator type (e.g., "jwt", "apikey").
var authenticatorFactories = make(map[string]AuthenticatorFactory)

// authorizerFactories stores all registered Authorizer factories.
// The key is a string identifier for the authorizer type (e.g., "rbac", "acl").
var authorizerFactories = make(map[string]AuthorizerFactory)

// RegisterAuthenticatorFactory registers an AuthenticatorFactory with a given name.
// This should be called during application initialization (e.g., in an init() function).
func RegisterAuthenticatorFactory(name string, factory AuthenticatorFactory) {
	if _, exists := authenticatorFactories[name]; exists {
		panic(fmt.Sprintf("Authenticator factory with name '%s' already registered", name))
	}
	authenticatorFactories[name] = factory
}

// RegisterAuthorizerFactory registers an AuthorizerFactory with a given name.
// This should be called during application initialization (e.g., in an init() function).
func RegisterAuthorizerFactory(name string, factory AuthorizerFactory) {
	if _, exists := authorizerFactories[name]; exists {
		panic(fmt.Sprintf("Authorizer factory with name '%s' already registered", name))
	}
	authorizerFactories[name] = factory
}

// GetAuthenticators creates an Authenticator instance based on the provided AuthnConfig Protobuf message.
// It inspects the oneof field in AuthnConfig to determine which specific authenticator to create.
func GetAuthenticators(authnConfigs *securityv1.AuthenticatorConfigs, opts ...options.Option) (map[string]declarative.Authenticator, error) {
	if authnConfigs == nil {
		return nil, nil // No config provided, so no authenticator is needed
	}

	configs := authnConfigs.GetConfigs()
	authenticators := make(map[string]declarative.Authenticator)
	for _, config := range configs {
		authenticator, err := GetAuthenticator(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticator '%s': %w", config.GetType(), err)
		}
		key := cmp.Or(config.GetName(), config.GetType())
		authenticators[key] = authenticator
	}
	return authenticators, nil
}

func GetAuthorizers(authzConfigs *securityv1.AuthorizerConfigs, opts ...options.Option) (map[string]declarative.Authorizer, error) {
	if authzConfigs == nil {
		return nil, nil // No config provided, so no authorizer is needed
	}

	configs := authzConfigs.GetConfigs()
	authorizers := make(map[string]declarative.Authorizer)
	for _, config := range configs {
		authorizer, err := GetAuthorizer(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create authorizer '%s': %w", config.GetType(), err)
		}
		key := cmp.Or(config.GetName(), config.GetType())
		authorizers[key] = authorizer
	}
	return authorizers, nil
}

func GetAuthenticator(authnConfig *authnv1.Authenticator, opts ...options.Option) (declarative.Authenticator, error) {
	if authnConfig == nil {
		return nil, nil // No config provided, so no authenticator is needed
	}
	factory, ok := authenticatorFactories[authnConfig.GetType()]
	if !ok {
		return nil, fmt.Errorf("authenticator factory '%s' not registered", authnConfig.GetType())
	}
	return factory(authnConfig)
}

func GetAuthorizer(authzConfig *authzv1.Authorizer, opts ...options.Option) (declarative.Authorizer, error) {
	if authzConfig == nil {
		return nil, nil // No config provided, so no authorizer is needed
	}
	factory, ok := authorizerFactories[authzConfig.GetType()]
	if !ok {
		return nil, fmt.Errorf("authorizer factory '%s' not registered", authzConfig.GetType())
	}
	return factory(authzConfig)
}
