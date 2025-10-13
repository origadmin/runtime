package middleware_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1/selector"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/middleware"
)

// Define all supported middleware types
var supportedMiddlewareTypes = map[string]struct{}{
	"metadata":        {},
	"logging":         {},
	"selector":        {},
	"rate_limiter":    {},
	"circuit_breaker": {},
	"jwt":             {},
	"cors":            {},
	"metrics":         {},
	"validator":       {},
}

func TestMiddleware(t *testing.T) {
	t.Run("LoadAndBuild", TestMiddleware_LoadAndBuild)
	t.Run("Selector", TestSelectorMiddleware)
	t.Run("Creation", TestMiddleware_Creation)
	t.Run("EdgeCases", TestMiddleware_EdgeCases)
}

// getMiddlewareConfigField 获取中间件配置字段
func getMiddlewareConfigField(mw *middlewarev1.MiddlewareConfig, fieldName string) interface{} {
	switch fieldName {
	case "metadata":
		return mw.Metadata
	case "logging":
		return mw.Logging
	case "selector":
		return mw.Selector
	case "rate_limiter":
		return mw.RateLimiter
	case "circuit_breaker":
		return mw.CircuitBreaker
	default:
		return nil
	}
}

// loadTestConfig 加载测试配置文件
func loadTestConfig(t *testing.T) *middlewarev1.Middlewares {
	// 获取当前测试文件的目录
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// 构建配置文件的完整路径
	configPath := filepath.Join(dir, "configs", "config.yaml")

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file not found: %s", configPath)
	}

	// 加载配置
	var configs middlewarev1.Middlewares
	cfg, err := config.Load(configPath, &configs)
	require.NoError(t, err, "Failed to load config")
	require.NotNil(t, cfg, "Config instance should not be nil")
	t.Cleanup(func() {
		_ = cfg.Close()
	})

	return &configs
}

