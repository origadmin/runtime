package helper

import (
	"encoding/json" // Add json import
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	filesource "github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/extensions/optionutil"
	"github.com/origadmin/runtime/interfaces/options"
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

	// Fix: Marshal to JSON as per the function comment.
	return json.Marshal(configMap)
}

// LoadMiddleware loads and parses a middleware configuration file using the framework's config.
func LoadMiddleware(t *testing.T, configPath string) (*middlewarev1.Middlewares, error) {
	t.Helper()

	source := filesource.NewSource(configPath)
	c := runtimeconfig.NewKConfig(runtimeconfig.WithKSource(source))
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

// LoadConfigFromFile loads a configuration from the specified file path into a proto.Message
// using the framework's config package.
func LoadConfigFromFile(t *testing.T, configPath string, v proto.Message) {
	t.Helper()
	source := filesource.NewSource(configPath)
	c := runtimeconfig.NewKConfig(runtimeconfig.WithKSource(source))
	if err := c.Load(); err != nil {
		t.Fatalf("Failed to load config from %s: %v", configPath, err)
	}
	defer c.Close()

	if err := c.Scan(v); err != nil {
		t.Fatalf("Failed to scan config from %s into struct: %v", configPath, err)
	}
}

// mockWatcher implements frameworkConfig.Watcher for testing purposes.
type mockWatcher struct {
	data map[string]string
}

func (mw *mockWatcher) Next() ([]*kratosconfig.KeyValue, error) {
	kvs := make([]*runtimeconfig.KKeyValue, 0, len(mw.data))
	for k, v := range mw.data {
		kvs = append(kvs, &runtimeconfig.KKeyValue{
			Key:    k,
			Value:  []byte(v),
			Format: "yaml", // Assuming YAML format for simplicity in mock
		})
	}
	return kvs, nil
}

// Stop is a no-op for the mock watcher.
func (mw *mockWatcher) Stop() error {
	return nil
}

// MockConsulSource is a mock implementation of runtimeconfig.Source for Consul.
type MockConsulSource struct {
	config *sourcev1.SourceConfig
	// mockEntries stores the data for each key, along with its intended format.
	mockEntries map[string]struct {
		Value  interface{} // Can be string, map[string]interface{}, etc.
		Format string      // "string", "json", "yaml"
	}
}

// WithMockData provides mock data for MockConsulSource.
// It expects string values, which will be treated as YAML content.
// This function maintains backward compatibility with the original signature.
func WithMockData(data map[string]string) func(*MockConsulSource) {
	return func(m *MockConsulSource) {
		if m.mockEntries == nil {
			m.mockEntries = make(map[string]struct {
				Value  interface{}
				Format string
			})
		}
		for k, v := range data {
			m.mockEntries[k] = struct {
				Value  interface{}
				Format string
			}{
				Value:  v,
				Format: "yaml", // Original behavior was to treat these as YAML strings
			}
		}
	}
}

// WithMockDataString provides mock data for MockConsulSource, allowing specification of the format for string values.
func WithMockDataString(data map[string]string, format string) func(*MockConsulSource) {
	return func(m *MockConsulSource) {
		if m.mockEntries == nil {
			m.mockEntries = make(map[string]struct {
				Value  interface{}
				Format string
			})
		}
		for k, v := range data {
			m.mockEntries[k] = struct {
				Value  interface{}
				Format string
			}{
				Value:  v,
				Format: format,
			}
		}
	}
}

// WithMockDataJSON provides mock data for MockConsulSource, marshaling values to JSON.
func WithMockDataJSON(data map[string]interface{}) options.Option {
	return optionutil.Update(func(m *MockConsulSource) {
		if m.mockEntries == nil {
			m.mockEntries = make(map[string]struct {
				Value  interface{}
				Format string
			})
		}
		for k, v := range data {
			m.mockEntries[k] = struct {
				Value  interface{}
				Format string
			}{
				Value:  v,
				Format: "json",
			}
		}
	})
}

// WithMockDataYAML provides mock data for MockConsulSource, marshaling values to YAML.
func WithMockDataYAML(data map[string]interface{}) options.Option {
	return optionutil.Update(func(o *MockConsulSource) {
		if o.mockEntries == nil {
			o.mockEntries = make(map[string]struct {
				Value  interface{}
				Format string
			})
		}
		for k, v := range data {
			o.mockEntries[k] = struct {
				Value  interface{}
				Format string
			}{
				Value:  v,
				Format: "yaml",
			}
		}
	})
}

// NewSource creates a new instance of MockConsulSource.
func (m *MockConsulSource) NewSource(config *sourcev1.SourceConfig, opts ...options.Option) (kratosconfig.Source,
	error) {
	optionutil.Apply(m, opts...)
	m.config = config
	return m, nil
}

// Load returns the mock data as KeyValue pairs.
func (m *MockConsulSource) Load() ([]*runtimeconfig.KKeyValue, error) {
	kvs := make([]*runtimeconfig.KKeyValue, 0, len(m.mockEntries))
	for k, entry := range m.mockEntries {
		var valueBytes []byte
		var err error

		switch entry.Format {
		case "json":
			valueBytes, err = json.Marshal(entry.Value)
		case "yaml":
			valueBytes, err = yaml.Marshal(entry.Value)
		default:
			// If format is unknown, or if entry.Value is already []byte or string, use it directly.
			if b, ok := entry.Value.([]byte); ok {
				valueBytes = b
			} else if s, ok := entry.Value.(string); ok {
				valueBytes = []byte(s)
			} else {
				// If it's not a string or []byte, and format is not json/yaml, default to json marshal
				valueBytes, err = json.Marshal(entry.Value)
				entry.Format = "json" // Update format if we marshaled as JSON
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to marshal mock data for key %s (format %s): %w", k, entry.Format, err)
		}

		kvs = append(kvs, &runtimeconfig.KKeyValue{
			Key:    k,
			Value:  valueBytes,
			Format: entry.Format, // Use the specified or inferred format
		})
	}
	return kvs, nil
}

// Watch is not implemented for the mock source.
func (m *MockConsulSource) Watch() (runtimeconfig.KWatcher, error) {
	return &noopWatcher{}, nil
}

// A no-op watcher for MockConsulSource
type noopWatcher struct{}

func (nw *noopWatcher) Next() ([]*kratosconfig.KeyValue, error) {
	select {} // Block indefinitely, simulating no new events
}

func (nw *noopWatcher) Stop() error {
	return nil
}

// String returns the name of the mock source.
func (m *MockConsulSource) String() string {
	return "mock-consul-source"
}

// var _ runtimeconfig.SourceFactory = (*MockConsulSource)(nil) // Removed this line
