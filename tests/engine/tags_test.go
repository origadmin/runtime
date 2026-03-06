/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package engine_test

import (
	"context"
	"fmt"
	"sync/atomic"
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
	reg.Register("middleware", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		if h.Name() == "common" {
			return "CommonMiddleware", nil
		}
		return nil, nil
	})

	reg.Register("middleware", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		if h.Name() == "authn" {
			return "GatewayMiddleware", nil
		}
		return nil, nil
	}, engine.WithTag("gateway"))

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
		h := reg.In("middleware", engine.WithInTags("gateway"))
		results := make(map[string]any)
		for name, inst := range h.Iter(ctx) {
			results[name] = inst
		}
		assert.Contains(t, results, "common")
		assert.Contains(t, results, "authn")
		assert.NotContains(t, results, "authz")
	})

	t.Run("FeaturePerspective", func(t *testing.T) {
		h := reg.In("middleware", engine.WithInTags("feature"))
		results := make(map[string]any)
		for name, inst := range h.Iter(ctx) {
			results[name] = inst
		}
		assert.Contains(t, results, "common")
		assert.Contains(t, results, "authz")
		assert.NotContains(t, results, "authn")
	})
}

func TestContainer_LazyInitializationWithTags(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()

	var gatewayCreated int32
	var featureCreated int32

	// 1. Register with side-effect counters
	reg.Register("lazy", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		atomic.AddInt32(&gatewayCreated, 1)
		return "GatewayInst", nil
	}, engine.WithTag("gateway"))

	reg.Register("lazy", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		atomic.AddInt32(&featureCreated, 1)
		return "FeatureInst", nil
	}, engine.WithTag("feature"))

	// 2. Load Configuration
	err := reg.Load(ctx, nil, engine.WithLoadResolver(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{
				{Name: "item", Value: nil},
			},
		}, nil
	}))
	assert.NoError(t, err)

	// VERIFY: Load does not trigger creation
	assert.Equal(t, int32(0), atomic.LoadInt32(&gatewayCreated))
	assert.Equal(t, int32(0), atomic.LoadInt32(&featureCreated))

	// 3. Request Gateway component
	hGateway := reg.In("lazy", engine.WithInTags("gateway"))
	gwInst, err := hGateway.Get(ctx, "item")
	assert.NoError(t, err)
	assert.Equal(t, "GatewayInst", gwInst)

	// VERIFY: Only Gateway is created
	assert.Equal(t, int32(1), atomic.LoadInt32(&gatewayCreated))
	assert.Equal(t, int32(0), atomic.LoadInt32(&featureCreated))

	// 4. Request Feature component - with the fix, this should now work
	// The feature provider should be able to recreate the instance with its own tag
	hFeature := reg.In("lazy", engine.WithInTags("feature"))
	featInst, err := hFeature.Get(ctx, "item")
	assert.NoError(t, err)
	assert.Equal(t, "FeatureInst", featInst)

	// VERIFY: Both providers are now called because each can create its own instance
	assert.Equal(t, int32(1), atomic.LoadInt32(&gatewayCreated))
	assert.Equal(t, int32(1), atomic.LoadInt32(&featureCreated))
}

func TestContainer_TagIsolationNoOverwrite(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()

	var gatewayCalls []string
	var featureCalls []string

	// 1. Register components that track their instances
	reg.Register("service", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		instanceID := fmt.Sprintf("gateway-%d", len(gatewayCalls))
		gatewayCalls = append(gatewayCalls, instanceID)
		return instanceID, nil
	}, engine.WithTag("gateway"))

	reg.Register("service", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		instanceID := fmt.Sprintf("feature-%d", len(featureCalls))
		featureCalls = append(featureCalls, instanceID)
		return instanceID, nil
	}, engine.WithTag("feature"))

	// 2. Load Configuration
	err := reg.Load(ctx, nil, engine.WithLoadResolver(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{
				{Name: "api", Value: nil},
			},
		}, nil
	}))
	assert.NoError(t, err)

	// 3. Request Gateway component multiple times
	hGateway := reg.In("service", engine.WithInTags("gateway"))
	inst1, err := hGateway.Get(ctx, "api")
	assert.NoError(t, err)
	assert.Equal(t, "gateway-0", inst1)

	inst2, err := hGateway.Get(ctx, "api")
	assert.NoError(t, err)
	assert.Equal(t, "gateway-0", inst2) // Should be the same instance

	// 4. Request Feature component multiple times
	hFeature := reg.In("service", engine.WithInTags("feature"))
	inst3, err := hFeature.Get(ctx, "api")
	assert.NoError(t, err)
	assert.Equal(t, "feature-0", inst3)

	inst4, err := hFeature.Get(ctx, "api")
	assert.NoError(t, err)
	assert.Equal(t, "feature-0", inst4) // Should be the same instance

	// 5. Verify no cross-contamination
	assert.Equal(t, 1, len(gatewayCalls), "Gateway provider should only be called once")
	assert.Equal(t, 1, len(featureCalls), "Feature provider should only be called once")

	// 6. Verify gateway instance is still accessible and unchanged
	inst5, err := hGateway.Get(ctx, "api")
	assert.NoError(t, err)
	assert.Equal(t, "gateway-0", inst5) // Should still be the original gateway instance
}