// TestMiddleware_LoadAndBuild 测试中间件加载和构建
func TestMiddleware_LoadAndBuild(t *testing.T) {
	// 加载测试配置
	configs := loadTestConfig(t)

	t.Run("VerifyConfig", func(t *testing.T) {
		// 检查配置中的中间件数量
		assert.GreaterOrEqual(t, len(configs.Middlewares), 2, "Should have at least 2 middlewares in config")

		// 检查每个中间件的配置
		for _, mw := range configs.Middlewares {
			t.Run(mw.Name, func(t *testing.T) {
				// 验证必填字段
				assert.NotEmpty(t, mw.Name, "Middleware name should not be empty")
				assert.True(t, mw.Enabled, "Middleware should be enabled by default")

				// 验证类型是否支持
				_, exists := supportedMiddlewareTypes[mw.Type]
				assert.True(t, exists, "Unsupported middleware type: %s", mw.Type)

				// 验证配置字段是否存在
				configValue := getMiddlewareConfigField(mw, mw.Type)
				assert.NotNil(t, configValue, "Config for %s should not be nil", mw.Type)
			})
		}
	})

	t.Run("BuildClientMiddleware", func(t *testing.T) {
		// 构建客户端中间件链
		clientMWs := middleware.BuildClient(configs)

		// 计算实际支持的客户端中间件数量（排除不支持的中间件类型）
		expectedCount := 0
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "rate_limiter" { // rate_limiter 在客户端不受支持
				expectedCount++
			}
		}

		// 验证中间件数量
		assert.Equal(t, expectedCount, len(clientMWs), "Number of client middlewares should match expected")
	})

	t.Run("BuildServerMiddleware", func(t *testing.T) {
		// 构建服务端中间件链
		serverMWs := middleware.BuildServer(configs)

		// 计算实际支持的服务端中间件数量（排除不支持的中间件类型）
		expectedCount := 0
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "circuit_breaker" { // circuit_breaker 在服务端不受支持
				expectedCount++
			}
		}

		// 验证中间件数量
		assert.Equal(t, expectedCount, len(serverMWs), "Number of server middlewares should match expected")
	})

	t.Run("MiddlewareNaming", func(t *testing.T) {
		// 创建一个自定义的中间件配置
		customConfig := &middlewarev1.Middlewares{
			Middlewares: []*middlewarev1.MiddlewareConfig{
				{
					Name:    "custom-name",
					Type:    "logging",
					Enabled: true,
					Logging: &middlewarev1.Logging{},
				},
				{
					Name:     "", // 测试未命名的中间件
					Type:     "metadata",
					Enabled:  true,
					Metadata: &middlewarev1.Metadata{},
				},
			},
		}

		// 测试客户端中间件
		clientMWs := middleware.BuildClient(customConfig)
		assert.NotEmpty(t, clientMWs, "Client middlewares should not be empty")

		// 测试服务端中间件
		serverMWs := middleware.BuildServer(customConfig)
		assert.NotEmpty(t, serverMWs, "Server middlewares should not be empty")
	})

	// 验证配置中的中间件
	t.Run("VerifyConfig", func(t *testing.T) {
		// 检查配置中的中间件数量
		assert.GreaterOrEqual(t, len(configs.Middlewares), 2, "Should have at least 2 middlewares in config")

		// 检查每个中间件的配置
		for _, mw := range configs.Middlewares {
			t.Run(mw.Name, func(t *testing.T) {
				// 验证必填字段
				assert.NotEmpty(t, mw.Name, "Middleware name should not be empty")
				assert.True(t, mw.Enabled, "Middleware should be enabled by default")

				// 验证类型是否支持
				_, exists := supportedMiddlewareTypes[mw.Type]
				assert.True(t, exists, "Unsupported middleware type: %s", mw.Type)

				// 验证配置字段是否存在
				configValue := getMiddlewareConfigField(mw, mw.Type)
				assert.NotNil(t, configValue, "Config for %s should not be nil", mw.Type)
			})
		}
	})

	// 测试客户端中间件构建
	t.Run("BuildClientMiddleware", func(t *testing.T) {
		// 构建客户端中间件链
		clientMWs := middleware.BuildClient(configs)

		// 计算实际支持的客户端中间件数量（排除不支持的中间件类型）
		expectedCount := 0
		var expectedOrder []string
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "rate_limiter" { // rate_limiter 在客户端不受支持
				expectedOrder = append(expectedOrder, mw.Name)
				expectedCount++
			}
		}

		// 验证中间件数量
		assert.Equal(t, expectedCount, len(clientMWs), "Number of client middlewares should match expected")

		// 验证每个中间件都被正确处理
		for i := 0; i < len(clientMWs); i++ {
			assert.NotNil(t, clientMWs[i], "Client middleware at index %d should not be nil", i)
		}

		// 验证 Selector 中间件存在
		selectorFound := false
		for _, mwName := range expectedOrder {
			if mwName == "selector" {
				selectorFound = true
				break
			}
		}
		assert.True(t, selectorFound, "Selector middleware should be present in client middlewares")
	})

	// 测试服务端中间件构建
	t.Run("BuildServerMiddleware", func(t *testing.T) {
		// 构建服务端中间件链
		serverMWs := middleware.BuildServer(configs)

		// 计算实际支持的服务端中间件数量（排除不支持的中间件类型）
		expectedCount := 0
		var expectedOrder []string
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "circuit_breaker" { // circuit_breaker 在服务端不受支持
				expectedOrder = append(expectedOrder, mw.Name)
				expectedCount++
			}
		}

		// 验证中间件数量
		assert.Equal(t, expectedCount, len(serverMWs), "Number of server middlewares should match expected")

		// 验证每个中间件都被正确处理
		for i := 0; i < len(serverMWs); i++ {
			assert.NotNil(t, serverMWs[i], "Server middleware at index %d should not be nil", i)
		}

		// 验证 Selector 中间件存在
		selectorFound := false
		for _, mwName := range expectedOrder {
			if mwName == "selector" {
				selectorFound = true
				break
			}
		}
		assert.True(t, selectorFound, "Selector middleware should be present in server middlewares")
	})
}

// TestSelectorMiddleware 测试Selector中间件的includes/excludes功能
func TestSelectorMiddleware(t *testing.T) {
	// 准备测试数据
	tests := []struct {
		name        string
		config      *middlewarev1.Middlewares
		expectCount int
		side        string // "client" or "server"
	}{
		{
			name: "selector_with_includes",
			config: &middlewarev1.Middlewares{
				Middlewares: []*middlewarev1.MiddlewareConfig{
					{
						Name:    "metadata",
						Type:    "metadata",
						Enabled: true,
					},
					{
						Name:    "logging",
						Type:    "logging",
						Enabled: true,
					},
					{
						Name:    "selector",
						Type:    "selector",
						Enabled: true,
						Selector: &selectorv1.Selector{
							Includes: []string{"logging"}, // 只包含logging中间件
						},
					},
				},
			},
			expectCount: 2, // selector + logging
			side:        "client",
		},
		{
			name: "selector_with_excludes",
			config: &middlewarev1.Middlewares{
				Middlewares: []*middlewarev1.MiddlewareConfig{
					{
						Name:    "metadata",
						Type:    "metadata",
						Enabled: true,
					},
					{
						Name:    "logging",
						Type:    "logging",
						Enabled: true,
					},
					{
						Name:    "selector",
						Type:    "selector",
						Enabled: true,
						Selector: &selectorv1.Selector{
							Excludes: []string{"metadata"}, // 排除metadata中间件
						},
					},
				},
			},
			expectCount: 2, // selector + logging
			side:        "server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.side {
			case "client":
				mw := middleware.BuildClient(tt.config)
				assert.Len(t, mw, tt.expectCount, "Unexpected number of client middlewares")
			case "server":
				mw := middleware.BuildServer(tt.config)
				assert.Len(t, mw, tt.expectCount, "Unexpected number of server middlewares")
			}
		})
	}
}

