/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides the core components for the declarative security framework.
package declarative

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// ErrPolicyNotFound is returned when a requested security policy is not found.
var ErrPolicyNotFound = errors.New(500, "SECURITY_POLICY_NOT_FOUND", "security policy not found")

// Manager is a thread-safe manager for security policies.
// It holds registered policy factories and created policy instances.
type Manager struct {
	mu        sync.RWMutex
	factories map[string]declarative.SecurityFactory
	policies  map[string]declarative.SecurityPolicy
	configs   map[string][]byte
}

// NewManager creates a new security policy manager.
func NewManager() *Manager {
	return &Manager{
		factories: make(map[string]declarative.SecurityFactory),
		policies:  make(map[string]declarative.SecurityPolicy),
		configs:   make(map[string][]byte),
	}
}

// RegisterFactory registers a SecurityFactory for a given policy name.
// This is typically called during application initialization.
func (m *Manager) RegisterFactory(name string, factory declarative.SecurityFactory) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.factories[name] = factory
}

// AddConfig adds the raw configuration for a policy.
// This is used to provide configuration data to the factory when the policy is first requested.
func (m *Manager) AddConfig(name string, config []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configs[name] = config
}

// GetPolicy retrieves a SecurityPolicy by its name.
// It uses a lazy-initialization approach: the policy is created by its factory
// on the first request and then cached for subsequent requests.
func (m *Manager) GetPolicy(name string) (declarative.SecurityPolicy, error) {
	m.mu.RLock()
	policy, ok := m.policies[name]
	m.mu.RUnlock()
	if ok {
		return policy, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check in case it was created while waiting for the lock.
	policy, ok = m.policies[name]
	if ok {
		return policy, nil
	}

	factory, ok := m.factories[name]
	if !ok {
		return nil, fmt.Errorf("%w: no factory registered for policy '%s'", ErrPolicyNotFound, name)
	}

	config, _ := m.configs[name] // config can be nil if not provided

	newPolicy, err := factory(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy '%s': %w", name, err)
	}

	m.policies[name] = newPolicy
	return newPolicy, nil
}
