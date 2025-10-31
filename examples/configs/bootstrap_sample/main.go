package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/log"

	// Import all contrib packages that need self-registration
	//_ "github.com/origadmin/contrib/config/env"   // Register env config source
	//_ "github.com/origadmin/contrib/config/file"  // Register file config source
	//_ "github.com/origadmin/contrib/registry/consul" // Register consul registry (if used in config)

	// 导入生成的配置
	conf "github.com/origadmin/runtime/examples/protos/bootstrap_sample"
)

// Version information injected at build time
var (
	Version = "v0.1.0"
	Name    = "bootstrap-sample"
	ID      = "bootstrap-sample"
)

func main() {
	// 1. Load .env file (optional, used for local development and environment variable presets)
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found or could not be loaded: %v\n", err)
	}

	// 2. Create and start Runtime instance
	// NewFromBootstrap encapsulates all bootstrap processes
	rt, err := runtime.NewFromBootstrap(
		"./configs/bootstrap.yaml",
		runtime.WithAppInfo(&runtime.AppInfo{
			ID:      ID,
			Name:    Name,
			Version: Version,
		}),
	)
	if err != nil {
		fmt.Println("Failed to create Runtime:", err)
		os.Exit(1)
	}
	defer rt.Cleanup()

	// 3. Get components from Runtime and use them
	logger := log.NewHelper(rt.Logger())
	appInfo := rt.AppInfo()

	logger.Infof("App %s (%s) is starting...", appInfo.Name, appInfo.Version)

	// Get the generated Bootstrap configuration
	var bc conf.Bootstrap
	if err := rt.Config().Decode("", &bc); err != nil {
		log.Fatalf("Failed to decode bootstrap config: %v", err)
	}

	// Print some information from the configuration to demonstrate successful loading

	logger.Infof("Logger Level: %s, Format: %s", bc.GetLogger().GetLevel(), bc.GetLogger().GetFormat())
	logger.Infof("HTTP Server Addr: %s, Timeout: %s", bc.GetServers().GetHttp().GetAddr(), bc.GetServers().GetHttp().GetTimeout().AsDuration())
	logger.Infof("gRPC Server Addr: %s, Timeout: %s", bc.GetServers().GetGrpc().GetAddr(), bc.GetServers().GetGrpc().GetTimeout().AsDuration())

	if bc.GetData() != nil && bc.GetData().GetDatabase() != nil {
		logger.Infof("Database Driver: %s, Source: %s", bc.GetData().GetDatabase().GetDriver(), bc.GetData().GetDatabase().GetSource())
	}
	if bc.GetData() != nil && bc.GetData().GetRedis() != nil {
		logger.Infof("Redis Addr: %s, DB: %d", bc.GetData().GetRedis().GetAddr(), bc.GetData().GetRedis().GetDb())
	}

	// The proto file for the Security component doesn't exist, so GetSecurity() will return nil
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

	// 4. Here you can create and run the Kratos App based on the configuration in bc
	// But to keep the example minimal and focused on config loading, we only demonstrate config loading
	log.Println("Bootstrap config loaded and runtime initialized successfully.")
	log.Println("Application will exit after cleanup.")
}
