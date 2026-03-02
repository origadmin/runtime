package main

import (
	"fmt"
	"log"

	rt "github.com/origadmin/runtime"
	conf "github.com/origadmin/runtime/examples/protos/custom_extension"
)

func main() {
	appInfo := rt.NewAppInfo("custom-extension-example", "1.0.0")
	appInfo.Env = "development"
	rtInstance := rt.NewWithAppInfo(appInfo)
	err := rtInstance.Load("examples/configs/load_with_custom_extension/config/bootstrap.yaml")
	if err != nil {
		return
	}
	decoder := rtInstance.Config()
	defer decoder.Close()

	var mw conf.MiddlewareStruct
	if err := decoder.Scan(&mw); err != nil {
		log.Fatalf("Failed to scan config: %v", err)
	}

	fmt.Println("�?Successfully loaded config into MiddlewareStruct.")
	fmt.Printf("   Name: %s\n   Type: %s\n", mw.Name, mw.Type)

	if mw.Customize != nil {
		// Proper way to access fields in google.protobuf.Struct
		fields := mw.Customize.GetFields()
		if p, ok := fields["policy"]; ok {
			fmt.Printf("   Policy: %s\n", p.GetStringValue())
		}
		if s, ok := fields["required_scope"]; ok {
			fmt.Printf("   Required Scope: %s\n", s.GetStringValue())
		}
	}
}
