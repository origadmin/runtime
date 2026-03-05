/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"context"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
)

// Provider defines the interface for a middleware service provider.
type Provider interface {
	Middleware(name string) (KMiddleware, error)
	DefaultMiddleware() (KMiddleware, error)
}

// providerImpl implements the Provider interface.
type providerImpl struct {
	handle component.Handle
}

func (p *providerImpl) Middleware(name string) (KMiddleware, error) {
	return engine.Get[KMiddleware](context.Background(), p.handle.In(component.CategoryMiddleware), name)
}

func (p *providerImpl) DefaultMiddleware() (KMiddleware, error) {
	return engine.GetDefault[KMiddleware](context.Background(), p.handle.In(component.CategoryMiddleware))
}

// NewProvider creates a new middleware provider instance.
func NewProvider(handle component.Handle) Provider {
	return &providerImpl{handle: handle}
}

// ServerProvider is the engine-compatible provider for server-side middleware components.
var ServerProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	if h.Scope() != component.ServerScope {
		return nil, nil // Let other providers handle it
	}
	cfg, err := engine.AsConfig[middlewarev1.Middleware](h)
	if err != nil {
		return nil, err
	}
	m, ok := NewServer(cfg, opts...)
	if !ok {
		return nil, nil
	}
	return m, nil
}

// ClientProvider is the engine-compatible provider for client-side middleware components.
var ClientProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	if h.Scope() != component.ClientScope {
		return nil, nil // Let other providers handle it
	}
	cfg, err := engine.AsConfig[middlewarev1.Middleware](h)
	if err != nil {
		return nil, err
	}
	m, ok := NewClient(cfg, opts...)
	if !ok {
		return nil, nil
	}
	return m, nil
}

// DefaultProvider acts as a catch-all dispatcher for legacy or non-scoped registrations.
var DefaultProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	cfg, err := engine.AsConfig[middlewarev1.Middleware](h)
	if err != nil {
		return nil, err
	}

	if h.Scope() == component.ClientScope {
		m, ok := NewClient(cfg, opts...)
		if ok {
			return m, nil
		}
	} else {
		m, ok := NewServer(cfg, opts...)
		if ok {
			return m, nil
		}
	}
	return nil, nil
}
