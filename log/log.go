/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Logging struct {
	Logger log.Logger
	source log.Logger
}

func (l *Logging) Init(kv ...any) {
	l.Logger = log.With(l.source, kv...)
	log.SetLogger(l.Logger)
}
