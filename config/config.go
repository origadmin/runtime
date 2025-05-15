/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"google.golang.org/protobuf/proto"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

type (
	// Builder is an interface that defines a method for registering a config Builder.
	Builder interface {
		factory.Registry[Factory]
		NewConfig(*configv1.SourceConfig, ...Option) (KConfig, error)
	}
	// Factory is an interface that defines a method for creating a new config.
	Factory interface {
		// NewSource creates a new config using the given KConfig and a list of Options.
		NewSource(*configv1.SourceConfig, *Options) (KSource, error)
	}

	// Syncer is an interface that defines a method for synchronizing a config.
	Syncer interface {
		SyncConfig(*configv1.SourceConfig, string, any, *Options) error
	}

	// ProtoSyncer is an interface that defines a method for synchronizing a protobuf message.
	ProtoSyncer interface {
		SyncConfig(*configv1.SourceConfig, string, proto.Message, *Options) error
	}
)
