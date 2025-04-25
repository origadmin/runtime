/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides functions for loading configurations and registering services.
package runtime

import (
	"sync"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/origadmin/runtime/application"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/toolkits/errors"
)

const (
	DefaultEnvPrefix = "ORIGADMIN_RUNTIME_SERVICE"
)

type Builder interface {
	config.Builder
	registry.Builder
	service.Builder
	MiddlewareBuilders
	middlewareBuildRegistry
}

// build is a global variable that holds an instance of the builder struct.
var (
	once    = &sync.Once{}
	runtime = &Runtime[any]{
		builder:   &builder{},
		EnvPrefix: DefaultEnvPrefix,
	}
)

// ErrNotFound is an error that is returned when a ConfigBuilder or RegistryBuilder is not found.
var ErrNotFound = errors.String("not found")

type Runtime[T any] struct {
	once        sync.Once
	builder     *builder
	Debug       bool
	EnvPrefix   string
	WorkDir     string
	Application application.Application
	Logging     log.Logging
	Config      config.Config
	Registry    registry.Registry
	Middleware  middleware.Middleware
	Service     service.Service
	bootstrap   T
}

func (r *Runtime[T]) Load(bs *bootstrap.Bootstrap) error {
	var rerr error
	r.once.Do(func() {
		sc := new(configv1.SourceConfig)
		rerr = bootstrap.LoadLocalConfig(bs, sc)
		if rerr != nil {
			return
		}

		//// todo: add init and check before load
		//// todo: load config
		if err := r.Config.LoadFromSource(sc); err != nil {
			rerr = errors.Wrap(err, "load config")
			return
		}
		//// todo: load registry
		//if err := r.Registry.Load(); err != nil {
		//	rerr = errors.Wrap(err, "load registry")
		//	return
		//}
		//// todo: load service
		//if err := r.Service.Load(); err != nil {
		//	rerr = errors.Wrap(err, "load service")
		//	return
		//}
		//// todo: load middleware
		//if err := r.Build.Load(); err != nil {
		//	rerr = errors.Wrap(err, "load middleware")
		//	return
		//}
	})
	return rerr
}

func (r *Runtime[T]) Build(rr registry.Registry, ss ...transport.Server) *kratos.App {
	// todo: add init and check before build

	return kratos.New(
		kratos.ID(r.Application.ID),
		kratos.Name(r.Application.Name),
		kratos.Version(r.Application.Version),
		kratos.Metadata(r.Application.Metadata),
		kratos.Logger(r.Logging.Logger),
		kratos.Server(ss...),
		kratos.Registrar(rr),
	)
}

func (r *Runtime[T]) CreateRegistrar(serviceName string, ss ...registry.Option) (registry.KRegistrar, error) {
	cfg, err := r.Config.Registry(serviceName)
	if err != nil {
		return nil, err
	}
	return r.builder.NewRegistrar(cfg, ss...)
}

func (r *Runtime[T]) CreateDiscovery(serviceName string, ss ...registry.Option) (registry.KDiscovery, error) {
	cfg, err := r.Config.Registry(serviceName)
	if err != nil {
		return nil, err
	}
	return r.builder.NewDiscovery(cfg, ss...)
}

func (r *Runtime[T]) CreateGRPCServer(serviceName string, ss ...service.GRPCOption) (*service.GRPCServer, error) {
	cfg, err := r.Config.Service(serviceName)
	if err != nil {
		return nil, err
	}
	return r.builder.NewGRPCServer(cfg, ss...)
}

func (r *Runtime[T]) CreateHTTPServer(serviceName string, ss ...service.HTTPOption) (*service.HTTPServer, error) {
	cfg, err := r.Config.Service(serviceName)
	if err != nil {
		return nil, err
	}
	return r.builder.NewHTTPServer(cfg, ss...)
}

func init() {
	once.Do(func() {
		runtime.builder.init()
	})
}

// GlobalBuilder returns the global instance of the builder.
func GlobalBuilder() Builder {
	return runtime.builder
}

// Global returns the global instance of the Runtime struct.
func Global() *Runtime[any] {
	return runtime
}

// New returns a new instance of the Runtime struct.
func New[T any]() Runtime[T] {
	return Runtime[T]{
		builder:   runtime.builder,
		EnvPrefix: DefaultEnvPrefix,
	}
}

// Load uses the global Runtime instance to load configurations and other resources
// with the default bootstrap settings. It returns an error if the loading process fails.
func Load() error {
	return runtime.Load(bootstrap.New())
}

// LoadWithBootstrap uses the global Runtime instance to load configurations and other resources
// with the provided bootstrap settings. It returns an error if the loading process fails.
func LoadWithBootstrap(bs *bootstrap.Bootstrap) error {
	return runtime.Load(bs)
}
