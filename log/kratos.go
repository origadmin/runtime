/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	_ "github.com/go-kratos/kratos/v2/log"
)

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/log klog
//go:adapter:package:type *
//go:adapter:package:type:prefix K
//go:adapter:package:func *
//go:adapter:package:func:prefix K
