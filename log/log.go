/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	kslog "github.com/origadmin/slog-kratos"
	"github.com/origadmin/toolkits/slogx"
)

type Logging struct {
	Logger log.Logger
	source log.Logger
}

func (l *Logging) Init(kv ...any) {
	l.Logger = log.With(l.source, kv...)
	log.SetLogger(l.Logger)
}

func New(cfg *configv1.Logger) log.Logger {
	if cfg == nil || cfg.GetDisabled() {
		return NewDiscard()
	}
	options := make([]slogx.Option, 0)
	if cfg.GetFile() != nil {
		options = append(options, slogx.WithFile(cfg.GetFile().GetPath()))
		_, file := filepath.Split(cfg.GetFile().GetPath())
		options = append(options, slogx.WithLumberjack(&slogx.LumberjackLogger{
			Filename:   file,
			MaxSize:    int(cfg.GetFile().GetMaxSize()),
			MaxAge:     int(cfg.GetFile().GetMaxAge()),
			MaxBackups: int(cfg.GetFile().GetMaxBackups()),
			LocalTime:  cfg.GetFile().GetLocalTime(),
			Compress:   cfg.GetFile().GetCompress(),
		}))
	}
	if cfg.GetStdout() {
		options = append(options, slogx.WithConsole(true))
	}
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
	options = append(options, LevelOption(cfg.GetLevel()))
	//l := log.With(kslog.NewLogger(),
	//	"ts", log.DefaultTimestamp,
	//	"caller", log.DefaultCaller,
	//	"service.id", flags.ServiceID(),
	//	"service.name", flags.ServiceName(),
	//	"service.version", flags.Version(),
	//	"trace.id", tracing.TraceID(),
	//	"span.id", tracing.SpanID(),
	//)
	//log.SetLogger(l)
	logger := slogx.New(options...)
	return kslog.NewLogger(logger)
}

func LevelOption(level configv1.LoggerLevel) slogx.Option {
	ll := slogx.LevelInfo
	switch level {
	case configv1.LoggerLevel_LOGGER_LEVEL_FATAL:
		ll = slogx.LevelFatal
	case configv1.LoggerLevel_LOGGER_LEVEL_DEBUG:
		ll = slogx.LevelDebug
	case configv1.LoggerLevel_LOGGER_LEVEL_ERROR:
		ll = slogx.LevelError
	case configv1.LoggerLevel_LOGGER_LEVEL_WARN:
		ll = slogx.LevelWarn
	case configv1.LoggerLevel_LOGGER_LEVEL_INFO:
		ll = slogx.LevelInfo
	}
	return slogx.WithLevel(ll)
}
