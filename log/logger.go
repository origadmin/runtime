/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	"strings"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	kslog "github.com/origadmin/slog-kratos"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1" // Added for AppInfo
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/toolkits/slogx"
)

// NewLogger creates a new kratos logger based on the provided configuration.
// It uses slog as the underlying logging library and slog-kratos as an adapter.
func NewLogger(cfg *loggerv1.Logger) Logger {
	if cfg == nil {
		return DefaultLogger
	}
	if cfg.GetDisabled() {
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

// LevelOption converts a string level to a slogx.Option.
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

// WithDecorate decorates a kratos logger with common service information.
// It adds fields like service ID, name, version, timestamp, caller, trace ID, and span ID.
func WithDecorate(l kratoslog.Logger, info *appv1.App) kratoslog.Logger {
	if info == nil {
		info = &appv1.App{} // Use empty App to avoid nil pointer dereference
	}
	return kratoslog.With(
		l,
		"service.id", info.GetId(),
		"service.name", info.GetName(),
		"service.version", info.GetVersion(),
		"ts", kratoslog.DefaultTimestamp,
		"caller", kratoslog.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
}
