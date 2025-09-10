/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"github.com/origadmin/runtime/api/gen/go/apierrors"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/factory"
)

// Factory is the interface for creating new registrar and discovery components.
type Factory interface {
	NewRegistrar(*configv1.Discovery, ...Option) (KRegistrar, error)
	NewDiscovery(*configv1.Discovery, ...Option) (KDiscovery, error)
}

// buildImpl is the concrete implementation of the Builder.
type buildImpl struct {
	factory.Registry[Factory]
}

func (b *buildImpl) NewRegistrar(cfg *configv1.Discovery, opts ...Option) (KRegistrar, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, errors.NewMessage(apierrors.ErrorReason_INVALID_REGISTRY_CONFIG, "registry configuration or type is missing")
	}

	f, ok := b.Get(cfg.Type)
	if !ok {
		return nil, errors.NewMessageWithMeta(apierrors.ErrorReason_NOT_FOUND, map[string]string{"type": cfg.Type}, "no registry factory found for type: %s", cfg.Type)
	}
	registrar, err := f.NewRegistrar(cfg, opts...)
	if err != nil {
		return nil, errors.WrapAndConvert(err, apierrors.ErrorReason_REGISTRY_CREATION_FAILURE, "failed to create registrar for type %s", cfg.Type)
	}
	return registrar, nil
}

func (b *buildImpl) NewDiscovery(cfg *configv1.Discovery, opts ...Option) (KDiscovery, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, errors.NewMessage(apierrors.ErrorReason_INVALID_REGISTRY_CONFIG, "registry configuration or type is missing")
	}

	f, ok := b.Get(cfg.Type)
	if !ok {
		return nil, errors.NewMessageWithMeta(apierrors.ErrorReason_NOT_FOUND, map[string]string{"type": cfg.Type}, "no registry factory found for type: %s", cfg.Type)
	}
	discovery, err := f.NewDiscovery(cfg, opts...)
	if err != nil {
		return nil, errors.WrapAndConvert(err, apierrors.ErrorReason_REGISTRY_CREATION_FAILURE, "failed to create discovery for type %s", cfg.Type)
	}
	return discovery, nil
}

// defaultBuilder is a private variable to prevent accidental modification from other packages.
var defaultBuilder = &buildImpl{
	Registry: factory.New[Factory](),
}

// DefaultBuilder returns the shared instance of the registry builder.
func DefaultBuilder() Builder {
	return defaultBuilder
}

// NewBuilder is an alias for DefaultBuilder for consistency, though it returns a shared instance.
// func NewBuilder() Builder {
// 	return DefaultBuilder()
// }
