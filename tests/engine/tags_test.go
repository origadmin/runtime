/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package engine_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/container"
)

func TestContainer_TagsAndCommon(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()

	// 1. Register components with different tags
	// Common provider (no tags) - handles "common" name
	reg.Register("middleware", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		if h.Name() == "common" {
			return "CommonMiddleware", nil
		}
		return nil, nil // Not handled by this provider
	})

	// Gateway specific provider - handles "authn" name
	reg.Register("middleware", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		if h.Name() == "authn" {
			return "GatewayMiddleware", nil
		}
		return nil, nil
	}, engine.WithTag("gateway"))

	// Feature specific provider - handles "authz" name
	reg.Register("middleware", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		if h.Name() == "authz" {
			return "FeatureMiddleware", nil
		}
		return nil, nil
	}, engine.WithTag("feature"))

	// 2. Load configurations
	err := reg.Load(ctx, nil, engine.WithLoadResolver(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{
				{Name: "common", Value: nil},
				{Name: "authn", Value: nil},
				{Name: "authz", Value: nil},
			},
		}, nil
	}))
	assert.NoError(t, err)

	t.Run("GatewayPerspective", func(t *testing.T) {
		// Should get Common + Gateway tags
		h := reg.In("middleware", engine.WithInTags("gateway"))

		results := make(map[string]any)
		for name, inst := range h.Iter(ctx) {
			results[name] = inst
		}

		assert.Contains(t, results, "common")
		assert.Contains(t, results, "authn")
		assert.NotContains(t, results, "authz")
		assert.Equal(t, "CommonMiddleware", results["common"])
		assert.Equal(t, "GatewayMiddleware", results["authn"])
	})

	t.Run("FeaturePerspective", func(t *testing.T) {
		// Should get Common + Feature tags
		h := reg.In("middleware", engine.WithInTags("feature"))

		results := make(map[string]any)
		for name, inst := range h.Iter(ctx) {
			results[name] = inst
		}

		assert.Contains(t, results, "common")
		assert.Contains(t, results, "authz")
		assert.NotContains(t, results, "authn")
		assert.Equal(t, "CommonMiddleware", results["common"])
		assert.Equal(t, "FeatureMiddleware", results["authz"])
	})

	t.Run("FullPerspective", func(t *testing.T) {
		// No tag filter, should get all
		h := reg.In("middleware")

		results := make(map[string]any)
		for name, inst := range h.Iter(ctx) {
			results[name] = inst
		}
		assert.Len(t, results, 3)
		assert.Equal(t, "CommonMiddleware", results["common"])
		assert.Equal(t, "GatewayMiddleware", results["authn"])
		assert.Equal(t, "FeatureMiddleware", results["authz"])
	})
}
