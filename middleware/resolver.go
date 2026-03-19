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
	"github.com/origadmin/runtime/helpers/configutil"
	"github.com/origadmin/runtime/log"
)

// Resolve resolves the middleware configuration.
func Resolve(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
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
			name := configutil.ExtractName(entry)
			if name == "" {
				continue
			}

			resEntry := component.ConfigEntry{
				Name:  name,
				Value: entry,
			}

			// SELECTIVE INJECTION: Only Selector needs the Carrier logic.
			if entry.GetType() == string(Selector) {
				resEntry.RequirementResolver = resolveSelectorRequirement
			} else {
				resEntry.RequirementResolver = resolveBaseRequirement
			}

			res.Entries = append(res.Entries, resEntry)
		}
		return res, nil
	}
	return nil, nil
}

// getBaseOptions collects generic options like Logger for all middlewares.
func getBaseOptions(ctx context.Context, h component.Handle) []Option {
	var resOpts []Option
	// Attach global Logger
	if l, err := comp.Get[log.Logger](ctx, h.Locator().In(CategoryLogger)); err == nil {
		resOpts = append(resOpts, log.WithLogger(l))
	}
	return resOpts
}

// getCarrierOptions collects Carrier logic specifically for Selectors.
func getCarrierOptions(ctx context.Context, h component.Handle) []Option {
	var resOpts []Option
	clients := make(map[string]KMiddleware)
	servers := make(map[string]KMiddleware)

	// h.Locator() is already scoped and automatically skips the requester (self).
	var it = h.Locator().Iter(ctx)
	for it.Next() {
		name, inst := it.Value()
		if m, ok := inst.(KMiddleware); ok {
			if h.Scope() == ClientScope {
				clients[name] = m
			} else {
				servers[name] = m
			}
		}
	}

	if h.Scope() == ClientScope && len(clients) > 0 {
		resOpts = append(resOpts, WithClientCarrier(clients))
	} else if h.Scope() == ServerScope && len(servers) > 0 {
		resOpts = append(resOpts, WithServerCarrier(servers))
	}
	return resOpts
}

// resolveBaseRequirement provides generic options for standard middlewares.
func resolveBaseRequirement(ctx context.Context, h component.Handle, purpose string) (any, error) {
	if purpose == RequirementOption {
		return getBaseOptions(ctx, h), nil
	}
	return nil, fmt.Errorf("middleware: unknown base requirement %s", purpose)
}

// resolveSelectorRequirement provides full options including Carrier for Selectors.
func resolveSelectorRequirement(ctx context.Context, h component.Handle, purpose string) (any, error) {
	if purpose == RequirementOption {
		opts := getBaseOptions(ctx, h)
		opts = append(opts, getCarrierOptions(ctx, h)...)
		return opts, nil
	}
	return nil, fmt.Errorf("middleware: unknown selector requirement %s", purpose)
}
