/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"errors"
)

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/registry
//go:adapter:package:type *
//go:adapter:package:type:prefix K
//go:adapter:package:func *
//go:adapter:package:func:prefix K

var (
	ErrRegistryNotFound = errors.New("registry not found")
)
