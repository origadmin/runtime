package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"google.golang.org/protobuf/proto"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/runtime/source/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
)

type Builder interface {
	factory.Registry[SourceFactory]
	NewConfig(*sourcev1.Sources, ...options.Option) (interfaces.Config, error)
	//SyncConfig(*sourcev1.Sources, any, ...Option) error // Add SyncConfig method
}

type SourceFunc func(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error)

func (c SourceFunc) NewSource(config *sourcev1.SourceConfig, options ...options.Option) (kratosconfig.Source, error) {
	return c(config, options...)
}

type SourceFactory interface {
	// NewSource creates a new config using the given KConfig and a list of Options.
	NewSource(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error)
}

type Syncer interface {
	SyncConfig(*sourcev1.SourceConfig, string, proto.Message, ...options.Option) error
}

type ProtoSyncer interface {
	SyncConfig(*sourcev1.SourceConfig, string, proto.Message, ...options.Option) error
}

type FileConfig func(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error)

func (f FileConfig) NewSource(sourceConfig *sourcev1.SourceConfig, opts ...options.Option) (kratosconfig.Source, error) {
	return f(sourceConfig, opts...)
}
