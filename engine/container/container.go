package container

import (
	"context"
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/origadmin/runtime/contracts/component"
)

type Status int

const (
	StatusNone Status = iota
	StatusInstantiating
	StatusReady
	StatusError

	defaultInstanceName = "_default"
	globalScopeName     = "_global"
)

type moduleKey struct {
	category component.Category
	scope    component.Scope
}

type componentMeta struct {
	config              any
	requirementResolver component.RequirementResolver
	status              Status
	inst                any
	err                 error
}

type moduleState struct {
	mu          sync.RWMutex
	instances   map[string]*componentMeta // Key: name + ":" + tag
	order       []string
	defaultName string
	bound       bool
}

// containerBackend defines the internal operations of the container.
// It decouples the implementation from handles and locators.
type containerBackend interface {
	register(cat component.Category, p component.Provider, opts ...component.RegisterOption)
	inject(cat component.Category, name string, inst any, opts ...component.RegisterOption)
	isRegistered(cat component.Category, opts ...component.RegisterOption) bool
	requirement(cat component.Category, purpose string, res component.RequirementResolver)
	getCategoryRequirementResolver(cat component.Category, purpose string) component.RequirementResolver
	instantiate(ctx context.Context, cat component.Category, scope component.Scope, name string, tags []string) (any, error)
	iter(ctx context.Context, cat component.Category, scope component.Scope, tags []string, skips []string) iter.Seq2[string, any]
	scopes(cat component.Category) []component.Scope
}

// makeInstanceKey defines the PHYSICAL COORDINATE of an instance.
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
	provider            component.Provider
	resolver            component.ConfigResolver
	requirementResolver component.RequirementResolver
	scopes              []component.Scope
	priority            component.Priority
	tag                 string
	defaultEntries      []string
}

func isProviderCompatible(providerTag, requestedTag string) bool {
	if providerTag == "" {
		return true
	}
	return providerTag == requestedTag
}

type Option func(*containerImpl)

func WithCategoryResolvers(res map[component.Category]component.ConfigResolver) Option {
	return func(c *containerImpl) {
		if res != nil {
			for k, v := range res {
				c.categoryResolvers[k] = v
			}
		}
	}
}

type containerImpl struct {
	mu                           sync.RWMutex
	modules                      map[moduleKey]*moduleState
	providers                    map[component.Category][]*providerEntry
	categoryResolvers            map[component.Category]component.ConfigResolver
	categoryRequirementResolvers map[component.Category]map[string]component.RequirementResolver
	isLoaded                     bool
}

func (c *containerImpl) Register(cat component.Category, p component.Provider, opts ...component.RegisterOption) {
	c.register(cat, p, opts...)
}

func (c *containerImpl) register(cat component.Category, p component.Provider, opts ...component.RegisterOption) {
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
		scopes = []component.Scope{globalScopeName}
	} else {
		// Map empty scope to globalScopeName internal alias
		for i, s := range scopes {
			if s == "" {
				scopes[i] = globalScopeName
			}
		}
	}
	entry := &providerEntry{
		provider:            p,
		resolver:            cfg.ConfigResolver,
		requirementResolver: cfg.RequirementResolver,
		scopes:              scopes,
		priority:            cfg.Priority,
		tag:                 cfg.Tag,
		defaultEntries:      cfg.DefaultEntries,
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
	c.inject(cat, name, inst, opts...)
}

func (c *containerImpl) inject(cat component.Category, name string, inst any, opts ...component.RegisterOption) {
	cfg := &component.RegistrationOptions{}
	for _, opt := range opts {
		opt(cfg)
	}
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []component.Scope{globalScopeName}
	}
	for _, s := range scopes {
		internalScope := s
		if s == "" {
			internalScope = globalScopeName
		}
		mKey := moduleKey{category: cat, scope: internalScope}
		state := c.getModuleState(mKey)
		state.mu.Lock()
		finalName := name
		if finalName == "" {
			finalName = "_injected_" + string(cat)
		}
		key := makeInstanceKey(finalName, cfg.Tag)
		state.instances[key] = &componentMeta{inst: inst, status: StatusReady}
		if !contains(state.order, finalName) {
			state.order = append(state.order, finalName)
		}
		if state.defaultName == "" || name == "" {
			state.defaultName = finalName
		}
		if state.defaultName != "" {
			if dMeta, ok := state.instances[makeInstanceKey(state.defaultName, cfg.Tag)]; ok {
				state.instances[makeInstanceKey(defaultInstanceName, "")] = dMeta
			}
		}
		state.mu.Unlock()
	}
}

