/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides functions for loading configurations and registering services.
package runtime

import (
	"os"
	"sync"
	"syscall"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
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
	signals     []os.Signal
	EnvPrefix   string
	WorkDir     string
	Logging     log.Logging
	logger      log.KLogger
	Config      config.KConfig
	Registry    registry.Registry
	Middleware  middleware.Middleware
	Service     service.Service
	source      *configv1.SourceConfig
	serviceInfo bootstrap.ServiceInfo
}

func (r *Runtime[T]) Load(bs *bootstrap.Bootstrap) error {
	var rerr error
	r.once.Do(func() {
		sourceConfig, err := bootstrap.LoadSourceConfig(bs)
		if err != nil {
			rerr = errors.Wrap(err, "load source config")
			return
		}
		r.source = sourceConfig
		kcfg, err := r.builder.NewConfig(r.source, config.WithServiceName(bs.ServiceName()))
		if err != nil {
			rerr = errors.Wrap(err, "new config")
		}
		r.Config = kcfg
		r.serviceInfo = bootstrap.ServiceInfo{
			ID:       bs.ServiceID(),
			Name:     bs.ServiceName(),
			Version:  bs.Version(),
			Metadata: bs.Metadata(),
		}
	})

	return rerr
}

func (r *Runtime[T]) Scan(v any) error {
	if err := r.Config.Scan(v); err != nil {
		return errors.Wrap(err, "scan config")
	}
	return nil
}

func (r *Runtime[T]) CreateApp(ctx context.Context, rr registry.Registry, ss ...transport.Server) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(r.serviceInfo.ID),
		kratos.Name(r.serviceInfo.Name),
		kratos.Version(r.serviceInfo.Version),
		kratos.Metadata(r.serviceInfo.Metadata),
		kratos.Context(ctx),
		kratos.Signal(r.signals...),
		kratos.Logger(r.logger),
	}

	if rr != nil {
		opts = append(opts, kratos.Registrar(rr))
	}
	if len(ss) > 0 {
		opts = append(opts, kratos.Server(ss...))
	}

	return kratos.New(opts...)
}

//
//func (r *Runtime[T]) CreateRegistrar(serviceName string, ss ...registry.Option) (registry.KRegistrar, error) {
//	err := r.Config.Scan(serviceName)
//	if err != nil {
//		return nil, err
//	}
//	return r.builder.NewRegistrar(cfg, ss...)
//}
//
//func (r *Runtime[T]) CreateDiscovery(serviceName string, ss ...registry.Option) (registry.KDiscovery, error) {
//	cfg, err := r.builder.NewDiscovery()
//	if err != nil {
//		return nil, err
//	}
//	return r.builder.NewDiscovery(cfg, ss...)
//}
//
//func (r *Runtime[T]) CreateGRPCServer(serviceName string, ss ...service.GRPCOption) (*service.GRPCServer, error) {
//	cfg, err := r.Config.Service(serviceName)
//	if err != nil {
//		return nil, err
//	}
//	return r.builder.NewGRPCServer(cfg, ss...)
//}
//
//func (r *Runtime[T]) CreateHTTPServer(serviceName string, ss ...service.HTTPOption) (*service.HTTPServer, error) {
//	cfg, err := r.Config.Service(serviceName)
//	if err != nil {
//		return nil, err
//	}
//	return r.builder.NewHTTPServer(cfg, ss...)
//}

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
		signals: []os.Signal{
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		},
	}
}

// Load uses the global Runtime instance to load configurations and other resources
// with the provided bootstrap settings. It returns an error if the loading process fails.
func Load(bs *bootstrap.Bootstrap) error {
	return runtime.Load(bs)
}
