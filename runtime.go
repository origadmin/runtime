/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides functions for loading configurations and registering services.
package runtime

import (
	"sync"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/toolkits/errors"
)

const (
	DefaultEnvPrefix = "ORIGADMIN_RUNTIME_SERVICE"
)

type Builder interface {
	ConfigBuilder
	registry.Builder
	service.Builder
	MiddlewareBuilders

	configBuildRegistry
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
	EnvPrefix  string
	Config     config.Config
	Registry   registry.Registry
	Middleware middleware.Middleware
	Service    service.Service
}

func New(prefix string) Runtime {
	if prefix == "" {
		prefix = DefaultEnvPrefix
	}
	return Runtime{
		EnvPrefix: prefix,
	}
}
