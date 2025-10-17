package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/config/file"
)

// LoadYAMLConfig loads a YAML configuration file and converts it to a JSON format byte array.
func LoadYAMLConfig(filename string) ([]byte, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(currentDir, filename)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var configMap map[string]interface{}
	if err := yaml.Unmarshal(data, &configMap); err != nil {
		return nil, err
	}

	return yaml.Marshal(configMap)
}

// LoadMiddleware loads and parses a middleware configuration file.
func LoadMiddleware(t *testing.T, configPath string) (*middlewarev1.Middlewares, error) {
	t.Helper()

	source := file.NewSource(configPath)
	c := config.New(config.WithSource(source))
	if err := c.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	defer c.Close()

	var middlewares middlewarev1.Middlewares
	if err := c.Scan(&middlewares); err != nil {
		return nil, fmt.Errorf("failed to unmarshal middleware config: %w", err)
	}

	return &middlewares, nil
}

// SetupIntegrationTest sets up the environment for integration tests.
// It changes the working directory to the runtime module's root, allowing
// tests to use consistent relative paths for configuration files.
// It returns a cleanup function that restores the original working directory.
func SetupIntegrationTest(t *testing.T) func() {
	t.Helper()

	// Get the directory of the test file that called this helper.
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("Failed to get caller info for test setup")
	}

	// Assuming the caller is in .../test/integration/config, we go up 3 levels to the runtime root.
	runtimeRoot := filepath.Join(filepath.Dir(filename), "..", "..", "..")

	// Save the original working directory.
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}

	// Change to the runtime module root.
	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	// Return a cleanup function to restore the original directory.
	return func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}
}

// LoadConfigFromFile loads a configuration from the specified file path into a proto.Message.
func LoadConfigFromFile(t *testing.T, filePath string, v proto.Message) {
	t.Helper()
	source := file.NewSource(filePath)
	c := config.New(config.WithSource(source))
	if err := c.Load(); err != nil {
		t.Fatalf("Failed to load config from %s: %v", filePath, err)
	}
	defer c.Close()

	if err := c.Scan(v); err != nil {
		t.Fatalf("Failed to scan config from %s into struct: %v", filePath, err)
	}
}

// MockConsulSource is a mock implementation of config.Source for Consul.
type MockConsulSource struct {
	data map[string]string
}

// NewMockConsulSource creates a new MockConsulSource.
func NewMockConsulSource(data map[string]string) *MockConsulSource {
	return &MockConsulSource{data: data}
}

// Load returns the mock data as KeyValue pairs.
func (m *MockConsulSource) Load() ([]*config.KeyValue, error) {
	kvs := make([]*config.KeyValue, 0, len(m.data))
	for k, v := range m.data {
		kvs = append(kvs, &config.KeyValue{
			Key:    k,
			Value:  []byte(v),
			Format: "yaml", // Assuming YAML format for simplicity in mock
		})
	}
	return kvs, nil
}

// Watch is not implemented for the mock source.
func (m *MockConsulSource) Watch() (config.Watcher, error) {
	return nil, fmt.Errorf("watch not implemented for MockConsulSource")
}

// String returns the name of the mock source.
func (m *MockConsulSource) String() string {
	return "mock-consul-source"
}
