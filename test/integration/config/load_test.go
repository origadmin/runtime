package config

import (
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/stretchr/testify/assert"

	// Import our test-specific, generated bootstrap proto package
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

func TestConfigLoading(t *testing.T) {
	// 1. Create configuration file source
	source := file.NewSource("./configs/full_config.yaml")

	// 2. Create Kratos config instance
	c := config.New(
		config.WithSource(source),
	)
	defer c.Close()

	// 3. Load configuration
	if err := c.Load(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 4. Scan configuration into our Bootstrap struct
	var bc testconfigs.Bootstrap
	if err := c.Scan(&bc); err != nil {
		t.Fatalf("Failed to scan config: %v", err)
	}

	// 5. --- Start assertion validation ---

	// Validate discovery pool
	assert.Len(t, bc.Discoveries, 2, "Should have 2 discovery configurations")
	assert.Equal(t, "internal-consul", bc.Discoveries[0].Name)
	assert.Equal(t, "my-test-app", bc.Discoveries[0].Config.ServiceName)
	assert.Equal(t, "consul.internal:8500", bc.Discoveries[0].Config.Consul.Address)
	assert.Equal(t, "legacy-etcd", bc.Discoveries[1].Name)
	assert.Equal(t, "etcd.legacy:2379", bc.Discoveries[1].Config.Etcd.Endpoints[0])

	// Validate service registration name
	assert.Equal(t, "internal-consul", bc.RegistrationDiscoveryName)

	// Validate service endpoints
	assert.Len(t, bc.GrpcServers, 1)
	assert.Equal(t, ":9001", bc.GrpcServers[0].Addr)
	assert.Len(t, bc.HttpServers, 1)
	assert.Equal(t, ":8001", bc.HttpServers[0].Addr)

	// Validate clients (most critical part)
	assert.Len(t, bc.Clients, 2, "Should have 2 client configurations")
	// Validate first client
	assert.Equal(t, "user-service", bc.Clients[0].Name)
	assert.Equal(t, "internal-consul", bc.Clients[0].DiscoveryName, "user-service client should use internal-consul")
	assert.Equal(t, "v1.5.0", bc.Clients[0].Selector.Version)
	// Validate second client
	assert.Equal(t, "stock-service", bc.Clients[1].Name)
	assert.Equal(t, "legacy-etcd", bc.Clients[1].DiscoveryName, "stock-service client should use legacy-etcd")
	assert.Equal(t, "v1.0.1", bc.Clients[1].Selector.Version)

	t.Log("Config loaded and verified successfully!")
}
