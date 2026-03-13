/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
	"github.com/origadmin/runtime/helpers/configutil"
)

// Resolve resolves the registry configuration.
func Resolve(source any, _ component.Category) (*component.ModuleConfig, error) {
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

		res := &component.ModuleConfig{Active: comp.ExtractName(def)}
		for _, cfg := range configs {
			if name := comp.ExtractName(cfg); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: cfg})
			}
		}
		return res, nil
	}
	return nil, nil
}
