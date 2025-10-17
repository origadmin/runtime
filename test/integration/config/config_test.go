package config_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite is the unified entry point for all configuration-related integration tests.
// Individual test suites are located in sub-packages under 'test_cases/' and are discovered
// and run automatically by 'go test ./...'.
type ConfigTestSuite struct {
	suite.Suite
}

// TestConfigTestSuite is a placeholder to ensure the main package has a test entry point.
// Actual test suites are run by 'go test ./...' in their respective sub-packages.
func TestConfigTestSuite(t *testing.T) {
	// This function can be used for top-level setup/teardown if needed, but does not run sub-suites directly.
	// suite.Run(t, new(ConfigTestSuite))
}
