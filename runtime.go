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
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/configure"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/internal/decoder"
	runtimeLog "github.com/origadmin/runtime/log"
	runtimeRegistry "github.com/origadmin/runtime/registry"
)

// Runtime defines the application's runtime environment, providing access to
// core components like configuration, logging, and service discovery/registration.
type Runtime interface {
	AppInfo() AppInfo
	Logger() klog.Logger
	NewApp(servers ...transport.Server) *kratos.App
	DefaultRegistrar() registry.Registrar
	Discovery(name string) (registry.Discovery, bool)
	Registrar(name string) (registry.Registrar, bool)
}

// runtime is the internal implementation of the Runtime interface.
type runtime struct {
	app              AppInfo
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
func WithDecoderProvider(p interfaces.ConfigDecoderProvider) Option {
	return func(o *options) {
		o.decoderProvider = p
	}
}

// WithAppInfo sets the application information for the Runtime.
// This is a required option.
func WithAppInfo(appInfo AppInfo) Option {
	return func(o *options) {
		o.appInfo = appInfo
	}
}

// NewFromBootstrap is the recommended, one-stop function to create a new Runtime instance.
// It loads configuration from the given path and initializes the full runtime environment.
func NewFromBootstrap(bootstrapPath string, opts ...Option) (Runtime, func(), error) {
	// Load config from the given path using the bootstrap package.
	kratosConfig, err := bootstrap.Load(bootstrapPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load bootstrap config from path %s: %w", bootstrapPath, err)
	}

	// Call the main New function with the loaded config.
	rt, cleanup, err := New(kratosConfig, opts...)
	if err != nil {
		// If New fails, we must close the config source we just opened.
		if e := kratosConfig.Close(); e != nil {
			klog.Errorf("failed to close config source after New() error: %v", e)
		}
		return nil, nil, err
	}

	// Chain the cleanup functions. The final cleanup must close the config source.
	finalCleanup := func() {
		cleanup()
		if e := kratosConfig.Close(); e != nil {
			klog.Errorf("failed to close config source during cleanup: %v", e)
		}
	}

	return rt, finalCleanup, nil
}

// New creates a new runtime instance from a pre-existing Kratos config instance.
// It requires that app info is provided via options.
func New(kratosConfig kratosconfig.Config, opts ...Option) (Runtime, func(), error) {
	// Apply options
	appliedOpts := configure.Apply(&options{
		decoderProvider: decoder.DefaultDecoder,
	}, opts)

	// --- 1. Validate Essential Options ---
	appInfo := appliedOpts.appInfo
	if appInfo.ID == "" || appInfo.Name == "" || appInfo.Version == "" || appInfo.Env == "" {
		return nil, nil, fmt.Errorf("app ID, Name, Version, and Env cannot be empty and must be provided via WithAppInfo option")
	}

	// --- 2. Initialize Decoder ---
	if appliedOpts.decoderProvider == nil {
		return nil, nil, fmt.Errorf("DecoderProvider must be provided")
	}
	configDecoder, err := appliedOpts.decoderProvider.GetConfigDecoder(kratosConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get config decoder: %w", err)
	}

	// --- 3. Create and Enrich Logger ---
	// This is the first point where we have both the config and the appInfo context.
	logger := newLogger(configDecoder, appInfo)
	klog.SetLogger(logger)

	// --- 4. Initialize all configured Service Registries & Discoveries ---
	registriesCfg := getRegistriesConfig(configDecoder)

	registrars := make(map[string]registry.Registrar)
	discoveries := make(map[string]registry.Discovery)

	for name, registryCfg := range registriesCfg.Registries {
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
			return nil, nil, fmt.Errorf("failed to create discovery for '%s': %w", name, err)
		}
		registrars[name] = r
		discoveries[name] = d
	}

	// --- 5. Identify the default registrar for self-registration ---
	var defaultReg registry.Registrar
	if registriesCfg.DefaultRegistry != "" {
		var ok bool
		defaultReg, ok = registrars[registriesCfg.DefaultRegistry]
		if !ok {
			return nil, nil, fmt.Errorf("default registry '%s' not found in configured registries", registriesCfg.DefaultRegistry)
		}
		klog.Infof("Default registry for self-registration set to: '%s'", registriesCfg.DefaultRegistry)
	} else if len(registrars) > 0 {
		klog.Warn("No default registry specified. The service will not register itself despite registries being configured.")
	} else {
		klog.Info("No service registry configured, running in local mode.")
	}

	rt := &runtime{
		app:              appInfo,
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

// newLogger creates the logger backend from config and enriches it with app info.
func newLogger(decoder interfaces.ConfigDecoder, appInfo AppInfo) klog.Logger {
	var loggerConfig *configv1.Logger

	// Fast path: If the decoder directly provides logger config, use it.
	if d, ok := decoder.(interfaces.LoggerConfig); ok {
		loggerConfig = d.GetLogger()
	} else {
		// Slow path: Fall back to generic decoding.
		loggerConfig = new(configv1.Logger) // Initialize if not from fast path
		if err := decoder.Decode("logger", loggerConfig); err != nil {
			klog.Warnf("Failed to decode logger config, using default: %v", err)
		}
	}

	loggerBackend := runtimeLog.NewLogger(loggerConfig)

	return klog.With(loggerBackend,
		"service.name", appInfo.Name,
		"service.version", appInfo.Version,
		"service.id", appInfo.ID,
		"service.env", appInfo.Env,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
}

// registriesConfig holds the configuration for all service registries.
type registriesConfig struct {
	Registries      map[string]*discoveryv1.Discovery
	DefaultRegistry string
}

// getRegistriesConfig encapsulates the logic for decoding the registries' configuration.
func getRegistriesConfig(decoder interfaces.ConfigDecoder) registriesConfig {
	var cfg registriesConfig
	if d := decoder.(interfaces.DiscoveryConfig); d != nil {
		cfg.Registries = d.GetDiscoveries()
	} else {
		cfg.Registries = make(map[string]*discoveryv1.Discovery)
		if err := decoder.Decode("registries", &cfg.Registries); err != nil {
			klog.Warnf("Failed to decode registries config, running in standalone mode: %v", err)
		}
	}
	//cfg.DefaultRegistry = decoder.GetDefaultDiscovery()
	return cfg
}

func (r *runtime) AppInfo() AppInfo {
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
	// Start with the application identity options.
	kratosOpts := r.app.Options()

	// Append runtime-specific options.
	kratosOpts = append(kratosOpts,
		kratos.Logger(r.logger),
		kratos.Server(servers...),
	)

	// Conditionally add the registrar.
	if r.defaultRegistrar != nil {
		kratosOpts = append(kratosOpts, kratos.Registrar(r.defaultRegistrar))
	}

	return kratos.New(kratosOpts...)
}
