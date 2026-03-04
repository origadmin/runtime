package container

import (
	"context"
	"fmt"
	"iter"
	"reflect"
	"sync"

	"github.com/origadmin/runtime/contracts/component"
)

type Status int

const (
	StatusNone Status = iota
	StatusResolving
	StatusReady
)

type moduleKey struct {
	category component.Category
	scope    component.Scope
}

type componentMeta struct {
	mu       sync.Mutex
	config   any
	status   Status
	instance any
}

type moduleState struct {
	mu          sync.RWMutex
	bound       bool
	order       []string
	defaultName string
	instances   map[string]*componentMeta
}

type providerEntry struct {
	scopes   map[component.Scope]bool
	provider component.Provider
	resolver component.Resolver
	priority component.Priority
}

type containerImpl struct {
	regMu    sync.RWMutex
	registry map[component.Category][]*providerEntry

	stateMu sync.RWMutex
	states  map[moduleKey]*moduleState

	mu             sync.RWMutex
	globalResolver component.Resolver
}

// NewContainer returns a clean, parameterless container instance.
func NewContainer() component.Registry {
	return &containerImpl{
		registry: make(map[component.Category][]*providerEntry),
		states:   make(map[moduleKey]*moduleState),
	}
}

func (c *containerImpl) SetResolver(res component.Resolver) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.globalResolver = res
}

func (c *containerImpl) getModuleState(mKey moduleKey) *moduleState {
	c.stateMu.RLock()
	s, ok := c.states[mKey]
	c.stateMu.RUnlock()
	if ok {
		return s
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if s, ok = c.states[mKey]; ok {
		return s
	}
	s = &moduleState{instances: make(map[string]*componentMeta)}
	c.states[mKey] = s
	return s
}

func (c *containerImpl) findProvider(cat component.Category, scope component.Scope) (*providerEntry, bool) {
	c.regMu.RLock()
	defer c.regMu.RUnlock()
	entries, ok := c.registry[cat]
	if !ok {
		return nil, false
	}
	var fallback *providerEntry
	for _, entry := range entries {
		if entry.scopes[scope] {
			return entry, true
		}
		if len(entry.scopes) == 0 {
			fallback = entry
		}
	}
	return fallback, fallback != nil
}

func (c *containerImpl) Config() any                  { return nil }
func (c *containerImpl) Scope() component.Scope       { return component.GlobalScope }
func (c *containerImpl) Category() component.Category { return "" }

func (c *containerImpl) In(cat component.Category, opts ...component.InOption) component.Handle {
	o := &component.InOptions{Scope: component.GlobalScope}
	for _, opt := range opts {
		opt(o)
	}
	return &scopedHandle{
		container: c,
		scope:     o.Scope,
		category:  cat,
	}
}

func (c *containerImpl) Get(ctx context.Context, name string) (any, error) {
	return nil, fmt.Errorf("engine: must use In(category) before Get()")
}

func (c *containerImpl) Iter(ctx context.Context) iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {}
}

func (c *containerImpl) BindConfig(target any) error {
	return fmt.Errorf("engine: BindConfig must be called from Provider handle")
}

func (c *containerImpl) Register(cat component.Category, p component.Provider, opts ...component.RegisterOption) {
	if component.IsReserved(string(cat)) {
		panic(fmt.Sprintf("engine: category name '%s' is reserved", cat))
	}

	o := &component.RegistrationOptions{
		Priority: 100,
	}
	for _, opt := range opts {
		opt(o)
	}

	entry := &providerEntry{
		scopes:   make(map[component.Scope]bool),
		provider: p,
		resolver: o.Resolver,
		priority: o.Priority,
	}
	for _, s := range o.Scopes {
		entry.scopes[s] = true
	}

	c.regMu.Lock()
	defer c.regMu.Unlock()
	c.registry[cat] = append(c.registry[cat], entry)
}

func (c *containerImpl) Has(cat component.Category, opts ...component.RegisterOption) bool {
	o := &component.RegistrationOptions{}
	for _, opt := range opts {
		opt(o)
	}
	scope := component.GlobalScope
	if len(o.Scopes) > 0 {
		scope = o.Scopes[0]
	}
	_, found := c.findProvider(cat, scope)
	return found
}

func (c *containerImpl) Load(ctx context.Context, source any, opts ...component.LoadOption) error {
	o := &component.LoadOptions{}
	for _, opt := range opts {
		opt(o)
	}

	c.regMu.RLock()
	var targets []component.Category
	if o.Category != "" {
		targets = append(targets, o.Category)
	} else {
		for cat := range c.registry {
			targets = append(targets, cat)
		}
	}
	c.regMu.RUnlock()

	standardScopes := []component.Scope{component.GlobalScope, "server", "client"}
	for _, cat := range targets {
		for _, scope := range standardScopes {
			entry, found := c.findProvider(cat, scope)
			if !found {
				continue
			}

			res := o.Resolver
			if res == nil {
				res = c.globalResolver
			}

			if err := c.bindWithSource(cat, scope, entry, source, res, o.Name); err != nil {
				continue
			}
		}
	}
	return nil
}

