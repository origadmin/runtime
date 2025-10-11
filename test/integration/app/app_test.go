package app_test

import (
	"fmt"
	"testing"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file"

	configs "github.com/origadmin/runtime/test/integration/app/proto"
)

// TestAppBootstrap demonstrates the advantages of defining a unified structure in the application-specific bootstrap.proto
func TestAppBootstrap(t *testing.T) {
    // 1. Use runtime/config to load YAML configuration file
    source := file.NewSource("config.yaml")
    c := config.NewKConfig(config.WithKSource(source))
    if err := c.Load(); err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }
    defer c.Close()

    // Parse configuration into Bootstrap struct
    var bootstrap configs.Bootstrap
    if err := c.Scan(&bootstrap); err != nil {
        t.Fatalf("Failed to unmarshal config: %v", err)
    }

    // 2. Proof: "transport division" issue has been resolved
    // We can iterate through a unified `servers` list defined in bootstrap
    fmt.Println("--- Processing Unified Server List ---")
    if len(bootstrap.Servers) != 2 {
        t.Errorf("Expected 2 servers, got %d", len(bootstrap.Servers))
    }
    for _, srv := range bootstrap.Servers {
        // Using a simple switch statement, we can handle different types of transport
        // This is much cleaner than handling two separate lists (grpc_servers, http_servers)
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
            t.Errorf("Unknown server type in unified list: %T", c)
        }
    }
    fmt.Println("--- Server Processing Complete ---")

    fmt.Println("") // Spacing

    // 3. Proof: "different clients need different Middleware" issue has been resolved
    // We can iterate through the client list, where each client has its own dedicated middleware chain
    fmt.Println("--- Processing Clients with Specific Middlewares ---")
    if len(bootstrap.Clients) != 2 {
        t.Errorf("Expected 2 clients, got %d", len(bootstrap.Clients))
    }
    for _, cli := range bootstrap.Clients {
        // Note: Changed how target is obtained because Discoveries is an array
        target := ""
        if len(cli.Discoveries) > 0 {
            target = cli.Discoveries[0].Name
        }
        fmt.Printf("Client for target '%s' has %d specific middlewares:\n", target, len(cli.Middlewares))

        // Assertions to prove we loaded the correct, dedicated data
        if target == "user-service" && len(cli.Middlewares) != 2 {
            t.Errorf("Expected 2 middlewares for user-service, got %d", len(cli.Middlewares))
        }
        if target == "order-service" && len(cli.Middlewares) != 2 {
            t.Errorf("Expected 2 middlewares for order-service, got %d", len(cli.Middlewares))
        }
    }
    fmt.Println("--- Client Processing Complete ---")
}