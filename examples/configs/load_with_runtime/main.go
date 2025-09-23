package main

import (
	"encoding/json"
	"errors"
	"fmt"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	// Import the generated Go code from the load_with_runtime proto definition.
	conf "github.com/origadmin/runtime/examples/protos/load_with_runtime"
	"github.com/origadmin/runtime/log" // Import the log package
)

// ProtoConfig is a custom implementation of interfaces.Config that handles Protobuf decoding.
type ProtoConfig struct {
	kratosCfg kratosconfig.Config // Keep for Raw() and Close()
	bootstrap *conf.Bootstrap     // The fully decoded protobuf config
}

// NewProtoConfig creates a new ProtoConfig instance and decodes the entire Kratos config into the Bootstrap proto.
func NewProtoConfig(c kratosconfig.Config) (*ProtoConfig, error) {
	var bc conf.Bootstrap
	if err := c.Load(); err != nil {
		return nil, err
	}
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	return &ProtoConfig{
		kratosCfg: c,
		bootstrap: &bc,
	}, nil
}

// Decode implements the interfaces.Config interface.
func (d *ProtoConfig) Decode(key string, target interface{}) error {
	var source interface{} // This will hold the Go representation of the config section

	if key == "" || key == "/" {
		source = d.bootstrap
	} else {
		// For specific top-level keys, extract the corresponding field from the bootstrap proto.
		// Note: GetX() methods return nil if the field is not set.
		switch key {
		case "servers":
			source = d.bootstrap.GetServers()
		case "clients":
			source = d.bootstrap.GetClients()
		default:
			// If the key is not a direct top-level field of Bootstrap, and not empty/root,
			// we can try to scan it from the original Kratos config as a fallback.
			// This handles cases where some config might not be explicitly defined in the proto.
			return d.kratosCfg.Value(key).Scan(target)
		}
	}

	if source == nil {
		return fmt.Errorf("config key '%s' not found or is nil in bootstrap proto", key)
	}

	// Now, marshal the extracted 'source' (which is a Go struct/slice/map, potentially a proto.Message)
	// to JSON and then unmarshal into the target. This is the most flexible way
	// to handle generic `interface{}` targets, especially when dealing with protobuf messages.
	jsonBytes, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal source data for key '%s' to JSON: %w", key, err)
	}

	if pm, ok := target.(proto.Message); ok {
		// If target is a proto.Message, use protojson.Unmarshal
		if err := protojson.Unmarshal(jsonBytes, pm); err != nil {
			return fmt.Errorf("failed to protojson unmarshal JSON to %T for key '%s': %w", pm, key, err)
		}
	} else {
		// If target is not a proto.Message (e.g., map[string]interface{}, struct), use standard json.Unmarshal
		if err := json.Unmarshal(jsonBytes, target); err != nil {
			return fmt.Errorf("failed to unmarshal JSON to %T for key '%s': %w", target, key, err)
		}
	}

	return nil
}

// Raw implements the interfaces.Config interface.
func (d *ProtoConfig) Raw() kratosconfig.Config {
	return d.kratosCfg
}

// Close implements the interfaces.Config interface.
func (d *ProtoConfig) Close() error {
	return d.kratosCfg.Close()
}

// DecodeLogger implements the interfaces.LoggerConfigDecoder interface.
func (d *ProtoConfig) DecodeLogger() (*loggerv1.Logger, error) {
	// The log config is directly available in d.bootstrap
	if d.bootstrap.GetLogger() == nil {
		return nil, errors.New("logger config not found in bootstrap proto")
	}

	return d.bootstrap.GetLogger(), nil
}

// DecodeDiscoveries implements the interfaces.DiscoveriesConfigDecoder interface.
func (d *ProtoConfig) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	discoveries := make(map[string]*discoveryv1.Discovery)

	// First, check if we have registries configuration with discoveries
	if registries := d.bootstrap.GetRegistries(); registries != nil {
		// Get discoveries from the registries configuration
		for name, discovery := range registries.GetDiscoveries() {
			discoveries[name] = discovery
		}
	}

	// Fallback to client configurations for backward compatibility
	if len(discoveries) == 0 {
		// Iterate through the clients defined in the bootstrap proto
		for clientName, clientCfg := range d.bootstrap.GetClients() {
			if clientCfg.GetDiscovery() != nil {
				// clientCfg.GetDiscovery() returns a *discoveryv1.Client.
				// The target is map[string]*discoveryv1.Discovery.
				// We need to convert discoveryv1.Client to discoveryv1.Discovery.
				jsonBytes, err := protojson.Marshal(clientCfg.GetDiscovery())
				if err != nil {
					return nil, fmt.Errorf("failed to marshal discoveryv1.Client for client '%s' to JSON: %w", clientName, err)
				}
				var disc discoveryv1.Discovery
				if err := protojson.Unmarshal(jsonBytes, &disc); err != nil {
					return nil, fmt.Errorf("failed to protojson unmarshal JSON to discoveryv1.Discovery for client '%s': %w", clientName, err)
				}
				discoveries[clientName] = &disc
			}
		}
	}

	return discoveries, nil
}

