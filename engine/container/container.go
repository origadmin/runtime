package container

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/origadmin/runtime/contracts/options"
	enginecontext "github.com/origadmin/runtime/engine/context"
	"github.com/origadmin/runtime/engine/protocol"
)

type Status int

const (
	StatusNone Status = iota
	StatusResolving
	StatusReady
)

type Handle interface {
	Get(ctx context.Context, name string) (any, error)
	In(category enginecontext.Category, opts ...RegisterOption) Handle
	BindConfig(target any) error
	Config() any
	Scope() enginecontext.Scope
	Category() enginecontext.Category
}

type Provider func(ctx context.Context, h Handle, opts ...options.Option) (any, error)

type Registry interface {
	Handle
	Register(c enginecontext.Category, e protocol.Extractor, p Provider, opts ...RegisterOption)
	BindRoot(root any)
	Init(ctx context.Context) error
}

type RegisterOption func(*regOpts)

type regOpts struct {
	scope    enginecontext.Scope
	priority int
}

func WithScope(s enginecontext.Scope) RegisterOption {
	return func(o *regOpts) {
		o.scope = s
	}
}

func WithPriority(p int) RegisterOption {
	return func(o *regOpts) {
		o.priority = p
	}
}

type instanceKey struct {
	category enginecontext.Category
	scope    enginecontext.Scope
	name     string
}

type moduleKey struct {
	category enginecontext.Category
	scope    enginecontext.Scope
}

type componentMeta struct {
	mu       sync.Mutex
	config   any
	status   Status
	instance any
}

type containerImpl struct {
	mu         sync.RWMutex
	rootConfig any
	providers  map[moduleKey]Provider
	extractors map[moduleKey]protocol.Extractor
	priorities map[moduleKey]int
	pool       map[instanceKey]*componentMeta
	poolOrder  map[moduleKey][]string
	defaults   map[moduleKey]string
}

func NewContainer(root any) Registry {
	return &containerImpl{
		rootConfig: root,
		providers:  make(map[moduleKey]Provider),
		extractors: make(map[moduleKey]protocol.Extractor),
		priorities: make(map[moduleKey]int),
		pool:       make(map[instanceKey]*componentMeta),
		poolOrder:  make(map[moduleKey][]string),
		defaults:   make(map[moduleKey]string),
	}
}

func (c *containerImpl) BindRoot(root any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rootConfig = root
}

func (c *containerImpl) Config() any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rootConfig
}

func (c *containerImpl) Scope() enginecontext.Scope       { return enginecontext.GlobalScope }
func (c *containerImpl) Category() enginecontext.Category { return "" }

func (c *containerImpl) In(cat enginecontext.Category, opts ...RegisterOption) Handle {
	o := &regOpts{scope: enginecontext.GlobalScope}
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

func (c *containerImpl) Register(cat enginecontext.Category, e protocol.Extractor, p Provider, opts ...RegisterOption) {
	o := &regOpts{
		scope:    enginecontext.GlobalScope,
		priority: enginecontext.PriorityInfrastructure,
	}
	for _, opt := range opts {
		opt(o)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	mKey := moduleKey{category: cat, scope: o.scope}
	c.providers[mKey] = p
	c.extractors[mKey] = e
	c.priorities[mKey] = o.priority
}

func (c *containerImpl) Init(ctx context.Context) error {
	c.mu.RLock()
	var keys []moduleKey
	for k := range c.providers {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return c.priorities[keys[i]] < c.priorities[keys[j]]
	})
	c.mu.RUnlock()

	for _, k := range keys {
		if err := c.bindCategory(k.category, k.scope, k); err != nil {
			return fmt.Errorf("engine: failed to bind %s in scope %s: %w", k.category, k.scope, err)
		}
		c.mu.RLock()
		order := c.poolOrder[k]
		c.mu.RUnlock()
		for _, name := range order {
			if _, err := c.getInternal(ctx, k.category, k.scope, name); err != nil {
				return fmt.Errorf("engine: failed to init %s/%s: %w", k.category, name, err)
			}
		}
	}
	return nil
}

// Internal Logic

func (c *containerImpl) getInternal(ctx context.Context, cat enginecontext.Category, scope enginecontext.Scope, name string) (any, error) {
	c.mu.RLock()
	mKey := moduleKey{category: cat, scope: scope}
	actualName := name
	if actualName == "" {
		actualName = c.defaults[mKey]
	}
	key := instanceKey{category: cat, scope: scope, name: actualName}
	meta, exists := c.pool[key]
	c.mu.RUnlock()

	if !exists || (actualName == "" && name == "") {
		if err := c.bindCategory(cat, scope, mKey); err != nil {
			if scope != enginecontext.GlobalScope {
				return c.getInternal(ctx, cat, enginecontext.GlobalScope, name)
			}
			return nil, err
		}
		c.mu.RLock()
		if actualName == "" {
			actualName = c.defaults[mKey]
		}
		key = instanceKey{category: cat, scope: scope, name: actualName}
		meta, exists = c.pool[key]
		c.mu.RUnlock()
		if !exists {
			if scope != enginecontext.GlobalScope {
				return c.getInternal(ctx, cat, enginecontext.GlobalScope, name)
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
	c.mu.RLock()
	provider, ok := c.providers[mKey]
	if !ok {
		provider, ok = c.providers[moduleKey{category: cat, scope: enginecontext.GlobalScope}]
	}
	c.mu.RUnlock()

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

func (c *containerImpl) bindCategory(cat enginecontext.Category, targetScope enginecontext.Scope, regKey moduleKey) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	mKey := moduleKey{category: cat, scope: targetScope}
	extractor, ok := c.extractors[regKey]
	if !ok {
		if regKey.scope != enginecontext.GlobalScope {
			extractor, ok = c.extractors[moduleKey{category: cat, scope: enginecontext.GlobalScope}]
		}
	}
	if !ok {
		return fmt.Errorf("engine: no extractor registered for %s", cat)
	}

	if c.rootConfig == nil {
		return fmt.Errorf("engine: root config is nil, cannot bind %s", cat)
	}

	mc, err := extractor(c.rootConfig)
	if err != nil || mc == nil {
		return err
	}

	var order []string
	for _, entry := range mc.Entries {
		key := instanceKey{category: cat, scope: targetScope, name: entry.Name}
		if _, exists := c.pool[key]; !exists {
			c.pool[key] = &componentMeta{config: entry.Value, status: StatusNone}
		}
		order = append(order, entry.Name)
	}
	c.poolOrder[mKey] = order

	if mc.Active != "" {
		c.defaults[mKey] = mc.Active
	} else if len(mc.Entries) > 0 && c.defaults[mKey] == "" {
		c.defaults[mKey] = mc.Entries[0].Name
	}

	return nil
}

type scopedHandle struct {
	container *containerImpl
	scope     enginecontext.Scope
	category  enginecontext.Category
	name      string
	config    any
}

func (h *scopedHandle) Scope() enginecontext.Scope       { return h.scope }
func (h *scopedHandle) Category() enginecontext.Category { return h.category }
func (h *scopedHandle) Config() any                      { return h.config }

func (h *scopedHandle) In(cat enginecontext.Category, opts ...RegisterOption) Handle {
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
