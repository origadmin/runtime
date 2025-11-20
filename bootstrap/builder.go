/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/bootstrap/internal/container"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/storage"

	data_storage "github.com/origadmin/runtime/data/storage" // Import the data storage package
)

// containerImpl implements the interfaces.Container interface.
type containerImpl struct {
	discoveries       map[string]registry.Discovery
	registrars        map[string]registry.Registrar
	defaultRegistrar  registry.Registrar
	serverMiddlewares map[string]middleware.Middleware
	clientMiddlewares map[string]middleware.Middleware
	components        map[string]interfaces.ComponentFactory

	storageProvider storage.Provider // Add storage provider
}

// Discoveries returns a map of all configured service discovery components.
func (c *containerImpl) Discoveries() map[string]registry.Discovery {
	return c.discoveries
}

// Discovery returns a discovery client by name.
func (c *containerImpl) Discovery(name string) (registry.Discovery, bool) {
	d, ok := c.discoveries[name]
	return d, ok
}

// Registrars returns a map of all configured service registrar components.
func (c *containerImpl) Registrars() map[string]registry.Registrar {
	return c.registrars
}

// Registrar returns a registrar by name.
func (c *containerImpl) Registrar(name string) (registry.Registrar, bool) {
	r, ok := c.registrars[name]
	return r, ok
}

// DefaultRegistrar returns the default service registrar, used for service self-registration.
// It may be nil if no default registry is configured.
func (c *containerImpl) DefaultRegistrar() registry.Registrar {
	return c.defaultRegistrar
}

// ServerMiddlewares returns a map of all configured server middlewares.
func (c *containerImpl) ServerMiddlewares() map[string]middleware.Middleware {
	return c.serverMiddlewares
}

// ServerMiddleware returns a server middleware by name.
func (c *containerImpl) ServerMiddleware(name string) (middleware.Middleware, bool) {
	m, ok := c.serverMiddlewares[name]
	return m, ok
}

// ClientMiddlewares returns a map of all configured client middlewares.
func (c *containerImpl) ClientMiddlewares() map[string]middleware.Middleware {
	return c.clientMiddlewares
}

// ClientMiddleware returns a client middleware by name.
func (c *containerImpl) ClientMiddleware(name string) (middleware.Middleware, bool) {
	m, ok := c.clientMiddlewares[name]
	return m, ok
}

// StorageProvider returns the configured storage provider.
func (c *containerImpl) StorageProvider() storage.Provider {
	return c.storageProvider
}

// Component retrieves a generic component by its registered name.
// This allows for future components to be added without changing the interface.
func (c *containerImpl) Component(name string) (component interface{}, ok bool) {
	factory, ok := c.components[name]
	if !ok {
		return nil, false
	}
	// Components are created on demand.
	comp, err := factory(nil, c) // Pass nil for component-specific config for now
	if err != nil {
		log.Errorf("failed to create component '%s': %v", name, err)
		return nil, false
	}
	return comp, true
}

// buildContainer builds and returns the component container.
func buildContainer(sc interfaces.StructuredConfig, providerOptions *ProviderOptions) (interfaces.Container,
	log.Logger, error) {
	factories := providerOptions.componentFactories
	// 1. Create the component provider implementation.
	builder := container.NewBuilder(factories).WithConfig(sc)

	// 2. Initialize core components by consuming the config.
	c, err := builder.Build(providerOptions.rawOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize components: %w", err)
	}

	return c, builder.Logger(), err
}

// buildStorageProvider builds the storage provider.
func buildStorageProvider(sc interfaces.StructuredConfig) (storage.Provider, error) {
	return data_storage.New(sc)
}