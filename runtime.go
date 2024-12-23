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
	"github.com/origadmin/runtime/config"
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
	once  = &sync.Once{}
	build = &builder{}
)

// ErrNotFound is an error that is returned when a ConfigBuilder or RegistryBuilder is not found.
var ErrNotFound = errors.String("not found")

type Runtime struct {
	once        sync.Once
	Debug       bool
	EnvPrefix   string
	WorkDir     string
	Application application.Application
	Logging     log.Logging
	Config      config.Config
	Registry    registry.Registry
	Middleware  middleware.Middleware
	Service     service.Service
}

func (r *Runtime) Load() error {
	var rerr error
	r.once.Do(func() {
		//// todo: add init and check before load
		//// todo: load config
		//if err := r.Config.Load(); err != nil {
		//	rerr = errors.Wrap(err, "load config")
		//	return
		//}
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
		//if err := r.Middleware.Load(); err != nil {
		//	rerr = errors.Wrap(err, "load middleware")
		//	return
		//}
	})
	return rerr
}

func (r *Runtime) Build(rr registry.Registry, servers ...transport.Server) *kratos.App {
	// todo: add init and check before build

	return kratos.New(
		kratos.ID(r.Application.ID),
		kratos.Name(r.Application.Name),
		kratos.Version(r.Application.Version),
		kratos.Metadata(r.Application.Metadata),
		kratos.Logger(r.Logging.Logger),
		kratos.Server(servers...),
		kratos.Registrar(rr),
	)
}
func New() Runtime {
	return Runtime{
		EnvPrefix: DefaultEnvPrefix,
	}
}
