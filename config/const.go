/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

const Type = "config"

// Define types from kratos config package
type (
	KDecoder  = kratosconfig.Decoder
	KKeyValue = kratosconfig.KeyValue
	KMerge    = kratosconfig.Merge
	KObserver = kratosconfig.Observer
	KReader   = kratosconfig.Reader
	KResolver = kratosconfig.Resolver
	KSource   = kratosconfig.Source
	KOption   = kratosconfig.Option
	KConfig   = kratosconfig.Config
	KValue    = kratosconfig.Value
	KWatcher  = kratosconfig.Watcher
)

var (
	// ErrNotFound defined error from kratos config package
	ErrNotFound = kratosconfig.ErrNotFound
)

// NewSourceConfig returns a new config instance
func NewSourceConfig(opts ...KOption) KConfig {
	return kratosconfig.New(opts...)
}

// WithDecoder sets the decoder
func WithDecoder(d KDecoder) KOption {
	return kratosconfig.WithDecoder(d)
}

// WithMergeFunc sets the merge function
func WithMergeFunc(m KMerge) KOption {
	return kratosconfig.WithMergeFunc(m)
}

// WithResolveActualTypes enables resolving actual types
func WithResolveActualTypes(enableConvertToType bool) KOption {
	return kratosconfig.WithResolveActualTypes(enableConvertToType)
}

// WithResolver sets the resolver
func WithResolver(r KResolver) KOption {
	return kratosconfig.WithResolver(r)
}

// WithSource sets the sourceConfig
func WithSource(s ...KSource) KOption {
	return kratosconfig.WithSource(s...)
}
