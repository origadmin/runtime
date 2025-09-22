/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides the foundational implementation for building distributed-ready services
// leveraging the Kratos microservice framework.
package runtime

import (
	"fmt"

	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log" // Import Kratos log package
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/configure"

	"github.com/origadmin/runtime/interfaces"
)

// Runtime defines the application's runtime environment, providing access to
// core components like configuration, logging, and service discovery/registration.
type Runtime interface {
	AppInfo() AppInfo
	Logger() kratoslog.Logger // Changed to kratoslog.Logger
	NewApp(servers ...transport.Server) *kratos.App
	DefaultRegistrar() registry.Registrar
	Discovery(name string) (registry.Discovery, bool)
	Registrar(name string) (registry.Registrar, bool)
	WithLogger(logger kratoslog.Logger) Runtime // Changed to kratoslog.Logger
}

// runtime is the internal implementation of the Runtime interface.
type runtime struct {
	app              AppInfo
	logger           kratoslog.Logger // Changed to kratoslog.Logger
	registrars       map[string]registry.Registrar
	discoveries      map[string]registry.Discovery
	defaultRegistrar registry.Registrar
}

// Option is a function type that applies a configuration option to the Runtime.
type Option func(*options)

// options holds the configuration options for creating a Runtime instance.
type options struct {
	appInfo AppInfo
}

// WithAppInfo sets the application information for the Runtime.
// This is a required option.
func WithAppInfo(appInfo AppInfo) Option {
	return func(o *options) {
		o.appInfo = appInfo
	}
}

// New creates a new runtime instance from a ComponentProvider.
// It requires that app info is provided via options.
func New(provider interfaces.ComponentProvider, opts ...Option) (Runtime, error) {
	// Apply options
	appliedOpts := configure.Apply(&options{}, opts)

	// --- 1. Validate Essential Options ---
	appInfo := appliedOpts.appInfo
	if appInfo.ID == "" || appInfo.Name == "" || appInfo.Version == "" || appInfo.Env == "" {
		return nil, fmt.Errorf("app ID, Name, Version, and Env cannot be empty and must be provided via WithAppInfo option")
	}

	// --- 2. Get Components from Provider ---
	logger := provider.GetLogger()
	if logger == nil {
		return nil, fmt.Errorf("logger component is missing from the provider")
	}

	registrars := provider.GetRegistrars()
	discoveries := provider.GetDiscoveries()
	defaultRegistrar := provider.GetDefaultRegistrar()

	// Log messages about registries/discoveries can be handled by runtime's logger
	if len(registrars) == 0 {
		logger.Info("No service registry configured, running in local mode.")
	} else if defaultRegistrar == nil {
		logger.Warn("No default registrar found. Service will not self-register.")
	} else {
		logger.Infof("Default registrar for self-registration set.")
	}

	rt := &runtime{
		app:              appInfo,
		logger:           logger,
		registrars:       registrars,
		discoveries:      discoveries,
		defaultRegistrar: defaultRegistrar,
	}

	return rt, nil
}

func (r *runtime) AppInfo() AppInfo {
	return r.app
}

func (r *runtime) Logger() kratoslog.Logger { // Changed to kratoslog.Logger
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

func (r *runtime) WithLogger(logger kratoslog.Logger) Runtime { // Changed to kratoslog.Logger
	newRt := *r
	newRt.logger = logger
	return &newRt
}

// NewApp creates a new Kratos application with the provided servers.
func (r *runtime) NewApp(servers ...transport.Server) *kratos.App {
	logger := kratoslog.NewHelper(r.logger) // Changed to kratoslog.NewHelper
	if len(servers) == 0 {
		logger.Fatal("No servers provided. Please provide at least one server.")
	}

	kratosOpts := r.app.Options()
	kratosOpts = append(kratosOpts, kratos.Logger(r.logger), kratos.Server(servers...))

	if r.defaultRegistrar != nil {
		kratosOpts = append(kratosOpts, kratos.Registrar(r.defaultRegistrar))
	}

	return kratos.New(kratosOpts...)
}
