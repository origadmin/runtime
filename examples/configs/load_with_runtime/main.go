package main

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1" // Added for transportv1.Server
	// Import the generated Go code from the api_gateway proto definition.
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
	"github.com/origadmin/runtime/log" // Import the log package
)

func main() {
	// 1. Create a new Runtime instance from the new api_gateway config.
	//    Path is now relative to the CWD (runtime directory), pointing to the bootstrap.yaml.
	rt, cleanup, err := runtime.NewFromBootstrap(
		"examples/configs/load_with_runtime/config/bootstrap.yaml", // Correctly load bootstrap.yaml
		bootstrap.WithAppInfo(interfaces.AppInfo{
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

	// 3. Decode the entire configuration into our Bootstrap struct.
	// Manually decode servers and clients using JSON marshal/unmarshal for Protobuf compatibility.
	var bc conf.Bootstrap

	// Decode Servers
	var rawServers []interface{}
	if err := decoder.Decode("servers", &rawServers); err != nil {
		appLogger.Errorf("Failed to decode raw servers config: %v", err)
		panic(err)
	}

	for i, rawServer := range rawServers {
		jsonServer, err := json.Marshal(rawServer)
		if err != nil {
			appLogger.Errorf("Failed to marshal server config %d to JSON: %v", i, err)
			panic(err)
		}
		var server transportv1.Server
		if err := protojson.Unmarshal(jsonServer, &server); err != nil {
			appLogger.Errorf("Failed to protojson unmarshal JSON to server %d: %v", i, err)
			panic(err)
		}
		bc.Servers = append(bc.Servers, &server)
	}

	// Decode Clients
	var rawClients map[string]interface{}
	if err := decoder.Decode("clients", &rawClients); err != nil {
		appLogger.Errorf("Failed to decode raw clients config: %v", err)
		panic(err)
	}

	bc.Clients = make(map[string]*conf.ClientConfig)
	for name, rawClient := range rawClients {
		jsonClient, err := json.Marshal(rawClient)
		if err != nil {
			appLogger.Errorf("Failed to marshal client config '%s' to JSON: %v", name, err)
			panic(err)
		}
		var client conf.ClientConfig
		if err := protojson.Unmarshal(jsonClient, &client); err != nil {
			appLogger.Errorf("Failed to protojson unmarshal JSON to client '%s': %v", name, err)
			panic(err)
		}
		bc.Clients[name] = &client
	}

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
