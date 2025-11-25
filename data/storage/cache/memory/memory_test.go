package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
)

func TestNewMemoryCache(t *testing.T) {
	ctx := context.Background()

	// Test with default configuration
	cache, err := New(&cachev1.CacheConfig{
		Driver: DriverName, // Set the driver name
		Memory: &cachev1.MemoryConfig{},
	})
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	defer cache.Close(ctx)

	memCache, ok := cache.(*Cache)
	assert.True(t, ok)
	assert.Equal(t, DefaultSize, memCache.size)
	assert.Equal(t, time.Duration(0)*time.Millisecond, memCache.defaultExpiry, "Default expiration should be 0ms (no expiration)") // Corrected assertion
	assert.Equal(t, time.Duration(-1)*time.Millisecond, memCache.cleanupInterval, "Default cleanup interval should be -1ms (disabled)") // Corrected assertion

	// Test with custom configuration
	customSize := int32(50)
	customCapacity := int32(100)
	customExpiration := int32(1000)     // 1 second
	customCleanupInterval := int32(500) // 0.5 second

	cache, err = New(&cachev1.CacheConfig{
		Driver: DriverName, // Set the driver name
		Memory: &cachev1.MemoryConfig{
			Size:            customSize,
			Capacity:        customCapacity,
			Expiration:      int64(customExpiration),
			CleanupInterval: int64(customCleanupInterval),
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	defer cache.Close(ctx)

	memCache, ok = cache.(*Cache)
	assert.True(t, ok)
	assert.Equal(t, int(customSize), memCache.size)
	assert.Equal(t, time.Duration(customExpiration)*time.Millisecond, memCache.defaultExpiry)
	assert.Equal(t, time.Duration(customCleanupInterval)*time.Millisecond, memCache.cleanupInterval)
}

func TestMemoryCache_SetAndGet(t *testing.T) {
	ctx := context.Background()
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{},
	})
	defer cache.Close(ctx)

	// Test Set and Get
	err := cache.Set(ctx, "key1", "value1")
	assert.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// Test Get non-existent key
	_, err = cache.Get(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	// Test update value
	err = cache.Set(ctx, "key1", "newValue")
	assert.NoError(t, err)
	val, err = cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "newValue", val)
}

func TestMemoryCache_Expiration(t *testing.T) {
	ctx := context.Background()
	// Set default expiration to 100ms
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{
			Expiration:      100, // 100 milliseconds
			CleanupInterval: 50,  // 50 milliseconds
		},
	})
	defer cache.Close(ctx)

	// Test default expiration
	err := cache.Set(ctx, "key1", "value1")
	assert.NoError(t, err)

	time.Sleep(150 * time.Millisecond) // Wait for item to expire and cleanup to run

	_, err = cache.Get(ctx, "key1")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	// Test custom expiration
	err = cache.Set(ctx, "key2", "value2", 50*time.Millisecond)
	assert.NoError(t, err)

	val, err := cache.Get(ctx, "key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", val)

	time.Sleep(100 * time.Millisecond) // Wait for item to expire and cleanup to run

	_, err = cache.Get(ctx, "key2")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	// Test no expiration (0)
	cacheNoExp, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{
			Expiration:      0, // No expiration
			CleanupInterval: 50,
		},
	})
	defer cacheNoExp.Close(ctx)

	err = cacheNoExp.Set(ctx, "key3", "value3")
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond) // Wait longer than cleanup interval

	val, err = cacheNoExp.Get(ctx, "key3")
	assert.NoError(t, err)
	assert.Equal(t, "value3", val)
}

