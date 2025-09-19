package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/origadmin/runtime"
	// Import the generated Go code from the api_gateway proto definition.
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
)

func main() {
	// 1. Create a new Runtime instance from the new api_gateway config.
	//    Path is now relative to the CWD (runtime directory), pointing to the bootstrap.yaml.
	rt, cleanup, err := runtime.NewFromBootstrap(
		"examples/configs/load_with_runtime/config/bootstrap.yaml", // Correctly load bootstrap.yaml
		runtime.WithAppInfo(runtime.AppInfo{
			ID:      "api-gateway-runtime-example",
			Name:    "ApiGatewayRuntimeExample",
			Version: "1.0.0",
			Env:     "dev",
		}),
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 2. Get the configuration decoder from the runtime instance.
	decoder := rt.Config()

	// --- DIAGNOSTIC PRINT Kratos Config Content START ---
	fmt.Printf("\n--- Debugging Kratos Config Content ---\n")

	// Get the raw config value from the decoder
	rawConfig := decoder.Config()
	fmt.Printf("Type of decoder.Config(): %v\n", reflect.TypeOf(rawConfig))
	fmt.Printf("Value of decoder.Config(): %+v\n", rawConfig)

	// Attempt to cast to map[string]any and then marshal
	rawConfigMap, ok := rawConfig.(map[string]any)
	if !ok {
		fmt.Println("Error: decoder.Config() is not a map[string]any.")
	} else {
		// Access "servers" from the map and marshal to JSON string
		if serversVal, exists := rawConfigMap["servers"]; exists {
			if serversStr, err := json.Marshal(serversVal); err == nil {
				fmt.Printf("Kratos Config 'servers' value (string): %s\n", serversStr)
			} else {
				fmt.Printf("Error marshalling 'servers' value to string: %v\n", err)
			}
		} else {
			fmt.Println("Kratos Config 'servers' value not found in raw map.")
		}

		// Access "clients" from the map and marshal to JSON string
		if clientsVal, exists := rawConfigMap["clients"]; exists {
			if clientsStr, err := json.Marshal(clientsVal); err == nil {
				fmt.Printf("Kratos Config 'clients' value (string): %s\n", clientsStr)
			} else {
				fmt.Printf("Error marshalling 'clients' value to string: %v\n", err)
			}
		} else {
			fmt.Println("Kratos Config 'clients' value not found in raw map.")
		}
	}
	fmt.Printf("--------------------------------------\n\n")
	// --- DIAGNOSTIC PRINT Kratos Config Content END ---

	// 3. Decode the entire configuration into our Bootstrap struct.
	var bc conf.Bootstrap
	if err := decoder.Decode("", &bc); err != nil {
		panic(err)
	}

	// --- DIAGNOSTIC PRINT Bootstrap struct content after Decode START ---
	fmt.Printf("\n--- Debugging Bootstrap struct content after Decode ---\n")
	fmt.Printf("%+v\n", bc)
	fmt.Printf("--------------------------------------\n\n")
	// --- DIAGNOSTIC PRINT Bootstrap struct content after Decode END ---

	// 4. Print the loaded configuration to verify.
	fmt.Printf("--- Loaded API Gateway config via runtime interface ---\n")

	// Verify server config
	if len(bc.Servers) > 0 && bc.Servers[0].GetHttp() != nil {
		fmt.Printf("Server HTTP Addr: %s\n", bc.Servers[0].GetHttp().GetAddr())
	} else {
		fmt.Println("No HTTP server configuration found.")
	}

	// Verify client config
	if userService, ok := bc.Clients["user-service"]; ok {
		fmt.Printf("Client 'user-service' Endpoint: %s\n", userService.GetDiscovery().GetEndpoint())
	} else {
		fmt.Println("No 'user-service' client configuration found.")
	}
}