func (c *containerImpl) bindWithSource(cat component.Category, scope component.Scope, entry *providerEntry, source any, resolver component.Resolver, filterName string) error {
	mKey := moduleKey{category: cat, scope: scope}
	s := c.getModuleState(mKey)
	s.mu.Lock()
	defer s.mu.Unlock()

	var mc *component.ModuleConfig
	var err error

	// Priority: Local Resolver (formerly Extractor) > Load Option Resolver > Container Global Resolver
	if entry.resolver != nil {
		mc, err = entry.resolver(source, cat)
	} else if resolver != nil {
		mc, err = resolver(source, cat)
	}

	if err != nil || mc == nil {
		return err
	}

	for _, cfgEntry := range mc.Entries {
		if filterName != "" && cfgEntry.Name != filterName {
			continue
		}
		if _, exists := s.instances[cfgEntry.Name]; !exists {
			s.instances[cfgEntry.Name] = &componentMeta{config: cfgEntry.Value, status: StatusNone}
			s.order = append(s.order, cfgEntry.Name)
		}
	}

	if mc.Active != "" && (filterName == "" || mc.Active == filterName) {
		s.defaultName = mc.Active
	} else if len(mc.Entries) > 0 && s.defaultName == "" {
		s.defaultName = mc.Entries[0].Name
	}
	s.bound = true
	return nil
}

func (c *containerImpl) getInternal(ctx context.Context, cat component.Category, scope component.Scope, name string) (any, error) {
	mKey := moduleKey{category: cat, scope: scope}
	s := c.getModuleState(mKey)
	s.mu.RLock()
	actualName := name
	if actualName == "" {
		actualName = s.defaultName
	}
	meta, exists := s.instances[actualName]
	s.mu.RUnlock()

	if exists {
		return c.resolveMeta(ctx, cat, scope, actualName, meta)
	}

	if scope != component.GlobalScope {
		return c.getInternal(ctx, cat, component.GlobalScope, name)
	}
	return nil, fmt.Errorf("engine: component %s/%s not found", cat, name)
}

func (c *containerImpl) resolveMeta(ctx context.Context, cat component.Category, scope component.Scope, actualName string, meta *componentMeta) (any, error) {
	if meta.status == StatusReady {
		return meta.instance, nil
	}
	meta.mu.Lock()
	defer meta.mu.Unlock()
	if meta.status == StatusReady {
		return meta.instance, nil
	}
	if meta.status == StatusResolving {
		return nil, fmt.Errorf("engine: circular dependency detected for %s/%s", cat, actualName)
	}
	meta.status = StatusResolving
	entry, found := c.findProvider(cat, scope)
	if !found {
		meta.status = StatusNone
		return nil, fmt.Errorf("engine: no provider for %s during resolution", cat)
	}
	h := &scopedHandle{
		container: c,
		scope:     scope,
		category:  cat,
		name:      actualName,
		config:    meta.config,
	}
	meta.mu.Unlock()
	inst, err := entry.provider(ctx, h)
	meta.mu.Lock()
	if err != nil {
		meta.status = StatusNone
		return nil, err
	}
	meta.instance = inst
	meta.status = StatusReady
	return inst, nil
}

type scopedHandle struct {
	container *containerImpl
	scope     component.Scope
	category  component.Category
	name      string
	config    any
}

func (h *scopedHandle) Scope() component.Scope       { return h.scope }
func (h *scopedHandle) Category() component.Category { return h.category }
func (h *scopedHandle) Config() any                  { return h.config }

func (h *scopedHandle) In(cat component.Category, opts ...component.InOption) component.Handle {
	o := &component.InOptions{Scope: h.scope}
	for _, opt := range opts {
		opt(o)
	}
	return &scopedHandle{
		container: h.container,
		scope:     o.Scope,
		category:  cat,
	}
}

func (h *scopedHandle) Get(ctx context.Context, name string) (any, error) {
	if h.category == "" {
		return nil, fmt.Errorf("engine: category context is required for Get()")
	}
	return h.container.getInternal(ctx, h.category, h.scope, name)
}

func (h *scopedHandle) Iter(ctx context.Context) iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		if h.category == "" {
			return
		}
		uniqueNames := make(map[string]bool)
		mKey := moduleKey{category: h.category, scope: h.scope}
		s := h.container.getModuleState(mKey)
		s.mu.RLock()
		for _, n := range s.order {
			uniqueNames[n] = true
		}
		s.mu.RUnlock()
		if h.scope != component.GlobalScope {
			gKey := moduleKey{category: h.category, scope: component.GlobalScope}
			gs := h.container.getModuleState(gKey)
			gs.mu.RLock()
			for _, n := range gs.order {
				uniqueNames[n] = true
			}
			gs.mu.RUnlock()
		}
		for name := range uniqueNames {
			inst, err := h.Get(ctx, name)
			if err == nil && inst != nil {
				if !yield(name, inst) {
					return
				}
			}
		}
	}
}

func (h *scopedHandle) BindConfig(target any) error {
	vTarget := reflect.ValueOf(target)
	if vTarget.Kind() != reflect.Ptr || vTarget.IsNil() {
		return fmt.Errorf("engine: BindConfig target must be a non-nil pointer")
	}
	vSrc := reflect.ValueOf(h.config)
	if h.config == nil {
		return fmt.Errorf("engine: config is nil")
	}
	if vSrc.Kind() == reflect.Ptr {
		vSrc = vSrc.Elem()
	}
	vTargetElem := vTarget.Elem()
	if !vSrc.Type().AssignableTo(vTargetElem.Type()) {
		return fmt.Errorf("engine: cannot assign %T to %T", h.config, target)
	}
	vTargetElem.Set(vSrc)
	return nil
}
