package main

import (
	"encoding/json"
	"reflect"

	"github.com/origadmin/runtime"
	// Import the generated Go code from the api_gateway proto definition.
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
	"github.com/origadmin/runtime/log" // Import the log package
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

	// Get the configured logger from the runtime instance
	appLogger := log.NewHelper(rt.Logger()) // Use log.NewHelper for convenience

	appLogger.Info("Application started successfully!") // Log a message using the configured logger

	// 2. Get the configuration decoder from the runtime instance.
	decoder := rt.Config()

	// --- DIAGNOSTIC PRINT Kratos Config Content START ---
	appLogger.Info("--- Debugging Kratos Config Content ---") // Use appLogger

	// Get the raw config value from the decoder
	rawConfig := decoder.Config()
	appLogger.Infof("Type of decoder.Config(): %v", reflect.TypeOf(rawConfig)) // Use appLogger
	appLogger.Infof("Value of decoder.Config(): %+v", rawConfig)               // Use appLogger

	// Attempt to cast to map[string]any and then marshal
	rawConfigMap, ok := rawConfig.(map[string]any)
	if !ok {
		appLogger.Error("Error: decoder.Config() is not a map[string]any.") // Use appLogger
	} else {
		// Access "servers" from the map and marshal to JSON string
		if serversVal, exists := rawConfigMap["servers"]; exists {
			if serversStr, err := json.Marshal(serversVal); err == nil {
				appLogger.Infof("Kratos Config 'servers' value (string): %s", serversStr) // Use appLogger
			} else {
				appLogger.Errorf("Error marshalling 'servers' value to string: %v", err) // Use appLogger
			}
		} else {
			appLogger.Info("Kratos Config 'servers' value not found in raw map.") // Use appLogger
		}

		// Access "clients" from the map and marshal to JSON string
		if clientsVal, exists := rawConfigMap["clients"]; exists {
			if clientsStr, err := json.Marshal(clientsVal); err == nil {
				appLogger.Infof("Kratos Config 'clients' value (string): %s", clientsStr) // Use appLogger
			} else {
				appLogger.Errorf("Error marshalling 'clients' value to string: %v", err) // Use appLogger
			}
		} else {
			appLogger.Info("Kratos Config 'clients' value not found in raw map.") // Use appLogger
		}
	}
	appLogger.Info("--------------------------------------") // Use appLogger
	// --- DIAGNOSTIC PRINT Kratos Config Content END ---

	// 3. Decode the entire configuration into our Bootstrap struct.
	var bc conf.Bootstrap
	if err := decoder.Decode("", &bc); err != nil {
		appLogger.Errorf("Failed to decode bootstrap config: %v", err) // Use appLogger
		panic(err)
	}

	// --- DIAGNOSTIC PRINT Bootstrap struct content after Decode START ---
	appLogger.Info("--- Debugging Bootstrap struct content after Decode ---") // Use appLogger
	appLogger.Infof("%+v", bc)                                                // Use appLogger
	appLogger.Info("--------------------------------------")                  // Use appLogger

	// 4. Print the loaded configuration to verify.
	appLogger.Info("--- Loaded API Gateway config via runtime interface ---") // Use appLogger

	// Verify server config
	if len(bc.Servers) > 0 && bc.Servers[0].GetHttp() != nil {
		appLogger.Infof("Server HTTP Addr: %s", bc.Servers[0].GetHttp().GetAddr()) // Use appLogger
	} else {
		appLogger.Info("No HTTP server configuration found.") // Use appLogger
	}

	// Verify client config
	if userService, ok := bc.Clients["user-service"]; ok {
		appLogger.Infof("Client 'user-service' Endpoint: %s", userService.GetDiscovery().GetEndpoint()) // Use appLogger
	} else {
		appLogger.Info("No 'user-service' client configuration found.") // Use appLogger
	}

	appLogger.Info("Application finished.") // Log a final message
}
