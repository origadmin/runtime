package interfaces

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

type ServiceConfig interface {
	GetType() string
	// Add other methods from configv1.Service that are needed by Resolved
}

type LoggerConfig interface {
	GetDisabled() bool
	GetFile() *LoggerFileConfig
	GetStdout() bool
	GetFormat() string
	GetLevel() string
	// Add other methods from configv1.Logger that are needed by Resolved
}

type LoggerFileConfig interface {
	GetPath() string
	GetLumberjack() bool
	GetMaxSize() int32
	GetMaxAge() int32
	GetMaxBackups() int32
	GetLocalTime() bool
	GetCompress() bool
	// Add other methods from configv1.Logger.File that are needed by Resolved
}

type MiddlewareConfig interface {
	GetEnabledMiddlewares() []string
	GetJwt() *MiddlewareJwtConfig
	GetSelector() *MiddlewareSelectorConfig
	GetMetadata() *MiddlewareMetadataConfig
	GetRateLimiter() *MiddlewareRateLimiterConfig
	GetValidator() *MiddlewareValidatorConfig
	// Add other methods from middlewarev1.Middleware that are needed by Resolved
}

type MiddlewareJwtConfig interface {
	GetEnabled() bool
	// Add other methods from middlewarev1.Middleware.Jwt that are needed by Resolved
}

type MiddlewareSelectorConfig interface {
	GetEnabled() bool
	// Add other methods from middlewarev1.Middleware.Selector that are needed by Resolved
}

type MiddlewareMetadataConfig interface {
	GetEnabled() bool
	// Add other methods from middlewarev1.Middleware.Metadata that are needed by Resolved
}

type MiddlewareRateLimiterConfig interface {
	GetEnabled() bool
	// Add other methods from middlewarev1.Middleware.RateLimiter that are needed by Resolved
}

type MiddlewareValidatorConfig interface {
	GetEnabled() bool
	// Add other methods from middlewarev1.Middleware.Validator that are needed by Resolved
}

type Resolved interface {
	WithDecode(name string, v any, decode func([]byte, any) error) error
	Value(name string) (any, error)
	Middleware() MiddlewareConfig
	Services() []ServiceConfig
	Logger() LoggerConfig
	Discovery() DiscoveryConfig
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
