package main

import (
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	rt "github.com/origadmin/runtime"
	conf "github.com/origadmin/runtime/examples/protos/custom_extension"
)

// TempMiddleware is a temporary struct used only for unmarshaling the config file.
// Its 'Customize' field is a generic `map[string]interface{}` to accept the schemaless YAML block.
type TempMiddleware struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Customize map[string]interface{} `json:"customize"`
}

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
	decoder := rtInstance.Config()
	defer decoder.Close()

	// --- 2. Scan into a TEMPORARY Go struct ---
	var tempMw TempMiddleware
	// Use Kratos Scan instead of Decode
	if err := decoder.Scan(&tempMw); err != nil {
		log.Fatalf("Failed to scan config into temporary struct: %v", err)
	}

	fmt.Println("✅ Successfully loaded config into a temporary Go struct.")
	fmt.Printf("   Middleware Name: %s\n", tempMw.Name)
	fmt.Printf("   Raw Customize Map: %v\n", tempMw.Customize)

	var authCfgDirect conf.CustomAuthConfig
	// Use Value().Scan() instead of Decode
	if err := decoder.Value("customize").Scan(&authCfgDirect); err != nil {
		log.Fatalf("Failed to scan customize config: %v", err)
	}
	fmt.Printf("   Direct Customize Config.")
	fmt.Printf("   Policy: %s\n", authCfgDirect.Policy)
	fmt.Printf("   Required Scope: %s\n", authCfgDirect.RequiredScope)

	// --- 3. Manually Convert the Map to the Target Proto ---
	var authCfg conf.CustomAuthConfig
	jsonBytes, err := json.Marshal(tempMw.Customize)
	if err != nil {
		log.Fatalf("Failed to marshal map to JSON: %v", err)
	}
	if err := protojson.Unmarshal(jsonBytes, &authCfg); err != nil {
		log.Fatalf("Failed to unmarshal JSON to CustomAuthConfig: %v", err)
	}

	fmt.Println("\n✅ Manually converted map to strongly-typed CustomAuthConfig.")
	fmt.Printf("   Policy: %s\n", authCfg.Policy)

	anyValue, err := anypb.New(&authCfg)
	if err != nil {
		log.Fatalf("Failed to pack CustomAuthConfig into Any: %v", err)
	}

	fmt.Println("\n✅ Packed the strongly-typed config into a google.protobuf.Any.")

	finalMw := &conf.MiddlewareAny{
		Name:      tempMw.Name,
		Type:      tempMw.Type,
		Customize: anyValue,
	}

	fmt.Println("\n👍 Final Middleware object is constructed and ready to be used.")
	fmt.Printf("   Final Object Name: %s\n", finalMw.Name)
	fmt.Printf("   Is 'Customize' field an 'Any'?: %t\n", finalMw.Customize != nil)
}
