package container

import (
	"context"
	"fmt"
	"iter"
	"reflect"
	"sync"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine/metadata"
)

type Status int

const (
	StatusNone Status = iota
	StatusResolving
	StatusReady
)

type moduleKey struct {
	category metadata.Category
	scope    metadata.Scope
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
	scopes    map[metadata.Scope]bool
	provider  component.Provider
	extractor component.Extractor
	priority  int
}

type containerImpl struct {
	regMu    sync.RWMutex
	registry map[metadata.Category][]*providerEntry

	stateMu sync.RWMutex
	states  map[moduleKey]*moduleState

	mu         sync.RWMutex
	rootConfig any
}

func NewContainer() component.Registry {
	return &containerImpl{
		registry: make(map[metadata.Category][]*providerEntry),
		states:   make(map[moduleKey]*moduleState),
	}
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

func (c *containerImpl) findProvider(cat metadata.Category, scope metadata.Scope) (*providerEntry, bool) {
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

func (c *containerImpl) Config() any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rootConfig
}

func (c *containerImpl) Scope() metadata.Scope       { return metadata.GlobalScope }
func (c *containerImpl) Category() metadata.Category { return "" }

func (c *containerImpl) In(cat metadata.Category, opts ...component.InOption) component.Handle {
	o := &component.InOptions{Scope: metadata.GlobalScope}
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

func (c *containerImpl) Register(cat metadata.Category, e component.Extractor, p component.Provider, opts ...component.RegisterOption) {
	o := &component.RegistrationOptions{
		Priority: metadata.PriorityInfrastructure,
	}
	for _, opt := range opts {
		opt(o)
	}
	entry := &providerEntry{
		scopes:    make(map[metadata.Scope]bool),
		provider:  p,
		extractor: e,
		priority:  o.Priority,
	}
	for _, s := range o.Scopes {
		entry.scopes[s] = true
	}
	c.regMu.Lock()
	defer c.regMu.Unlock()
	c.registry[cat] = append(c.registry[cat], entry)
}

func (c *containerImpl) Has(cat metadata.Category, opts ...component.RegisterOption) bool {
	o := &component.RegistrationOptions{}
	for _, opt := range opts {
		opt(o)
	}
	scope := metadata.GlobalScope
	if len(o.Scopes) > 0 {
		scope = o.Scopes[0]
	}
	_, found := c.findProvider(cat, scope)
	return found
}

func (c *containerImpl) Init(ctx context.Context, root any) error {
	if root == nil {
		return fmt.Errorf("engine: cannot init with nil root config")
	}
	c.mu.Lock()
	c.rootConfig = root
	c.mu.Unlock()

	c.regMu.RLock()
	var cats []metadata.Category
	for cat := range c.registry {
		cats = append(cats, cat)
	}
	c.regMu.RUnlock()

	standardScopes := []metadata.Scope{metadata.GlobalScope, metadata.ServerScope, metadata.ClientScope}
	for _, cat := range cats {
		for _, scope := range standardScopes {
			if _, found := c.findProvider(cat, scope); found {
				if _, err := c.getInternal(ctx, cat, scope, ""); err != nil {
					continue
				}
			}
		}
	}
	return nil
}

func (c *containerImpl) getInternal(ctx context.Context, cat metadata.Category, scope metadata.Scope, name string) (any, error) {
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

	if entry, found := c.findProvider(cat, scope); found {
		if err := c.bindWithEntry(cat, scope, entry); err == nil {
			s.mu.RLock()
			if actualName == "" {
				actualName = s.defaultName
			}
			meta, exists = s.instances[actualName]
			s.mu.RUnlock()
			if exists {
				return c.resolveMeta(ctx, cat, scope, actualName, meta)
			}
		}
	}

	if scope != metadata.GlobalScope {
		return c.getInternal(ctx, cat, metadata.GlobalScope, name)
	}
	return nil, fmt.Errorf("engine: component %s/%s not found", cat, name)
}

func (c *containerImpl) resolveMeta(ctx context.Context, cat metadata.Category, scope metadata.Scope, actualName string, meta *componentMeta) (any, error) {
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
		return nil, fmt.Errorf("engine: no provider for %s", cat)
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

func (c *containerImpl) bindWithEntry(cat metadata.Category, scope metadata.Scope, entry *providerEntry) error {
	mKey := moduleKey{category: cat, scope: scope}
	s := c.getModuleState(mKey)
	s.mu.Lock()
	defer s.mu.Unlock()
	root := c.Config()
	if root == nil {
		return fmt.Errorf("engine: root config is nil")
	}
	mc, err := entry.extractor(root)
	if err != nil || mc == nil {
		return err
	}
	for _, cfgEntry := range mc.Entries {
		if _, exists := s.instances[cfgEntry.Name]; !exists {
			s.instances[cfgEntry.Name] = &componentMeta{config: cfgEntry.Value, status: StatusNone}
			s.order = append(s.order, cfgEntry.Name)
		}
	}
	if mc.Active != "" {
		s.defaultName = mc.Active
	} else if len(mc.Entries) > 0 && s.defaultName == "" {
		s.defaultName = mc.Entries[0].Name
	}
	s.bound = true
	return nil
}

type scopedHandle struct {
	container *containerImpl
	scope     metadata.Scope
	category  metadata.Category
	name      string
	config    any
}

func (h *scopedHandle) Scope() metadata.Scope       { return h.scope }
func (h *scopedHandle) Category() metadata.Category { return h.category }
func (h *scopedHandle) Config() any                 { return h.config }

func (h *scopedHandle) In(cat metadata.Category, opts ...component.InOption) component.Handle {
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
		if entry, found := h.container.findProvider(h.category, h.scope); found {
			_ = h.container.bindWithEntry(h.category, h.scope, entry)
		}
		s.mu.RLock()
		for _, n := range s.order {
			uniqueNames[n] = true
		}
		s.mu.RUnlock()
		if h.scope != metadata.GlobalScope {
			gKey := moduleKey{category: h.category, scope: metadata.GlobalScope}
			gs := h.container.getModuleState(gKey)
			if entry, found := h.container.findProvider(h.category, metadata.GlobalScope); found {
				_ = h.container.bindWithEntry(h.category, metadata.GlobalScope, entry)
			}
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
