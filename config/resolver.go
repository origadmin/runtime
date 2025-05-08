/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
)

type Resolver interface {
	Observer(string, KValue)
	Resolve(KConfig) (Resolved, error)
}

type Resolved interface {
	Middleware() (*middlewarev1.Middleware, error)
	Service() (*configv1.Service, error)
}
