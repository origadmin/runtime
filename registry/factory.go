/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"fmt"
	"net/http"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/factory"
)

// Factory is the interface for creating new registrar and discovery components.
type Factory interface {
	NewRegistrar(cfg *configv1.Discovery, opts ...Option) (KRegistrar, error)
	NewDiscovery(cfg *configv1.Discovery, opts ...Option) (KDiscovery, error)
}

// buildImpl is the concrete implementation of the Builder.
type buildImpl struct {
	factory.Registry[Factory]
}

func (b *buildImpl) NewRegistrar(cfg *configv1.Discovery, opts ...Option) (KRegistrar, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, errors.New(http.StatusBadRequest, ReasonInvalidConfig, "registry configuration or type is missing")
	}
	f, ok := b.Get(cfg.Type)
	if !ok {
		err := errors.New(http.StatusNotFound, ReasonRegistryNotFound, fmt.Sprintf("no registry factory found for type: %s", cfg.Type))
		return nil, err.WithMetadata(map[string]string{"type": cfg.Type})
	}
	registrar, err := f.NewRegistrar(cfg, opts...)
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, ReasonCreationFailure, fmt.Sprintf("failed to create registrar for type %s: %v", cfg.Type, err))
	}
	return registrar, nil
}

func (b *buildImpl) NewDiscovery(cfg *configv1.Discovery, opts ...Option) (KDiscovery, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, errors.New(http.StatusBadRequest, ReasonInvalidConfig, "registry configuration or type is missing")
	}
	f, ok := b.Get(cfg.Type)
	if !ok {
		err := errors.New(http.StatusNotFound, ReasonRegistryNotFound, fmt.Sprintf("no registry factory found for type: %s", cfg.Type))
		return nil, err.WithMetadata(map[string]string{"type": cfg.Type})
	}
	discovery, err := f.NewDiscovery(cfg, opts...)
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, ReasonCreationFailure, fmt.Sprintf("failed to create discovery for type %s: %v", cfg.Type, err))
	}
	return discovery, nil
}

// defaultBuilder is a private variable to prevent accidental modification from other packages.
var defaultBuilder = &buildImpl{
	Registry: factory.New[Factory](),
}

// GetDefaultBuilder returns the shared instance of the registry builder.
func GetDefaultBuilder() Builder {
	return defaultBuilder
}

// NewBuilder is an alias for GetDefaultBuilder for consistency, though it returns a shared instance.
func NewBuilder() Builder {
	return GetDefaultBuilder()
}
