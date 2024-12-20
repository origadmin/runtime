/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

// Define types from kratos config package
type (
	Decoder      = kratosconfig.Decoder
	KeyValue     = kratosconfig.KeyValue
	Merge        = kratosconfig.Merge
	Observer     = kratosconfig.Observer
	Reader       = kratosconfig.Reader
	Resolver     = kratosconfig.Resolver
	Source       = kratosconfig.Source
	SourceOption = kratosconfig.Option
	SourceConfig = kratosconfig.Config
	Value        = kratosconfig.Value
	Watcher      = kratosconfig.Watcher
)

var (
	// ErrNotFound defined error from kratos config package
	ErrNotFound = kratosconfig.ErrNotFound
)

// NewSourceConfig returns a new config instance
func NewSourceConfig(opts ...SourceOption) SourceConfig {
	return kratosconfig.New(opts...)
}

// WithDecoder sets the decoder
func WithDecoder(d Decoder) SourceOption {
	return kratosconfig.WithDecoder(d)
}

// WithMergeFunc sets the merge function
func WithMergeFunc(m Merge) SourceOption {
	return kratosconfig.WithMergeFunc(m)
}

// WithResolveActualTypes enables resolving actual types
func WithResolveActualTypes(enableConvertToType bool) SourceOption {
	return kratosconfig.WithResolveActualTypes(enableConvertToType)
}

// WithResolver sets the resolver
func WithResolver(r Resolver) SourceOption {
	return kratosconfig.WithResolver(r)
}

// WithSource sets the source
func WithSource(s ...Source) SourceOption {
	return kratosconfig.WithSource(s...)
}
