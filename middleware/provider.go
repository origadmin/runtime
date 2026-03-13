/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"context"
	"errors"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
)

// Provider defines the interface for a middleware service provider.
type Provider interface {
	Middleware(name string) (KMiddleware, error)
	DefaultMiddleware() (KMiddleware, error)
}

// providerImpl implements the Provider interface.
type providerImpl struct {
	locator component.Locator
}

func (p *providerImpl) Middleware(name string) (KMiddleware, error) {
	return comp.Get[KMiddleware](context.Background(), p.locator.In(component.CategoryMiddleware), name)
}

func (p *providerImpl) DefaultMiddleware() (KMiddleware, error) {
	return comp.GetDefault[KMiddleware](context.Background(), p.locator.In(component.CategoryMiddleware))
}

// GetMiddlewareList collects all middlewares from the given locator as a slice.
func GetMiddlewareList(ctx context.Context, locator component.Locator) ([]KMiddleware, error) {
	var mws []KMiddleware
	for _, inst := range locator.Iter(ctx) {
		if m, ok := inst.(KMiddleware); ok {
			mws = append(mws, m)
		}
	}
	return mws, nil
}

// GetMiddlewares collects all available middlewares from the given locator as a map.
func GetMiddlewares(ctx context.Context, locator component.Locator) (map[string]KMiddleware, error) {
	mws := make(map[string]KMiddleware)
	for name, inst := range locator.Iter(ctx) {
		if m, ok := inst.(KMiddleware); ok {
			mws[name] = m
		}
	}
	return mws, nil
}

// NewProvider creates a new middleware provider instance.
func NewProvider(locator component.Locator) Provider {
	return &providerImpl{locator: locator}
}

// collectOptions gathers all required options for middleware creation via Require calls.
func collectOptions(h component.Handle) (*middlewarev1.Middleware, []Option, error) {
	// 1. Get base configuration directly from handle using centralized helper
	cfg, err := comp.AsConfig[middlewarev1.Middleware](h)
	if err != nil {
		return nil, nil, err
	}

	// 2. Resolve dynamic creation options (Carrier, Logger, etc.) via Require
	// This is where silent logic like WithCarrier for Selectors is injected.
	opts, err := comp.RequireTyped[[]Option](h, component.RequirementOption)
	if err != nil && !errors.Is(err, component.ErrRequirementNotFound) {
		return nil, nil, err
	}

	return cfg, opts, nil
}

// ServerProvider is the engine-compatible provider for server-side middleware components.
var ServerProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	if h.Scope() != component.ServerScope {
		return nil, nil
	}
	cfg, opts, err := collectOptions(h)
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
var ClientProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	if h.Scope() != component.ClientScope {
		return nil, nil
	}
	cfg, opts, err := collectOptions(h)
	if err != nil {
		return nil, err
	}
	m, ok := NewClient(cfg, opts...)
	if !ok {
		return nil, nil
	}
	return m, nil
}

// DefaultProvider acts as a catch-all dispatcher.
var DefaultProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	cfg, opts, err := collectOptions(h)
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
