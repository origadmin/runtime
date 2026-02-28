package config_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite is the unified entry point for all configuration-related integration tests.
type ConfigTestSuite struct {
	suite.Suite
}

// TestConfigTestSuite is the main test entry point for the config package.
// It runs all test suites in the test_cases directory.
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