func (c *containerImpl) IsRegistered(cat component.Category, opts ...component.RegisterOption) bool {
	return c.isRegistered(cat, opts...)
}

func (c *containerImpl) isRegistered(cat component.Category, opts ...component.RegisterOption) bool {
	cfg := &component.RegistrationOptions{}
	for _, opt := range opts {
		opt(cfg)
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	entries, ok := c.providers[cat]
	if !ok {
		return false
	}
	if cfg.Tag == "" && len(cfg.Scopes) == 0 {
		return true
	}
	for _, e := range entries {
		if cfg.Tag != "" && e.tag != cfg.Tag {
			continue
		}
		if len(cfg.Scopes) > 0 {
			match := false
			for _, s := range cfg.Scopes {
				target := s
				if s == "" {
					target = globalScopeName
				}
				if matchScope(e.scopes, target) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		return true
	}
	return false
}

func (c *containerImpl) Requirement(cat component.Category, purpose string, res component.RequirementResolver) {
	c.requirement(cat, purpose, res)
}

func (c *containerImpl) requirement(cat component.Category, purpose string, res component.RequirementResolver) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.categoryRequirementResolvers == nil {
		c.categoryRequirementResolvers = make(map[component.Category]map[string]component.RequirementResolver)
	}
	if c.categoryRequirementResolvers[cat] == nil {
		c.categoryRequirementResolvers[cat] = make(map[string]component.RequirementResolver)
	}
	c.categoryRequirementResolvers[cat][purpose] = res
}

func (c *containerImpl) getCategoryRequirementResolver(cat component.Category, purpose string) component.RequirementResolver {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.categoryRequirementResolvers == nil {
		return nil
	}
	if resMap, ok := c.categoryRequirementResolvers[cat]; ok {
		return resMap[purpose]
	}
	return nil
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
				registeredScopes[globalScopeName] = true
			} else {
				for _, s := range entry.scopes {
					registeredScopes[s] = true
				}
			}
		}
		for s := range registeredScopes {
			// Filter by Scope if requested
			if loadOpts.Scope != "" {
				target := loadOpts.Scope
				if target == "" {
					target = globalScopeName
				}
				if s != target {
					continue
				}
			}
			// CLONE opts for this specific category and scope
			currentOpts := *loadOpts
			currentOpts.Category = cat
			currentOpts.Scope = s

			if err := c.bindWithSource(ctx, primaryEntry, source, &currentOpts); err != nil {
				return err
			}
		}
		for _, entry := range entries {
			if len(entry.defaultEntries) > 0 {
				for _, name := range entry.defaultEntries {
					for s := range registeredScopes {
						// Apply scope filtering for seeding too
						if loadOpts.Scope != "" {
							target := loadOpts.Scope
							if target == "" {
								target = globalScopeName
							}
							if s != target {
								continue
							}
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
							state.instances[configKey(defaultInstanceName)] = state.instances[configKey(name)]
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

func (c *containerImpl) bindWithSource(ctx context.Context, entry *providerEntry, source any, opts *component.LoadOptions) error {
	internalScope := opts.Scope
	if internalScope == "" {
		internalScope = globalScopeName
	}
	mKey := moduleKey{category: opts.Category, scope: internalScope}
	s := c.getModuleState(mKey)
	s.mu.Lock()
	defer s.mu.Unlock()

	var mc *component.ModuleConfig
	var err error
	// Priority: Load side > Registration side > Global default
	effectiveResolver := opts.Resolver
	if effectiveResolver == nil {
		effectiveResolver = entry.resolver
	}
	if effectiveResolver == nil {
		c.mu.RLock()
		effectiveResolver = c.categoryResolvers[opts.Category]
		c.mu.RUnlock()
	}

	if effectiveResolver != nil {
		mc, err = effectiveResolver(ctx, source, opts)
	}

	if err == nil && mc == nil {
		name := string(opts.Category)
		mc = &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: name, Value: source}}, Active: name}
	}
	if err != nil {
		return err
	}
	for _, cfgEntry := range mc.Entries {
		if opts.Name != "" && cfgEntry.Name != opts.Name {
			continue
		}
		key := configKey(cfgEntry.Name)
		if _, exists := s.instances[key]; !exists {
			// Resolve RequirementResolver: Entry > Module
			res := cfgEntry.RequirementResolver
			if res == nil {
				res = mc.RequirementResolver
			}
			s.instances[key] = &componentMeta{
				config:              cfgEntry.Value,
				requirementResolver: res,
				status:              StatusNone,
			}
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
	if newDefault != "" && (opts.Name == "" || newDefault == opts.Name) {
		s.defaultName = newDefault
	}
	if s.defaultName != "" {
		if meta, ok := s.instances[configKey(s.defaultName)]; ok {
			s.instances[configKey(defaultInstanceName)] = meta
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

func (c *containerImpl) scopes(cat component.Category) []component.Scope {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var res []component.Scope
	for k := range c.modules {
		if k.category == cat {
			// Map internal alias back to empty string for public API
			s := k.scope
			if s == globalScopeName {
				s = ""
			}
			res = append(res, s)
		}
	}
	return res
}

func (c *containerImpl) iter(ctx context.Context, cat component.Category, scope component.Scope, tags []string, skips []string) iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		internalScope := scope
		if scope == "" {
			internalScope = globalScopeName
		}
		mKey := moduleKey{category: cat, scope: internalScope}
		s := c.getModuleState(mKey)
		s.mu.RLock()
		order := make([]string, len(s.order))
		copy(order, s.order)
		s.mu.RUnlock()
		for _, name := range order {
			if contains(skips, name) {
				continue
			}
			inst, err := c.instantiate(ctx, cat, internalScope, name, tags)
			if err != nil {
				if isCircularDependencyError(err) {
					continue
				}
				continue
			}
			if !yield(name, inst) {
				return
			}
		}
	}
}

func parseInstanceName(name string) (string, []string) {
	parts := strings.Split(name, ":")
	if len(parts) <= 1 {
		return name, nil
	}
	return parts[0], parts[1:]
}

func (c *containerImpl) In(cat component.Category, opts ...component.InOption) component.Registry {
	var res component.Registry = &locatorHandle{c: c, category: cat, scope: ""}
	for _, opt := range opts {
		if opt != nil {
			res = opt(res)
		}
	}
	return res
}

func (c *containerImpl) instantiate(ctx context.Context, cat component.Category, scope component.Scope, name string, tags []string) (any, error) {
	reqName := name
	if reqName == "" {
		reqName = defaultInstanceName
	}
	internalScope := scope
	if scope == "" {
		internalScope = globalScopeName
	}
	mKey := moduleKey{category: cat, scope: internalScope}
	c.mu.RLock()
	s, exists := c.modules[mKey]
	c.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("engine: scope %s not initialized for %s", scope, cat)
	}
	realName := reqName
	if reqName == defaultInstanceName {
		realName = s.defaultName
	}
	s.mu.RLock()
	if meta, ok := s.instances[makeInstanceKey(realName, "")]; ok && meta.status == StatusReady {
		s.mu.RUnlock()
		return meta.inst, nil
	}
	if reqName == defaultInstanceName {
		if meta, ok := s.instances[makeInstanceKey(defaultInstanceName, "")]; ok && meta.status == StatusReady {
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
	tagsToTry := tags
	if len(tagsToTry) == 0 {
		tagsToTry = []string{""}
	}
	realReqName, demandTags := parseInstanceName(realName)
	if len(demandTags) > 0 {
		tagsToTry = demandTags
		realName = realReqName
	}
	c.mu.RLock()
	entries := c.providers[cat]
	c.mu.RUnlock()
	var lastErr error
	for _, curTag := range tagsToTry {
		for _, entry := range entries {
			if !matchScope(entry.scopes, internalScope) || !isProviderCompatible(entry.tag, curTag) || entry.provider == nil {
				continue
			}
			iKey := makeInstanceKey(realName, curTag)
			s.mu.Lock()
			meta, exists := s.instances[iKey]
			if !exists {
				res := cfgMeta.requirementResolver
				if res == nil {
					res = entry.requirementResolver
				}
				meta = &componentMeta{config: cfgMeta.config, requirementResolver: res, status: StatusNone}
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
			h := &entryHandle{
				category:  cat,
				scope:     scope,
				name:      realName,
				meta:      meta,
				activeTag: curTag,
				l:         (&locatorHandle{c: c, category: cat, scope: scope, tags: demandTags}).Skip(realName),
				c:         c,
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

type locatorHandle struct {
	c        containerBackend
	category component.Category
	scope    component.Scope
	tags     []string
	skips    []string
}

func (l *locatorHandle) Get(ctx context.Context, name string) (any, error) {
	return l.c.instantiate(ctx, l.category, l.scope, name, l.tags)
}
func (l *locatorHandle) Iter(ctx context.Context) iter.Seq2[string, any] {
	return l.c.iter(ctx, l.category, l.scope, l.tags, l.skips)
}
func (l *locatorHandle) In(cat component.Category, opts ...component.InOption) component.Registry {
	var res component.Registry = &locatorHandle{c: l.c, category: cat, scope: l.scope, tags: l.tags}
	for _, opt := range opts {
		if opt != nil {
			res = opt(res)
		}
	}
	return res
}
func (l *locatorHandle) WithInScope(s component.Scope) component.Locator {
	if l.scope == s {
		return l
	}
	newHandle := *l
	newHandle.scope = s
	return &newHandle
}
func (l *locatorHandle) WithInTags(tags ...string) component.Locator {
	var validTags []string
	for _, t := range tags {
		if t != "" {
			validTags = append(validTags, t)
		}
	}
	if equalTags(l.tags, validTags) {
		return l
	}
	newHandle := *l
	newHandle.tags = validTags
	return &newHandle
}
func (l *locatorHandle) Skip(names ...string) component.Locator {
	if len(names) == 0 {
		return l
	}
	newHandle := *l
	newHandle.skips = make([]string, len(l.skips)+len(names))
	copy(newHandle.skips, l.skips)
	copy(newHandle.skips[len(l.skips):], names)
	return &newHandle
}
func (l *locatorHandle) Category() component.Category { return l.category }
func (l *locatorHandle) Scope() component.Scope       { return l.scope }
func (l *locatorHandle) Scopes() []component.Scope    { return l.c.scopes(l.category) }
func (l *locatorHandle) Tags() []string               { return l.tags }

func (l *locatorHandle) Register(p component.Provider, opts ...component.RegisterOption) {
	l.c.register(l.category, p, opts...)
}
func (l *locatorHandle) Inject(name string, inst any, opts ...component.RegisterOption) {
	l.c.inject(l.category, name, inst, opts...)
}
func (l *locatorHandle) IsRegistered(opts ...component.RegisterOption) bool {
	return l.c.isRegistered(l.category, opts...)
}
func (l *locatorHandle) Requirement(purpose string, resolver component.RequirementResolver) {
	l.c.requirement(l.category, purpose, resolver)
}

func equalTags(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type entryHandle struct {
	category  component.Category
	scope     component.Scope
	name      string
	meta      *componentMeta
	activeTag string
	l         component.Locator
	c         containerBackend
}

func (e *entryHandle) Name() string                 { return e.name }
func (e *entryHandle) Category() component.Category { return e.category }
func (e *entryHandle) Scope() component.Scope       { return e.scope }
func (e *entryHandle) Config() any {
	if e.meta == nil {
		return nil
	}
	return e.meta.config
}
func (e *entryHandle) Locator() component.Locator { return e.l }
func (e *entryHandle) Tag() string                { return e.activeTag }
func (e *entryHandle) Require(purpose string) (any, error) {
	if e.meta != nil && e.meta.requirementResolver != nil {
		return e.meta.requirementResolver(context.Background(), e, purpose)
	}
	if res := e.c.getCategoryRequirementResolver(e.category, purpose); res != nil {
		return res(context.Background(), e, purpose)
	}
	return nil, fmt.Errorf("%w: %s (no resolver provided)", component.ErrRequirementNotFound, purpose)
}

func isCircularDependencyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "circular dependency")
}

func NewContainer(opts ...Option) component.Container {
	c := &containerImpl{
		modules:           make(map[moduleKey]*moduleState),
		providers:         make(map[component.Category][]*providerEntry),
		categoryResolvers: make(map[component.Category]component.ConfigResolver),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}
