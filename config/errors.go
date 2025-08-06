/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"github.com/origadmin/toolkits/errors"
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

var (
	ErrInvalidConfigType = errors.New("invalid config type")
	ErrNotFound = kratosconfig.ErrNotFound
)
