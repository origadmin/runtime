package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/origadmin/runtime/api/gen/go/runtime/config/v1"
	"github.com/origadmin/runtime/test/helper"
)

// 测试服务器配置的加载和解析
func TestServerConfigLoading(t *testing.T) {
	// Test server configurations in different formats
	testCases := []struct {
		name     string
		filePath string
	}{{
		name:     "YAML",
		filePath: "configs/server_config.yaml",
	}, {
		name:     "JSON",
		filePath: "configs/server_config.json",
	}, {
		name:     "TOML",
		filePath: "configs/server_config.toml",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var serverConfig configv1.Server
			helper.LoadConfigFromFile(t, tc.filePath, &serverConfig)

			// Verify the GRPC configuration
			assert.NotNil(t, serverConfig.Grpc)
			assert.Equal(t, "tcp", serverConfig.Grpc.Network)
			assert.Equal(t, ":9000", serverConfig.Grpc.Addr)
			assert.Equal(t, "1s", serverConfig.Grpc.Timeout.AsDuration().String())
			assert.True(t, serverConfig.Grpc.EnableReflection)

			// Verify HTTP configuration
			assert.NotNil(t, serverConfig.Http)
			assert.Equal(t, "tcp", serverConfig.Http.Network)
			assert.Equal(t, ":8000", serverConfig.Http.Addr)
			assert.Equal(t, "2s", serverConfig.Http.Timeout.AsDuration().String())
		})
	}
}
