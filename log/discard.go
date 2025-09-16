/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	kratoslog "github.com/go-kratos/kratos/v2/log"
)

var discardInstance DiscardLogger

func NewDiscard() kratoslog.Logger {
	return discardInstance
}

type DiscardLogger struct{}

func (d DiscardLogger) Log(level kratoslog.Level, keyvals ...interface{}) error {
	return nil
}

func (d DiscardLogger) Debug(msg string, keyvals ...interface{}) {
}

func (d DiscardLogger) Info(msg string, keyvals ...interface{}) {
}

func (d DiscardLogger) Warn(msg string, keyvals ...interface{}) {
}

func (d DiscardLogger) Error(msg string, keyvals ...interface{}) {
}
