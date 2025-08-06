/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config provides adapters for Kratos config types and functions.
package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

// Kratos Config Type Adapters
type (
	KConfig   = kratosconfig.Config
	KDecoder  = kratosconfig.Decoder
	KKeyValue = kratosconfig.KeyValue
	KMerge    = kratosconfig.Merge
	KObserver = kratosconfig.Observer
	KReader   = kratosconfig.Reader
	KResolver = kratosconfig.Resolver
	KSource   = kratosconfig.Source
	KOption   = kratosconfig.Option
	KWatcher  = kratosconfig.Watcher
)

// Kratos Config Option Helpers

// NewKratosConfig returns a new config instance
func NewKratosConfig(opts ...KOption) KConfig {
	return kratosconfig.New(opts...)
}

// WithKratosDecoder sets the decoder
func WithKratosDecoder(d KDecoder) KOption {
	return kratosconfig.WithDecoder(d)
}

// WithKratosMergeFunc sets the merge function
func WithKratosMergeFunc(m KMerge) KOption {
	return kratosconfig.WithMergeFunc(m)
}

// WithKratosResolveActualTypes enables resolving actual types
func WithKratosResolveActualTypes(enableConvertToType bool) KOption {
	return kratosconfig.WithResolveActualTypes(enableConvertToType)
}

// WithKratosResolver sets the resolver
func WithKratosResolver(r KResolver) KOption {
	return kratosconfig.WithResolver(r)
}

// WithKratosSource sets the sourceConfig
func WithKratosSource(s ...KSource) KOption {
	return kratosconfig.WithSource(s...)
}
