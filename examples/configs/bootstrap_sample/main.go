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

	// 2. Create and start App instance
	// NewFromBootstrap encapsulates all bootstrap processes
	rt, err := runtime.NewFromBootstrap(
		"./configs/bootstrap.yaml",
		runtime.WithAppInfo(runtime.NewAppInfo(Name, Version, runtime.WithAppInfoID(ID))),
	)
	if err != nil {
		fmt.Println("Failed to create App:", err)
		os.Exit(1)
	}
	// The Cleanup method has been removed. The application manages its own lifecycle.

	// 3. Get components from App and use them
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
	for _, srv := range bc.GetServers().GetConfigs() {
		switch srv.GetProtocol() {
		case "http":
			logger.Infof("HTTP Server Addr: %s, Timeout: %s", srv.GetHttp().GetAddr(), srv.GetHttp().GetTimeout().AsDuration())
		case "grpc":
			logger.Infof("gRPC Server Addr: %s, Timeout: %s", srv.GetGrpc().GetAddr(), srv.GetGrpc().GetTimeout().AsDuration())
		default:
			logger.Warnf("Unknown server protocol type %s", srv.GetProtocol())
			continue
		}
	}

	if bc.GetData() != nil && bc.GetData().GetDatabases() != nil {
		for _, db := range bc.GetData().GetDatabases().GetConfigs() {
			logger.Infof("Database Driver: %s, Source: %s", db.GetDialect(), db.GetSource())
		}

	}
	if bc.GetData() != nil && bc.GetData().GetCaches() != nil {
		for _, cache := range bc.GetData().GetCaches().GetConfigs() {
			logger.Infof("Cache Type: %s, Address: %s", cache.GetDriver(), cache.GetRedis().GetAddr())
		}
	}

	// The proto file for the Security component doesn't exist, so GetSecurity() will return nil
	// if bc.GetSecurity() != nil && len(bc.GetSecurity().GetAuthenticators()) > 0 {
	// 	logger.Infof("Security Authenticator Type: %s", bc.GetSecurity().GetAuthenticators()[0].GetType())
	// }

	if bc.GetDiscoveries() != nil && len(bc.GetDiscoveries().GetConfigs()) > 0 {
		logger.Infof("Discovery Type: %s, Address: %s", bc.GetDiscoveries().GetConfigs()[0].GetType(),
			bc.GetDiscoveries().GetConfigs()[0].GetConsul().GetAddress())
	}

	if bc.GetMiddlewares() != nil && len(bc.GetMiddlewares().GetConfigs()) > 0 {
		logger.Infof("Middlewares configured: %d", len(bc.GetMiddlewares().GetConfigs()))
	}

	if bc.GetBrokers() != nil && bc.GetBrokers().GetBrokers() != nil {
		for _, broker := range bc.GetBrokers().GetBrokers() {
			logger.Infof("Broker Type: Kafka, Brokers: %v", broker.GetKafka())
		}

	}

	// 4. Here you can create and run the Kratos App based on the configuration in bc
	// But to keep the example minimal and focused on config loading, we only demonstrate config loading
	logger.Info("Bootstrap config loaded and runtime initialized successfully.")
	logger.Info("Application will exit after cleanup.")
}
