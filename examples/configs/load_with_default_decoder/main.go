package main

import (
	"fmt"

	"github.com/origadmin/runtime/config"
	// Import the generated Go code from the api_gateway proto definition.
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
)

func main() {
	var bc conf.Bootstrap

	// Corrected path relative to the CWD (runtime directory).
	configPath := "examples/configs/load_with_default_decoder/config/config.yaml"

	// Use the Load function to read the new api_gateway config.
	c, err := config.Load(configPath, &bc)
	if err != nil {
		panic(fmt.Errorf("failed to load bootstrap config from path %s: %w", configPath, err))
	}
	defer c.Close()

	// Print the loaded configuration to verify.
	fmt.Printf("--- Loaded API Gateway config via default decoder ---\n")

	// Verify server config
	if len(bc.Servers) > 0 && bc.Servers[0].GetHttp() != nil {
		fmt.Printf("Server HTTP Addr: %s\n", bc.Servers[0].GetHttp().GetAddr())
	} else {
		fmt.Println("No HTTP server configuration found.")
	}

	// Verify client config
	if userService, ok := bc.Clients["user-service"]; ok {
		fmt.Printf("Client 'user-service' Endpoint: %s\n", userService.GetClient().GetEndpoint())
	} else {
		fmt.Println("No 'user-service' client configuration found.")
	}
}
