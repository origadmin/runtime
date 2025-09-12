package resolver

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	jwt "github.com/origadmin/runtime/api/gen/go/middleware/v1/jwt"
	metrics "github.com/origadmin/runtime/api/gen/go/middleware/v1/metrics"
	ratelimit "github.com/origadmin/runtime/api/gen/go/middleware/v1/ratelimit"
	selector "github.com/origadmin/runtime/api/gen/go/middleware/v1/selector"
	validator "github.com/origadmin/runtime/api/gen/go/middleware/v1/validator"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
)

type ResolveFunc func(config kratosconfig.Config) (interfaces.Resolved, error)

func (r ResolveFunc) Resolve(config kratosconfig.Config) (interfaces.Resolved, error) {
	return r(config)
}

type resolver struct {
	values map[string]any

	// Cached decoded values
	servicesOnce   sync.Once
	cachedServices interfaces.ServiceConfig

	discoveryOnce   sync.Once
	cachedDiscovery interfaces.DiscoveryConfig

	middlewareOnce   sync.Once
	cachedMiddleware interfaces.MiddlewareConfig

	loggerOnce   sync.Once
	cachedLogger interfaces.LoggerConfig
}

func (r *resolver) Services() interfaces.ServiceConfig {
	r.servicesOnce.Do(func() {
		var ss configv1.Service
		if !r.decodeConfig("service", &ss) {
			log.Warnf("Failed to load 'service' configuration or it does not exist.")
		}
		r.cachedServices = &ss
	})
	return r.cachedServices
}

func (r *resolver) Discovery() interfaces.DiscoveryConfig {
	r.discoveryOnce.Do(func() {
		var discovery configv1.Discovery
		if !r.decodeConfig("discovery", &discovery) {
			log.Warnf("Failed to load 'discovery' configuration or it does not exist.")
		}
		r.cachedDiscovery = &discovery
	})
	return r.cachedDiscovery
}

func (r *resolver) WithDecode(name string, v any, decode func([]byte, any) error) error {
	if v == nil {
		return fmt.Errorf("value %s is nil", name)
	}
	data, err := r.Value(name)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("value %s is nil", name)
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return decode(marshal, v)
}

func (r *resolver) Value(name string) (any, error) {
	v, ok := r.values[name]
	if !ok {
		return nil, fmt.Errorf("value %s not found", name)
	}
	return v, nil
}

func (r *resolver) Middleware() interfaces.MiddlewareConfig {
	r.middlewareOnce.Do(func() {
		var m middlewarev1.Middleware
		if !r.decodeConfig("middleware", &m) {
			log.Warnf("Failed to load 'middleware' configuration or it does not exist.")
		}
		r.cachedMiddleware = &middlewareConfigAdapter{&m}
	})
	return r.cachedMiddleware
}

// middlewareConfigAdapter adapts middlewarev1.Middleware to interfaces.MiddlewareConfig
type middlewareConfigAdapter struct {
	m *middlewarev1.Middleware
}

func (a *middlewareConfigAdapter) GetEnabledMiddlewares() []string {
	return a.m.GetEnabledMiddlewares()
}

func (a *middlewareConfigAdapter) GetJwt() interfaces.MiddlewareJwtConfig {
	return &middlewareJwtConfigAdapter{a.m.GetJwt()}
}

func (a *middlewareConfigAdapter) GetSelector() interfaces.MiddlewareSelectorConfig {
	return &middlewareSelectorConfigAdapter{a.m.GetSelector()}
}

func (a *middlewareConfigAdapter) GetMetadata() interfaces.MiddlewareMetadataConfig {
	return &middlewareMetadataConfigAdapter{a.m.GetMetadata()}
}

func (a *middlewareConfigAdapter) GetRateLimiter() interfaces.MiddlewareRateLimiterConfig {
	return &middlewareRateLimiterConfigAdapter{a.m.GetRateLimiter()}
}

func (a *middlewareConfigAdapter) GetMetrics() interfaces.MiddlewareMetricsConfig {
	return &middlewareMetricsConfigAdapter{a.m.GetMetrics()}
}

