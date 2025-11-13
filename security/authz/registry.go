package authz

import (
	"cmp"
	"fmt"

	authzv1 "github.com/origadmin/runtime/api/gen/go/config/security/authz/v1"
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// AuthorizerFactory is a function type that creates an Authorizer instance
// from a specific Protobuf configuration message.
type AuthorizerFactory func(config *authzv1.Authorizer, opts ...options.Option) (declarative.Authorizer, error)

// authorizerFactories stores all registered Authorizer factories.
// The key is a string identifier for the authorizer type (e.g., "rbac", "acl").
var authorizerFactories = make(map[string]AuthorizerFactory)

// RegisterAuthorizerFactory registers an AuthorizerFactory with a given name.
// This should be called during application initialization (e.g., in an init() function).
func RegisterAuthorizerFactory(name string, factory AuthorizerFactory) {
	if _, exists := authorizerFactories[name]; exists {
		panic(fmt.Sprintf("Authorizer factory with name '%s' already registered", name))
	}
	authorizerFactories[name] = factory
}

func NewAuthorizers(authzConfigs *securityv1.AuthorizerConfigs, opts ...options.Option) (map[string]declarative.Authorizer, error) {
	if authzConfigs == nil {
		return nil, nil // No config provided, so no authorizer is needed
	}

	configs := authzConfigs.GetConfigs()
	authorizers := make(map[string]declarative.Authorizer)
	for _, config := range configs {
		authorizer, err := NewAuthorizer(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create authorizer '%s': %w", config.GetType(), err)
		}
		key := cmp.Or(config.GetName(), config.GetType())
		authorizers[key] = authorizer
	}
	return authorizers, nil
}

func NewAuthorizer(authzConfig *authzv1.Authorizer, opts ...options.Option) (declarative.Authorizer, error) {
	if authzConfig == nil {
		return nil, nil // No config provided, so no authorizer is needed
	}
	factory, ok := authorizerFactories[authzConfig.GetType()]
	if !ok {
		return nil, fmt.Errorf("authorizer factory '%s' not registered", authzConfig.GetType())
	}
	return factory(authzConfig)
}
