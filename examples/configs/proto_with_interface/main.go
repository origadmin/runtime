package main

import (
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

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

func NewProtoDecoder(c kratosconfig.Config) interfaces.ConfigDecoder {
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
			file.NewSource("examples/configs/proto_with_interface/config/config.yaml"),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	defer c.Close()

	// Create a new decoder.
	decoder := NewProtoDecoder(c)

	// Decode the 'servers' key into a slice of Server structs.
	var servers []*transportv1.Server
	if err := decoder.Decode("servers", &servers); err != nil {
		panic(err)
	}

	// Decode the 'clients' key into a map of ClientConfig structs.
	var clients map[string]*conf.ClientConfig
	// The key is empty because we are decoding the entire config into the Bootstrap struct.
	if err := decoder.Decode("clients", &clients); err != nil {
		panic(err)
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
		fmt.Printf("Client 'user-service' Endpoint: %s\n", userService.GetDiscovery().GetEndpoint())
	} else {
		fmt.Println("No 'user-service' client configuration found.")
	}
}
