/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log acts as a bridge to the Kratos logging library.
// The adptool directives below generate wrapper code, making Kratos's logging
// functions and types directly available in this package without non-idiomatic prefixes.
package log

import (
	// This import ensures that the Kratos log package is a dependency for the generator.
	_ "github.com/go-kratos/kratos/v2/log"
)

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/log kratoslog
////go:adapter:package:type *
////go:adapter:package:type:prefix K
////go:adapter:package:func *
////go:adapter:package:func:prefix K
