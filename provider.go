/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
)

// --- Default Global Resolver (The Dispatcher) ---

// DefaultGlobalResolver is the primary config dispatcher for the framework.
var DefaultGlobalResolver component.Resolver = func(source any, cat component.Category) (*component.ModuleConfig, error) {
	// Centralized dispatching based on category.
	// This reduces redundant assertions across multiple extractors.
	switch cat {
	case CategoryLogger:
		return resolveLogger(source)
	case CategoryRegistry:
		return resolveRegistry(source)
	case CategoryMiddleware:
		return resolveMiddleware(source)
	case CategoryDatabase:
		return resolveDatabase(source)
	case CategoryCache:
		return resolveCache(source)
	case CategoryObjectStore:
		return resolveObjectStore(source)
	default:
		// Try generic resolution for other categories.
		res, err := resolveGeneric(source, cat)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
		// Unknown categories return an empty object to satisfy the dispatcher.
		return &component.ModuleConfig{}, nil
	}
}

// resolveGeneric attempts to extract configuration using reflection.
func resolveGeneric(source any, cat component.Category) (*component.ModuleConfig, error) {
	if source == nil {
		return nil, nil
	}
	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, nil
	}

	// Try standard naming conventions: Get<Category>s, Get<Category>, etc.
	name := string(cat)
	titleName := strings.ToUpper(name[:1]) + name[1:]
	methodNames := []string{
		"Get" + titleName + "s",
		"Get" + titleName,
	}

	for _, methodName := range methodNames {
		method := v.MethodByName(methodName)
		if !method.IsValid() {
			// Try on the pointer if it was a struct value
			if v.CanAddr() {
				method = v.Addr().MethodByName(methodName)
			}
		}

		if method.IsValid() && method.Type().NumIn() == 0 && method.Type().NumOut() == 1 {
			results := method.Call(nil)
			return resolveFromSource(results[0].Interface())
		}
	}

	// Fallback: Check if the source itself is a container or item for this category
	return resolveFromSource(source)
}

// resolveFromSource attempts to convert a generic source into a ModuleConfig.
func resolveFromSource(val any) (*component.ModuleConfig, error) {
	if val == nil {
		return nil, nil
	}

	v := reflect.ValueOf(val)
	// originalV is used for method calls which might be on pointer receivers
	originalV := v
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, nil
	}

	res := &component.ModuleConfig{}
	uniqueEntries := make(map[string]any)

	// 1. Get Active string
	if a, ok := val.(interface{ GetActive() string }); ok {
		res.Active = a.GetActive()
	} else if activeMethod := originalV.MethodByName("GetActive"); activeMethod.IsValid() {
		results := activeMethod.Call(nil)
		res.Active = results[0].String()
	} else if activeField := v.FieldByName("Active"); activeField.IsValid() {
		if activeField.Kind() == reflect.Ptr && !activeField.IsNil() {
			res.Active = activeField.Elem().String()
		} else if activeField.Kind() == reflect.String {
			res.Active = activeField.String()
		}
	}

	// 2. Get Default object
	var defaultItem any
	if defaultMethod := originalV.MethodByName("GetDefault"); defaultMethod.IsValid() {
		results := defaultMethod.Call(nil)
		if len(results) > 0 && !results[0].IsNil() {
			defaultItem = results[0].Interface()
			uniqueEntries["default"] = defaultItem
			// Also add by its own name
			name := extractName(defaultItem)
			if name != "" {
				uniqueEntries[name] = defaultItem
			}
		}
	}

	// 3. Get Configs slice
	var configs []any
	if configsMethod := originalV.MethodByName("GetConfigs"); configsMethod.IsValid() {
		results := configsMethod.Call(nil)
		if len(results) > 0 && results[0].Kind() == reflect.Slice {
			resVals := results[0]
			for i := 0; i < resVals.Len(); i++ {
				configs = append(configs, resVals.Index(i).Interface())
			}
		}
	} else if configsField := v.FieldByName("Configs"); configsField.IsValid() && configsField.Kind() == reflect.Slice {
		for i := 0; i < configsField.Len(); i++ {
			configs = append(configs, configsField.Index(i).Interface())
		}
	}

	// Add all configs to unique map
	for _, item := range configs {
		name := extractName(item)
		if name != "" {
			uniqueEntries[name] = item
		}
	}

	// 4. Special Case: No Default and No Configs found?
	// Check if the source ITSELF is a config item (has a name/type).
	if defaultItem == nil && len(configs) == 0 {
		selfName := extractName(val)
		if selfName != "" {
			// Treat source as a single item container
			uniqueEntries[selfName] = val
			uniqueEntries["default"] = val
		}
	} else if defaultItem == nil && len(configs) == 1 {
		// Single config in list promotion
		uniqueEntries["default"] = configs[0]
	}

	// Convert map to entries
	for name, item := range uniqueEntries {
		res.Entries = append(res.Entries, component.ConfigEntry{
			Name:  name,
			Value: item,
		})
	}

	// 5. Final Active check
	if res.Active == "" && uniqueEntries["default"] != nil {
		res.Active = "default"
	}

	return res, nil
}

