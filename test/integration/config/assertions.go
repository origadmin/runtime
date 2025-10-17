package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// AssertTestConfig contains all the validation logic for the unified TestConfig struct.
// It is reused across all configuration loading tests to ensure consistency.
func AssertTestConfig(t *testing.T, cfg *testconfigs.TestConfig) {
	assert := assert.New(t)

	// App configuration assertions
	assert.NotNil(cfg.App)
	assert.Equal("test-app-id", cfg.App.GetId())
	assert.Equal("TestApp", cfg.App.GetName())
	assert.Equal("1.0.0", cfg.App.GetVersion())
	assert.Equal("test", cfg.App.GetEnv())
	assert.Contains(cfg.App.GetMetadata(), "key1")
	assert.Contains(cfg.App.GetMetadata(), "key2")
	assert.Equal("value1", cfg.App.GetMetadata()["key1"])
	assert.Equal("value2", cfg.App.GetMetadata()["key2"])

	// Server configuration assertions
	assert.Len(cfg.GrpcServers, 1)
	assert.Equal("tcp", cfg.GrpcServers[0].GetNetwork())
	assert.Equal(":9000", cfg.GrpcServers[0].GetAddr())
	assert.Equal("1s", cfg.GrpcServers[0].GetTimeout().AsDuration().String())

	assert.Len(cfg.HttpServers, 1)
	assert.Equal("tcp", cfg.HttpServers[0].GetNetwork())
	assert.Equal(":8000", cfg.HttpServers[0].GetAddr())
	assert.Equal("2s", cfg.HttpServers[0].GetTimeout().AsDuration().String())

	// Client configuration assertions
	assert.NotNil(cfg.Client)
	assert.Equal("discovery:///user-service", cfg.Client.GetEndpoint())
	assert.Equal("3s", cfg.Client.GetTimeout().AsDuration().String())
	assert.NotNil(cfg.Client.GetSelector())
	assert.Equal("v1.0.0", cfg.Client.GetSelector().GetVersion())

	// Discovery configuration assertions
	assert.Len(cfg.Discoveries, 2)
	assert.Equal("internal-consul", cfg.Discoveries[0].GetName())
	assert.NotNil(cfg.Discoveries[0].GetConfig())
	assert.Equal("consul", cfg.Discoveries[0].GetConfig().GetType())
	assert.NotNil(cfg.Discoveries[0].GetConfig().GetConsul())
	assert.Equal("consul.internal:8500", cfg.Discoveries[0].GetConfig().GetConsul().GetAddress())

	assert.Equal("legacy-etcd", cfg.Discoveries[1].GetName())
	assert.NotNil(cfg.Discoveries[1].GetConfig())
	assert.Equal("etcd", cfg.Discoveries[1].GetConfig().GetType())
	assert.NotNil(cfg.Discoveries[1].GetConfig().GetEtcd())
	assert.Len(cfg.Discoveries[1].GetConfig().GetEtcd().GetEndpoints(), 1)
	assert.Equal("etcd.legacy:2379", cfg.Discoveries[1].GetConfig().GetEtcd().GetEndpoints()[0])

	// Registration discovery name assertion
	assert.Equal("internal-consul", cfg.GetRegistrationDiscoveryName())

	// Logger configuration assertions
	assert.NotNil(cfg.Logger)
	assert.Equal("info", cfg.Logger.GetLevel())
	assert.Equal("json", cfg.Logger.GetFormat())
	assert.True(cfg.Logger.GetStdout())

	// Tracer configuration assertions
	assert.NotNil(cfg.Tracer)
	assert.Equal("jaeger", cfg.Tracer.GetName())
	assert.Equal("http://jaeger:14268/api/traces", cfg.Tracer.GetEndpoint())

	// Middleware configuration assertions
	assert.NotNil(cfg.Middleware)
	assert.Len(cfg.Middleware.GetMiddlewares(), 1, "Should have 1 middleware configured")

	corsMiddleware := cfg.Middleware.GetMiddlewares()[0]
	assert.Equal("cors", corsMiddleware.GetType())
	assert.True(corsMiddleware.GetEnabled())
	assert.NotNil(corsMiddleware.GetCors(), "CORS config should not be nil for middleware of type cors")
	assert.Len(corsMiddleware.GetCors().GetAllowOrigins(), 1)
	assert.Equal("*", corsMiddleware.GetCors().GetAllowOrigins()[0])
	assert.Len(corsMiddleware.GetCors().GetAllowMethods(), 2)
	assert.Equal("GET", corsMiddleware.GetCors().GetAllowMethods()[0])
	assert.Equal("POST", corsMiddleware.GetCors().GetAllowMethods()[1])
}
