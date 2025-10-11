package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/test/helper"
	// Import our test-specific, generated bootstrap proto package
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"

	// Import Kratos's XML codec to enable XML format support for loading
	_ "github.com/go-kratos/kratos/v2/encoding/xml"
)

// runAssertions contains all the validation logic for the Bootstrap config.
// It is reused across all configuration format tests to ensure consistency.
func runAssertions(t *testing.T, bc *testconfigs.Bootstrap) {
	// Validate discovery pool
	assert.Len(t, bc.Discoveries, 2, "Should have 2 discovery configurations")
	assert.Equal(t, "internal-consul", bc.Discoveries[0].Name)
	assert.Equal(t, "my-test-app", bc.Discoveries[0].Config.Name)
	assert.Equal(t, "consul.internal:8500", bc.Discoveries[0].Config.Consul.Address)
	assert.Equal(t, "legacy-etcd", bc.Discoveries[1].Name)
	assert.Equal(t, "etcd.legacy:2379", bc.Discoveries[1].Config.Etcd.Endpoints[0])

	// Validate service registration name
	assert.Equal(t, "internal-consul", bc.RegistrationDiscoveryName)

	// Validate service endpoints
	assert.Len(t, bc.GrpcServers, 1)
	assert.Equal(t, ":9001", bc.GrpcServers[0].Addr)
	assert.Len(t, bc.HttpServers, 1)
	assert.Equal(t, ":8001", bc.HttpServers[0].Addr)

	// Validate clients (most critical part)
	assert.Len(t, bc.Clients, 2, "Should have 2 client configurations")
	// Validate first client
	//assert.Equal(t, "user-service", bc.Clients[0].Name)
	//assert.Equal(t, "internal-consul", bc.Clients[0].DiscoveryName, "user-service client should use internal-consul")
	//assert.Equal(t, "v1.5.0", bc.Clients[0].Selector.Version)
	//// Validate second client
	//assert.Equal(t, "stock-service", bc.Clients[1].Name)
	//assert.Equal(t, "legacy-etcd", bc.Clients[1].DiscoveryName, "stock-service client should use legacy-etcd")
	//assert.Equal(t, "v1.0.1", bc.Clients[1].Selector.Version)
}

// TestMultiFormatConfigLoading uses a table-driven approach to test loading
// configurations from all supported file formats.
func TestMultiFormatConfigLoading(t *testing.T) {
	// 1. Define the test cases for each format
	testCases := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: "./configs/full_config.yaml"},
		{name: "JSON", filePath: "./configs/full_config.json"},
		{name: "TOML", filePath: "./configs/full_config.toml"},
		//{name: "INI", filePath: "./configs/full_config.ini"},
		//{name: "HCL", filePath: "./configs/full_config.hcl"},
		//{name: "ENV", filePath: "./configs/full_config.env"},
		//{name: "PROPERTIES", filePath: "./configs/full_config.properties"},
		//{name: "XML", filePath: "./configs/full_config.xml"},
		//{name: "ProtoText", filePath: "./configs/full_config.pb.txt"},
	}

	// 2. Loop through all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bc testconfigs.Bootstrap
			helper.LoadConfigFromFile(t, tc.filePath, &bc)
			// 3. --- Run the exact same assertion validation for every format ---
			runAssertions(t, &bc)
			t.Logf("%s config loaded and verified successfully!", tc.name)
		})
	}
}

// TestConfigInteroperability provides the ultimate proof of robustness by testing
// that a configuration loaded from any format can be saved to any other format
// without data loss.
func TestConfigInteroperability(t *testing.T) {
	// 1. Define all supported formats and their file paths.
	formats := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: "./configs/full_config.yaml"},
		{name: "JSON", filePath: "./configs/full_config.json"},
		{name: "TOML", filePath: "./configs/full_config.toml"},
		//{name: "INI", filePath: "./configs/full_config.ini"},
		//{name: "HCL", filePath: "./configs/full_config.hcl"},
		//{name: "ENV", filePath: "./configs/full_config.env"},
		//{name: "PROPERTIES", filePath: "./configs/full_config.properties"},
		//{name: "XML", filePath: "./configs/full_config.xml"},
		//{name: "ProtoJSON", filePath: "./configs/full_config.protojson"},
	}
	// 仅 yaml/json/toml/hcl 当前受支持，其他作为占位并在运行时跳过
	//supportedInterop := map[string]bool{"YAML": true, "JSON": true, "TOML": true, "HCL": true}

	// 2. Outer loop: Iterate through each format as the SOURCE
	for _, src := range formats {
		// 3. Inner loop: Iterate through each format as the TARGET
		for _, dst := range formats {
			// Define a unique name for this specific conversion test
			testName := fmt.Sprintf("LoadFrom%s_SaveTo%s", src.name, dst.name)
			name := strings.ToLower(dst.name)
			t.Run(testName, func(t *testing.T) {
				// STEP A: Load the original source file into a Go struct.
				var originalBC testconfigs.Bootstrap
				helper.LoadConfigFromFile(t, src.filePath, &originalBC)

				// STEP B: Save the Go struct into the target format in a temporary file.
				tempDir := t.TempDir()
				tempFilePath := filepath.Join(tempDir, "temp_config."+name)

				// Only save to formats that have encoders implemented
				if !(name == "yaml" || name == "json" || name == "toml") {
					t.Skipf("skip saving to unsupported target format: %s", dst.name)
				}
				helper.SaveConfigToFile(t, &originalBC, tempFilePath, name)

				// STEP C: Load the newly created target file back into another Go struct.
				var convertedBC testconfigs.Bootstrap
				helper.LoadConfigFromFile(t, tempFilePath, &convertedBC)

				// STEP D: The ultimate proof. Assert that the original struct and the
				// converted struct are semantically identical using proto.Equal.
				assert.True(t, proto.Equal(&originalBC, &convertedBC),
					"Struct loaded from %s should be identical to the one saved to %s and loaded back", src.name, dst.name)
			})
		}
	}
}

func init() {
	if err := generateAllFormatsFromYAML(); err != nil {
		panic(fmt.Errorf("failed to generate test files: %v", err))
	}
}

func generateAllFormatsFromYAML() error {
	// 1. Load YAML configuration using Kratos
	source := file.NewSource("full_config_source.yaml")
	c := config.NewKConfig(config.WithKSource(source))
	if err := c.Load(); err != nil {
		return fmt.Errorf("failed to load YAML config: %v", err)
	}
	defer c.Close()

	// 2. Parse configuration into a struct
	var configBootstrap testconfigs.Bootstrap
	if err := c.Scan(&configBootstrap); err != nil {
		return fmt.Errorf("failed to scan config: %v", err)
	}

	// 3. Define the supported formats and their encoders
	formats := []struct {
		name string
	}{
		{name: "yaml"},
		{name: "json"},
		{name: "toml"},
		//{name: "ini"},
		//{name: "hcl"},
		//{name: "env"},
		//{name: "properties"},
	}

	// 4. Generate and save profiles in various formats
	for _, format := range formats {
		//if !supportedGen[format.name] {
		//	continue
		//}
		// Create a minimal testing.T object for logging
	t := &testing.T{}
	helper.SaveConfigToFile(t, &configBootstrap, "configs/full_config."+format.name, format.name)
	}
	return nil
}


