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
}

type moduleState struct {
	mu          sync.RWMutex
	instances   map[string]*componentMeta // Key: name + ":" + tag
	order       []string
	defaultName string
	bound       bool
}

// makeInstanceKey defines the PHYSICAL COORDINATE of an instance.
// name + tag = Identity.
func makeInstanceKey(name, tag string) string {
	if tag == "" {
		return name
	}
	return name + ":" + tag
}

func configKey(name string) string {
	return name + "@_config"
}

type providerEntry struct {
	provider       component.Provider
	resolver       component.Resolver
	scopes         []component.Scope
	priority       component.Priority
	tag            string // Producer's identity tag
	defaultEntries []string
}

// isProviderCompatible checks if a provider can perform work for the requested tag.
func isProviderCompatible(providerTag, requestedTag string) bool {
	if providerTag == "" {
		return true // Common provider serves all tags
	}
	return providerTag == requestedTag // Specific provider serves only its tag
}

type Option func(*containerImpl)

func WithCategoryResolvers(res map[component.Category]component.Resolver) Option {
	return func(c *containerImpl) {
		if res != nil {
			for k, v := range res {
				c.categoryResolvers[k] = v
			}
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
		panic(fmt.Sprintf("engine: cannot register category %s after Load()", cat))
	}
	cfg := &component.RegistrationOptions{}
	for _, opt := range opts {
		opt(cfg)
	}
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []component.Scope{component.GlobalScope}
	}
	entry := &providerEntry{
		provider:       p,
		resolver:       cfg.Resolver,
		scopes:         scopes,
		priority:       cfg.Priority,
		tag:            cfg.Tag,
		defaultEntries: cfg.DefaultEntries,
	}
	entries := c.providers[cat]
	inserted := false
	for i, e := range entries {
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

func (c *containerImpl) Inject(cat component.Category, name string, inst any, opts ...component.RegisterOption) {
	cfg := &component.RegistrationOptions{}
	for _, opt := range opts {
		opt(cfg)
	}
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []component.Scope{component.GlobalScope}
	}
	for _, s := range scopes {
		mKey := moduleKey{category: cat, scope: s}
		state := c.getModuleState(mKey)
		state.mu.Lock()
		finalName := name
		isExplicitDefault := false
		if finalName == "" || finalName == component.DefaultName {
			finalName = "_injected_" + string(cat)
			isExplicitDefault = true
		}

		// Injected instance identity is determined by its name and its own tag
		key := makeInstanceKey(finalName, cfg.Tag)
		state.instances[key] = &componentMeta{inst: inst, status: StatusReady}

		if !contains(state.order, finalName) {
			state.order = append(state.order, finalName)
		}
		if isExplicitDefault || state.defaultName == "" {
			state.defaultName = finalName
		}

		// Update _default marker
		if state.defaultName != "" {
			if dMeta, ok := state.instances[makeInstanceKey(state.defaultName, cfg.Tag)]; ok {
				state.instances[makeInstanceKey(component.DefaultName, "")] = dMeta
			}
		}
		state.mu.Unlock()
	}
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
		// SEEDING
		for _, entry := range entries {
			if len(entry.defaultEntries) > 0 {
				for _, name := range entry.defaultEntries {
					for s := range registeredScopes {
						if loadOpts.Scope != "" && s != loadOpts.Scope {
							continue
						}
						mKey := moduleKey{category: cat, scope: s}
						state := c.getModuleState(mKey)
						state.mu.Lock()
						if _, exists := state.instances[configKey(name)]; !exists {
							state.instances[configKey(name)] = &componentMeta{config: nil, status: StatusNone}
							state.order = append(state.order, name)
						}
						if state.defaultName == "" {
							state.defaultName = name
							state.instances[configKey(component.DefaultName)] = state.instances[configKey(name)]
						}
						state.mu.Unlock()
					}
				}
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
		name := string(cat)
		mc = &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: name, Value: source}}, Active: name}
	}
	if err != nil {
		return err
	}
	for _, cfgEntry := range mc.Entries {
		if filterName != "" && cfgEntry.Name != filterName {
			continue
		}
		key := configKey(cfgEntry.Name)
		if _, exists := s.instances[key]; !exists {
			s.instances[key] = &componentMeta{config: cfgEntry.Value, status: StatusNone}
			if !contains(s.order, cfgEntry.Name) {
				s.order = append(s.order, cfgEntry.Name)
			}
		}
	}
	newDefault := ""
	for _, e := range mc.Entries {
		if e.Name == "default" {
			newDefault = e.Name
			break
		}
	}
	if newDefault == "" && mc.Active != "" {
		newDefault = mc.Active
	}
	if newDefault == "" && len(mc.Entries) == 1 {
		newDefault = mc.Entries[0].Name
	}
	if newDefault != "" && (filterName == "" || newDefault == filterName) {
		s.defaultName = newDefault
	}
	if s.defaultName != "" {
		if meta, ok := s.instances[configKey(s.defaultName)]; ok {
			s.instances[configKey(component.DefaultName)] = meta
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
	s := &moduleState{instances: make(map[string]*componentMeta)}
	c.modules[key] = s
	return s
}

func (c *containerImpl) In(cat component.Category, opts ...component.InOption) component.Locator {
	var res component.Locator = &locatorHandle{c: c, category: cat, scope: component.GlobalScope}
	for _, opt := range opts {
		if opt != nil {
			res = opt(res)
		}
	}
	return res
}

func (c *containerImpl) Config() any                  { return nil }
func (c *containerImpl) Name() string                 { return "" }
func (c *containerImpl) Scope() component.Scope       { return component.GlobalScope }
func (c *containerImpl) Category() component.Category { return "" }
func (c *containerImpl) Tags() []string               { return nil }
func (c *containerImpl) Tag() string                  { return "" }

func (c *containerImpl) instantiate(ctx context.Context, cat component.Category, scope component.Scope, name string, demandTags []string) (any, error) {
	reqName := name
	if reqName == "" {
		reqName = component.DefaultName
	}
	mKey := moduleKey{category: cat, scope: scope}
	c.mu.RLock()
	s, exists := c.modules[mKey]
	c.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("engine: scope %s not initialized for %s", scope, cat)
	}

	realName := reqName
	if reqName == component.DefaultName {
		realName = s.defaultName
	}

	// 1. Check for directly ready instance (Inject or agnostic cached)
	s.mu.RLock()
	if meta, ok := s.instances[makeInstanceKey(realName, "")]; ok && meta.status == StatusReady {
		s.mu.RUnlock()
		return meta.inst, nil
	}
	if reqName == component.DefaultName {
		if meta, ok := s.instances[makeInstanceKey(component.DefaultName, "")]; ok && meta.status == StatusReady {
			s.mu.RUnlock()
			return meta.inst, nil
		}
	}
	cfgMeta, ok := s.instances[configKey(reqName)]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("engine: component %s/%s not found", cat, reqName)
	}
	if cfgMeta.status == StatusReady {
		return cfgMeta.inst, nil
	}

	// 2. RETRIEVAL (Package): Find results for any tag in the package
	tagsToTry := demandTags
	if len(tagsToTry) == 0 {
		tagsToTry = []string{""}
	}

	c.mu.RLock()
	entries := c.providers[cat]
	c.mu.RUnlock()

	var lastErr error
	for _, curTag := range tagsToTry {
		for _, entry := range entries {
			if !matchScope(entry.scopes, scope) {
				continue
			}
			if !isProviderCompatible(entry.tag, curTag) {
				continue
			}

			// 3. WORK (Creation): Create the instance for this specific identity
			iKey := makeInstanceKey(realName, curTag)
			s.mu.Lock()
			meta, exists := s.instances[iKey]
			if !exists {
				meta = &componentMeta{config: cfgMeta.config, status: StatusNone}
				s.instances[iKey] = meta
			}
			if meta.status == StatusReady {
				inst := meta.inst
				s.mu.Unlock()
				return inst, nil
			}
			if meta.status == StatusInstantiating {
				s.mu.Unlock()
				return nil, fmt.Errorf("engine: circular dependency %s", iKey)
			}
			meta.status = StatusInstantiating
			s.mu.Unlock()

			// CALLBACK: Single Work Tag
			h := &entryHandle{
				name:      realName,
				meta:      meta,
				activeTag: curTag,
				l:         &locatorHandle{c: c, category: cat, scope: scope, tags: demandTags},
			}
			inst, err := entry.provider(ctx, h)
			s.mu.Lock()
			if err == nil && inst != nil {
				meta.inst = inst
				meta.status = StatusReady
				s.mu.Unlock()
				return inst, nil
			}
			meta.status = StatusNone
			if err != nil {
				lastErr = err
			}
			s.mu.Unlock()
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("engine: no provider found for %s/%s with tags %v", cat, realName, demandTags)
}

func matchScope(ss []component.Scope, t component.Scope) bool {
	if len(ss) == 0 {
		return true
	}
	for _, s := range ss {
		if s == t {
			return true
		}
	}
	return false
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func (c *containerImpl) iter(ctx context.Context, cat component.Category, scope component.Scope, tags []string) iter.Seq2[string, any] {
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

// containerReader defines the internal read-only operations of the container.
// It ensures that locators can only trigger instantiation and iteration, 
// without having access to registration or loading methods.
type containerReader interface {
	instantiate(ctx context.Context, cat component.Category, scope component.Scope, name string, tags []string) (any, error)
	iter(ctx context.Context, cat component.Category, scope component.Scope, tags []string) iter.Seq2[string, any]
}

type locatorHandle struct {
	c        containerReader
	category component.Category
	scope    component.Scope
	tags     []string
}

func (l *locatorHandle) Get(ctx context.Context, name string) (any, error) {
	return l.c.instantiate(ctx, l.category, l.scope, name, l.tags)
}
func (l *locatorHandle) Iter(ctx context.Context) iter.Seq2[string, any] {
	return l.c.iter(ctx, l.category, l.scope, l.tags)
}
func (l *locatorHandle) In(cat component.Category, opts ...component.InOption) component.Locator {
	var res component.Locator = &locatorHandle{c: l.c, category: cat, scope: l.scope, tags: l.tags}
	for _, opt := range opts {
		if opt != nil {
			res = opt(res)
		}
	}
	return res
}
func (l *locatorHandle) WithInScope(s component.Scope) component.Locator {
	return &locatorHandle{c: l.c, category: l.category, scope: s, tags: l.tags}
}
func (l *locatorHandle) WithInTags(tags ...string) component.Locator {
	return &locatorHandle{c: l.c, category: l.category, scope: l.scope, tags: tags}
}
func (l *locatorHandle) Scope() component.Scope       { return l.scope }
func (l *locatorHandle) Category() component.Category { return l.category }
func (l *locatorHandle) Tags() []string               { return l.tags }
func (l *locatorHandle) Tag() string {
	if len(l.tags) > 0 {
		return l.tags[0]
	}
	return ""
}

type entryHandle struct {
	name      string
	meta      *componentMeta
	activeTag string
	l         component.Locator
}

func (e *entryHandle) Name() string { return e.name }
func (e *entryHandle) Config() any {
	if e.meta == nil {
		return nil
	}
	return e.meta.config
}
func (e *entryHandle) Locator() component.Locator { return e.l }
func (e *entryHandle) Tag() string                { return e.activeTag }

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
