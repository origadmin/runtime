/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides functions for loading configurations and registering services.
package runtime

import (
	"sync"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/toolkits/errors"
)

type Builder interface {
	ConfigBuilder
	RegistryBuilder
	ServiceBuilder
	MiddlewareBuilders

	configBuildRegistry
	registryBuildRegistry
	serviceBuildRegistry
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
	Config   config.Config
	Registry registry.Registry
}
