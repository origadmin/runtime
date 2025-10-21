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
	assertions := assert.New(t)

	// App configuration assertions
	assertions.NotNil(cfg.App)
	assertions.Equal("test-app-id", cfg.App.GetId())
	assertions.Equal("TestApp", cfg.App.GetName())
	assertions.Equal("1.0.0", cfg.App.GetVersion())
	assertions.Equal("test", cfg.App.GetEnv())
	assertions.Contains(cfg.App.GetMetadata(), "key1")
	assertions.Contains(cfg.App.GetMetadata(), "key2")
	assertions.Equal("value1", cfg.App.GetMetadata()["key1"])
	assertions.Equal("value2", cfg.App.GetMetadata()["key2"])

	// Server configuration assertions
	assertions.Len(cfg.GetServers().GetServers(), 2, "Should have 1 Servers message")
	serverConfigs := cfg.GetServers().GetServers()
	assertions.Len(serverConfigs, 2, "Should have 2 Server configurations (gRPC and HTTP)")

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

	assertions.NotNil(grpcServer, "gRPC server configuration not found")
	assertions.Equal("tcp", grpcServer.GetGrpc().GetNetwork())
	assertions.Equal(":9000", grpcServer.GetGrpc().GetAddr())
	assertions.Equal("1s", grpcServer.GetGrpc().GetTimeout().AsDuration().String())

	assertions.NotNil(httpServer, "HTTP server configuration not found")
	assertions.Equal("tcp", httpServer.GetHttp().GetNetwork())
	assertions.Equal(":8000", httpServer.GetHttp().GetAddr())
	assertions.Equal("2s", httpServer.GetHttp().GetTimeout().AsDuration().String())

	// Client configuration assertions
	assertions.NotNil(cfg.Client)
	assertions.Equal("discovery:///user-service", cfg.Client.GetEndpoint())
	assertions.Equal("3s", cfg.Client.GetTimeout().AsDuration().String())
	assertions.NotNil(cfg.Client.GetSelector())
	assertions.Equal("v1.0.0", cfg.Client.GetSelector().GetVersion())

	// Discovery configuration assertions
	assertions.Len(cfg.GetDiscoveries().GetDiscoveries(), 2)
	assertions.Equal("internal-consul", cfg.GetDiscoveries().GetDiscoveries()[0].GetName())
	assertions.NotNil(cfg.GetDiscoveries().GetDiscoveries()[0])
	assertions.Equal("consul", cfg.GetDiscoveries().GetDiscoveries()[0].GetType())
	assertions.NotNil(cfg.GetDiscoveries().GetDiscoveries()[0].GetConsul())
	assertions.Equal("consul.internal:8500", cfg.GetDiscoveries().GetDiscoveries()[0].GetConsul().GetAddress())

	assertions.Equal("legacy-etcd", cfg.GetDiscoveries().GetDiscoveries()[1].GetName())
	assertions.NotNil(cfg.GetDiscoveries().GetDiscoveries()[1])
	assertions.Equal("etcd", cfg.GetDiscoveries().GetDiscoveries()[1].GetType())
	assertions.NotNil(cfg.GetDiscoveries().GetDiscoveries()[1].GetEtcd())
	assertions.Len(cfg.GetDiscoveries().GetDiscoveries()[1].GetEtcd().GetEndpoints(), 1)
	assertions.Equal("etcd.legacy:2379", cfg.GetDiscoveries().GetDiscoveries()[1].GetEtcd().GetEndpoints()[0])

	// Registration discovery name assertion
	assertions.Equal("internal-consul", cfg.GetRegistrationDiscoveryName())

	// Logger configuration assertions
	assertions.NotNil(cfg.Logger)
	assertions.Equal("info", cfg.Logger.GetLevel())
	assertions.Equal("json", cfg.Logger.GetFormat())
	assertions.True(cfg.Logger.GetStdout())

	// Tracer configuration assertions
	assertions.NotNil(cfg.Tracer)
	assertions.Equal("jaeger", cfg.Tracer.GetName())
	assertions.Equal("http://jaeger:14268/api/traces", cfg.Tracer.GetEndpoint())

	// Middleware configuration assertions
	assertions.NotNil(cfg.Middlewares)
	assertions.Len(cfg.Middlewares.GetMiddlewares(), 1, "Should have 1 middleware configured")

	corsMiddleware := cfg.Middlewares.GetMiddlewares()[0]
	assertions.Equal("cors", corsMiddleware.GetType())
	assertions.True(corsMiddleware.GetEnabled())
	assertions.NotNil(corsMiddleware.GetCors(), "CORS config should not be nil for middleware of type cors")
	assertions.Len(corsMiddleware.GetCors().GetAllowedOrigins(), 1)
	assertions.Equal("*", corsMiddleware.GetCors().GetAllowedOrigins()[0])
	assertions.Len(corsMiddleware.GetCors().GetAllowedMethods(), 2)
	assertions.Equal("GET", corsMiddleware.GetCors().GetAllowedMethods()[0])
	assertions.Equal("POST", corsMiddleware.GetCors().GetAllowedMethods()[1])
}
