package container

import (
	"context"
	"fmt"
	"reflect"
	"sort"
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

type regOpts struct {
	scope    metadata.Scope
	priority int
}

func (o *regOpts) SetScope(s metadata.Scope) { o.scope = s }
func (o *regOpts) SetPriority(p int)         { o.priority = p }

type instanceKey struct {
	category metadata.Category
	scope    metadata.Scope
	name     string
}

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

type containerImpl struct {
	regMu      sync.RWMutex
	providers  map[moduleKey]component.Provider
	extractors map[moduleKey]component.Extractor
	priorities map[moduleKey]int

	stateMu sync.RWMutex
	states  map[moduleKey]*moduleState

	mu         sync.RWMutex
	rootConfig any
}

func NewContainer() component.Registry {
	return &containerImpl{
		providers:  make(map[moduleKey]component.Provider),
		extractors: make(map[moduleKey]component.Extractor),
		priorities: make(map[moduleKey]int),
		states:     make(map[moduleKey]*moduleState),
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

func (c *containerImpl) Config() any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rootConfig
}

func (c *containerImpl) Scope() metadata.Scope       { return metadata.GlobalScope }
func (c *containerImpl) Category() metadata.Category { return "" }

func (c *containerImpl) In(cat metadata.Category, opts ...component.RegisterOption) component.Handle {
	o := &regOpts{scope: metadata.GlobalScope}
	for _, opt := range opts {
		opt(o)
	}
	return &scopedHandle{
		container: c,
		scope:     o.scope,
		category:  cat,
	}
}

func (c *containerImpl) Get(ctx context.Context, name string) (any, error) {
	return nil, fmt.Errorf("engine: must use In(category) before Get()")
}

func (c *containerImpl) BindConfig(target any) error {
	return fmt.Errorf("engine: BindConfig must be called from Provider handle")
}

func (c *containerImpl) Register(cat metadata.Category, e component.Extractor, p component.Provider, opts ...component.RegisterOption) {
	o := &regOpts{
		scope:    metadata.GlobalScope,
		priority: metadata.PriorityInfrastructure,
	}
	for _, opt := range opts {
		opt(o)
	}
	c.regMu.Lock()
	defer c.regMu.Unlock()
	mKey := moduleKey{category: cat, scope: o.scope}
	c.providers[mKey] = p
	c.extractors[mKey] = e
	c.priorities[mKey] = o.priority
}

func (c *containerImpl) Has(cat metadata.Category, opts ...component.RegisterOption) bool {
	o := &regOpts{scope: metadata.GlobalScope}
	for _, opt := range opts {
		opt(o)
	}
	c.regMu.RLock()
	defer c.regMu.RUnlock()
	mKey := moduleKey{category: cat, scope: o.scope}
	_, exists := c.providers[mKey]
	return exists
}

func (c *containerImpl) Init(ctx context.Context, root any) error {
	if root == nil {
		return fmt.Errorf("engine: cannot init with nil root config")
	}
	c.mu.Lock()
	c.rootConfig = root
	c.mu.Unlock()

	c.regMu.RLock()
	var keys []moduleKey
	for k := range c.providers {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return c.priorities[keys[i]] < c.priorities[keys[j]]
	})
	c.regMu.RUnlock()

	for _, k := range keys {
		if err := c.bindCategory(k.category, k.scope, k); err != nil {
			return fmt.Errorf("engine: failed to bind %s in scope %s: %w", k.category, k.scope, err)
		}
		s := c.getModuleState(k)
		s.mu.RLock()
		order := s.order
		s.mu.RUnlock()
		for _, name := range order {
			if _, err := c.getInternal(ctx, k.category, k.scope, name); err != nil {
				return fmt.Errorf("engine: failed to init %s/%s: %w", k.category, name, err)
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
	isBound := s.bound
	s.mu.RUnlock()

	if !isBound || (!exists && name != "") || (actualName == "" && name == "") {
		if err := c.bindCategory(cat, scope, mKey); err != nil {
			if scope != metadata.GlobalScope {
				return c.getInternal(ctx, cat, metadata.GlobalScope, name)
			}
			return nil, err
		}
		s.mu.RLock()
		if actualName == "" {
			actualName = s.defaultName
		}
		meta, exists = s.instances[actualName]
		s.mu.RUnlock()
		if !exists {
			if scope != metadata.GlobalScope {
				return c.getInternal(ctx, cat, metadata.GlobalScope, name)
			}
			return nil, fmt.Errorf("engine: component %s/%s not found", cat, name)
		}
	}

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
	c.regMu.RLock()
	provider, ok := c.providers[mKey]
	if !ok {
		provider, ok = c.providers[moduleKey{category: cat, scope: metadata.GlobalScope}]
	}
	c.regMu.RUnlock()

	if !ok {
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
	inst, err := provider(ctx, h)
	meta.mu.Lock()

	if err != nil {
		meta.status = StatusNone
		return nil, err
	}

	meta.instance = inst
	meta.status = StatusReady
	return inst, nil
}

func (c *containerImpl) bindCategory(cat metadata.Category, targetScope metadata.Scope, regKey moduleKey) error {
	mKey := moduleKey{category: cat, scope: targetScope}
	s := c.getModuleState(mKey)

	s.mu.Lock()
	if s.bound {
		s.mu.Unlock()
		return nil
	}
	defer s.mu.Unlock()

	c.regMu.RLock()
	extractor, ok := c.extractors[regKey]
	if !ok {
		if regKey.scope != metadata.GlobalScope {
			extractor, ok = c.extractors[moduleKey{category: cat, scope: metadata.GlobalScope}]
		}
	}
	c.regMu.RUnlock()

	if !ok {
		return fmt.Errorf("no extractor")
	}

	root := c.Config()
	if root == nil {
		return fmt.Errorf("engine: root config is nil")
	}

	mc, err := extractor(root)
	if err != nil || mc == nil {
		return err
	}

	for _, entry := range mc.Entries {
		if _, exists := s.instances[entry.Name]; !exists {
			s.instances[entry.Name] = &componentMeta{config: entry.Value, status: StatusNone}
		}
		s.order = append(s.order, entry.Name)
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

func (h *scopedHandle) In(cat metadata.Category, opts ...component.RegisterOption) component.Handle {
	o := &regOpts{scope: h.scope}
	for _, opt := range opts {
		opt(o)
	}
	return &scopedHandle{
		container: h.container,
		scope:     o.scope,
		category:  cat,
	}
}

func (h *scopedHandle) Get(ctx context.Context, name string) (any, error) {
	if h.category == "" {
		return nil, fmt.Errorf("engine: category context is required for Get()")
	}
	return h.container.getInternal(ctx, h.category, h.scope, name)
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