func main() {
	// Define the ConfigTransformFunc to create our custom ProtoConfig.
	configTransformer := bootstrap.ConfigTransformFunc(func(kc kratosconfig.Config) (interfaces.Config, error) {
		protoCfg, err := NewProtoConfig(kc)
		if err != nil {
			return nil, fmt.Errorf("failed to create ProtoConfig: %w", err)
		}
		return protoCfg, nil
	})

	// 1. Create a new Runtime instance from the new bootstrap config.
	//    Path is now relative to the CWD (runtime directory), pointing to the bootstrap.yaml.
	rt, cleanup, err := runtime.NewFromBootstrap(
		"examples/configs/load_with_runtime/config/bootstrap.yaml", // Correctly load bootstrap.yaml
		bootstrap.WithAppInfo(interfaces.AppInfo{
			ID:      "rich-config-runtime-example",
			Name:    "RichConfigRuntimeExample",
			Version: "1.0.0",
			Env:     "dev",
		}),
		bootstrap.WithDecoderOptions(bootstrap.WithConfigTransformer(configTransformer)), // Inject custom transformer
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Get the configured logger from the runtime instance
	appLogger := log.NewHelper(rt.Logger()) // Use log.NewHelper for convenience

	appLogger.Info("Application started successfully!") // Log a message using the configured logger

	// 2. Get the configuration decoder from the runtime instance and assert it to our ProtoConfig type
	decoder := rt.Config()

	// Type assert the decoder to our ProtoConfig
	protoCfg, ok := decoder.(*ProtoConfig)
	if !ok {
		appLogger.Error("Failed to assert config to ProtoConfig type")
		panic("config decoder is not of type *ProtoConfig")
	}

	// 3. Get the bootstrap config directly from our ProtoConfig
	bc := protoCfg.bootstrap
	if bc == nil {
		appLogger.Error("Bootstrap config is nil")
		panic("bootstrap config is nil")
	}

	// 4. Print the loaded configuration to verify
	appLogger.Info("--- Loaded Rich Config via runtime interface ---")

	// Verify logger config
	if loggerCfg := bc.GetLogger(); loggerCfg != nil {
		appLogger.Infof("Logger Level: %s, Format: %s, Stdout: %v",
			loggerCfg.GetLevel(),
			loggerCfg.GetFormat(),
			loggerCfg.GetStdout())

		if fileCfg := loggerCfg.GetFile(); fileCfg != nil {
			appLogger.Infof("Log File: %s, MaxSize: %dMB, MaxBackups: %d, MaxAge: %dd, Compress: %v",
				fileCfg.GetPath(),
				fileCfg.GetMaxSize(),
				fileCfg.GetMaxBackups(),
				fileCfg.GetMaxAge(),
				fileCfg.GetCompress())
		}
	} else {
		appLogger.Info("No logger configuration found")
	}

	// Verify server configs
	servers := bc.GetServers()
	if len(servers) > 0 {
		appLogger.Infof("Found %d server configurations", len(servers))
		for i, srv := range servers {
			if httpSrv := srv.GetHttp(); httpSrv != nil {
				appLogger.Infof("Server[%d] HTTP: Addr=%s, Network=%s, Timeout=%s",
					i, httpSrv.GetAddr(), httpSrv.GetNetwork(), httpSrv.GetTimeout().AsDuration())
			}
			if grpcSrv := srv.GetGrpc(); grpcSrv != nil {
				appLogger.Infof("Server[%d] gRPC: Addr=%s, Network=%s, Timeout=%s",
					i, grpcSrv.GetAddr(), grpcSrv.GetNetwork(), grpcSrv.GetTimeout().AsDuration())
			}
		}
	} else {
		appLogger.Info("No server configurations found")
	}

	// Verify client configs
	clients := bc.GetClients()
	if len(clients) > 0 {
		appLogger.Infof("Found %d client configurations", len(clients))
		for name, client := range clients {
			appLogger.Infof("Client '%s':", name)

			if discovery := client.GetDiscovery(); discovery != nil {
				appLogger.Infof("  Discovery: Name=%s, Endpoint=%s, Selector=%s, Version=%s",
					discovery.GetName(),
					discovery.GetEndpoint(),
					discovery.GetSelector().GetType(),
					discovery.GetSelector().GetVersion())
			}

			if transport := client.GetTransport(); transport != nil {
				if grpcCfg := transport.GetGrpc(); grpcCfg != nil {
					appLogger.Infof("  Transport gRPC: Target=%s, Timeout=%s",
						grpcCfg.GetTarget(),
						grpcCfg.GetDialTimeout().AsDuration())
				}
			}
		}
	} else {
		appLogger.Info("No client configurations found")
	}

	// Verify registries and discoveries
	if registries := bc.GetRegistries(); registries != nil {
		discoveries := registries.GetDiscoveries()
		if len(discoveries) > 0 {
			appLogger.Infof("Found %d discovery configurations", len(discoveries))
			for name, disc := range discoveries {
				appLogger.Infof("Discovery '%s': Type=%s, ServiceName=%s",
					name, disc.GetType(), disc.GetServiceName())
			}
		} else {
			appLogger.Info("No discovery configurations found in registries")
		}
	} else {
		appLogger.Info("No registries configuration found")
	}

	appLogger.Info("Application finished.") // Log a final message
}
