package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	rt "github.com/origadmin/runtime"
	conf "github.com/origadmin/runtime/examples/protos/custom_extension"
)

func main() {
	// Initialize runtime with bootstrap configuration
	rtInstance, err := rt.New("Custom Extension Example",
		"1.0.0",
		rt.WithID("custom-extension-example"), // Changed from rt.WithAppInfoID
		rt.WithEnv("development"),             // Changed from rt.WithAppInfoEnv
	)
	rtInstance.Load("examples/configs/load_with_custom_extension/config/bootstrap.yaml")
	if err != nil {
		log.Fatalf("Failed to initialize runtime: %v", err)
	}
	// Removed defer rtInstance.Cleanup() as it's no longer available

	// Get config decoder
	decoder := rtInstance.Config()

	// Decode the entire config into our custom application config
	var appConfig conf.ApplicationConfig
	if err := decoder.Decode("", &appConfig); err != nil {
		log.Fatalf("Failed to decode configuration: %v", err)
	}

	// Process extensions
	if appConfig.CustomizeConfig != nil {
		fmt.Printf("Processing extension: %s\n", appConfig.GetCustomizeConfig())

		// Example of how to unmarshal Any to a specific type
		customCfg, err := getCustomConfig(appConfig.GetCustomizeConfig())
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

func getCustomConfig(customize *structpb.Struct) (*conf.CustomConfig, error) {
	if customize == nil {
		return nil, fmt.Errorf("customize config is nil")
	}

	// Debug: Print the raw data
	fmt.Printf("Raw config data: %+v\n", customize)

	// Convert the data to JSON bytes
	jsonBytes, err := protojson.Marshal(customize)
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
