package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/config/decoder" // Import the new public decoder package
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
	// Import the generated Go code from the api_gateway proto definition.
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
)

// CustomSettings represents the structure of our custom configuration section.
type CustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"endpoints"`
}

// customConfigDecoder implements the interfaces.ConfigDecoder interface
// It embeds decoder.Decoder to inherit default behaviors.
type customConfigDecoder struct {
	*decoder.Decoder // Embed the Decoder from config/decoder
}

// customDecoderProvider implements the interfaces.ConfigDecoderProvider interface
type customDecoderProvider struct{}

// DefaultCustomDecoder is the default instance of customDecoderProvider.
var DefaultCustomDecoder = &customDecoderProvider{}

// GetConfigDecoder returns a new customConfigDecoder.
func (p *customDecoderProvider) GetConfigDecoder(kratosConfig kratosconfig.Config) (interfaces.ConfigDecoder, error) {
	// Initialize the embedded Decoder
	return &customConfigDecoder{
		Decoder: decoder.NewDecoder(kratosConfig),
	}, nil
}

// Note: We are not implementing Decode, DecodeLogger, DecodeDiscoveries here.
// This customConfigDecoder will rely on the embedded Decoder's implementations
// for these methods. If specific custom logic were needed, we would override them here.

func main() {
	// 1. Create a new Runtime instance from the new api_gateway config.
	//    Path is now relative to the CWD (runtime directory), pointing to the bootstrap.yaml.
	rt, cleanup, err := runtime.NewFromBootstrap(
		"examples/configs/load_with_custom_parser/config/bootstrap.yaml",
		runtime.WithAppInfo(runtime.AppInfo{
			ID:      "api-gateway-custom-parser-example",
			Name:    "ApiGatewayCustomParserExample",
			Version: "1.0.0",
			Env:     "dev",
		}),
		runtime.WithDecoderProvider(DefaultCustomDecoder), // Use our custom decoder provider
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	appLogger := log.NewHelper(rt.Logger())
	appLogger.Info("Application started successfully with custom parser!")

	// 2. Get the configuration decoder from the runtime instance.
	decoder := rt.Config()

	// --- DIAGNOSTIC PRINT Kratos Config Content START ---
	appLogger.Info("--- Debugging Kratos Config Content ---")

	// Use the Decode method to get the raw config map, as Config() is no longer available.
	var rawConfig map[string]any
	if err := decoder.Decode("", &rawConfig); err != nil {
		appLogger.Errorf("Error decoding raw config: %v", err)
		// Depending on desired behavior, you might panic or handle more gracefully.
		panic(err)
	}

	appLogger.Infof("Type of rawConfig: %v", reflect.TypeOf(rawConfig))
	appLogger.Infof("Value of rawConfig: %+v", rawConfig)

	rawConfigMap, ok := rawConfig.(map[string]any)
	if !ok {
		appLogger.Error("Error: rawConfig is not a map[string]any.")
	} else {
		if serversVal, exists := rawConfigMap["servers"]; exists {
			if serversStr, err := json.Marshal(serversVal); err == nil {
				appLogger.Infof("Kratos Config 'servers' value (string): %s", serversStr)
			} else {
				appLogger.Errorf("Error marshalling 'servers' value to string: %v", err)
			}
		} else {
			appLogger.Info("Kratos Config 'servers' value not found in raw map.")
		}

		if clientsVal, exists := rawConfigMap["clients"]; exists {
			if clientsStr, err := json.Marshal(clientsVal); err == nil {
				appLogger.Infof("Kratos Config 'clients' value (string): %s", clientsStr)
			} else {
				appLogger.Errorf("Error marshalling 'clients' value to string: %v", err)
			}
		} else {
			appLogger.Info("Kratos Config 'clients' value not found in raw map.")
		}
	}
	appLogger.Info("--------------------------------------")

	// --- Decode Custom Settings ---
	var customSettings CustomSettings
	if err := decoder.Decode("custom_settings", &customSettings); err != nil {
		appLogger.Errorf("Failed to decode custom settings: %v", err)
		panic(err)
	}
	appLogger.Infof("Decoded Custom Settings: %+v", customSettings)
	appLogger.Infof("Custom Setting - Feature Enabled: %t", customSettings.FeatureEnabled)
	appLogger.Infof("Custom Setting - API Key: %s", customSettings.APIKey)
	appLogger.Infof("Custom Setting - Rate Limit: %d", customSettings.RateLimit)
	for i, ep := range customSettings.Endpoints {
		appLogger.Infof("Custom Setting - Endpoint %d: Name=%s, Path=%s", i, ep.Name, ep.Path)
	}

	// 3. Decode the entire configuration into our Bootstrap struct.
	var bc conf.Bootstrap
	if err := decoder.Decode("", &bc); err != nil {
		appLogger.Errorf("Failed to decode bootstrap config: %v", err)
		panic(err)
	}

	// --- DIAGNOSTIC PRINT Bootstrap struct content after Decode START ---
	appLogger.Info("--- Debugging Bootstrap struct content after Decode ---")
	appLogger.Infof("%+v", bc)
	appLogger.Info("--------------------------------------")

	// 4. Print the loaded configuration to verify.
	appLogger.Info("--- Loaded API Gateway config via runtime interface ---")

	// Verify server config
	if len(bc.Servers) > 0 && bc.Servers[0].GetHttp() != nil {
		appLogger.Infof("Server HTTP Addr: %s", bc.Servers[0].GetHttp().GetAddr())
	} else {
		appLogger.Info("No HTTP server configuration found.")
	}

	// Verify client config
	if userService, ok := bc.Clients["user-service"]; ok {
		appLogger.Infof("Client 'user-service' Endpoint: %s", userService.GetDiscovery().GetEndpoint())
	} else {
		appLogger.Info("No 'user-service' client configuration found.")
	}

	appLogger.Info("Application finished.")
}
