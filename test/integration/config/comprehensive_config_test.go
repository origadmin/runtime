// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/test/helper"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// 运行所有配置的通用断言
func runCompleteConfigAssertions(t *testing.T, config *testconfigs.AllConfigs) {
	assert := assert.New(t)

	// App配置断言
	assert.Equal("test-app-id", config.App.Id)
	assert.Equal("TestApp", config.App.Name)
	assert.Equal("1.0.0", config.App.Version)
	assert.Equal("test", config.App.Env)
	assert.Contains(config.App.Metadata, "key1")
	assert.Contains(config.App.Metadata, "key2")

	// Bootstrap配置断言
	assert.NotNil(config.Bootstrap)
	assert.Equal(2, len(config.Bootstrap.Sources))

	// Client配置断言
	assert.Equal("discovery:///user-service", config.Client.Endpoint)
	assert.Equal("v1.0.0", config.Client.Selector.Version)

	// Server配置断言
	assert.NotNil(config.Server.Grpc)
	assert.Equal("tcp", config.Server.Grpc.Network)
	assert.Equal(":9000", config.Server.Grpc.Addr)
	assert.NotNil(config.Server.Http)
	assert.Equal("tcp", config.Server.Http.Network)
	assert.Equal(":8000", config.Server.Http.Addr)

	// Logger配置断言
	assert.Equal("info", config.Logger.Level)
	assert.Equal("json", config.Logger.Format)
	assert.True(config.Logger.Stdout)
}

// 测试所有格式的完整配置加载
func TestCompleteConfigMultiFormatLoading(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		filePath string
	}{{
		name:     "YAML",
		filePath: "configs/complete_config.yaml",
	}, {
		name:     "JSON",
		filePath: "configs/complete_config.json",
	}, {
		name:     "TOML",
		filePath: "configs/complete_config.toml",
	}}

	// 遍历所有测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var config testconfigs.AllConfigs
			helper.LoadConfigFromFile(t, tc.filePath, &config)

			// 运行断言
			runCompleteConfigAssertions(t, &config)
			t.Logf("%s format complete config loaded and verified successfully!", tc.name)
		})
	}
}

// 测试配置格式互操作性
func TestCompleteConfigInteroperability(t *testing.T) {
	// 定义支持的格式
	formats := []struct {
		name     string
		filePath string
		ext      string
	}{{
		name:     "YAML",
		filePath: "configs/complete_config.yaml",
		ext:      "yaml",
	}, {
		name:     "JSON",
		filePath: "configs/complete_config.json",
		ext:      "json",
	}, {
		name:     "TOML",
		filePath: "configs/complete_config.toml",
		ext:      "toml",
	}}

	// 外循环：源格式
	for _, src := range formats {
		// 内循环：目标格式
		for _, dst := range formats {
			testName := src.name + "_to_" + dst.name
			t.Run(testName, func(t *testing.T) {
				// 加载原始配置
				var originalConfig testconfigs.AllConfigs
				helper.LoadConfigFromFile(t, src.filePath, &originalConfig)

				// 保存为目标格式
				tempDir := t.TempDir()
				tempFilePath := filepath.Join(tempDir, "temp_config."+dst.ext)
				helper.SaveConfigToFile(t, &originalConfig, tempFilePath, dst.ext)

				// 重新加载目标格式配置
				var convertedConfig testconfigs.AllConfigs
				helper.LoadConfigFromFile(t, tempFilePath, &convertedConfig)

				// 断言配置一致
				assert.True(t, proto.Equal(&originalConfig, &convertedConfig),
					"Struct loaded from %s should be identical to the one saved to %s and loaded back",
					src.name, dst.name)
			})
		}
	}
}

// 测试使用Runtime加载完整配置
func TestRuntimeLoadCompleteConfig(t *testing.T) {
	assert := assert.New(t)

	// 获取当前目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file info")
	}
	currentDir := filepath.Dir(filename)

	// 保存原始工作目录并在测试结束后恢复
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	defer func() {
		err := os.Chdir(originalCwd)
		if err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	// 更改工作目录到runtime模块根目录
	runtimeRoot := filepath.Join(currentDir, "../../..")
	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	// 为每种格式运行测试
	formats := []string{"yaml", "json", "toml"}
	for _, format := range formats {
		t.Run("Runtime_"+format, func(t *testing.T) {
			// 创建临时引导配置文件
			tempDir := t.TempDir()
			tempBootstrapPath := filepath.Join(tempDir, "bootstrap."+format)

			// 写入临时引导配置
			var bootstrapContent string
			switch format {
			case "yaml":
				bootstrapContent = "sources:\n  - type: \"file\"\n    name: \"complete-config\"\n    file:\n      path: \"test/integration/config/configs/complete_config." + format + "\"\n    priority: 100"
			case "json":
				bootstrapContent = `{\"sources\": [{\"type\": \"file\", \"name\": \"complete-config\", \"file\": {\"path\": \"test/integration/config/configs/complete_config.` + format + `\"}, \"priority\": 100}]}`
			case "toml":
				bootstrapContent = `[[sources]]\ntype = \"file\"\nname = \"complete-config\"\nfile.path = \"test/integration/config/configs/complete_config.` + format + `\"\npriority = 100`
			}

			if err := os.WriteFile(tempBootstrapPath, []byte(bootstrapContent), 0644); err != nil {
				t.Fatalf("Failed to write temp bootstrap file: %v", err)
			}

			// 初始化Runtime
			rt, cleanup, err := rt.NewFromBootstrap(
				tempBootstrapPath,
				bootstrap.WithAppInfo(&interfaces.AppInfo{
					ID:      "test-complete-config",
					Name:    "TestCompleteConfig",
					Version: "1.0.0",
				}),
			)
			if err != nil {
				t.Fatalf("Failed to initialize runtime with %s config: %v", format, err)
			}
			defer cleanup()

			// 获取配置解码器
			configDecoder := rt.Config()
			assert.NotNil(configDecoder)

			// 解码为完整配置结构
			var config testconfigs.AllConfigs
			err = configDecoder.Decode("", &config)
			assert.NoError(err)

			// 运行断言
			runCompleteConfigAssertions(t, &config)
			t.Logf("Runtime loaded and verified %s format complete config successfully!", format)
		})
	}
}
