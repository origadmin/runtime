/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package authz

import (
	"fmt"

	authzv1 "github.com/origadmin/runtime/api/gen/go/config/security/authz/v1"
	"github.com/origadmin/runtime/interfaces/options"
	internalfactory "github.com/origadmin/runtime/internal/factory"
)

var (
	defaultFactories = internalfactory.New[Factory]()
)

// Register registers a new authorizer provider factory.
// This function is intended to be called from the init() function of each provider implementation.
func Register(name string, factory Factory) {
	if _, ok := defaultFactories.Get(name); ok {
		panic(fmt.Sprintf("authorizer factory %q already registered", name))
	}
	defaultFactories.Register(name, factory)
}

// Create creates a new authorizer provider instance based on the given configuration.
// It looks up the appropriate factory using the type specified in the config and invokes it.
// The returned Provider instance is NOT stored globally; it is the caller's responsibility
// to manage its lifecycle and inject it where needed.
func Create(cfg *authzv1.Authorizer, opts ...options.Option) (Provider, error) {
	factory, ok := defaultFactories.Get(cfg.GetType())
	if !ok {
		return nil, fmt.Errorf("authorizer factory %q not found", cfg.GetType())
	}
	return factory(cfg, opts...)
}
