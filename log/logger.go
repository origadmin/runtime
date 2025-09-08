/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	kslog "github.com/origadmin/slog-kratos"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/toolkits/slogx"
)

// NewLogger creates a new kratos logger based on the provided configuration.
// It uses slog as the underlying logging library and slog-kratos as an adapter.
func NewLogger(cfg *configv1.Logger) log.Logger {
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
			// Assuming slogx.WithFile uses the service name for the filename if path is not absolute
			options = append(options, slogx.WithFile(cfg.GetName()))
		}
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
		log.SetLogger(kratosLogger)
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
