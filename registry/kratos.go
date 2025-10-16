/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	commonv1 "github.com/origadmin/runtime/api/gen/go/runtime/common/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
)

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/registry
//go:adapter:package:type *
//go:adapter:package:type:prefix K
//go:adapter:package:func *
//go:adapter:package:func:regex New([A-Z])=NewK$1
//go:adapter:package:func *
//go:adapter:package:func:prefix K

var (
	ErrRegistryNotFound = runtimeerrors.WithReason(runtimeerrors.NewStructured("registry", "registry not found").WithCaller(), commonv1.ErrorReason_NOT_FOUND)
)
