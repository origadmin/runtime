/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"context"

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
	handle component.Handle
}

func (p *providerImpl) Middleware(name string) (KMiddleware, error) {
	return comp.Get[KMiddleware](context.Background(), p.handle.In(component.CategoryMiddleware), name)
}

func (p *providerImpl) DefaultMiddleware() (KMiddleware, error) {
	return comp.GetDefault[KMiddleware](context.Background(), p.handle.In(component.CategoryMiddleware))
}

// GetMiddlewares collects all middlewares from the given handle as a slice.
func GetMiddlewares(ctx context.Context, h component.Handle) ([]KMiddleware, error) {
	var mws []KMiddleware
	for _, inst := range h.Iter(ctx) {
		if m, ok := inst.(KMiddleware); ok {
			mws = append(mws, m)
		}
	}
	return mws, nil
}

// GetMiddlewareMap collects all available middlewares from the given handle as a map.
func GetMiddlewareMap(ctx context.Context, h component.Handle) (map[string]KMiddleware, error) {
	mws := make(map[string]KMiddleware)
	for name, inst := range h.Iter(ctx) {
		if m, ok := inst.(KMiddleware); ok {
			mws[name] = m
		}
	}
	return mws, nil
}

// NewProvider creates a new middleware provider instance.
func NewProvider(handle component.Handle) Provider {
	return &providerImpl{handle: handle}
}

// ServerProvider is the engine-compatible provider for server-side middleware components.
var ServerProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	if h.Scope() != component.ServerScope {
		return nil, nil
	}
	cfg, err := comp.AsConfig[middlewarev1.Middleware](h)
	if err != nil {
		return nil, err
	}
	m, ok := NewServer(cfg)
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
	cfg, err := comp.AsConfig[middlewarev1.Middleware](h)
	if err != nil {
		return nil, err
	}
	m, ok := NewClient(cfg)
	if !ok {
		return nil, nil
	}
	return m, nil
}

// DefaultProvider acts as a catch-all dispatcher.
var DefaultProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	cfg, err := comp.AsConfig[middlewarev1.Middleware](h)
	if err != nil {
		return nil, err
	}

	if h.Scope() == component.ClientScope {
		m, ok := NewClient(cfg)
		if ok {
			return m, nil
		}
	} else {
		m, ok := NewServer(cfg)
		if ok {
			return m, nil
		}
	}
	return nil, nil
}

