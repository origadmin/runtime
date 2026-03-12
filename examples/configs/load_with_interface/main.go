package main

import (
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/origadmin/runtime"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
)

func main() {
	// Create a new runtime App instance.
	app := runtime.New("api-gateway", "1.0.0")

	// Load configuration using the new Bootstrap engine.
	// Path is relative to the runtime directory.
	err := app.Load("examples/configs/load_with_interface/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	defer func() {
		if app.Decoder() != nil {
			_ = app.Decoder().Close()
		}
	}()

	// Access the KConfig decoder directly from the App.
	decoder := app.Decoder()

	// 1. Decode 'servers' into proto messages.
	// KConfig native scanning returns raw map/slice for complex objects if not typed.
	var rawServers []interface{}
	if err := decoder.Value("servers").Scan(&rawServers); err != nil {
		log.Fatalf("Failed to decode raw servers config: %v", err)
	}

	var servers []*transportv1.Server
	for i, rawServer := range rawServers {
		jsonServer, err := json.Marshal(rawServer)
		if err != nil {
			log.Fatalf("Failed to marshal server config %d to JSON: %v", i, err)
		}
		var server transportv1.Server
		if err := protojson.Unmarshal(jsonServer, &server); err != nil {
			log.Fatalf("Failed to protojson unmarshal JSON to server %d: %v", i, err)
		}
		servers = append(servers, &server)
	}

	// 2. Decode 'clients' into a map of proto messages.
	var rawClients map[string]interface{}
	if err := decoder.Value("clients").Scan(&rawClients); err != nil {
		log.Fatalf("Failed to decode raw clients config: %v", err)
	}

	clients := make(map[string]*conf.ClientConfig)
	for name, rawClient := range rawClients {
		jsonClient, err := json.Marshal(rawClient)
		if err != nil {
			log.Fatalf("Failed to marshal client config '%s' to JSON: %v", name, err)
		}
		var client conf.ClientConfig
		if err := protojson.Unmarshal(jsonClient, &client); err != nil {
			log.Fatalf("Failed to protojson unmarshal JSON to client '%s': %v", name, err)
		}
		clients[name] = &client
	}

	// Print the loaded configuration to verify.
	fmt.Printf("--- Loaded API Gateway config using latest runtime App ---\n")

	if len(servers) > 0 && servers[0].GetHttp() != nil {
		fmt.Printf("Server HTTP Addr: %s\n", servers[0].GetHttp().GetAddr())
	}

	if userService, ok := clients["user-service"]; ok {
		if userService.GetClient() != nil {
			fmt.Printf("Client 'user-service' Endpoint: %s\n", userService.GetClient())
		}
	}
}
