/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"google.golang.org/protobuf/proto"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces" // Import interfaces package
)

type (
	SyncFunc func(*configv1.SourceConfig, string, any, *interfaces.Options) error

	// Syncer is an interface that defines a method for synchronizing a config.
	Syncer interface {
		SyncConfig(*configv1.SourceConfig, string, any, *interfaces.Options) error
	}

	// ProtoSyncer is an interface that defines a method for synchronizing a protobuf message.
	ProtoSyncer interface {
		SyncConfig(*configv1.SourceConfig, string, proto.Message, *interfaces.Options) error
	}
)

func (fn SyncFunc) SyncConfig(cfg *configv1.SourceConfig, key string, value any, opts *interfaces.Options) error {
	return fn(cfg, key, value, opts)
}
