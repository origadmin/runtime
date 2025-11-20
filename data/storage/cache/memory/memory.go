/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package memory

import (
	"context"
	"sync"
	"time"

	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"
)

const (
	DriverName             = "memory" // Name of the memory cache driver
	DefaultSize            = 1024
	DefaultCapacity        = 1024
	DefaultExpiration      = 0
	DefaultCleanupInterval = 300
)

// Item represents an item stored in the cache.
type Item struct {
	value      string
	expiryTime time.Time
}

// Cache is an in-memory cache implementation.
type Cache struct {
	items           map[string]Item
	mu              sync.RWMutex
	size            int           // Maximum number of cache entries
	defaultExpiry   time.Duration // Default expiration time for cache items
	cleanupInterval time.Duration // Interval for cleaning up expired items
	stopCleanup     chan struct{} // Channel to stop the background cleanup goroutine
	closeOnce       sync.Once     // Ensures stopCleanup is closed only once
}

var (
	ErrKeyNotFound           = errors.New("key not found")
	ErrCacheSizeLimitReached = errors.New("cache size limit reached")
)

// Get retrieves the value associated with the given key.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return "", ErrKeyNotFound
	}

	// Check if the item has expired
	if !item.expiryTime.IsZero() && time.Now().After(item.expiryTime) {
		return "", ErrKeyNotFound // Item expired, let background cleanup handle deletion
	}

	return item.value, nil
}

// GetAndDelete retrieves the value associated with the given key and deletes it.
func (c *Cache) GetAndDelete(ctx context.Context, key string) (string, error) {
	value, err := c.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if err := c.Delete(ctx, key); err != nil {
		return "", err
	}

	return value, nil
}

// Exists checks if a value exists for the given key.
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	_, found := c.items[key]
	c.mu.RUnlock()

	return found, nil
}

// Set sets the value for the given key.
func (c *Cache) Set(ctx context.Context, key string, value string, exp ...time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Enforce size limit
	if c.size > 0 && len(c.items) >= c.size {
		return ErrCacheSizeLimitReached
	}

	var ttl time.Duration
	if len(exp) > 0 {
		ttl = exp[0]
	} else {
		ttl = c.defaultExpiry
	}

	var expiryTime time.Time
	if ttl > 0 {
		expiryTime = time.Now().Add(ttl)
	}

	c.items[key] = Item{
		value:      value,
		expiryTime: expiryTime,
	}

	return nil
}

// Delete deletes the value associated with the given key.
func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Deleting a non-existent key is not an error.
	// This makes the Delete operation idempotent.
	delete(c.items, key)
	return nil
}

// Clear clears the cache.
func (c *Cache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
	return nil
}

// Close closes the cache.
func (c *Cache) Close(ctx context.Context) error {
	c.closeOnce.Do(func() {
		close(c.stopCleanup)
	})
	return nil
}

// startCleanup starts a background goroutine to periodically clean up expired items.
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, item := range c.items {
				if !item.expiryTime.IsZero() && time.Now().After(item.expiryTime) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		case <-c.stopCleanup:
			return
		}
	}
}

func New(cfg *cachev1.CacheConfig, _ ...options.Option) (storageiface.Cache, error) {
	if cfg == nil || cfg.GetDriver() != DriverName {
		return nil, errors.New("invalid cache config")
	}
	memoryCfg := cfg.GetMemory()
	if memoryCfg == nil {
		// If no specific memory config is provided, create a default one to proceed.
		memoryCfg = &cachev1.MemoryConfig{}
	}

	// Use local variables for configuration, applying defaults without modifying the input cfg.
	size := memoryCfg.GetSize()
	if size == 0 {
		size = DefaultSize
	}

	capacity := memoryCfg.GetCapacity()
	if capacity == 0 {
		capacity = DefaultCapacity
	}

	expiration := memoryCfg.GetExpiration()
	if expiration == 0 {
		expiration = DefaultExpiration
	}

	cleanupInterval := memoryCfg.GetCleanupInterval()
	if cleanupInterval == 0 {
		// If cleanup interval is 0, disable background cleanup
		cleanupInterval = -1 // Use a negative value to indicate no cleanup
	}

	cache := &Cache{
		items:           make(map[string]Item, capacity),
		size:            int(size),
		defaultExpiry:   time.Duration(expiration) * time.Millisecond,
		cleanupInterval: time.Duration(cleanupInterval) * time.Millisecond,
		stopCleanup:     make(chan struct{}),
	}

	// Start background cleanup goroutine if cleanupInterval is set and positive
	if cache.cleanupInterval > 0 {
		go cache.startCleanup()
	}

	return cache, nil
}
