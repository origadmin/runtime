package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/origadmin/runtime"

	// 导入所有需要自注册的 contrib 包
	_ "github.com/origadmin/contrib/config/env"   // 注册 env 配置源
	_ "github.com/origadmin/contrib/config/file"  // 注册 file 配置源
	_ "github.com/origadmin/contrib/registry/consul" // 注册 consul 注册中心 (如果配置中使用了)

	// 导入生成的配置
	"github.com/origadmin/runtime/examples/configs/bootstrap_sample/conf"
)

// 编译时注入版本信息
var (
	Version = "v0.1.0"
	Name    = "bootstrap-sample"
	ID      = "bootstrap-sample"
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

	// 获取生成的 Bootstrap 配置
	var bc conf.Bootstrap
	if err := rt.Config().Decode("", &bc); err != nil {
		log.Fatalf("Failed to decode bootstrap config: %v", err)
	}

	// 打印一些从配置中获取的信息，以演示配置加载成功
	logger.Infof("Logger Level: %s, Format: %s", bc.GetLogger().GetLevel(), bc.GetLogger().GetFormat())
	logger.Infof("HTTP Server Addr: %s, Timeout: %s", bc.GetServers().GetHttp().GetAddr(), bc.GetServers().GetHttp().GetTimeout().AsDuration())
	logger.Infof("gRPC Server Addr: %s, Timeout: %s", bc.GetServers().GetGrpc().GetAddr(), bc.GetServers().GetGrpc().GetTimeout().AsDuration())

	if bc.GetData() != nil && bc.GetData().GetDatabase() != nil {
		logger.Infof("Database Driver: %s, Source: %s", bc.GetData().GetDatabase().GetDriver(), bc.GetData().GetDatabase().GetSource())
	}
	if bc.GetData() != nil && bc.GetData().GetRedis() != nil {
		logger.Infof("Redis Addr: %s, DB: %d", bc.GetData().GetRedis().GetAddr(), bc.GetData().GetRedis().GetDb())
	}

	// Security 组件的 proto 文件不存在，所以 GetSecurity() 会返回 nil
	// if bc.GetSecurity() != nil && len(bc.GetSecurity().GetAuthenticators()) > 0 {
	// 	logger.Infof("Security Authenticator Type: %s", bc.GetSecurity().GetAuthenticators()[0].GetType())
	// }

	if bc.GetDiscoveries() != nil && len(bc.GetDiscoveries().GetDiscoveries()) > 0 {
		logger.Infof("Discovery Type: %s, Address: %s", bc.GetDiscoveries().GetDiscoveries()[0].GetType(), bc.GetDiscoveries().GetDiscoveries()[0].GetConsul().GetAddress())
	}

	if bc.GetMiddlewares() != nil && len(bc.GetMiddlewares().GetMiddlewares()) > 0 {
		logger.Infof("Middlewares configured: %d", len(bc.GetMiddlewares().GetMiddlewares()))
	}

	if bc.GetBroker() != nil && bc.GetBroker().GetKafka() != nil {
		logger.Infof("Broker Type: Kafka, Brokers: %v", bc.GetBroker().GetKafka().GetBrokers())
	}

	if bc.GetWebsocket() != nil && bc.GetWebsocket().GetServer() != nil {
		logger.Infof("Websocket Server Addr: %s, Path: %s", bc.GetWebsocket().GetServer().GetAddr(), bc.GetWebsocket().GetServer().GetPath())
	}

	// 4. 在这里可以根据 bc 中的配置来创建和运行 Kratos App
	// 但为了保持示例的最小化和聚焦于配置加载，我们只演示配置加载
	log.Println("Bootstrap config loaded and runtime initialized successfully.")
	log.Println("Application will exit after cleanup.")
}
