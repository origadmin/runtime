package main

import (
	"context"
	"fmt"
	"time"

	kratosregistry "github.com/go-kratos/kratos/v2/registry"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	"github.com/origadmin/runtime/registry"
)

// --- Temporary Local Registry Implementation for Example --- START

type localDiscovery struct{}

func (d *localDiscovery) GetService(ctx context.Context, serviceName string) ([]*kratosregistry.ServiceInstance, error) {
	return []*kratosregistry.ServiceInstance{
		{
			ID:        fmt.Sprintf("%s-local-0", serviceName),
			Name:      serviceName,
			Version:   "1.0.0",
			Endpoints: []string{"http://localhost:8080"}, // Dummy endpoint
		},
	}, nil
}

func (d *localDiscovery) Watch(ctx context.Context, serviceName string) (kratosregistry.Watcher, error) {
	return &localWatcher{serviceName: serviceName, ctx: ctx}, nil
}

type localRegistrar struct{}

func (r *localRegistrar) Register(ctx context.Context, service *kratosregistry.ServiceInstance) error {
	fmt.Printf("Local Registrar: Registered service %s - %s\n", service.Name, service.ID)
	return nil
}

func (r *localRegistrar) Deregister(ctx context.Context, service *kratosregistry.ServiceInstance) error {
	fmt.Printf("Local Registrar: Deregistered service %s - %s\n", service.Name, service.ID)
	return nil
}

type localWatcher struct {
	serviceName string
	ctx         context.Context
}

func (w *localWatcher) Next() ([]*kratosregistry.ServiceInstance, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-time.After(1 * time.Second): // Simulate some delay
		return []*kratosregistry.ServiceInstance{
			{
				ID:        fmt.Sprintf("%s-local-0", w.serviceName),
				Name:      w.serviceName,
				Version:   "1.0.0",
				Endpoints: []string{"http://localhost:8080"}, // Dummy endpoint
			},
		}, nil
	}
}

func (w *localWatcher) Stop() error {
	return nil
}

type localFactory struct{}

func (f *localFactory) NewDiscovery(cfg *discoveryv1.Discovery, opts ...registry.Option) (kratosregistry.Discovery, error) {
	fmt.Printf("Creating Local Discovery for service: %s\n", cfg.GetServiceName())
	return &localDiscovery{}, nil
}

func (f *localFactory) NewRegistrar(cfg *discoveryv1.Discovery, opts ...registry.Option) (kratosregistry.Registrar, error) {
	fmt.Printf("Creating Local Registrar for service: %s\n", cfg.GetServiceName())
	return &localRegistrar{}, nil
}

// Register the local factory with the default builder.
func init() {
	fmt.Println("DEBUG: Initializing local registry package in example...") // Added debug print
	registry.Register("local", &localFactory{})
}

// --- Temporary Local Registry Implementation for Example --- END
