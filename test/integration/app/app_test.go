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
		switch srv.GetProtocol() {
		case "grpc":
			c := srv.GetGrpc()
			fmt.Printf("Found gRPC Server '%s' at address: %s\n", srv.Name, c.Addr)
			if c.Addr != ":9000" {
				t.Errorf("Unexpected gRPC addr: %s", c.Addr)
			}
		case "http":
			c := srv.GetHttp()
			fmt.Printf("Found HTTP Server '%s' at address: %s\n", srv.Name, c.Addr)
			if c.Addr != ":8000" {
				t.Errorf("Unexpected HTTP addr: %s", c.Addr)
			}
		default:
			t.Errorf("Unknown server type in unified list: %T", srv.GetProtocol())
		}
	}
	fmt.Println("--- Server Processing Complete ---")

	fmt.Println("") // Spacing

	// 3. Proof: "different clients need different Middleware" issue has been resolved
	// We can iterate through the client list, where each client has its own dedicated middleware chain
	fmt.Println("--- Processing Clients with Specific Middlewares ---")
	if len(bootstrap.GetClients().GetClients()) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(bootstrap.GetClients().GetClients()))
	}
	for _, cli := range bootstrap.GetClients().GetClients() {
		// Note: Client configuration is now nested under specific protocols (grpc, http)
		// We need to check the protocol and then access the specific client config.
		target := cli.GetProtocol()
		var middlewares []string

		switch cli.GetProtocol() {
		case "grpc":
			c := cli.GetGrpc()
			if c != nil {
				middlewares = c.GetMiddlewares()
				target = c.GetEndpoint()
			}
		case "http":
			c := cli.GetHttp()
			if c != nil {
				middlewares = c.GetMiddlewares()
				target = c.GetEndpoint()
			}
		default:
			t.Errorf("Unknown client protocol: %s", cli.GetProtocol())
		}

		fmt.Printf("Client for target '%s' (%s) has %d specific middlewares:\n", target, cli.GetName(), len(middlewares))

		// Assertions to prove we loaded the correct, dedicated data
		if cli.GetName() == "user-service" && len(middlewares) != 2 {
			t.Errorf("Expected 2 middlewares for user-service, got %d", len(middlewares))
		}
		if cli.GetName() == "order-service" && len(middlewares) != 2 {
			t.Errorf("Expected 2 middlewares for order-service, got %d", len(middlewares))
		}
	}
	fmt.Println("--- Client Processing Complete ---")

	fmt.Println("") // Spacing

	// 4. Proof: Top-level middlewares are loaded correctly
	fmt.Println("--- Processing Top-level Middlewares ---")
	if bootstrap.Middlewares == nil || len(bootstrap.Middlewares.GetMiddlewares()) != 2 {
		t.Errorf("Expected 2 top-level middlewares, got %d", len(bootstrap.Middlewares.GetMiddlewares()))
	}
	for i, mw := range bootstrap.Middlewares.GetMiddlewares() {
		fmt.Printf("Top-level Middleware %d: Name='%s', Type='%s', Enabled=%t\n", i+1, mw.GetName(), mw.GetType(), mw.GetEnabled())
	}
	fmt.Println("--- Top-level Middleware Processing Complete ---")
}
