/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides the core components for the declarative security framework.
package declarative

import (
	"sync"

	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// Provider is a thread-safe implementation of declarative.PolicyProvider.
// It holds registered policy factories and created policy instances.
type Provider struct {
	mu        sync.RWMutex
	factories map[string]declarative.SecurityFactory
	policies  map[string]declarative.SecurityPolicy
	configs   map[string][]byte
}

// NewProvider creates a new security policy provider.
func NewProvider() *Provider {
	return &Provider{
		factories: make(map[string]declarative.SecurityFactory),
		policies:  make(map[string]declarative.SecurityPolicy),
		configs:   make(map[string][]byte),
	}
}

// RegisterFactory registers a SecurityFactory for a given policy name.
// This is typically called during application initialization.
func (p *Provider) RegisterFactory(name string, factory declarative.SecurityFactory) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.factories[name] = factory
}

// AddConfig adds the raw configuration for a policy.
// This is used to provide configuration data to the factory when the policy is first requested.
func (p *Provider) AddConfig(name string, config []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.configs[name] = config
}

// GetPolicy retrieves a SecurityPolicy by its name.
// It uses a lazy-initialization approach: the policy is created by its factory
// on the first request and then cached for subsequent requests.
func (p *Provider) GetPolicy(name string) (declarative.SecurityPolicy, error) {
	p.mu.RLock()
	policy, ok := p.policies[name]
	p.mu.RUnlock()
	if ok {
		return policy, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check in case it was created while waiting for the lock.
	policy, ok = p.policies[name]
	if ok {
		return policy, nil
	}

	factory, ok := p.factories[name]
	if !ok {
		return nil, runtimeerrors.NewStructured(
			"SECURITY_POLICY_NOT_FOUND",
			"security policy factory not found for: %s", name,
		)
	}

	config, _ := p.configs[name] // config can be nil if not provided

	newPolicy, err := factory(config)
	if err != nil {
		return nil, runtimeerrors.NewStructured(
			"SECURITY_POLICY_CREATION_FAILED",
			"failed to create security policy '%s': %s", name, err.Error(),
		)
	}

	p.policies[name] = newPolicy
	return newPolicy, nil
}
