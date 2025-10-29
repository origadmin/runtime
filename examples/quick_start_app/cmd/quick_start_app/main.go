package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/origadmin/runtime"

	// 导入所有需要自注册的 contrib 包
	_ "github.com/origadmin/contrib/config/env"
	_ "github.com/origadmin/contrib/config/file"

	// 导入应用内部的 wire 生成文件
	"github.com/origadmin/runtime/examples/quick_start_app/internal/conf"
	"github.com/origadmin/runtime/examples/quick_start_app/internal/server"
	"github.com/origadmin/runtime/examples/quick_start_app/internal/service"
)

// 编译时注入版本信息
var (
	Version = "v0.1.0"
	Name    = "quick-start-app"
	ID      = "quick-start-app"
)

func main() {
	// 1. 加载 .env 文件 (可选，用于本地开发和环境变量预设)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded:", err)
	}

	// 2. 创建并启动 Runtime 实例
	// NewFromBootstrap 封装了所有引导流程
	rt, err := runtime.NewFromBootstrap(
		"./configs/bootstrap.yaml",
		runtime.WithAppInfo(&runtime.AppInfo{
			ID:      ID,
			Name:    Name,
			Version: Version,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create runtime: %v", err)
	}
	defer rt.Cleanup()

	// 3. 从 Runtime 获取组件并使用
	logger := rt.Logger()
	appInfo := rt.AppInfo()

	logger.Infof("App %s (%s) is starting...", appInfo.Name, appInfo.Version)

	// 4. 使用 wire 组装应用，并运行 Kratos App
	// newApp 函数由 wire 生成，它会接收 rt 实例，并返回 kratos.App
	app, cleanupApp, err := newApp(rt)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	defer cleanupApp()

	// 运行应用
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
