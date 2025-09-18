package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	// Import our test-specific, generated bootstrap proto package
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"

	// Import all necessary Kratos codecs to enable format support
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/toml"
	_ "github.com/go-kratos/kratos/v2/encoding/xml"
	_ "github.com/go-kratos/kratos/v2/encoding/yaml"
)

// runAssertions contains all the validation logic for the Bootstrap config.
// It is reused across all configuration format tests to ensure consistency.
func runAssertions(t *testing.T, bc *testconfigs.Bootstrap) {
	// Validate discovery pool
	assert.Len(t, bc.Discoveries, 2, "Should have 2 discovery configurations")
	assert.Equal(t, "internal-consul", bc.Discoveries[0].Name)
	assert.Equal(t, "my-test-app", bc.Discoveries[0].Config.ServiceName)
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
	assert.Equal(t, "user-service", bc.Clients[0].Name)
	assert.Equal(t, "internal-consul", bc.Clients[0].DiscoveryName, "user-service client should use internal-consul")
	assert.Equal(t, "v1.5.0", bc.Clients[0].Selector.Version)
	// Validate second client
	assert.Equal(t, "stock-service", bc.Clients[1].Name)
	assert.Equal(t, "legacy-etcd", bc.Clients[1].DiscoveryName, "stock-service client should use legacy-etcd")
	assert.Equal(t, "v1.0.1", bc.Clients[1].Selector.Version)
}

// TestMultiFormatConfigLoading uses a table-driven approach to test loading
// configurations from all supported file formats.
func TestMultiFormatConfigLoading(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: "./configs/full_config.yaml"},
		{name: "JSON", filePath: "./configs/full_config.json"},
		{name: "TOML", filePath: "./configs/full_config.toml"},
		{name: "XML", filePath: "./configs/full_config.xml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create Kratos config instance
			source := file.NewSource(filepath.Clean(tc.filePath))
			c := config.New(config.WithSource(source))
			defer c.Close()

			// Load and Scan
			assert.NoError(t, c.Load(), "Failed to load config")
			var bc testconfigs.Bootstrap
			assert.NoError(t, c.Scan(&bc), "Failed to scan config")

			// Run assertions
			runAssertions(t, &bc)
			t.Logf("%s config loaded and verified successfully!", tc.name)
		})
	}
}

// TestConfigInteroperability is the definitive test for format conversion robustness.
// It now correctly uses the Kratos config instance and its codecs for the entire process.
func TestConfigInteroperability(t *testing.T) {
	formats := []struct {
		name      string
		filePath  string
		codecName string // The name registered in Kratos's encoding package
	}{
		{name: "YAML", filePath: "./configs/full_config.yaml", codecName: "yaml"},
		{name: "JSON", filePath: "./configs/full_config.json", codecName: "json"},
		{name: "TOML", filePath: "./configs/full_config.toml", codecName: "toml"},
		{name: "XML", filePath: "./configs/full_config.xml", codecName: "xml"},
	}

	for _, src := range formats {
		for _, dst := range formats {
			testName := fmt.Sprintf("LoadFrom%s_ConvertTo_%s", src.name, dst.name)

			t.Run(testName, func(t *testing.T) {
				// 1. Load the source file using Kratos config
				srcSource := file.NewSource(filepath.Clean(src.filePath))
				srcConf := config.New(config.WithSource(srcSource))
				assert.NoError(t, srcConf.Load())
				defer srcConf.Close()

				// 2. Get the raw, format-agnostic value from the source config
				var rawValue map[string]interface{}
				assert.NoError(t, srcConf.Scan(&rawValue))

				// 3. Get the target format's codec from Kratos's encoding registry
				codec := encoding.GetCodec(dst.codecName)
				assert.NotNil(t, codec)

				// 4. Marshal the raw value into the target format using the codec
				dstData, err := codec.Marshal(rawValue)
				assert.NoError(t, err)

				// 5. Write the marshaled data to a temporary file
				tempDir := t.TempDir()
				tempFilePath := filepath.Join(tempDir, "temp_config."+dst.codecName)
				err = os.WriteFile(tempFilePath, dstData, 0644)
				assert.NoError(t, err)

				// 6. Create a new Kratos config instance from the temporary file
				dstSource := file.NewSource(tempFilePath)
				dstConf := config.New(config.WithSource(dstSource))
				assert.NoError(t, dstConf.Load())
				defer dstConf.Close()

				// 7. Scan both configs into structs and compare them
				var originalBC, convertedBC testconfigs.Bootstrap
				assert.NoError(t, srcConf.Scan(&originalBC))
				assert.NoError(t, dstConf.Scan(&convertedBC))

				// 8. The ultimate proof: ensure they are semantically identical.
				assert.True(t, proto.Equal(&originalBC, &convertedBC),
					"Struct from %s should be identical to struct converted to %s and back", src.name, dst.name)
			})
		}
	}
}
