/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides the foundational implementation for building distributed-ready services
// leveraging the Kratos microservice framework.
package runtime

import (
	"fmt"

	"github.com/go-kratos/kratos/v2"
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/configure"
	"github.com/prometheus/common/log"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/internal/decoder"
	runtimeLog "github.com/origadmin/runtime/log"
	runtimeRegistry "github.com/origadmin/runtime/registry"
)

// AppInfo is an alias for configv1.App, representing the application's configured information.
type AppInfo struct {
	ID       string
	Name     string
	Version  string
	Metadata map[string]string
}

// Runtime defines the application's runtime environment, providing access to
// core components like configuration, logging, and service discovery/registration.
type Runtime interface {
	// AppInfo returns the application's configured information (ID, name, version).
	AppInfo() *AppInfo
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
	app              *AppInfo
	logger           klog.Logger
	registrars       map[string]registry.Registrar
	discoveries      map[string]registry.Discovery
	defaultRegistrar registry.Registrar
}

// Option is a function type that applies a configuration option to the Runtime.
type Option func(*options)

// options holds the configuration options for creating a Runtime instance.
type options struct {
	appInfo         AppInfo
	decoderProvider interfaces.ConfigDecoderProvider
}

// WithDecoderProvider sets the DecoderProvider for the Runtime.
// This allows the Runtime to access application-specific configurations.
func WithDecoderProvider(p interfaces.ConfigDecoderProvider) Option {
	return func(o *options) {
		o.decoderProvider = p
	}
}

func WithAppInfo(appInfo AppInfo) Option {
	return func(o *options) {
		o.appInfo = appInfo
	}
}

// New creates a new runtime instance from the given configuration.
// It initializes the logger and the service registry/discovery components.
// A cleanup function is returned to release resources, which should be deferred by the caller.
func New(kratosConfig kratosconfig.Config, opts ...Option) (Runtime, func(), error) {
	// Apply options
	appliedOpts := configure.Apply(&options{
		decoderProvider: decoder.DefaultDecoder,
	}, opts)

	// Ensure a DecoderProvider is set
	if appliedOpts.decoderProvider == nil {
		return nil, nil, fmt.Errorf("DecoderProvider must be provided via WithDecoderProvider option")
	}

	// Get the Decoder instance from the provider
	configDecoder, err := appliedOpts.decoderProvider.GetConfigDecoder(kratosConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get config decoder: %w", err)
	}

	// --- 1. Initialize AppInfo ---
	appInfo := appliedOpts.appInfo

	// --- 2. Initialize Logger ---
	var loggerConfig *configv1.Logger
	if d, ok := configDecoder.(interfaces.LoggerConfig); ok {
		loggerConfig = d.GetLogger()
	} else {
		loggerConfig = new(configv1.Logger)
		if err := configDecoder.Decode("logger", loggerConfig); err != nil {
			log.Warnf("Failed to decode logger config, using default: %v", err)
		}
	}
	logger := runtimeLog.NewLogger(loggerConfig)
	klog.SetLogger(logger) // Set global logger for Kratos's internal logging

	// --- 3. Initialize all configured Service Registries & Discoveries ---
	var registriesConfig struct {
		Registries      map[string]*discoveryv1.Discovery
		DefaultRegistry string
	}
	if d, ok := configDecoder.(interfaces.DiscoveryConfig); ok {
		registriesConfig.Registries = d.GetDiscoveries()
	} else {
		if err := configDecoder.Decode("registries", &registriesConfig); err != nil {
			log.Warnf("Failed to decode registries config, running in standalone mode: %v", err)
		}
	}

	registrars := make(map[string]registry.Registrar)
	discoveries := make(map[string]registry.Discovery)

	for name, registryCfg := range registriesConfig.Registries {
		if registryCfg == nil || registryCfg.GetType() == "" || registryCfg.GetType() == "none" {
			klog.Infof("Skipping registry '%s' due to missing or 'none' type.", name)
			continue
		}

		klog.Infof("Initializing service registry and discovery '%s' with type: %s", name, registryCfg.GetType())
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

	// 4. Identify the default registrar for self-registration
	var defaultReg registry.Registrar
	if registriesConfig.DefaultRegistry != "" {
		var ok bool
		defaultReg, ok = registrars[registriesConfig.DefaultRegistry]
		if !ok {
			return nil, nil, fmt.Errorf("default registry '%s' not found in configured registries", registriesConfig.DefaultRegistry)
		}
		klog.Infof("Default registry for self-registration set to: '%s'", registriesConfig.DefaultRegistry)
	} else if len(registrars) > 0 {
		klog.Warn("No default registry specified. The service will not register itself despite registries being configured.")
	} else {
		klog.Info("No service registry configured, running in local mode.")
	}

	rt := &runtime{
		app:              &appInfo,
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

func (r *runtime) AppInfo() *AppInfo {
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
		kratos.ID(r.app.ID),
		kratos.Name(r.app.Name),
		kratos.Version(r.app.Version),
		kratos.Metadata(r.app.Metadata), // Can be extended with more metadata from config
		kratos.Logger(r.logger),
		kratos.Server(servers...),
	}

	// If a default registrar is available, use it to register the service.
	if r.defaultRegistrar != nil {
		kratosOpts = append(kratosOpts, kratos.Registrar(r.defaultRegistrar))
	}

	return kratos.New(kratosOpts...)
}
