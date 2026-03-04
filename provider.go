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

	return nil, nil
}

// resolveFromSource attempts to convert a generic source into a ModuleConfig.
func resolveFromSource(val any) (*component.ModuleConfig, error) {
	if val == nil {
		return nil, nil
	}

	v := reflect.ValueOf(val)
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
	} else if activeField := v.FieldByName("Active"); activeField.IsValid() {
		if activeField.Kind() == reflect.Ptr && !activeField.IsNil() {
			res.Active = activeField.Elem().String()
		} else if activeField.Kind() == reflect.String {
			res.Active = activeField.String()
		}
	}

	// 2. Get Default object
	var defaultItem any
	if defaultMethod := v.MethodByName("GetDefault"); defaultMethod.IsValid() {
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
	if configsMethod := v.MethodByName("GetConfigs"); configsMethod.IsValid() {
		results := configsMethod.Call(nil)
		if len(results) > 0 && results[0].Kind() == reflect.Slice {
			for i := 0; i < results[0].Len(); i++ {
				configs = append(configs, results[0].Index(i).Interface())
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

	// 4. Special Case: No Default but only ONE config exists
	if defaultItem == nil && len(configs) == 1 {
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
	if c, ok := source.(component.RegistryConfig); ok {
		discoveries := c.GetDiscoveries()
		if discoveries == nil {
			return &component.ModuleConfig{}, nil
		}
		res := &component.ModuleConfig{Active: discoveries.GetActive()}
		uniqueEntries := make(map[string]any)

		if d := discoveries.GetDefault(); d != nil {
			uniqueEntries["default"] = d
			if name := extractName(d); name != "" {
				uniqueEntries[name] = d
			}
		}

		configs := discoveries.GetConfigs()
		for _, entry := range configs {
			if name := extractName(entry); name != "" {
				uniqueEntries[name] = entry
			}
		}

		if uniqueEntries["default"] == nil && len(configs) == 1 {
			uniqueEntries["default"] = configs[0]
		}

		for name, val := range uniqueEntries {
			res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: val})
		}
		if res.Active == "" && uniqueEntries["default"] != nil {
			res.Active = "default"
		}
		return res, nil
	}
	return resolveGeneric(source, CategoryRegistry)
}

// resolveMiddleware handles extraction for Middleware components.
func resolveMiddleware(source any) (*component.ModuleConfig, error) {
	if c, ok := source.(component.MiddlewareConfig); ok {
		mws := c.GetMiddlewares()
		if mws == nil {
			return &component.ModuleConfig{}, nil
		}
		res := &component.ModuleConfig{}
		uniqueEntries := make(map[string]any)

		configs := mws.GetConfigs()
		for _, entry := range configs {
			if name := extractName(entry); name != "" {
				uniqueEntries[name] = entry
			}
		}

		// Protocol: If only one config exists and no default, promote it to "default"
		if len(configs) == 1 {
			uniqueEntries["default"] = configs[0]
		}

		for name, val := range uniqueEntries {
			res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: val})
		}
		if res.Active == "" && uniqueEntries["default"] != nil {
			res.Active = "default"
		}
		return res, nil
	}
	return resolveGeneric(source, CategoryMiddleware)
}

// resolveDatabase handles extraction for Database components.
func resolveDatabase(source any) (*component.ModuleConfig, error) {
	if c, ok := source.(component.DataConfig); ok {
		data := c.GetData()
		if data == nil || data.GetDatabases() == nil {
			return &component.ModuleConfig{}, nil
		}
		dbs := data.GetDatabases()
		res := &component.ModuleConfig{Active: dbs.GetActive()}
		uniqueEntries := make(map[string]any)

		if d := dbs.GetDefault(); d != nil {
			uniqueEntries["default"] = d
			if name := extractName(d); name != "" {
				uniqueEntries[name] = d
			}
		}

		configs := dbs.GetConfigs()
		for _, entry := range configs {
			if name := extractName(entry); name != "" {
				uniqueEntries[name] = entry
			}
		}

		if uniqueEntries["default"] == nil && len(configs) == 1 {
			uniqueEntries["default"] = configs[0]
		}

		for name, val := range uniqueEntries {
			res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: val})
		}
		if res.Active == "" && uniqueEntries["default"] != nil {
			res.Active = "default"
		}
		return res, nil
	}
	return resolveGeneric(source, CategoryDatabase)
}

// resolveCache handles extraction for Cache components.
func resolveCache(source any) (*component.ModuleConfig, error) {
	if c, ok := source.(component.DataConfig); ok {
		data := c.GetData()
		if data == nil || data.GetCaches() == nil {
			return &component.ModuleConfig{}, nil
		}
		caches := data.GetCaches()
		res := &component.ModuleConfig{Active: caches.GetActive()}
		uniqueEntries := make(map[string]any)

		if d := caches.GetDefault(); d != nil {
			uniqueEntries["default"] = d
			if name := extractName(d); name != "" {
				uniqueEntries[name] = d
			}
		}

		configs := caches.GetConfigs()
		for _, entry := range configs {
			if name := extractName(entry); name != "" {
				uniqueEntries[name] = entry
			}
		}

		if uniqueEntries["default"] == nil && len(configs) == 1 {
			uniqueEntries["default"] = configs[0]
		}

		for name, val := range uniqueEntries {
			res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: val})
		}
		if res.Active == "" && uniqueEntries["default"] != nil {
			res.Active = "default"
		}
		return res, nil
	}
	return resolveGeneric(source, CategoryCache)
}

// resolveObjectStore handles extraction for ObjectStore components.
func resolveObjectStore(source any) (*component.ModuleConfig, error) {
	if c, ok := source.(component.DataConfig); ok {
		data := c.GetData()
		if data == nil || data.GetObjectStores() == nil {
			return &component.ModuleConfig{}, nil
		}
		oss := data.GetObjectStores()
		res := &component.ModuleConfig{Active: oss.GetActive()}
		uniqueEntries := make(map[string]any)

		if d := oss.GetDefault(); d != nil {
			uniqueEntries["default"] = d
			if name := extractName(d); name != "" {
				uniqueEntries[name] = d
			}
		}

		configs := oss.GetConfigs()
		for _, entry := range configs {
			if name := extractName(entry); name != "" {
				uniqueEntries[name] = entry
			}
		}

		if uniqueEntries["default"] == nil && len(configs) == 1 {
			uniqueEntries["default"] = configs[0]
		}

		for name, val := range uniqueEntries {
			res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: val})
		}
		if res.Active == "" && uniqueEntries["default"] != nil {
			res.Active = "default"
		}
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
