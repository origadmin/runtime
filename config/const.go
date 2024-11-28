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
	Config   = kratosconfig.Config
	Decoder  = kratosconfig.Decoder
	KeyValue = kratosconfig.KeyValue
	Merge    = kratosconfig.Merge
	Observer = kratosconfig.Observer
	Option   = kratosconfig.Option
	Reader   = kratosconfig.Reader
	Resolver = kratosconfig.Resolver
	Source   = kratosconfig.Source
	Value    = kratosconfig.Value
	Watcher  = kratosconfig.Watcher
)

var (
	// ErrNotFound defined error from kratos config package
	ErrNotFound = kratosconfig.ErrNotFound
)

// New returns a new config instance
func New(opts ...Option) Config {
	return kratosconfig.New(opts...)
}

// WithDecoder sets the decoder
func WithDecoder(d Decoder) Option {
	return kratosconfig.WithDecoder(d)
}

// WithMergeFunc sets the merge function
func WithMergeFunc(m Merge) Option {
	return kratosconfig.WithMergeFunc(m)
}

// WithResolveActualTypes enables resolving actual types
func WithResolveActualTypes(enableConvertToType bool) Option {
	return kratosconfig.WithResolveActualTypes(enableConvertToType)
}

// WithResolver sets the resolver
func WithResolver(r Resolver) Option {
	return kratosconfig.WithResolver(r)
}

// WithSource sets the source
func WithSource(s ...Source) Option {
	return kratosconfig.WithSource(s...)
}
