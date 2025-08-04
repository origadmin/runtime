package interfaces

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
)

type Resolver interface {
	Resolve(config kratosconfig.Config) (Resolved, error)
}

type Resolved interface {
	WithDecode(name string, v any, decode func([]byte, any) error) error
	Value(name string) (any, error)
	Middleware() *middlewarev1.Middleware
	Services() []*configv1.Service
	Logger() *configv1.Logger
	Discovery() *configv1.Discovery
}

type ResolveObserver interface {
	Observer(string, kratosconfig.Value)
}

type Options struct {
	ConfigName    string
	ServiceName   string
	ResolverName  string
	EnvPrefixes   []string
	Sources       []kratosconfig.Source
	ConfigOptions []kratosconfig.Option
	Decoder       kratosconfig.Decoder
	Encoder       Encoder
	ForceReload   bool
}

// Option is a function that takes a pointer to a KOption struct and modifies it.
type Option = func(s *Options)
