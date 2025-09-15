package interfaces

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

// Resolved is the main interface for accessing resolved configuration values.
// It provides type-safe, lazy-loaded, and cached access to different sections of the application's configuration.
type Resolved interface {
	// Resolve method is implemented by concrete resolvers to process the raw Kratos config.
	// Resolve(config kratosconfig.Config) (Resolved, error)

	// Services() ServiceConfig
	// Discovery() DiscoveryConfig
	// Middleware() MiddlewareConfig
	// Logger() LoggerConfig
	// Security() SecurityConfig
	// Add other top-level config accessors here (e.g., Data(), AppInfo())
}

// ConfigLoader defines the interface for loading application configuration.
// DEPRECATED: This interface is being phased out in favor of the new bootstrap.Load mechanism.
type ConfigLoader interface {
	Load(configPath string, bootstrapConfig interface{}) (kratosconfig.Config, error)
}

type ServiceConfig interface {
	// Add methods from configv1.Service that are needed by Resolved
	GetName() string
	GetProtocol() string
	// GetGrpc() *configv1.GRPC
	// GetHttp() *configv1.HTTP
	// GetWebsocket() *configv1.WebSocket
	// GetMessage() *configv1.Message
	// GetTask() *configv1.Task
	// GetMiddleware() *middlewarev1.Middleware
	// GetSelector() *configv1.Service_Selector
}

type LoggerConfig interface {
	GetDisabled() bool
	GetFile() LoggerFileConfig
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

// DiscoveryConfig provides type-safe access to discovery configuration.
type DiscoveryConfig interface {
	GetType() string
	GetServiceName() string
	GetDebug() bool
	// Add other methods from discoveryv1.Discovery that are needed by Resolved
}

// SecurityConfig provides type-safe access to security configuration.
type SecurityConfig interface {
	GetAuthn() AuthNConfig
	GetAuthz() AuthZConfig
	// Add other methods from securityv1.Security that are needed by Resolved
}

// AuthNConfig provides type-safe access to authentication configuration.
type AuthNConfig interface {
	GetDisabled() bool
	GetType() string
	GetPublicPaths() []string
	GetJwt() AuthNJwtConfig
	// Add other methods from securityv1.AuthNConfig that are needed by Resolved
}

// AuthNJwtConfig provides type-safe access to JWT authentication configuration.
type AuthNJwtConfig interface {
	GetAlgorithm() string
	GetSigningKey() string
	GetExpireTime() int64
	// Add other methods from securityv1.AuthNConfig_JWTConfig that are needed by Resolved
}

// AuthZConfig provides type-safe access to authorization configuration.
type AuthZConfig interface {
	GetDisabled() bool
	GetType() string
	GetPublicPaths() []string
	GetCasbin() AuthZCasbinConfig
	// Add other methods from securityv1.AuthZConfig that are needed by Resolved
}

// AuthZCasbinConfig provides type-safe access to Casbin authorization configuration.
type AuthZCasbinConfig interface {
	GetPolicyFile() string
	GetModelFile() string
	// Add other methods from securityv1.AuthZConfig_CasbinConfig that are needed by Resolved
}
