/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"context"
	"fmt"

	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
	"github.com/origadmin/runtime/log"
)

// Resolve resolves the middleware configuration.
func Resolve(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(contracts.MiddlewareConfig); ok {
		mws := c.GetMiddlewares()
		if mws == nil {
			return nil, nil
		}
		res := &component.ModuleConfig{}
		for _, entry := range mws.GetConfigs() {
			if !entry.GetEnabled() {
				continue
			}
			name := comp.ExtractName(entry)
			if name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{
					Name:  name,
					Value: entry,
					// Requirements are no longer used as a map
					RequirementResolver: resolveRequirement,
				})
			}
		}
		// Fallback to first if only one exists
		if len(res.Entries) == 1 {
			res.Active = res.Entries[0].Name
		}
		return res, nil
	}
	return nil, nil
}

// resolveRequirement provides silent resolution for middleware dependencies.
func resolveRequirement(ctx context.Context, h component.Handle, purpose string) (any, error) {
	if purpose == component.RequirementOption {
		// Logic: collect all other components in the current perspective to build options

		// If it's a selector, it needs the carrier of all OTHER middlewares in its scope
		// The engine handles Skip(self) automatically.
		clients := make(map[string]KMiddleware)
		servers := make(map[string]KMiddleware)

		for name, inst := range h.Locator().Iter(ctx) {
			if m, ok := inst.(KMiddleware); ok {
				if h.Scope() == component.ClientScope {
					clients[name] = m
				} else {
					servers[name] = m
				}
			}
		}

		var resOpts []Option
		if h.Scope() == component.ClientScope && len(clients) > 0 {
			resOpts = append(resOpts, WithClientCarrier(clients))
		} else if h.Scope() == component.ServerScope && len(servers) > 0 {
			resOpts = append(resOpts, WithServerCarrier(servers))
		}

		// Also attach global options like Logger
		if l, err := comp.GetDefault[log.Logger](ctx, h.Locator().In(component.CategoryLogger)); err == nil {
			resOpts = append(resOpts, log.WithLogger(l))
		}

		return resOpts, nil
	}
	return nil, fmt.Errorf("middleware: unknown requirement %s", purpose)
}
