package main

import (
	"fmt"
	"log"

	rt "github.com/origadmin/runtime"
	conf "github.com/origadmin/runtime/examples/protos/custom_extension"
)

func main() {
	// Create AppInfo using the new functional options pattern
	appInfo := rt.NewAppInfo(
		"custom-extension-example",
		"1.0.0",
	)
	appInfo.Env = "development"

	// --- 1. Load Configuration ---
	rtInstance := rt.NewWithAppInfo(appInfo)

	err := rtInstance.Load("examples/configs/load_with_custom_extension/config/bootstrap.yaml")
	if err != nil {
		return
	}

	// Get config decoder (KConfig)
	decoder := rtInstance.Decoder()
	defer decoder.Close()

	// --- 2. Directly scan into the CustomAuthConfig proto ---
	var authCfg conf.CustomAuthConfig
	// Use Kratos native API: Value("customize").Scan() instead of Decode
	if err := decoder.Value("customize").Scan(&authCfg); err != nil {
		log.Fatalf("Failed to scan config into CustomAuthConfig struct: %v", err)
	}

	fmt.Println("�?Successfully loaded config into CustomAuthConfig.")
	fmt.Printf("   Policy: %s\n", authCfg.Policy)
	fmt.Printf("   Required Scope: %s\n", authCfg.RequiredScope)
}
