/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides the foundational implementation for building distributed-ready services
// leveraging the Kratos microservice framework.
package runtime

import (
	"fmt"

	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	runtimeLog "github.com/origadmin/runtime/log"
	runtimeRegistry "github.com/origadmin/runtime/registry"
)

// Runtime defines the application's runtime environment, providing access to
// core components like configuration, logging, and service discovery/registration.
type Runtime interface {
	// AppInfo returns the application's configured information (ID, name, version).
	AppInfo() *configv1.App
	// Logger returns the configured Kratos logger.
	Logger() klog.Logger
	// NewApp creates a new Kratos application instance. It wires together the runtime's
	// configured components (like the default registrar) with the provided transport servers.
	NewApp(servers ...transport.Server) *kratos.App
	// DefaultRegistrar returns the default service registrar, used for self-registration.
	// It may be nil if no default registry is configured.
	DefaultRegistrar() registry.Registrar
	// Discovery returns a service discovery component by its configured name.
	// It returns the component and a boolean indicating if it was found.
	Discovery(name string) (registry.Discovery, bool)
	// Registrar returns a service registrar component by its configured name.
	// This is useful if a service needs to interact with a non-default registry.
	Registrar(name string) (registry.Registrar, bool)
}

// runtime is the internal implementation of the Runtime interface.
type runtime struct {
	app              *configv1.App
	logger           klog.Logger
	registrars       map[string]registry.Registrar
	discoveries      map[string]registry.Discovery
	defaultRegistrar registry.Registrar
}

// Config is the configuration for creating a new Runtime.
type Config struct {
	App    *configv1.App
	Logger *configv1.Logger
	// Registries holds the configuration for all service discovery/registration
	// components. The key is a unique name for the registry instance.
	Registries map[string]*configv1.Discovery
	// DefaultRegistry specifies the name of the registry to use for
	// self-registration. This key must exist in the Registries map.
	DefaultRegistry string
}

// New creates a new runtime instance from the given configuration.
// It initializes the logger and the service registry/discovery components.
// A cleanup function is returned to release resources, which should be deferred by the caller.
func New(cfg *Config) (Runtime, func(), error) {
	if cfg == nil || cfg.App == nil {
		return nil, nil, fmt.Errorf("app config must be provided")
	}

	// 1. Initialize Logger
	// The logger is fundamental and should be created first.
	logger := runtimeLog.NewLogger(cfg.Logger)
	klog.SetLogger(logger) // Set global logger for Kratos's internal logging

	// 2. Initialize all configured Service Registries & Discoveries
	// These are optional. If no registry config is provided, the app runs in standalone/local mode.
	registrars := make(map[string]registry.Registrar)
	discoveries := make(map[string]registry.Discovery)

	for name, registryCfg := range cfg.Registries {
		if registryCfg == nil || registryCfg.Type == "" || registryCfg.Type == "none" {
			klog.Infof("Skipping registry '%s' due to missing or 'none' type.", name)
			continue
		}

		klog.Infof("Initializing service registry and discovery '%s' with type: %s", name, registryCfg.Type)
		r, err := runtimeRegistry.NewRegistrar(registryCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create registrar for '%s': %w", name, err)
		}
		d, err := runtimeRegistry.NewDiscovery(registryCfg)
		if err != nil {
			// If registrar was created, we should ideally have a way to clean it up.
			// But since NewRegistrar doesn't return a cleanup func, we'll just log and continue.
			return nil, nil, fmt.Errorf("failed to create discovery for '%s': %w", name, err)
		}
		registrars[name] = r
		discoveries[name] = d
	}

	// 3. Identify the default registrar for self-registration
	var defaultReg registry.Registrar
	if cfg.DefaultRegistry != "" {
		var ok bool
		defaultReg, ok = registrars[cfg.DefaultRegistry]
		if !ok {
			return nil, nil, fmt.Errorf("default registry '%s' not found in configured registries", cfg.DefaultRegistry)
		}
		klog.Infof("Default registry for self-registration set to: '%s'", cfg.DefaultRegistry)
	} else if len(registrars) > 0 {
		klog.Warn("No default registry specified. The service will not register itself despite registries being configured.")
	} else {
		klog.Info("No service registry configured, running in local mode.")
	}

	rt := &runtime{
		app:              cfg.App,
		logger:           logger,
		registrars:       registrars,
		discoveries:      discoveries,
		defaultRegistrar: defaultReg,
	}

	cleanup := func() {
		klog.Info("Runtime cleanup executed.")
	}

	return rt, cleanup, nil
}

func (r *runtime) AppInfo() *configv1.App {
	return r.app
}

func (r *runtime) Logger() klog.Logger {
	return r.logger
}

func (r *runtime) DefaultRegistrar() registry.Registrar {
	return r.defaultRegistrar
}

func (r *runtime) Discovery(name string) (registry.Discovery, bool) {
	d, ok := r.discoveries[name]
	return d, ok
}

func (r *runtime) Registrar(name string) (registry.Registrar, bool) {
	reg, ok := r.registrars[name]
	return reg, ok
}

func (r *runtime) NewApp(servers ...transport.Server) *kratos.App {
	kratosOpts := []kratos.Option{
		kratos.ID(r.app.GetId()),
		kratos.Name(r.app.GetName()),
		kratos.Version(r.app.GetVersion()),
		kratos.Metadata(map[string]string{}), // Can be extended with more metadata from config
		kratos.Logger(r.logger),
		kratos.Server(servers...),
	}

	// If a default registrar is available, use it to register the service.
	if r.defaultRegistrar != nil {
		kratosOpts = append(kratosOpts, kratos.Registrar(r.defaultRegistrar))
	}

	return kratos.New(kratosOpts...)
}
