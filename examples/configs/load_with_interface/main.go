package main

import (
	"encoding/json"
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"google.golang.org/protobuf/encoding/protojson"
	// Import the transportv1 package which contains the Server message definition.
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	// Import the generated Go code from the api_gateway proto for the ClientConfig message.
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
	"github.com/origadmin/runtime/interfaces"
)

// ProtoDecoder remains the same as it's a generic wrapper.
type ProtoDecoder struct {
	c kratosconfig.Config
}

func (d *ProtoDecoder) Raw() kratosconfig.Config {
	return d.c
}

func (d *ProtoDecoder) Close() error {
	return d.c.Close()
}

func NewProtoDecoder(c kratosconfig.Config) interfaces.Config {
	return &ProtoDecoder{c: c}
}

func (d *ProtoDecoder) Config() any {
	return d.c
}

func (d *ProtoDecoder) Decode(key string, target interface{}) error {
	return d.c.Value(key).Scan(target)
}

func main() {
	// Create a Kratos config instance from the YAML file.
	// Path is now relative to the CWD (runtime directory).
	c := kratosconfig.New(
		kratosconfig.WithSource(
			file.NewSource("examples/configs/load_with_interface/config/config.yaml"),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	defer c.Close()

	// Create a new decoder.
	decoder := NewProtoDecoder(c)

	// Decode the 'servers' key into a slice of Server structs using JSON marshal/unmarshal
	var rawServers []interface{}
	if err := decoder.Decode("servers", &rawServers); err != nil {
		panic(fmt.Errorf("failed to decode raw servers config: %w", err))
	}

	var servers []*transportv1.Server
	for i, rawServer := range rawServers {
		jsonServer, err := json.Marshal(rawServer)
		if err != nil {
			panic(fmt.Errorf("failed to marshal server config %d to JSON: %w", i, err))
		}
		var server transportv1.Server
		if err := protojson.Unmarshal(jsonServer, &server); err != nil {
			panic(fmt.Errorf("failed to protojson unmarshal JSON to server %d: %w", i, err))
		}
		servers = append(servers, &server)
	}

	// Decode the 'clients' key into a map of ClientConfig structs.
	var rawClients map[string]interface{}
	if err := decoder.Decode("clients", &rawClients); err != nil {
		panic(fmt.Errorf("failed to decode raw clients config: %w", err))
	}

	var clients map[string]*conf.ClientConfig = make(map[string]*conf.ClientConfig)
	for name, rawClient := range rawClients {
		jsonClient, err := json.Marshal(rawClient)
		if err != nil {
			panic(fmt.Errorf("failed to marshal client config '%s' to JSON: %w", name, err))
		}
		var client conf.ClientConfig
		if err := protojson.Unmarshal(jsonClient, &client); err != nil {
			panic(fmt.Errorf("failed to protojson unmarshal JSON to client '%s': %w", name, err))
		}
		clients[name] = &client
	}

	// Print the loaded configuration to verify.
	fmt.Printf("--- Loaded API Gateway config via interface implementation ---\n")

	// Verify server config
	if len(servers) > 0 && servers[0].GetHttp() != nil {
		fmt.Printf("Server HTTP Addr: %s\n", servers[0].GetHttp().GetAddr())
	} else {
		fmt.Println("No HTTP server configuration found.")
	}

	// Verify client config
	if userService, ok := clients["user-service"]; ok {
		fmt.Printf("Client 'user-service' Endpoint: %s\n", userService.GetEndpoint())
	} else {
		fmt.Println("No 'user-service' client configuration found.")
	}
}
