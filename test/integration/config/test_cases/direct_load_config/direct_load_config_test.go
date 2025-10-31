package direct_load_config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/durationpb"

	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	corsv1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/cors/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/runtime/selector/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/runtime/trace/v1"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	"github.com/origadmin/runtime/test/helper"
	parentconfig "github.com/origadmin/runtime/test/integration/config" // Import the parent package for AssertTestConfig
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

var defaultConfig *testconfigs.TestConfig

func init() {
	defaultConfig = &testconfigs.TestConfig{
		App: &appv1.App{
			Id:      "test-app-id",
			Name:    "TestApp",
			Version: "1.0.0",
			Env:     "test",
			Metadata: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		Servers: &transportv1.Servers{
			Servers: []*transportv1.Server{
				{
					Name:     "grpc_server",
					Protocol: "grpc",
					Grpc: &grpcv1.Server{
						Network: "tcp",
						Addr:    ":9000",
						Timeout: durationpb.New(1000000000), // 1s
					},
				},
				{
					Name:     "http_server",
					Protocol: "http",
					Http: &httpv1.Server{
						Network: "tcp",
						Addr:    ":8000",
						Timeout: durationpb.New(2000000000), // 2s
					},
				},
			},
		},
		Client: &transportv1.Client{
			Grpc: &grpcv1.Client{
				Endpoint: "discovery:///user-service",
				Timeout:  durationpb.New(3000000000), // 3s
				Selector: &selectorv1.SelectorConfig{
					Version: "v1.0.0",
				},
			},
		},
		Logger: &loggerv1.Logger{
			Level:  "info",
			Format: "json",
			Stdout: true,
		},
		Discoveries: &discoveryv1.Discoveries{
			Discoveries: []*discoveryv1.Discovery{
				{
					Name: "internal-consul",
					Type: "consul",
					Consul: &discoveryv1.Consul{
						Address: "consul.internal:8500",
					},
				},
				{
					Name: "legacy-etcd",
					Type: "etcd", Etcd: &discoveryv1.ETCD{
					Endpoints: []string{"etcd.legacy:2379"},
				},
				},
			},
		},
		RegistrationDiscoveryName: "internal-consul",
		Trace: &tracev1.Trace{
			Name:     "jaeger",
			Endpoint: "http://jaeger:14268/api/traces",
		},
		Middlewares: &middlewarev1.Middlewares{
			Middlewares: []*middlewarev1.Middleware{
				{
					Name:    "cors-middleware",
					Type:    "cors",
					Enabled: true,
					Cors: &corsv1.Cors{
						AllowedOrigins: []string{"*"},
						AllowedMethods: []string{"GET", "POST"},
					},
				},
			},
		},
	}
}

// DirectLoadConfigTestSuite tests the direct loading of the unified config.yaml.
type DirectLoadConfigTestSuite struct {
	suite.Suite
}

func TestDirectLoadConfigTestSuite(t *testing.T) {
	suite.Run(t, new(DirectLoadConfigTestSuite))
}

// TestDirectConfigLoading verifies that the raw config.yaml file is well-formed and parsable
// into the unified TestConfig struct.
func (s *DirectLoadConfigTestSuite) TestDirectConfigLoading() {
	t := s.T()
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	// Define test cases for different formats of the unified config.
	// For now, we only have YAML, but this structure allows easy expansion.
	testCases := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: "config/test_cases/direct_load_config/config.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// If the config file does not exist, create it from a default TestConfig struct.
			// This ensures the test has a valid config and provides a "live" template for developers.
			if _, err := os.Stat(tc.filePath); os.IsNotExist(err) {
				helper.SaveConfigToFileWithViper(t, defaultConfig, tc.filePath, tc.name)
			}

			var cfg testconfigs.TestConfig
			helper.LoadConfigFromFile(t, tc.filePath, &cfg)
			parentconfig.AssertTestConfig(t, &cfg)
			t.Logf("%s unified config loaded and verified successfully!", tc.name)
		})
	}
}
