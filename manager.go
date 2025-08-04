/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/go-kratos/kratos/v2"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces"
)

type Manager struct {
	Middleware interfaces.MiddlewareProvider
	Service    interfaces.ServiceProvider
	//Discovery   RegistryProvider
	//Config     ConfigProvider
}
