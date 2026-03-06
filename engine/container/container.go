package container

import (
	"context"
	"fmt"
	"iter"
	"sync"

	"github.com/origadmin/runtime/contracts/component"
)

type Status int

const (
	StatusNone Status = iota
	StatusResolving
	StatusResolved
	StatusInstantiating
	StatusReady
	StatusError
)

type moduleKey struct {
	category component.Category
	scope    component.Scope
}

type componentMeta struct {
	config any
	status Status
	inst   any
	err    error
	tags   []string // Identity: Tags of the provider that created this instance
}

type moduleState struct {
	mu          sync.RWMutex
	instances   map[string]*componentMeta
	order       []string
	defaultName string
	bound       bool
}

type providerEntry struct {
	provider component.Provider
	resolver component.Resolver
	scopes   []component.Scope
	priority component.Priority
	tags     []string
}

// isProviderVisible checks if a provider with the given tags is visible in the current requested perspective.
func isProviderVisible(providerTags, requestedTags []string) bool {
	// 1. Full perspective (no tags requested) can see everything
	if len(requestedTags) == 0 {
		return true
	}
	// 2. Common providers (no tags) are visible in any perspective
	if len(providerTags) == 0 {
		return true
	}
	// 3. Tagged providers are visible if their tags are a subset of requested tags
	for _, pt := range providerTags {
		found := false
		for _, rt := range requestedTags {
			if pt == rt {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// tagsEqual checks if two tag sets are identical (order-independent).
func tagsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		found := false
		for _, x := range b {
			if v == x {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Option defines the internal option for container initialization.
type Option func(*containerImpl)

// WithCategoryResolvers injects category-specific resolvers into the container.
func WithCategoryResolvers(res map[component.Category]component.Resolver) Option {
	return func(c *containerImpl) {
		if res == nil {
			return
		}
		for k, v := range res {
			c.categoryResolvers[k] = v
		}
	}
}

type containerImpl struct {
	mu                sync.RWMutex
	modules           map[moduleKey]*moduleState
	providers         map[component.Category][]*providerEntry
	categoryResolvers map[component.Category]component.Resolver
	isLoaded          bool
}

func (c *containerImpl) Register(cat component.Category, p component.Provider, opts ...component.RegisterOption) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isLoaded {
		panic(fmt.Sprintf("engine: cannot register category %s after Load() has been called", cat))
	}

	cfg := &component.RegistrationOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	entry := &providerEntry{
		provider: p,
		resolver: cfg.Resolver,
		scopes:   cfg.Scopes,
		priority: cfg.Priority,
		tags:     cfg.Tags,
	}

	entries := c.providers[cat]
	inserted := false
	for i, e := range entries {
		// Newest registration with same priority takes precedence (descending priority)
		if entry.priority >= e.priority {
			entries = append(entries[:i], append([]*providerEntry{entry}, entries[i:]...)...)
			inserted = true
			break
		}
	}
	if !inserted {
		entries = append(entries, entry)
	}
	c.providers[cat] = entries
}

func (c *containerImpl) Has(cat component.Category, opts ...component.RegisterOption) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.providers[cat]
	return ok
}

func (c *containerImpl) Load(ctx context.Context, source any, opts ...component.LoadOption) error {
	c.mu.Lock()
	c.isLoaded = true
	c.mu.Unlock()

	loadOpts := &component.LoadOptions{}
	for _, opt := range opts {
		opt(loadOpts)
	}

	c.mu.RLock()
	var cats []component.Category
	if loadOpts.Category != "" {
		if _, ok := c.providers[loadOpts.Category]; ok {
			cats = append(cats, loadOpts.Category)
		}
	} else {
		for cat := range c.providers {
			cats = append(cats, cat)
		}
	}
	c.mu.RUnlock()

	for _, cat := range cats {
		entries := c.getProviderEntries(cat)
		if len(entries) == 0 {
			continue
		}

		primaryEntry := entries[0]
		registeredScopes := make(map[component.Scope]bool)
		for _, entry := range entries {
			if len(entry.scopes) == 0 {
				registeredScopes[component.GlobalScope] = true
			} else {
				for _, s := range entry.scopes {
					registeredScopes[s] = true
				}
			}
		}

		for s := range registeredScopes {
			if loadOpts.Scope != "" && s != loadOpts.Scope {
				continue
			}
			if err := c.bindWithSource(cat, s, primaryEntry, source, loadOpts.Resolver, loadOpts.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *containerImpl) getProviderEntries(cat component.Category) []*providerEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.providers[cat]
}

func (c *containerImpl) bindWithSource(cat component.Category, scope component.Scope, entry *providerEntry, source any, resolver component.Resolver, filterName string) error {
	mKey := moduleKey{category: cat, scope: scope}
	s := c.getModuleState(mKey)
	s.mu.Lock()
	defer s.mu.Unlock()

	var mc *component.ModuleConfig
	var err error

	if resolver != nil {
		mc, err = resolver(source, cat)
	} else if entry.resolver != nil {
		mc, err = entry.resolver(source, cat)
	} else {
		c.mu.RLock()
		r := c.categoryResolvers[cat]
		c.mu.RUnlock()
		if r != nil {
			mc, err = r(source, cat)
		}
	}

	if err == nil && mc == nil {
		mc = &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: component.DefaultName, Value: source}},
			Active:  component.DefaultName,
		}
	}

	if err != nil {
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
	} else if s.defaultName == "" {
		foundDefault := false
		for _, e := range mc.Entries {
			if e.Name == "default" {
				s.defaultName = "default"
				foundDefault = true
				break
			}
		}
		if !foundDefault && len(mc.Entries) == 1 {
			s.defaultName = mc.Entries[0].Name
		}
	}

	if s.defaultName != "" {
		if meta, ok := s.instances[s.defaultName]; ok {
			s.instances[component.DefaultName] = meta
		}
	}

	s.bound = true
	return nil
}

func (c *containerImpl) getModuleState(key moduleKey) *moduleState {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s, ok := c.modules[key]; ok {
		return s
	}
	s := &moduleState{
		instances: make(map[string]*componentMeta),
	}
	c.modules[key] = s
	return s
}

func (c *containerImpl) Get(ctx context.Context, name string) (any, error) {
	return c.instantiate(ctx, "", component.GlobalScope, name, nil)
}

func (c *containerImpl) Iter(ctx context.Context) iter.Seq2[string, any] {
	return c.iterInternal(ctx, "", component.GlobalScope, nil)
}

func (c *containerImpl) iterInternal(ctx context.Context, cat component.Category, scope component.Scope, tags []string) iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		mKey := moduleKey{category: cat, scope: scope}
		s := c.getModuleState(mKey)
		s.mu.RLock()
		order := make([]string, len(s.order))
		copy(order, s.order)
		s.mu.RUnlock()

		for _, name := range order {
			inst, err := c.instantiate(ctx, cat, scope, name, tags)
			if err == nil {
				if !yield(name, inst) {
					return
				}
			}
		}
	}
}

func (c *containerImpl) In(cat component.Category, opts ...component.InOption) component.Handle {
	inOpts := &component.InOptions{Scope: component.GlobalScope}
	for _, opt := range opts {
		opt(inOpts)
	}
	return &handleAdapter{c: c, category: cat, scope: inOpts.Scope, tags: inOpts.Tags}
}

func (c *containerImpl) Config() any                  { return nil }
func (c *containerImpl) Name() string                 { return "" }
func (c *containerImpl) Scope() component.Scope       { return component.GlobalScope }
func (c *containerImpl) Category() component.Category { return "" }

func (c *containerImpl) instantiate(ctx context.Context, cat component.Category, scope component.Scope, name string, requestedTags []string) (any, error) {
	if name == "" {
		name = component.DefaultName
	}

	mKey := moduleKey{category: cat, scope: scope}

	c.mu.RLock()
	s, exists := c.modules[mKey]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("engine: scope %s not initialized for category %s", scope, cat)
	}

	s.mu.RLock()
	meta, ok := s.instances[name]
	if !ok {
		s.mu.RUnlock()
		return nil, fmt.Errorf("engine: component %s/%s not found in scope %s", cat, name, scope)
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if meta.status == StatusInstantiating {
		return nil, fmt.Errorf("engine: circular dependency %s/%s", cat, name)
	}

	c.mu.RLock()
	entries := c.providers[cat]
	c.mu.RUnlock()
	if len(entries) == 0 {
		return nil, fmt.Errorf("engine: no provider for %s", cat)
	}

	var lastErr error
	for _, entry := range entries {
		// 1. Match Scope
		scopeMatch := false
		if len(entry.scopes) == 0 {
			scopeMatch = true
		} else {
			for _, s := range entry.scopes {
				if s == scope {
					scopeMatch = true
					break
				}
			}
		}
		if !scopeMatch {
			continue
		}

		// 2. Match Perspective Visibility
		// Can this perspective see this provider?
		if !isProviderVisible(entry.tags, requestedTags) {
			continue
		}

		// 3. Handle Ready Instance with Strict Compatibility
		if meta.status == StatusReady {
			// An instance is ONLY claimable if:
			// a) We are in Full Perspective (requestedTags is empty) - we see everything
			// b) The instance was created by a provider with exactly the same identity (tags)
			if len(requestedTags) == 0 || tagsEqual(meta.tags, entry.tags) {
				return meta.inst, nil
			}
			// Otherwise, this provider cannot "see" or "claim" this instance.
			continue
		}

		// 4. Try Instantiate
		meta.status = StatusInstantiating
		h := &handleAdapter{c: c, category: cat, scope: scope, name: name, meta: meta, tags: requestedTags}
		inst, err := entry.provider(ctx, h)
		if err == nil && inst != nil {
			meta.inst = inst
			meta.status = StatusReady
			meta.tags = entry.tags // Record creator's identity
			return inst, nil
		}

		meta.status = StatusNone
		if err != nil {
			lastErr = err
		}
	}

	meta.status = StatusError
	if lastErr != nil {
		meta.err = lastErr
		return nil, lastErr
	}
	return nil, fmt.Errorf("engine: no compatible provider found for %s/%s in scope %s with requested tags %v", cat, name, scope, requestedTags)
}

type handleAdapter struct {
	c        *containerImpl
	category component.Category
	scope    component.Scope
	name     string
	meta     *componentMeta
	tags     []string
}

func (h *handleAdapter) Get(ctx context.Context, name string) (any, error) {
	return h.c.instantiate(ctx, h.category, h.scope, name, h.tags)
}

func (h *handleAdapter) Iter(ctx context.Context) iter.Seq2[string, any] {
	return h.c.iterInternal(ctx, h.category, h.scope, h.tags)
}

func (h *handleAdapter) In(cat component.Category, opts ...component.InOption) component.Handle {
	inOpts := &component.InOptions{
		Scope: component.GlobalScope, // Perspective switch defaults to GlobalScope
		Tags:  h.tags,                // Inherit tags by default
	}
	for _, opt := range opts {
		opt(inOpts)
	}
	return &handleAdapter{c: h.c, category: cat, scope: inOpts.Scope, tags: inOpts.Tags}
}

func (h *handleAdapter) Config() any {
	if h.meta == nil {
		return nil
	}
	return h.meta.config
}

func (h *handleAdapter) Name() string {
	return h.name
}

func (h *handleAdapter) Scope() component.Scope       { return h.scope }
func (h *handleAdapter) Category() component.Category { return h.category }

func NewContainer(opts ...Option) component.Registry {
	c := &containerImpl{
		modules:           make(map[moduleKey]*moduleState),
		providers:         make(map[component.Category][]*providerEntry),
		categoryResolvers: make(map[component.Category]component.Resolver),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}
