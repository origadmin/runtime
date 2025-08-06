/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces"
)

const Type = "config"

// NewConfig creates a new config instance.
func NewConfig(cfg *configv1.SourceConfig, opts ...interfaces.Option) (kratosconfig.Config, error) {
	return defaultConfigFactory.NewConfig(cfg, opts...)
}

// Register registers a config factory.
func Register(name string, factory interfaces.ConfigFactory) {
	defaultConfigFactory.Register(name, factory)
}
