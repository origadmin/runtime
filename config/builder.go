package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"google.golang.org/protobuf/proto"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

type Builder interface {
	factory.Registry[SourceFactory]
	NewConfig(*configv1.Sources, ...Option) (kratosconfig.Config, error)
	//SyncConfig(*configv1.Sources, any, ...Option) error // Add SyncConfig method
}

type SourceFunc func(*configv1.SourceConfig, *Options) (kratosconfig.Source, error)

func (c SourceFunc) NewSource(config *configv1.SourceConfig, options *Options) (kratosconfig.Source, error) {
	return c(config, options)
}

type SourceFactory interface {
	// NewSource creates a new config using the given KConfig and a list of Options.
	NewSource(*configv1.SourceConfig, *Options) (kratosconfig.Source, error)
}

type Syncer interface {
	SyncConfig(*configv1.SourceConfig, string, any, *Options) error
}

type ProtoSyncer interface {
	SyncConfig(*configv1.SourceConfig, string, proto.Message, *Options) error
}

type FileConfig func(*configv1.SourceConfig, *Options) (kratosconfig.Source, error)

func (f FileConfig) NewSource(sourceConfig *configv1.SourceConfig, opts *Options) (kratosconfig.Source, error) {
	return f(sourceConfig, opts)
}
