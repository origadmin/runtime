package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	cachev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/cache/v1"
)

// CacheItem represents an item stored in the cache.
type CacheItem struct {
	value      string
	expiryTime time.Time
}

// Cache is an in-memory cache implementation.
type Cache struct {
	items           map[string]CacheItem
	mu              sync.RWMutex
	size            int           // Maximum number of cache entries
	defaultExpiry   time.Duration // Default expiration time for cache items
	cleanupInterval time.Duration // Interval for cleaning up expired items
	stopCleanup     chan struct{} // Channel to stop the background cleanup goroutine
}

// NewCache creates a new instance of Cache with the given configuration.
func NewCache(config *cachev1.MemoryConfig) *Cache {
	cache := &Cache{
		items:           make(map[string]CacheItem),
		size:            int(config.Size),
		defaultExpiry:   time.Duration(config.Expiration) * time.Millisecond,
		cleanupInterval: time.Duration(config.CleanupInterval) * time.Millisecond,
		stopCleanup:     make(chan struct{}),
	}

	// Start background cleanup goroutine if cleanupInterval is set
	if cache.cleanupInterval > 0 {
		go cache.startCleanup()
	}

	return cache
}

// Get retrieves the value associated with the given key.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return "", errors.New("key not found")
	}

	// Check if the item has expired
	if !item.expiryTime.IsZero() && time.Now().After(item.expiryTime) {
		c.Delete(ctx, key)
		return "", errors.New("key expired")
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
		return errors.New("cache size limit reached")
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

	c.items[key] = CacheItem{
		value:      value,
		expiryTime: expiryTime,
	}

	return nil
}

// Delete deletes the value associated with the given key.
func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)
	return nil
}

// Clear clears the cache.
func (c *Cache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
	return nil
}

// Close closes the cache.
func (c *Cache) Close(ctx context.Context) error {
	close(c.stopCleanup)
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
