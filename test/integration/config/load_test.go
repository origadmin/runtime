package config

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	"github.com/origadmin/toolkits/codec/toml"
	// Import our test-specific, generated bootstrap proto package
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"

	// Import Kratos's XML codec to enable XML format support for loading
	_ "github.com/go-kratos/kratos/v2/encoding/xml"
)

func init() {
	encoding.RegisterCodec(toml.Codec)
}

// loadFromFile is a helper to load a config from a given path into a Bootstrap struct.
func loadFromFile(t *testing.T, path string, bc any) {
	t.Helper()
	cleanPath := filepath.Clean(path)
	ext := filepath.Ext(cleanPath)

	source := file.NewSource(cleanPath)
	c := config.New(config.WithSource(source))
	defer c.Close()

	err := c.Load()
	if !assert.NoError(t, err, "Failed to load config from %s", cleanPath) {
		t.FailNow()
	}

	err = c.Scan(bc)
	if !assert.NoError(t, err, "Failed to scan config into Bootstrap struct for format %s", ext) {
		t.FailNow()
	}
}

// saveToFile is a helper to save a Bootstrap struct to a given path in a specific format.
func saveToFile(t *testing.T, bc *testconfigs.Bootstrap, path string, formatName string) {
	t.Helper()
	var (
		data []byte
		err  error
	)

	switch strings.ToUpper(formatName) {
	case "YAML":
		data, err = yaml.Marshal(bc)
	case "JSON":
		data, err = protojson.Marshal(bc)
	case "TOML":
		data, err = toml.Marshal(bc)
	case "XML":
		data, err = xml.Marshal(bc)
	case "ProtoText":
		opts := prototext.MarshalOptions{Multiline: true, Indent: "  "}
		data, err = opts.Marshal(bc)
	default:
		t.Fatalf("Unsupported format for saving: %s", formatName)
	}

	assert.NoError(t, err, "Failed to marshal data to %s", formatName)
	err = os.WriteFile(path, data, 0644)
	assert.NoError(t, err, "Failed to write temporary file for %s", formatName)
}

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
		//{name: "XML", filePath: "./configs/full_config.xml"},
		//{name: "ProtoText", filePath: "./configs/full_config.pb.txt"},
	}

	// 2. Loop through all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bc testconfigs.Bootstrap
			loadFromFile(t, tc.filePath, &bc)
			//var mapConfig map[string]interface{}
			//loadFromFile(t, tc.filePath, &mapConfig)
			//t.Logf("Loaded %s config: %+v", tc.name, mapConfig)
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
		//{name: "XML", filePath: "./configs/full_config.xml"},
		//{name: "ProtoJSON", filePath: "./configs/full_config.protojson"},
	}

	// 2. Outer loop: Iterate through each format as the SOURCE
	for _, src := range formats {
		// 3. Inner loop: Iterate through each format as the TARGET
		for _, dst := range formats {
			// Define a unique name for this specific conversion test
			testName := fmt.Sprintf("LoadFrom%s_SaveTo%s", src.name, dst.name)

			t.Run(testName, func(t *testing.T) {
				// STEP A: Load the original source file into a Go struct.
				var originalBC testconfigs.Bootstrap
				loadFromFile(t, src.filePath, &originalBC)

				// STEP B: Save the Go struct into the target format in a temporary file.
				tempDir := t.TempDir()
				tempFilePath := filepath.Join(tempDir, "temp_config."+dst.name)

				saveToFile(t, &originalBC, tempFilePath, dst.name)

				// STEP C: Load the newly created target file back into another Go struct.
				var convertedBC testconfigs.Bootstrap
				loadFromFile(t, tempFilePath, &convertedBC)

				// STEP D: The ultimate proof. Assert that the original struct and the
				// converted struct are semantically identical using proto.Equal.
				assert.True(t, proto.Equal(&originalBC, &convertedBC),
					"Struct loaded from %s should be identical to the one saved to %s and loaded back", src.name, dst.name)
			})
		}
	}
}

func generateAllFormatsFromYAML(t *testing.T) {
	// 1. Load YAML configuration using Kratos
	source := file.NewSource("configs/full_config.yaml")
	c := config.New(config.WithSource(source))
	if err := c.Load(); err != nil {
		fmt.Printf("Failed to load YAML config: %v\n", err)
		return
	}
	defer c.Close()

	// 2. Parse configuration into a struct
	var configBootstrap testconfigs.Bootstrap
	if err := c.Scan(&configBootstrap); err != nil {
		fmt.Printf("Failed to scan config: %v\n", err)
		return
	}
	var configsMap map[string]any
	if err := c.Scan(&configsMap); err != nil {
		fmt.Printf("Failed to scan config: %v\n", err)
		return
	}
	// 3. define the supported formats and their encoders
	formats := []struct {
		name string
	}{
		{
			name: "yaml",
		},
		{
			name: "json",
		},
		{
			name: "toml",
		},
		{
			name: "xml",
		},
	}

	// 4. generate and save profiles in various formats
	for _, format := range formats {
		saveToFile(t, &configBootstrap, "full_config."+format.name, format.name)
		t.Logf("Successfully generated %s\n", format.name)
	}
}

func TestGenerateAllFormats(t *testing.T) {
	generateAllFormatsFromYAML(t)
}
