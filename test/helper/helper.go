package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"gopkg.in/yaml.v3"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/config/file"
)

// LoadYAMLConfig loads YAML configuration file and converts it to JSON format byte array
func LoadYAMLConfig(filename string) ([]byte, error) {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Build complete configuration file path
	configPath := filepath.Join(currentDir, filename)

	// Read configuration file content
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML to map for conversion to JSON
	var configMap map[string]interface{}
	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		return nil, err
	}

	// Convert to JSON format
	return yaml.Marshal(configMap) // YAML library can directly convert map to JSON
}

// LoadMiddlewareConfig loads and parses middleware configuration file
func LoadMiddlewareConfig(t *testing.T, configPath string) (*middlewarev1.Middlewares, error) {
	// Set to continue execution when test fails
	t.Helper()

	// Load configuration file
	source := file.NewSource(configPath)
	c := config.New(config.WithSource(source))
	if err := c.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	defer c.Close()

	// Parse to Middlewares struct
	var middlewares middlewarev1.Middlewares
	if err := c.Scan(&middlewares); err != nil {
		return nil, fmt.Errorf("failed to unmarshal middleware config: %w", err)
	}

	return &middlewares, nil
}
