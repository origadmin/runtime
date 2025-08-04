package interfaces

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"google.golang.org/protobuf/proto"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

type ConfigBuilder interface {
	factory.Registry[ConfigFactory]
	NewConfig(*configv1.SourceConfig, ...Option) (kratosconfig.Config, error)
}

type ConfigFactory interface {
	// NewSource creates a new config using the given KConfig and a list of Options.
	NewSource(*configv1.SourceConfig, *Options) (kratosconfig.Source, error)
}

type ConfigSyncer interface {
	SyncConfig(*configv1.SourceConfig, string, any, *Options) error
}

type ConfigProtoSyncer interface {
	SyncConfig(*configv1.SourceConfig, string, proto.Message, *Options) error
}
