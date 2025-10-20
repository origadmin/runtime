package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
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
	assert.Len(cfg.GetServers(), 2, "Should have 1 Servers message")
	serverConfigs := cfg.GetServers()[0].GetServers()
	assert.Len(serverConfigs, 2, "Should have 2 Server configurations (gRPC and HTTP)")

	var grpcServer *transportv1.Server
	var httpServer *transportv1.Server

	for _, s := range serverConfigs {
		if s.GetGrpc() != nil {
			grpcServer = s
		}
		if s.GetHttp() != nil {
			httpServer = s
		}
	}

	assert.NotNil(grpcServer, "gRPC server configuration not found")
	assert.Equal("tcp", grpcServer.GetGrpc().GetNetwork())
	assert.Equal(":9000", grpcServer.GetGrpc().GetAddr())
	assert.Equal("1s", grpcServer.GetGrpc().GetTimeout().AsDuration().String())

	assert.NotNil(httpServer, "HTTP server configuration not found")
	assert.Equal("tcp", httpServer.GetHttp().GetNetwork())
	assert.Equal(":8000", httpServer.GetHttp().GetAddr())
	assert.Equal("2s", httpServer.GetHttp().GetTimeout().AsDuration().String())

	// Client configuration assertions
	assert.NotNil(cfg.Client)
	assert.Equal("discovery:///user-service", cfg.Client.GetEndpoint())
	assert.Equal("3s", cfg.Client.GetTimeout().AsDuration().String())
	assert.NotNil(cfg.Client.GetSelector())
	assert.Equal("v1.0.0", cfg.Client.GetSelector().GetVersion())

	// Discovery configuration assertions
	assert.Len(cfg.GetDiscoveries().GetDiscoveries(), 2)
	assert.Equal("internal-consul", cfg.GetDiscoveries().GetDiscoveries()[0].GetName())
	assert.NotNil(cfg.GetDiscoveries().GetDiscoveries()[0])
	assert.Equal("consul", cfg.GetDiscoveries().GetDiscoveries()[0].GetType())
	assert.NotNil(cfg.GetDiscoveries().GetDiscoveries()[0].GetConsul())
	assert.Equal("consul.internal:8500", cfg.GetDiscoveries().GetDiscoveries()[0].GetConsul().GetAddress())

	assert.Equal("legacy-etcd", cfg.GetDiscoveries().GetDiscoveries()[1].GetName())
	assert.NotNil(cfg.GetDiscoveries().GetDiscoveries()[1])
	assert.Equal("etcd", cfg.GetDiscoveries().GetDiscoveries()[1].GetType())
	assert.NotNil(cfg.GetDiscoveries().GetDiscoveries()[1].GetEtcd())
	assert.Len(cfg.GetDiscoveries().GetDiscoveries()[1].GetEtcd().GetEndpoints(), 1)
	assert.Equal("etcd.legacy:2379", cfg.GetDiscoveries().GetDiscoveries()[1].GetEtcd().GetEndpoints()[0])

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
	assert.NotNil(cfg.Middlewares)
	assert.Len(cfg.Middlewares.GetMiddlewares(), 1, "Should have 1 middleware configured")

	corsMiddleware := cfg.Middlewares.GetMiddlewares()[0]
	assert.Equal("cors", corsMiddleware.GetType())
	assert.True(corsMiddleware.GetEnabled())
	assert.NotNil(corsMiddleware.GetCors(), "CORS config should not be nil for middleware of type cors")
	assert.Len(corsMiddleware.GetCors().GetAllowedOrigins(), 1)
	assert.Equal("*", corsMiddleware.GetCors().GetAllowedOrigins()[0])
	assert.Len(corsMiddleware.GetCors().GetAllowedMethods(), 2)
	assert.Equal("GET", corsMiddleware.GetCors().GetAllowedMethods()[0])
	assert.Equal("POST", corsMiddleware.GetCors().GetAllowedMethods()[1])
}
