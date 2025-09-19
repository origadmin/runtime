package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/mitchellh/mapstructure"

	"github.com/origadmin/runtime"
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
type customConfigDecoder struct {
	kratosConfig kratosconfig.Config
	values       map[string]any // Store the entire config as a map
}

// customDecoderProvider implements the interfaces.ConfigDecoderProvider interface
type customDecoderProvider struct{}

// DefaultCustomDecoder is the default instance of customDecoderProvider.
var DefaultCustomDecoder = &customDecoderProvider{}

// GetConfigDecoder returns a new customConfigDecoder.
func (p *customDecoderProvider) GetConfigDecoder(kratosConfig kratosconfig.Config) (interfaces.ConfigDecoder, error) {
	decoder := &customConfigDecoder{kratosConfig: kratosConfig}
	// Scan the entire config into the internal map upon initialization
	if err := kratosConfig.Scan(&decoder.values); err != nil {
		return nil, fmt.Errorf("failed to scan config into custom decoder values: %w", err)
	}
	// Ensure that after scanning, decoder.values is not empty, indicating successful load
	if len(decoder.values) == 0 {
		return nil, fmt.Errorf("custom decoder values are empty after scanning config")
	}
	return decoder, nil
}

// Config returns the raw configuration data.
func (d *customConfigDecoder) Config() any {
	return d.values // Return the stored map
}

// Decode decodes the configuration into the given target.
func (d *customConfigDecoder) Decode(key string, target any) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	var dataToDecode any
	if key == "" {
		// If key is empty, decode the entire config
		dataToDecode = d.values
	} else {
		// Navigate through the map using the dot-separated key
		var currentValue any = d.values
		keys := strings.Split(key, ".")

		for i, k := range keys {
			currentMap, isMap := currentValue.(map[string]any)
			if !isMap {
				pathSegment := strings.Join(keys[:i], ".")
				return fmt.Errorf("config path '%s' is not a map at segment '%s'", pathSegment, keys[i-1])
			}

			val, ok := currentMap[k]
			if !ok {
				pathSegment := strings.Join(keys[:i+1], ".")
				return fmt.Errorf("config key '%s' not found at path '%s'", k, pathSegment)
			}
			currentValue = val
		}
		dataToDecode = currentValue
	}

	// Configure mapstructure to use "json" tags, allow weakly typed input.
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           target,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.TextUnmarshallerHookFunc(),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(dataToDecode)
}

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
	rawConfig := decoder.Config()
	appLogger.Infof("Type of decoder.Config(): %v", reflect.TypeOf(rawConfig))
	appLogger.Infof("Value of decoder.Config(): %+v", rawConfig)

	rawConfigMap, ok := rawConfig.(map[string]any)
	if !ok {
		appLogger.Error("Error: decoder.Config() is not a map[string]any.")
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
