package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
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
		r.cachedMiddleware = &m
	})
	return r.cachedMiddleware
}

func (r *resolver) Logger() interfaces.LoggerConfig {
	r.loggerOnce.Do(func() {
		var l configv1.Logger
		if !r.decodeConfig("logger", &l) {
			log.Warnf("Failed to load 'logger' configuration or it does not exist.")
		}
		r.cachedLogger = &l
	})
	return r.cachedLogger
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