func TestMemoryCache_GetAndDelete(t *testing.T) {
	ctx := context.Background()
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{},
	})
	defer cache.Close(ctx)

	err := cache.Set(ctx, "key1", "value1")
	assert.NoError(t, err)

	val, err := cache.GetAndDelete(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// Should not be found after deletion
	_, err = cache.Get(ctx, "key1")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	// GetAndDelete non-existent key
	_, err = cache.GetAndDelete(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestMemoryCache_Exists(t *testing.T) {
	ctx := context.Background()
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{},
	})
	defer cache.Close(ctx)

	exists, err := cache.Exists(ctx, "key1")
	assert.NoError(t, err)
	assert.False(t, exists)

	err = cache.Set(ctx, "key1", "value1")
	assert.NoError(t, err)

	exists, err = cache.Exists(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, exists)

	err = cache.Delete(ctx, "key1")
	assert.NoError(t, err)

	exists, err = cache.Exists(ctx, "key1")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestMemoryCache_Delete(t *testing.T) {
	ctx := context.Background()
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{},
	})
	defer cache.Close(ctx)

	err := cache.Set(ctx, "key1", "value1")
	assert.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	err = cache.Delete(ctx, "key1")
	assert.NoError(t, err)

	_, err = cache.Get(ctx, "key1")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	// Deleting a non-existent key should not return an error
	err = cache.Delete(ctx, "nonexistent")
	assert.NoError(t, err)
}

func TestMemoryCache_Clear(t *testing.T) {
	ctx := context.Background()
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{},
	})
	defer cache.Close(ctx)

	err := cache.Set(ctx, "key1", "value1")
	assert.NoError(t, err)
	err = cache.Set(ctx, "key2", "value2")
	assert.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	err = cache.Clear(ctx)
	assert.NoError(t, err)

	_, err = cache.Get(ctx, "key1")
	assert.ErrorIs(t, err, ErrKeyNotFound)
	_, err = cache.Get(ctx, "key2")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	memCache := cache.(*Cache)
	assert.Len(t, memCache.items, 0)
}

func TestMemoryCache_SizeLimit(t *testing.T) {
	ctx := context.Background()
	customSize := int32(2)
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{
			Size: customSize,
		},
	})
	defer cache.Close(ctx)

	err := cache.Set(ctx, "key1", "value1")
	assert.NoError(t, err)
	err = cache.Set(ctx, "key2", "value2")
	assert.NoError(t, err)

	// Attempt to add a third item, should fail
	err = cache.Set(ctx, "key3", "value3")
	assert.ErrorIs(t, err, ErrCacheSizeLimitReached)

	// Verify existing items are still there
	_, err = cache.Get(ctx, "key1")
	assert.NoError(t, err)
	_, err = cache.Get(ctx, "key2")
	assert.NoError(t, err)

	// Verify the third item was not added
	_, err = cache.Get(ctx, "key3")
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestMemoryCache_Concurrency(t *testing.T) {
	ctx := context.Background()
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{
			Size: 100,
		},
	})
	defer cache.Close(ctx)

	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "key_" + fmt.Sprintf("%d", goroutineID) + "_" + fmt.Sprintf("%d", j)
				value := "value_" + fmt.Sprintf("%d", goroutineID) + "_" + fmt.Sprintf("%d", j)

				// Set
				_ = cache.Set(ctx, key, value)

				// Get
				_, _ = cache.Get(ctx, key)

				// Exists
				_, _ = cache.Exists(ctx, key)

				// Delete (randomly)
				if j%10 == 0 {
					_ = cache.Delete(ctx, key)
				}
			}
		}(i)
	}
	wg.Wait()

	// After all operations, ensure no panics and basic functionality holds
	_, err := cache.Get(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	memCache := cache.(*Cache)
	t.Logf("Final cache size: %d", len(memCache.items))
}

func TestMemoryCache_Close(t *testing.T) {
	ctx := context.Background()
	cache, _ := New(&cachev1.CacheConfig{
		Driver: DriverName,
		Memory: &cachev1.MemoryConfig{
			CleanupInterval: 10, // Small interval to ensure goroutine starts
		},
	})
	memCache := cache.(*Cache)

	// Ensure cleanup goroutine is running
	time.Sleep(20 * time.Millisecond)

	err := cache.Close(ctx)
	assert.NoError(t, err)

	// Verify that the stopCleanup channel is closed
	select {
	case _, ok := <-memCache.stopCleanup:
		assert.False(t, ok, "stopCleanup channel should be closed")
	default:
		assert.Fail(t, "stopCleanup channel should be closed")
	}

	// Attempting to close again should not panic
	err = cache.Close(ctx)
	assert.NoError(t, err)
}
