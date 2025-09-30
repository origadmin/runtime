package config_test

import (
	"os"
	"path/filepath"
	stdruntime "runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto" // Import generated proto package
)

func TestConfigProtoIntegration(t *testing.T) {
	assert := assert.New(t)

	// Get the current file's directory to construct an absolute path for the config.
	_, filename, _, ok := stdruntime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file info")
	}
	currentDir := filepath.Dir(filename)

	// Calculate the runtime module root directory.
	// From .../runtime/test/integration/config, go up 3 levels to .../runtime
	runtimeRoot := filepath.Join(currentDir, "../../..")

	// Store original CWD and defer its restoration.
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	defer func() {
		err := os.Chdir(originalCwd)
		if err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	// Change CWD to the runtime module root.
	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	// The bootstrapPath is now relative to the runtime module root.
	bootstrapPath := "test/integration/config/configs/test_bootstrap.yaml"

	// 1. Initialize Runtime with the default decoder provider.
	rt, cleanup, err := runtime.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-proto-config",
			Name:    "TestProtoConfig",
			Version: "1.0.0",
			//Env:     "test",
		}),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer cleanup()

	// 2. Get the ConfigDecoder from the runtime.
	configDecoder := rt.Config()
	assert.NotNil(configDecoder)

	// 3. Decode the entire configuration into the generated Bootstrap struct.
	var bootstrapConfig testconfigs.Bootstrap
	err = configDecoder.Decode("", &bootstrapConfig)
	assert.NoError(err)

	// 4. Assert decoded values.
	// Verify logger (from test_config.yaml)
	logger := rt.Logger()
	assert.NotNil(logger)
	// Further assertions on logger can be added if needed, e.g., level, format.

	// Verify registration_discovery_name
	assert.Equal("test-discovery", bootstrapConfig.RegistrationDiscoveryName)

	// Verify servers
	assert.Len(bootstrapConfig.GrpcServers, 1) // Check GrpcServers length
	assert.Len(bootstrapConfig.HttpServers, 1) // Check HttpServers length

	// GRPC Server
	grpcServer := bootstrapConfig.GrpcServers[0]
	assert.NotNil(grpcServer)
	assert.Equal("tcp", grpcServer.Network)
	assert.Equal(":9000", grpcServer.Addr)
	assert.Equal("1s", grpcServer.Timeout.AsDuration().String())

	// HTTP Server
	httpServer := bootstrapConfig.HttpServers[0]
	assert.NotNil(httpServer)
	assert.Equal("tcp", httpServer.Network)
	assert.Equal(":8000", httpServer.Addr)
	assert.Equal("2s", httpServer.Timeout.AsDuration().String())

	// Verify clients (not defined in test_config.yaml, so should be empty)
	assert.Empty(bootstrapConfig.Clients)

	// Verify discoveries (not defined in test_config.yaml, so should be empty)
	assert.Empty(bootstrapConfig.Discoveries)
}
