/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	"strings"

	kratoslog "github.com/go-kratos/kratos/v2/log"

	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	kslog "github.com/origadmin/slog-kratos"
	"github.com/origadmin/toolkits/slogx"
)

// NewLogger creates a new kratos logger based on the provided configuration.
// It uses slog as the underlying logging library and slog-kratos as an adapter.
func NewLogger(cfg *loggerv1.Logger) kratoslog.Logger {
	if cfg == nil || cfg.GetDisabled() {
		return NewDiscard()
	}

	var options []slogx.Option

	// Configure output writers
	if cfg.GetStdout() {
		options = append(options, slogx.WithConsole(true))
	}

	if fileConfig := cfg.GetFile(); fileConfig != nil {
		if fileConfig.GetLumberjack() {
			options = append(options, slogx.WithLumberjack(&slogx.LumberjackLogger{
				Filename:   fileConfig.GetPath(),
				MaxSize:    int(fileConfig.GetMaxSize()),
				MaxAge:     int(fileConfig.GetMaxAge()),
				MaxBackups: int(fileConfig.GetMaxBackups()),
				LocalTime:  fileConfig.GetLocalTime(),
				Compress:   fileConfig.GetCompress(),
			}))
		} else {
			// Use the path from fileConfig, not the logger's name
			options = append(options, slogx.WithOutputFile(fileConfig.GetPath()))
		}
	}

	// If no output is configured, default to console output
	if len(options) == 0 {
		options = append(options, slogx.WithConsole(true))
	}

	// Configure log format
	switch cfg.GetFormat() {
	case "dev":
		options = append(options, slogx.WithFormat(slogx.FormatDev))
	case "json":
		options = append(options, slogx.WithFormat(slogx.FormatJSON))
	case "tint":
		options = append(options, slogx.WithFormat(slogx.FormatTint))
	default:
		options = append(options, slogx.WithFormat(slogx.FormatText))
	}

	// Configure log level
	options = append(options, LevelOption(cfg.GetLevel()))

	// Create the underlying slog logger
	slogLogger := slogx.New(options...)

	// Adapt the slog logger to the kratos logger interface
	kratosLogger := kslog.NewLogger(kslog.WithLogger(slogLogger))

	if cfg.GetDefault() {
		kratoslog.SetLogger(kratosLogger)
	}

	return kratosLogger
}

// LevelOption converts a string level to an slogx.Option.
func LevelOption(level string) slogx.Option {
	var ll slogx.Level
	switch strings.ToLower(level) {
	case "debug":
		ll = slogx.LevelDebug
	case "info":
		ll = slogx.LevelInfo
	case "warn":
		ll = slogx.LevelWarn
	case "error":
		ll = slogx.LevelError
	case "fatal":
		ll = slogx.LevelFatal
	default:
		ll = slogx.LevelInfo // Default to Info
	}
	return slogx.WithLevel(ll)
}
