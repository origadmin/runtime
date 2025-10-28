package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"

	rt "github.com/origadmin/runtime"
	extensionv1 "github.com/origadmin/runtime/api/gen/go/runtime/extension/v1"
	"github.com/origadmin/runtime/bootstrap"
	conf "github.com/origadmin/runtime/examples/protos/custom_extension_example"
	"github.com/origadmin/runtime/interfaces"
)

func main() {
	// Initialize runtime with bootstrap configuration
	rtInstance, err := rt.NewFromBootstrap(
		"examples/configs/load_with_custom_extension/config/bootstrap.yaml",
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "custom-extension-example",
			Name:    "Custom Extension Example",
			Version: "1.0.0",
			Env:     "development",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer rtInstance.Cleanup()

	// Get config decoder
	decoder := rtInstance.Config()

	// Decode the entire config into our custom application config
	var appConfig conf.ApplicationConfig
	if err := decoder.Decode("", &appConfig); err != nil {
		log.Fatalf("Failed to decode configuration: %v", err)
	}

	// Process extensions
	if appConfig.Customize != nil {
		fmt.Printf("Processing extension: %s\n", appConfig.GetCustomize().GetName())

		// Example of how to unmarshal Any to a specific type
		customCfg, err := getCustomConfig(appConfig.GetCustomize())
		if err != nil {
			log.Printf("Failed to get custom config: %v", err)
			return
		}

		printCustomConfig(customCfg)
	} else {
		log.Println("No customize configuration found")
	}

	// Direct custom config (if not using extensions map)
	//if appConfig.Custom != nil {
	//	printCustomConfig(appConfig.Custom)
	//}
}

func getCustomConfig(customize *extensionv1.CustomizeConfig) (*conf.CustomConfig, error) {
	if customize == nil {
		return nil, fmt.Errorf("customize config is nil")
	}

	// Debug: Print the raw data
	fmt.Printf("Raw config data: %+v\n", customize.GetData())

	// Convert the data to JSON bytes
	jsonBytes, err := protojson.Marshal(customize.GetData())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config data: %w", err)
	}
	fmt.Printf("JSON bytes: %s\n", string(jsonBytes))

	// Unmarshal to CustomConfig
	var cfg conf.CustomConfig
	if err := protojson.Unmarshal(jsonBytes, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to CustomConfig: %w", err)
	}

	return &cfg, nil
}

func printCustomConfig(cfg *conf.CustomConfig) {
	fmt.Printf("\nCustom Configuration:\n")
	fmt.Printf("Custom Field: %s\n", cfg.GetCustomField())
	fmt.Printf("Custom Number: %d\n", cfg.GetCustomNumber())

	fmt.Println("\nItems:")
	for i, item := range cfg.GetItems() {
		fmt.Printf("  %d: %s\n", i+1, item)
	}

	if nested := cfg.GetNested(); nested != nil {
		fmt.Printf("\nNested Config:\n")
		fmt.Printf("  Enabled: %v\n", nested.GetEnabled())
		fmt.Printf("  Name: %s\n", nested.GetName())
	}
}
