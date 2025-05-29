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

//type Logging struct {
//	Logger log.Logger
//	source log.Logger
//}
//
//func (l *Logging) Init(kv ...any) {
//	l.Logger = log.With(l.source, kv...)
//	log.SetLogger(l.Logger)
//}

func New(cfg *configv1.Logger) log.Logger {
	if cfg == nil || cfg.GetDisabled() {
		return NewDiscard()
	}
	options := make([]slogx.Option, 0)
	fileConfig := cfg.File
	if fileConfig != nil {
		log.Infof("show log file config: %+v", fileConfig)
		//path, _ := filepath.Abs(fileConfig.GetPath())
		//log.Infof("log file path: %s", path)
		//_, err := os.Stat(path)
		//if err != nil && os.IsNotExist(err) {
		//	err = os.MkdirAll(path, os.ModePerm)
		//	if err != nil {
		//		log.Errorf("create log file path error: %v", err)
		//		return NewDiscard()
		//	}
		//}
		//pathname := filepath.Join(fileConfig.GetPath(), cfg.GetName())

		if fileConfig.Lumberjack {
			options = append(options, slogx.WithLumberjack(&slogx.LumberjackLogger{
				Filename:   fileConfig.GetPath(),
				MaxSize:    int(fileConfig.GetMaxSize()),
				MaxAge:     int(fileConfig.GetMaxAge()),
				MaxBackups: int(fileConfig.GetMaxBackups()),
				LocalTime:  fileConfig.GetLocalTime(),
				Compress:   fileConfig.GetCompress(),
			}))
		} else {
			options = append(options, slogx.WithFile(cfg.GetName()))
		}
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
	l := NewKSLoggerWithOption(options)
	if cfg.GetDefault() {
		log.SetLogger(l)
	}
	return l
}

func NewKSLoggerWithOption(options []slogx.Option) *kslog.Logger {
	base := slogx.New(options...)
	return kslog.NewLogger(kslog.WithLogger(base))
}

func LevelOption(level string) slogx.Option {
	ll := slogx.LevelInfo
	switch strings.ToLower(level) {
	case "fatal":
		ll = slogx.LevelFatal
	case "debug":
		ll = slogx.LevelDebug
	case "error":
		ll = slogx.LevelError
	case "warn":
		ll = slogx.LevelWarn
	case "info":
		ll = slogx.LevelInfo
	default:

	}
	//case configv1.LoggerLevel_LOGGER_LEVEL_FATAL:
	//	ll = slogx.LevelFatal
	//case configv1.LoggerLevel_LOGGER_LEVEL_DEBUG:
	//	ll = slogx.LevelDebug
	//case configv1.LoggerLevel_LOGGER_LEVEL_ERROR:
	//	ll = slogx.LevelError
	//case configv1.LoggerLevel_LOGGER_LEVEL_WARN:
	//	ll = slogx.LevelWarn
	//case configv1.LoggerLevel_LOGGER_LEVEL_INFO:
	//	ll = slogx.LevelInfo
	return slogx.WithLevel(ll)
}

//func New(cfg *configv1.Logger) log.Logger {
//	// ... 原有代码 ...
//
//	// 补充采样配置
//	if cfg.GetSampling() != nil {
//		options = append(options, slogx.WithSampling(
//			cfg.GetSampling().GetEnabled(),
//			cfg.GetSampling().GetRate(),
//		))
//	}
//
//	// 补充异步写入
//	if cfg.GetAsync() {
//		options = append(options, slogx.WithAsync(true))
//	}
//
//	// 补充编码器配置
//	if ec := cfg.GetEncoder(); ec != nil {
//		encoderOpts := []slogx.EncoderOption{
//			slogx.WithTimeKey(ec.GetTimeKey()),
//			slogx.WithLevelKey(ec.GetLevelKey()),
//			slogx.WithNameKey(ec.GetNameKey()),
//			slogx.WithCallerKey(ec.GetCallerKey()),
//			slogx.WithFunctionKey(ec.GetFunctionKey()),
//			slogx.WithStacktraceKey(ec.GetStacktraceKey()),
//			slogx.WithTimeFormat(ec.GetTimeFormat()),
//			slogx.WithUTC(ec.GetUtc()),
//		}
//		options = append(options, slogx.WithEncoderOptions(encoderOpts...))
//	}
//
//	// 补充时间间隔切割
//	if fileConfig != nil && fileConfig.GetRotateInterval() > 0 {
//		options = append(options, slogx.WithRotateInterval(
//			time.Duration(fileConfig.GetRotateInterval())*time.Hour,
//		))
//	}
//
//	// 补充最小级别过滤
//	if cfg.GetMinLevel() != configv1.LoggerLevel_LOGGER_LEVEL_UNSPECIFIED {
//		options = append(options, LevelFilterOption(cfg.GetMinLevel()))
//	}
//
//	// ... 后续原有代码 ...
//}
//
//// 新增最小级别过滤函数
//func LevelFilterOption(level configv1.LoggerLevel) slogx.Option {
//	ll := slogx.LevelInfo
//	switch level {
//	// ... 复用已有的LevelOption转换逻辑 ...
//	}
//	return slogx.WithLevelFilter(ll)
//}

//message Logger {
//// ... 已有字段 ...
//
//// 采样配置
//message Sampling {
//bool enabled = 1;
//float rate = 2; // 采样率 0.0-1.0
//}
//Sampling sampling = 10;
//
//// 异步写入
//bool async = 11;
//
//// 编码器配置
//message Encoder {
//string time_key = 1;
//string level_key = 2;
//string name_key = 3;
//string caller_key = 4;
//string function_key = 5;
//string stacktrace_key = 6;
//string time_format = 7;
//bool utc = 8;
//}
//Encoder encoder = 12;
//
//// 最小日志级别
//LoggerLevel min_level = 13;
//
//// 文件切割时间间隔 (小时)
//int32 rotate_interval = 14;
//}