// TestMiddleware_Creation 测试中间件创建
func TestMiddleware_Creation(t *testing.T) {
	// 测试不同类型的中间件创建
	tests := []struct {
		name     string
		config   *middlewarev1.MiddlewareConfig
		exists   bool
		clientMW bool
		serverMW bool
	}{
		{
			name: "valid logging middleware",
			config: &middlewarev1.MiddlewareConfig{
				Name:    "test-logging",
				Type:    "logging",
				Enabled: true,
				Logging: &middlewarev1.Logging{},
			},
			exists:   true,
			clientMW: true,
			serverMW: true,
		},
		{
			name: "disabled middleware",
			config: &middlewarev1.MiddlewareConfig{
				Name:    "disabled-mw",
				Type:    "logging",
				Enabled: false,
			},
			exists:   false,
			clientMW: false,
			serverMW: false,
		},
		{
			name: "unknown middleware type",
			config: &middlewarev1.MiddlewareConfig{
				Name:    "unknown",
				Type:    "nonexistent",
				Enabled: true,
			},
			exists:   false,
			clientMW: false,
			serverMW: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试客户端中间件
			clientMW, _ := middleware.NewClient(tt.config)
			if tt.clientMW {
				assert.NotNil(t, clientMW, "Expected client middleware to be created")
			} else {
				assert.Nil(t, clientMW, "Expected no client middleware to be created")
			}

			// 测试服务端中间件
			serverMW, _ := middleware.NewServer(tt.config)
			if tt.serverMW {
				assert.NotNil(t, serverMW, "Expected server middleware to be created")
			} else {
				assert.Nil(t, serverMW, "Expected no server middleware to be created")
			}

			// 验证中间件是否被正确存储
			if tt.exists {
				// 创建中间件
				if tt.clientMW {
					mw, _ := middleware.NewClient(tt.config)
					assert.NotNil(t, mw, "Client middleware should be created")
				}

				if tt.serverMW {
					mw, _ := middleware.NewServer(tt.config)
					assert.NotNil(t, mw, "Server middleware should be created")
				}
			}
		})
	}
}

// TestMiddleware_EdgeCases 测试中间件的边缘情况
func TestMiddleware_EdgeCases(t *testing.T) {
	t.Run("NilConfig", func(t *testing.T) {
		// 测试nil配置
		nilMWs := middleware.BuildClient(nil)
		assert.Empty(t, nilMWs, "BuildClient with nil config should return empty middlewares")
		nilMWs = middleware.BuildServer(nil)
		assert.Empty(t, nilMWs, "BuildServer with nil config should return empty middlewares")
	})

	t.Run("EmptyConfig", func(t *testing.T) {
		emptyConfigs := &middlewarev1.Middlewares{}

		// 测试空配置的客户端中间件
		clientMWs := middleware.BuildClient(emptyConfigs)
		assert.Empty(t, clientMWs, "BuildClient with empty config should return empty middlewares")

		// 测试空配置的服务端中间件
		serverMWs := middleware.BuildServer(emptyConfigs)
		assert.Empty(t, serverMWs, "BuildServer with empty config should return empty middlewares")
	})

	t.Run("DisabledMiddleware", func(t *testing.T) {
		disabledConfig := &middlewarev1.Middlewares{
			Middlewares: []*middlewarev1.MiddlewareConfig{
				{
					Name:    "disabled-mw",
					Type:    "logging",
					Enabled: false,
					Logging: &middlewarev1.Logging{},
				},
			},
		}

		// 测试禁用的中间件不会被创建
		clientMWs := middleware.BuildClient(disabledConfig)
		assert.Empty(t, clientMWs, "Disabled middleware should not be created")

		serverMWs := middleware.BuildServer(disabledConfig)
		assert.Empty(t, serverMWs, "Disabled middleware should not be created")
	})

	t.Run("UnknownMiddlewareType", func(t *testing.T) {
		unknownConfig := &middlewarev1.Middlewares{
			Middlewares: []*middlewarev1.MiddlewareConfig{
				{
					Name:    "unknown-mw",
					Type:    "nonexistent",
					Enabled: true,
				},
			},
		}

		// 测试未知类型的中间件不会被创建
		clientMWs := middleware.BuildClient(unknownConfig)
		assert.Empty(t, clientMWs, "Unknown middleware type should not be created")

		serverMWs := middleware.BuildServer(unknownConfig)
		assert.Empty(t, serverMWs, "Unknown middleware type should not be created")
	})
}
