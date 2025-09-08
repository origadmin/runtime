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
	GetJwt() MiddlewareJwtConfig
	GetSelector() MiddlewareSelectorConfig
	GetMetadata() MiddlewareMetadataConfig
	GetRateLimiter() MiddlewareRateLimiterConfig
	GetMetrics() MiddlewareMetricsConfig
	GetValidator() MiddlewareValidatorConfig
}

type MiddlewareJwtConfig interface {
	GetEnabled() bool
}

type MiddlewareSelectorConfig interface {
	GetEnabled() bool
}

type MiddlewareMetadataConfig interface {
	GetEnabled() bool
	GetPrefixes() []string
	GetData() map[string]string
}

type MiddlewareRateLimiterConfig interface {
	GetEnabled() bool
}

type MiddlewareMetricsConfig interface {
	GetEnabled() bool
}

type MiddlewareValidatorConfig interface {
	GetEnabled() bool
}

// ConfigLoader defines the interface for loading application configuration.
type ConfigLoader interface {
	Load(configPath string, bootstrapConfig interface{}) (kratosconfig.Config, error)
}
