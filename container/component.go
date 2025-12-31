package container

import (
	"sync"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
)

// DefaultComponentPriority defines the default priority for components.
// Components with lower priority values are initialized first.
// It is recommended to use values in increments of 100 to leave space for future adjustments.
const DefaultComponentPriority = 1000

// ComponentFunc is an adapter to allow the use of ordinary functions as ComponentFactory.
type ComponentFunc func(cfg interfaces.ConfigObject, container Container, opts ...options.Option) (interfaces.Component, error)

// NewComponent calls the wrapped function.
func (c ComponentFunc) NewComponent(cfg interfaces.ConfigObject, container Container, opts ...options.Option) (interfaces.Component, error) {
	return c(cfg, container, opts...)
}

// Priority returns the default priority for a ComponentFunc.
// This ensures that function-based components have a predictable, lower priority.
func (c ComponentFunc) Priority() int {
	return DefaultComponentPriority
}

// ComponentFactory defines the interface for creating generic components.
// It includes a priority system to manage initialization order.
type ComponentFactory interface {
	// Priority determines the initialization order of the component.
	// Components with lower priority values are created and initialized first.
	Priority() int
	// NewComponent creates a new component instance.
	// It receives a component-specific configuration and the container instance,
	// allowing it to register other components or access other services.
	NewComponent(cfg interfaces.ConfigObject, container Container, opts ...options.Option) (interfaces.Component, error)
}

// componentStore is a concurrency-safe storage for component factories and instances.
// It acts as a pure storage mechanism, without any component creation logic.
type componentStore struct {
	mu         sync.RWMutex
	factories  map[string]ComponentFactory
	components map[string]interfaces.Component
	logger     runtimelog.Logger
}

// newComponentStore creates a new, initialized component store.
func newComponentStore(logger runtimelog.Logger) *componentStore {
	return &componentStore{
		factories:  make(map[string]ComponentFactory),
		components: make(map[string]interfaces.Component),
		logger:     logger,
	}
}

// RegisterFactory stores a factory function for a component.
func (s *componentStore) RegisterFactory(name string, factory ComponentFactory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, loaded := s.factories[name]; loaded {
		runtimelog.NewHelper(s.logger).Warnf("component factory with name '%s' is being overwritten", name)
	}
	// A new factory registration should invalidate any existing component instance.
	if _, loaded := s.components[name]; loaded {
		runtimelog.NewHelper(s.logger).Warnf("component instance with name '%s' is being invalidated by a new factory registration", name)
		delete(s.components, name)
	}
	s.factories[name] = factory
}

// RegisterInstance stores a pre-built component instance.
func (s *componentStore) RegisterInstance(name string, comp interfaces.Component) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, loaded := s.factories[name]; loaded {
		runtimelog.NewHelper(s.logger).Warnf("component factory with name '%s' is being overwritten by an instance registration", name)
		delete(s.factories, name)
	}
	if _, loaded := s.components[name]; loaded {
		runtimelog.NewHelper(s.logger).Warnf("component instance with name '%s' is being overwritten", name)
	}
	s.components[name] = comp
}

// GetInstance retrieves a stored component instance.
func (s *componentStore) GetInstance(name string) (interfaces.Component, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	comp, ok := s.components[name]
	return comp, ok
}

// GetFactory retrieves a stored component factory.
func (s *componentStore) GetFactory(name string) (ComponentFactory, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	factory, ok := s.factories[name]
	return factory, ok
}

// Has checks if an instance or a factory is stored for the given name.
func (s *componentStore) Has(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, okInstance := s.components[name]
	_, okFactory := s.factories[name]
	return okInstance || okFactory
}

// List returns a slice of names for all stored instances and factories.
func (s *componentStore) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.components)+len(s.factories))
	seen := make(map[string]struct{})
	for name := range s.components {
		names = append(names, name)
		seen[name] = struct{}{}
	}
	for name := range s.factories {
		if _, ok := seen[name]; !ok {
			names = append(names, name)
		}
	}
	return names
}
