/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config provides adapters for Kratos config types and functions.
package config

import (
	_ "github.com/go-kratos/kratos/v2/config"
)

//go:generate adptool ./kratos.go
//go:adapter:package github.com/go-kratos/kratos/v2/config
