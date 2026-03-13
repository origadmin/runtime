/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"context"

	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/configutil"
)

const (
	// CategoryRegistrar is the category for registrar components.
	CategoryRegistrar component.Category = "registrar"
	// CategoryDiscovery is the category for discovery components.
	CategoryDiscovery component.Category = "discovery"
)

// Resolve resolves the registry configuration.
func Resolve(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
	if c, ok := source.(contracts.DiscoveryConfig); ok {
		discoveries := c.GetDiscoveries()
		if discoveries == nil {
			return nil, nil
		}

		// Authorization flow: Default -> Active -> First
		def, configs, err := configutil.Normalize(discoveries.GetActive(), discoveries.GetDefault(), discoveries.GetConfigs())
		if err != nil {
			return nil, err
		}

		res := &component.ModuleConfig{Active: configutil.ExtractName(def)}
		for _, cfg := range configs {
			if name := configutil.ExtractName(cfg); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: cfg})
			}
		}
		return res, nil
	}
	return nil, nil
}
