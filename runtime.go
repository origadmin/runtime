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
	Config() config.Builder
	Registry() registry.Builder
	Service() service.Builder
	Middleware() middleware.Builder
}

// build is a global variable that holds an instance of the builder struct.
var (
	once    = &sync.Once{}
	runtime = &Runtime{
		builder:   &builder{},
		EnvPrefix: DefaultEnvPrefix,
	}
)

// ErrNotFound is an error that is returned when a ConfigBuilder or RegistryBuilder is not found.
var ErrNotFound = errors.String("not found")

type Provider interface {
	Load(*bootstrap.Bootstrap) error
	Scan(v any) error
	CreateApp(context.Context, registry.Registry, ...transport.Server) *kratos.App
	Signals() []os.Signal
	SetSignals([]os.Signal)
}

type Runtime struct {
	once        sync.Once
	builder     *builder
	env         string
	signals     []os.Signal
	Config      *config.Config
	EnvPrefix   string
	WorkDir     string
	Logging     log.Logging
	logger      log.KLogger
	Registry    registry.Registry
	Middleware  middleware.Middleware
	Service     service.Service
	source      *configv1.SourceConfig
	serviceInfo bootstrap.ServiceInfo
}

func (r *Runtime) Signals() []os.Signal {
	return r.signals
}

func (r *Runtime) SetSignals(signals []os.Signal) {
	r.signals = signals
}

func (r *Runtime) Load(bs *bootstrap.Bootstrap) error {
	var rerr error
	r.once.Do(func() {
		sourceConfig, err := bootstrap.LoadSourceConfig(bs)
		if err != nil {
			rerr = errors.Wrap(err, "load source config")
			return
		}
		r.source = sourceConfig
		r.Config = config.New(bs.ConfigFilePath())
		if err := r.Config.Load(bs.ServiceName(), sourceConfig); err != nil {
			return
		}
		r.serviceInfo = bs.ServiceInfo()
	})

	return rerr
}

func (r *Runtime) CreateApp(ctx context.Context, rr registry.Registry, ss ...transport.Server) *kratos.App {
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
//func (r *Runtime) CreateRegistrar(serviceName string, ss ...registry.Option) (registry.KRegistrar, error) {
//	err := r.Config.Scan(serviceName)
//	if err != nil {
//		return nil, err
//	}
//	return r.builder.NewRegistrar(cfg, ss...)
//}
//
//func (r *Runtime) CreateDiscovery(serviceName string, ss ...registry.Option) (registry.KDiscovery, error) {
//	cfg, err := r.builder.NewDiscovery()
//	if err != nil {
//		return nil, err
//	}
//	return r.builder.NewDiscovery(cfg, ss...)
//}
//
//func (r *Runtime) CreateGRPCServer(serviceName string, ss ...service.GRPCOption) (*service.GRPCServer, error) {
//	cfg, err := r.Config.Service(serviceName)
//	if err != nil {
//		return nil, err
//	}
//	return r.builder.NewGRPCServer(cfg, ss...)
//}
//
//func (r *Runtime) CreateHTTPServer(serviceName string, ss ...service.HTTPOption) (*service.HTTPServer, error) {
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

// Global returns the global instance of the Runtime struct.
func Global() *Runtime {
	return runtime
}

// New returns a new instance of the Runtime struct.
func New() Runtime {
	return Runtime{
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
