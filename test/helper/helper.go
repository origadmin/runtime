package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"gopkg.in/yaml.v3"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/config/file"
)

// LoadYAMLConfig 加载 YAML 配置文件并转换为 JSON 格式的字节数组
func LoadYAMLConfig(filename string) ([]byte, error) {
	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 构建完整的配置文件路径
	configPath := filepath.Join(currentDir, filename)

	// 读取配置文件内容
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// 解析 YAML 到 map 以转换为 JSON
	var configMap map[string]interface{}
	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		return nil, err
	}

	// 转换为 JSON 格式
	return yaml.Marshal(configMap) // YAML 库可以直接将 map 转换为 JSON
}

// LoadMiddlewareConfig 加载并解析中间件配置文件
func LoadMiddlewareConfig(t *testing.T, configPath string) (*middlewarev1.Middlewares, error) {
	// 设置测试失败时继续执行
	t.Helper()

	// 加载配置文件
	source := file.NewSource(configPath)
	c := config.New(config.WithSource(source))
	if err := c.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	defer c.Close()

	// 解析为Middlewares结构体
	var middlewares middlewarev1.Middlewares
	if err := c.Scan(&middlewares); err != nil {
		return nil, fmt.Errorf("failed to unmarshal middleware config: %w", err)
	}

	return &middlewares, nil
}
