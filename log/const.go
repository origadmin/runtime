/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	"github.com/go-kratos/kratos/v2/log"
)

const (
	LevelKey   = log.LevelKey
	LevelDebug = log.LevelDebug
	LevelInfo  = log.LevelInfo
	LevelWarn  = log.LevelWarn
	LevelError = log.LevelError
	LevelFatal = log.LevelFatal
)

type (
	KLogger         = log.Logger
	KLevel          = log.Level
	KValuer         = log.Valuer
	KFilterOption   = log.FilterOption
	KFilter         = log.Filter
	KOption         = log.Option
	KHelper         = log.Helper
	KWriterOptionFn = log.WriterOptionFn
)

var (
	DefaultLogger     = log.DefaultLogger
	DefaultCaller     = log.DefaultCaller
	DefaultTimestamp  = log.DefaultTimestamp
	DefaultMessageKey = log.DefaultMessageKey
)

var (
	With                = log.With
	WithContext         = log.WithContext
	NewStdLogger        = log.NewStdLogger
	ParseLevel          = log.ParseLevel
	Value               = log.Value
	Caller              = log.Caller
	Timestamp           = log.Timestamp
	FilterLevel         = log.FilterLevel
	FilterKey           = log.FilterKey
	FilterValue         = log.FilterValue
	FilterFunc          = log.FilterFunc
	NewFilter           = log.NewFilter
	SetLogger           = log.SetLogger
	GetLogger           = log.GetLogger
	Log                 = log.Log
	Context             = log.Context
	Debug               = log.Debug
	Debugf              = log.Debugf
	Debugw              = log.Debugw
	Info                = log.Info
	Infof               = log.Infof
	Infow               = log.Infow
	Warn                = log.Warn
	Warnf               = log.Warnf
	Warnw               = log.Warnw
	Error               = log.Error
	Errorf              = log.Errorf
	Errorw              = log.Errorw
	Fatal               = log.Fatal
	Fatalf              = log.Fatalf
	Fatalw              = log.Fatalw
	WithMessageKey      = log.WithMessageKey
	WithSprint          = log.WithSprint
	WithSprintf         = log.WithSprintf
	NewHelper           = log.NewHelper
	WithWriterLevel     = log.WithWriterLevel
	WithWriteMessageKey = log.WithWriteMessageKey
	NewWriter           = log.NewWriter
)
