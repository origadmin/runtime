/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	commonv1 "github.com/origadmin/runtime/api/gen/go/common/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
)

// Factory is the interface for creating new registrar and discovery components.
type Factory interface {
	NewRegistrar(*discoveryv1.Discovery, ...options.Option) (KRegistrar, error)
	NewDiscovery(*discoveryv1.Discovery, ...options.Option) (KDiscovery, error)
}

// buildImpl is the concrete implementation of the Builder.
type buildImpl struct {
	factory.Registry[Factory]
}

func (b *buildImpl) NewRegistrar(cfg *discoveryv1.Discovery, opts ...options.Option) (KRegistrar, error) {
	if cfg == nil || cfg.Type == "" {
		// TODO: Refactor to use a registry-specific error reason instead of a generic one.
		return nil, errors.NewMessage(commonv1.ErrorReason_VALIDATION_ERROR,
			"registry configuration or type is missing")
	}

	f, ok := b.Get(cfg.Type)
	if !ok {
		return nil, errors.NewMessageWithMeta(commonv1.ErrorReason_NOT_FOUND, map[string]string{"type": cfg.Type}, "no registry factory found for type: %s", cfg.Type)
	}
	registrar, err := f.NewRegistrar(cfg, opts...)
	if err != nil {
		// TODO: Refactor to use a registry-specific error reason instead of a generic one.
		return nil, errors.WrapAndConvert(err, commonv1.ErrorReason_INTERNAL_SERVER_ERROR, "failed to create registrar for type %s", cfg.Type)
	}
	return registrar, nil
}

func (b *buildImpl) NewDiscovery(cfg *discoveryv1.Discovery, opts ...options.Option) (KDiscovery, error) {
	if cfg == nil || cfg.Type == "" {
		// TODO: Refactor to use a registry-specific error reason instead of a generic one.
		return nil, errors.NewMessage(commonv1.ErrorReason_VALIDATION_ERROR, "registry configuration or type is missing")
	}

	f, ok := b.Get(cfg.Type)
	if !ok {
		return nil, errors.NewMessageWithMeta(commonv1.ErrorReason_NOT_FOUND, map[string]string{"type": cfg.Type}, "no registry factory found for type: %s", cfg.Type)
	}
	discovery, err := f.NewDiscovery(cfg, opts...)
	if err != nil {
		// TODO: Refactor to use a registry-specific error reason instead of a generic one.
		return nil, errors.WrapAndConvert(err, commonv1.ErrorReason_INTERNAL_SERVER_ERROR, "failed to create discovery for type %s", cfg.Type)
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
