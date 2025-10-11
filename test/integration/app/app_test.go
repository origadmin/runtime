package app_test

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"

	configs "github.com/origadmin/runtime/test/integration/app/proto"
)

// TestAppBootstrap demonstrates the advantages of defining a unified structure in the application-specific bootstrap.proto
func TestAppBootstrap(t *testing.T) {
	// 1. 加载我们新的、与图纸完全对应的 YAML 配置文件
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		t.Fatalf("Failed to read config.yaml: %v", err)
	}

	// 将 YAML 直接解析到我们测试专属的 Bootstrap 结构体中
	var bootstrap configs.Bootstrap
	if err := yaml.Unmarshal(yamlFile, &bootstrap); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// 2. 证明: "transport不分" 的问题已解决
	// 我们可以遍历一个在 bootstrap 中定义的、统一的 `servers` 列表
	fmt.Println("--- Processing Unified Server List ---")
	if len(bootstrap.Servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(bootstrap.Servers))
	}
	for _, srv := range bootstrap.Servers {
		// 使用简单的 switch 语句，我们就可以处理不同类型的 transport
		// 这比处理两个独立的列表 (grpc_servers, http_servers) 要清晰得多
		switch c := srv.Config.(type) {
		case *configs.Server_Grpc:
			fmt.Printf("Found gRPC Server '%s' at address: %s\n", srv.Name, c.Grpc.Addr)
			if c.Grpc.Addr != ":9000" {
				t.Errorf("Unexpected gRPC addr: %s", c.Grpc.Addr)
			}
		case *configs.Server_Http:
			fmt.Printf("Found HTTP Server '%s' at address: %s\n", srv.Name, c.Http.Addr)
			if c.Http.Addr != ":8000" {
				t.Errorf("Unexpected HTTP addr: %s", c.Http.Addr)
			}
		default:
			t.Errorf("Unknown server type in unified list")
		}
	}
	fmt.Println("--- Server Processing Complete ---")

	fmt.Println("") // 间隔

	// 3. 证明: "不同client需要不同Middleware" 的问题已解决
	// 我们可以遍历客户端列表，其中每个客户端都拥有自己专属的中间件链
	fmt.Println("--- Processing Clients with Specific Middlewares ---")
	if len(bootstrap.Clients) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(bootstrap.Clients))
	}
	for _, cli := range bootstrap.Clients {
		target := cli.Discoveries[0].Name
		fmt.Printf("Client for target '%s' has %d specific middlewares:\n", target, len(cli.Middlewares))

		// 断言以证明我们加载了正确的、专属的数据
		if target == "user-service" && len(cli.Middlewares) != 2 {
			t.Errorf("Expected 2 middlewares for user-service, got %d", len(cli.Middlewares))
		}
		if target == "order-service" && len(cli.Middlewares) != 2 {
			t.Errorf("Expected 2 middlewares for order-service, got %d", len(cli.Middlewares))
		}
	}
	fmt.Println("--- Client Processing Complete ---")
}
