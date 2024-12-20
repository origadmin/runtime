/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides functions for loading configurations and registering services.
package runtime

import (
	"sync"

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
	EnvPrefix   string
	Application application.Application
	Logging     log.Logging
	Config      config.Config
	Registry    registry.Registry
	Middleware  middleware.Middleware
	Service     service.Service
}

func New() Runtime {
	return Runtime{
		EnvPrefix: DefaultEnvPrefix,
	}
}
