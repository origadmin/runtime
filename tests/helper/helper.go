package helper

import (
	"encoding/json"
	"fmt"
	"testing"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/helpers/optionutil"
)

// LoadConfigFromFile loads configuration from a file and scans into target.
func LoadConfigFromFile(t *testing.T, path string, target any) (runtimeconfig.KConfig, error) {
	c := kratosconfig.New(
		kratosconfig.WithSource(
			file.NewSource(path),
		),
	)
	err := c.Load()
	require.NoError(t, err)
	if target != nil {
		err = c.Scan(target)
		require.NoError(t, err)
	}
	return c, nil
}

// MockConsulSource is a mock implementation of a config source for testing purposes.
type MockConsulSource struct {
	config *sourcev1.SourceConfig
	// mockEntries stores the data for each key, along with its intended format.
	mockEntries map[string]struct {
		Value  interface{}
		Format string
	}
}

// WithMockDataString provides mock data for MockConsulSource, allowing specification of the format for string values.
func WithMockDataString(data map[string]string, format string) options.Option {
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
				Format: format,
			}
		}
	})
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

var _ runtimeconfig.SourceFactory = (*MockConsulSource)(nil)

// Register MockConsulSource as a config source factory
func init() {
	runtimeconfig.RegisterSourceFactory("consul", &MockConsulSource{})
}
