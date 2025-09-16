/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides the foundational implementation for building distributed-ready services
// leveraging the Kratos microservice framework.
package runtime

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2"
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/configure"
	"github.com/prometheus/common/log"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	"github.com/origadmin/runtime/bootstrap" // 引入 bootstrap 包
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/internal/decoder"
	runtimeLog "github.com/origadmin/runtime/log"
	runtimeRegistry "github.com/origadmin/runtime/registry"
)

// AppInfo represents the application's configured information.
type AppInfo struct {
	ID        string
	Name      string
	Version   string
	Env       string
	StartTime time.Time
	Metadata  map[string]string
}

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
func WithAppInfo(appInfo AppInfo) Option {
	return func(o *options) {
		o.appInfo = appInfo
	}
}

// NewFromBootstrap is the recommended way to create a new Runtime instance.
// It takes a path to a bootstrap configuration file (e.g., "configs/bootstrap.yaml"),
// loads it, and then initializes the full runtime environment.
func NewFromBootstrap(bootstrapPath string, opts ...Option) (Runtime, func(), error) {
	// 1. Load config from the given path using the bootstrap package.
	kratosConfig, err := bootstrap.Load(bootstrapPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load bootstrap config from path %s: %w", bootstrapPath, err)
	}

	// 2. Call the original New function with the loaded config.
	rt, cleanup, err := New(kratosConfig, opts...)
	if err != nil {
		// If New fails, we must close the config source we just opened.
		if e := kratosConfig.Close(); e != nil {
			klog.Errorf("failed to close config source after New() error: %v", e)
		}
		return nil, nil, err
	}

	// 3. Chain the cleanup functions. The final cleanup must close the config source.
	finalCleanup := func() {
		cleanup()
		if e := kratosConfig.Close(); e != nil {
			klog.Errorf("failed to close config source during cleanup: %v", e)
		}
	}

	return rt, finalCleanup, nil
}

// New creates a new runtime instance from a pre-existing Kratos config instance.
// This is useful for advanced scenarios or testing. For standard usage, prefer NewFromBootstrap.
func New(kratosConfig kratosconfig.Config, opts ...Option) (Runtime, func(), error) {
	// Apply options
	appliedOpts := configure.Apply(&options{
		decoderProvider: decoder.DefaultDecoder,
	}, opts)

	// --- 1. Initialize and Validate AppInfo ---
	appInfo := appliedOpts.appInfo
	if appInfo.ID == "" || appInfo.Name == "" || appInfo.Version == "" || appInfo.Env == "" {
		return nil, nil, fmt.Errorf("app ID, Name, Version, and Env cannot be empty and must be provided via WithAppInfo option")
	}

	// Ensure a DecoderProvider is set
	if appliedOpts.decoderProvider == nil {
		return nil, nil, fmt.Errorf("DecoderProvider must be provided")
	}

	configDecoder, err := appliedOpts.decoderProvider.GetConfigDecoder(kratosConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get config decoder: %w", err)
	}

	// --- 2. Initialize Logger ---
	loggerConfig := getLoggerConfig(configDecoder)
	logger := runtimeLog.NewLogger(loggerConfig)
	klog.SetLogger(logger)

	// --- 3. Initialize all configured Service Registries & Discoveries ---
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

	// 4. Identify the default registrar for self-registration
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
		// Here you would add cleanup logic for components created by the runtime,
		// for example, closing registry connections if they had a Close() method.
		klog.Info("Runtime cleanup executed.")
	}

	return rt, cleanup, nil
}

// getLoggerConfig encapsulates the logic for decoding the logger configuration.
// It prioritizes the fast path via the LoggerConfig interface,
// falling back to a generic decode operation.
func getLoggerConfig(decoder interfaces.ConfigDecoder) *configv1.Logger {
	// Fast path: If the decoder directly provides logger config, use it.
	if d, ok := decoder.(interfaces.LoggerConfig); ok {
		return d.GetLogger()
	}

	// Slow path: Fall back to generic decoding.
	loggerConfig := new(configv1.Logger)
	if err := decoder.Decode("logger", loggerConfig); err != nil {
		log.Warnf("Failed to decode logger config, using default: %v", err)
		// On error, we still return the allocated (but empty) struct,
		// allowing the NewLogger function to apply its own defaults.
	}
	return loggerConfig
}

// registriesConfig holds the configuration for all service registries.
type registriesConfig struct {
	Registries      map[string]*discoveryv1.Discovery
	DefaultRegistry string
}

// getRegistriesConfig encapsulates the logic for decoding the registries' configuration.
func getRegistriesConfig(decoder interfaces.ConfigDecoder) registriesConfig {
	var cfg registriesConfig

	// Fast path: Check if the decoder can provide the discoveries map directly.
	if d, ok := decoder.(interfaces.DiscoveryConfig); ok {
		cfg.Registries = d.GetDiscoveries()
		// Attempt to decode the rest of the struct to get DefaultRegistry.
		if err := decoder.Decode("registries", &cfg); err != nil {
			log.Warnf("Could not decode 'registries' block for DefaultRegistry: %v", err)
		}
	} else {
		// Slow path: Decode the entire struct from the config.
		if err := decoder.Decode("registries", &cfg); err != nil {
			log.Warnf("Failed to decode registries config, running in standalone mode: %v", err)
		}
	}
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
	kratosOpts := []kratos.Option{
		kratos.ID(r.app.ID),
		kratos.Name(r.app.Name),
		kratos.Version(r.app.Version),
		kratos.Metadata(r.app.Metadata),
		kratos.Logger(r.logger),
		kratos.Server(servers...),
	}

	if r.defaultRegistrar != nil {
		kratosOpts = append(kratosOpts, kratos.Registrar(r.defaultRegistrar))
	}

	return kratos.New(kratosOpts...)
}
