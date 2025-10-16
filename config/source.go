/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package config

import (
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	commonv1 "github.com/origadmin/runtime/api/gen/go/runtime/common/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/runtime/source/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

const Type = "config"

var (
	ErrInvalidConfigType = runtimeerrors.WithReason(runtimeerrors.NewStructured(Type, "invalid config type").WithCaller(), commonv1.ErrorReason_VALIDATION_ERROR)
)

// NewConfig creates a new config instance.
func NewConfig(cfg *sourcev1.Sources, opts ...options.Option) (interfaces.Config, error) {
	return defaultBuilder.NewConfig(cfg, opts...)
}

// Register registers a config factory.
func Register(name string, sourceFactory any) {
	var factory SourceFactory
	switch fty := sourceFactory.(type) {
	case SourceFactory:
		factory = fty
	case SourceFunc:
		factory = fty
	case func(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error):
		factory = SourceFunc(fty)
	default:
		panic(ErrInvalidConfigType)
	}
	defaultBuilder.Register(name, factory)
}

// RegisterSourceFactory registers a source factory for a specific config type.
func RegisterSourceFactory(name string, factory SourceFactory) {
	defaultBuilder.Register(name, factory)
}

// RegisterSourceFunc registers a source function for a specific config type.
func RegisterSourceFunc(name string, factory SourceFunc) {
	defaultBuilder.Register(name, factory)
}

// Load loads configuration from the specified file path and scans it into the target struct.
// It returns the Kratos config instance, which should be closed by the caller when no longer needed.
func Load(configPath string, target interface{}) (kratosconfig.Config, error) {
	c := kratosconfig.New(
		kratosconfig.WithSource(
			file.NewSource(configPath),
		),
	)

	if err := c.Load(); err != nil {
		// Ensure config is closed on load error to prevent resource leaks
		c.Close()
		return nil, runtimeerrors.WithReason(runtimeerrors.WrapStructured(err, Type, fmt.Sprintf("failed to load config from %s", configPath)).WithCaller(), commonv1.ErrorReason_INTERNAL_SERVER_ERROR)
	}

	if err := c.Scan(target); err != nil {
		// Ensure config is closed on scan error
		c.Close()
		return nil, runtimeerrors.WithReason(runtimeerrors.WrapStructured(err, Type, "failed to scan config into target").WithCaller(), commonv1.ErrorReason_INTERNAL_SERVER_ERROR)
	}

	return c, nil
}
