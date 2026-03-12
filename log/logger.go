/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and contracts for the module.
package log

import (
	"strings"
	"sync"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	kslog "github.com/origadmin/slog-kratos"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/toolkits/slogx"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.RWMutex
	*slogx.Logger
}

var (
	DefaultSlogLogger = slogx.New(slogx.WithFormat(slogx.FormatDev), slogx.WithConsole(true))
)

func init() {
	global.SetLogger(DefaultSlogLogger)
	kratoslog.SetLogger(kslog.NewLogger(kslog.WithLogger(DefaultSlogLogger)))
}

func (a *loggerAppliance) SetLogger(in *slogx.Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetSlogLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetSlogLogger(logger *slogx.Logger) {
	global.SetLogger(logger)
}

// GetSlogLogger returns global logger appliance as logger in current process.
func GetSlogLogger() *slogx.Logger {
	global.lock.RLock()
	defer global.lock.RUnlock()
	return global.Logger
}

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
		SetSlogLogger(slogLogger)
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
		ll = slogx.LevelInfo
	}
	return slogx.WithLevel(ll)
}

// WithDecorate decorates a kratos logger with common service information.
// It adds fields like service ID, name, version, timestamp, caller, trace ID, and span ID.
func WithDecorate(l kratoslog.Logger, info *appv1.App) kratoslog.Logger {
	var (
		id      = ""
		name    = ""
		version = ""
	)
	if info != nil {
		id = info.GetId()
		name = info.GetName()
		version = info.GetVersion()
	}
	return kratoslog.With(
		l,
		"service.id", id,
		"service.name", name,
		"service.version", version,
		"ts", kratoslog.DefaultTimestamp,
		"caller", kratoslog.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
}
