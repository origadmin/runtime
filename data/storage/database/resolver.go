/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package database

import (
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
	"github.com/origadmin/runtime/helpers/configutil"
)

// Resolve resolves the database configuration.
func Resolve(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(contracts.DataConfig); ok {
		data := c.GetData()
		if data == nil || data.GetDatabases() == nil {
			return nil, nil
		}
		dbs := data.GetDatabases()

		def, configs, err := configutil.Normalize(dbs.GetActive(), dbs.GetDefault(), dbs.GetConfigs())
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