// extractName attempts to get a name from a config item.
func extractName(item any) string {
	if item == nil {
		return ""
	}

	// 1. Try instance name (GetName)
	if n, ok := item.(interface{ GetName() string }); ok {
		if name := n.GetName(); name != "" {
			return name
		}
	}

	// 2. Try implementation identifiers (Dialect, Type, Driver)
	if d, ok := item.(interface{ GetDialect() string }); ok {
		if name := d.GetDialect(); name != "" {
			return name
		}
	}
	if t, ok := item.(interface{ GetType() string }); ok {
		if name := t.GetType(); name != "" {
			return name
		}
	}
	if d, ok := item.(interface{ GetDriver() string }); ok {
		if name := d.GetDriver(); name != "" {
			return name
		}
	}

	// 3. Fallback to type name
	t := reflect.TypeOf(item)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// resolveLogger handles extraction for Logger components.
func resolveLogger(source any) (*component.ModuleConfig, error) {
	if c, ok := source.(component.LoggerConfig); ok {
		logger := c.GetLogger()
		if logger == nil {
			return &component.ModuleConfig{}, nil
		}
		// Treating the single logger segment as the default instance.
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "default", Value: logger}},
		}, nil
	}
	return resolveGeneric(source, CategoryLogger)
}

// resolveRegistry handles extraction for Registry/Discovery components.
func resolveRegistry(source any) (*component.ModuleConfig, error) {
	// 1. Try parent-wrapped mode (RegistryConfig -> Discoveries)
	if c, ok := source.(component.RegistryConfig); ok {
		if d := c.GetDiscoveries(); d != nil {
			return resolveFromSource(d)
		}
	}
	// 2. Try direct container or single item mode
	res, err := resolveFromSource(source)
	if err == nil && res != nil && len(res.Entries) > 0 {
		return res, nil
	}
	return resolveGeneric(source, CategoryRegistry)
}

// resolveMiddleware handles extraction for Middleware components.
func resolveMiddleware(source any) (*component.ModuleConfig, error) {
	// 1. Try parent-wrapped mode (MiddlewareConfig -> Middlewares)
	if c, ok := source.(component.MiddlewareConfig); ok {
		if m := c.GetMiddlewares(); m != nil {
			return resolveFromSource(m)
		}
	}
	// 2. Try direct container or single item mode
	res, err := resolveFromSource(source)
	if err == nil && res != nil && len(res.Entries) > 0 {
		return res, nil
	}
	return resolveGeneric(source, CategoryMiddleware)
}

// resolveDatabase handles extraction for Database components.
func resolveDatabase(source any) (*component.ModuleConfig, error) {
	// 1. Try parent-wrapped mode (DataConfig -> Databases)
	if c, ok := source.(component.DataConfig); ok {
		if data := c.GetData(); data != nil {
			if dbs := data.GetDatabases(); dbs != nil {
				return resolveFromSource(dbs)
			}
		}
	}
	// 2. Try direct container or single item mode
	res, err := resolveFromSource(source)
	if err == nil && res != nil && len(res.Entries) > 0 {
		return res, nil
	}
	return resolveGeneric(source, CategoryDatabase)
}

// resolveCache handles extraction for Cache components.
func resolveCache(source any) (*component.ModuleConfig, error) {
	// 1. Try parent-wrapped mode (DataConfig -> Caches)
	if c, ok := source.(component.DataConfig); ok {
		if data := c.GetData(); data != nil {
			if caches := data.GetCaches(); caches != nil {
				return resolveFromSource(caches)
			}
		}
	}
	// 2. Try direct container or single item mode
	res, err := resolveFromSource(source)
	if err == nil && res != nil && len(res.Entries) > 0 {
		return res, nil
	}
	return resolveGeneric(source, CategoryCache)
}

// resolveObjectStore handles extraction for ObjectStore components.
func resolveObjectStore(source any) (*component.ModuleConfig, error) {
	// 1. Try parent-wrapped mode (DataConfig -> ObjectStores)
	if c, ok := source.(component.DataConfig); ok {
		if data := c.GetData(); data != nil {
			if oss := data.GetObjectStores(); oss != nil {
				return resolveFromSource(oss)
			}
		}
	}
	// 2. Try direct container or single item mode
	res, err := resolveFromSource(source)
	if err == nil && res != nil && len(res.Entries) > 0 {
		return res, nil
	}
	return resolveGeneric(source, CategoryObjectStore)
}

// --- Default Providers (Component Factory) ---

// DefaultLoggerProvider creates a default logger instance.
var DefaultLoggerProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	return log.DefaultLogger, nil
}

// DefaultRegistryProvider creates a default registry instance.
var DefaultRegistryProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	return nil, nil
}

// --- Wire Providers ---

// ProvideLogger is a Wire provider function that extracts the logger from the App.
func ProvideLogger(rt *App) log.Logger {
	return rt.Logger()
}

// ProvideDefaultRegistrar is a Wire provider function that extracts the registrar from the App.
func ProvideDefaultRegistrar(rt *App) (registry.Registrar, error) {
	return rt.DefaultRegistrar()
}