func (a *middlewareConfigAdapter) GetValidator() interfaces.MiddlewareValidatorConfig {
	return &middlewareValidatorConfigAdapter{a.m.GetValidator()}
}

// Adapters for sub-configs
type middlewareJwtConfigAdapter struct {
	m *jwt.JWT
}

func (a *middlewareJwtConfigAdapter) GetEnabled() bool {
	return a.m.GetEnabled()
}

type middlewareSelectorConfigAdapter struct {
	m *selector.Selector
}

func (a *middlewareSelectorConfigAdapter) GetEnabled() bool {
	return a.m.GetEnabled()
}

type middlewareMetadataConfigAdapter struct {
	m *middlewarev1.Middleware_Metadata
}

func (a *middlewareMetadataConfigAdapter) GetEnabled() bool {
	return a.m.GetEnabled()
}

func (a *middlewareMetadataConfigAdapter) GetPrefixes() []string {
	return a.m.GetPrefixes()
}

func (a *middlewareMetadataConfigAdapter) GetData() map[string]string {
	return a.m.GetData()
}

type middlewareRateLimiterConfigAdapter struct {
	m *ratelimit.RateLimiter
}

func (a *middlewareRateLimiterConfigAdapter) GetEnabled() bool {
	return a.m.GetEnabled()
}

type middlewareMetricsConfigAdapter struct {
	m *metrics.Metrics
}

func (a *middlewareMetricsConfigAdapter) GetEnabled() bool {
	return a.m.GetEnabled()
}

type middlewareValidatorConfigAdapter struct {
	m *validator.Validator
}

func (a *middlewareValidatorConfigAdapter) GetEnabled() bool {
	return a.m.GetEnabled()
}

func (r *resolver) Logger() interfaces.LoggerConfig {
	r.loggerOnce.Do(func() {
		var l configv1.Logger
		if !r.decodeConfig("logger", &l) {
			log.Warnf("Failed to load 'logger' configuration or it does not exist.")
		}
		r.cachedLogger = &loggerConfigAdapter{&l}
	})
	return r.cachedLogger
}

// loggerConfigAdapter adapts configv1.Logger to interfaces.LoggerConfig
type loggerConfigAdapter struct {
	l *configv1.Logger
}

func (a *loggerConfigAdapter) GetDisabled() bool {
	return a.l.GetDisabled()
}

func (a *loggerConfigAdapter) GetFile() interfaces.LoggerFileConfig {
	return &loggerFileConfigAdapter{a.l.GetFile()}
}

func (a *loggerConfigAdapter) GetStdout() bool {
	return a.l.GetStdout()
}

func (a *loggerConfigAdapter) GetFormat() string {
	return a.l.GetFormat()
}

func (a *loggerConfigAdapter) GetLevel() string {
	return a.l.GetLevel()
}

type loggerFileConfigAdapter struct {
	l *configv1.Logger_File
}

func (a *loggerFileConfigAdapter) GetPath() string {
	return a.l.GetPath()
}

func (a *loggerFileConfigAdapter) GetLumberjack() bool {
	return a.l.GetLumberjack()
}

func (a *loggerFileConfigAdapter) GetMaxSize() int32 {
	return a.l.GetMaxSize()
}

func (a *loggerFileConfigAdapter) GetMaxAge() int32 {
	return a.l.GetMaxAge()
}

func (a *loggerFileConfigAdapter) GetMaxBackups() int32 {
	return a.l.GetMaxBackups()
}

func (a *loggerFileConfigAdapter) GetLocalTime() bool {
	return a.l.GetLocalTime()
}

func (a *loggerFileConfigAdapter) GetCompress() bool {
	return a.l.GetCompress()
}

func (r *resolver) decodeConfig(key string, target interface{}) bool {
	v, ok := r.values[key]
	if !ok {
		return false
	}
	if err := mapstructure.Decode(v, target); err != nil {
		log.Errorf("Failed to decode config key '%s': %v", key, err)
		return false
	}
	return true
}

var DefaultResolver interfaces.Resolver = ResolveFunc(func(config kratosconfig.Config) (interfaces.Resolved, error) {
	var r resolver
	err := config.Scan(&r.values)
	if err != nil {
		return nil, err
	}
	return &r, nil // Return pointer to resolver
})
