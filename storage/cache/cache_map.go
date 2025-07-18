package cache

import (
	"context"
	"sync"
	"time"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

const (
	DefaultCleanupInterval = 5 * time.Minute
)

var (
	ErrClosed         error = &cacheError{msg: "cache closed"}
	ErrNotFound       error = &cacheError{msg: "cache not found"}
	ErrInvalidElement error = &cacheError{msg: "invalid cache element"}
	ErrExpired        error = &cacheError{msg: "cache expired"}
)

type element struct {
	value    string
	expireAt time.Time
}

type mapCache struct {
	elems *sync.Pool
	maps  *sync.Map

	cleanupInterval time.Duration // How often to run cleanup
	stopCleanup     chan struct{} // Channel to signal cleanup goroutine to stop
}

func (s *mapCache) Clear(ctx context.Context) error {
	s.maps.Clear()
	return nil
}

func (s *mapCache) Exists(ctx context.Context, key string) (bool, error) {
	value, ok := s.maps.Load(key)
	if !ok {
		return false, nil
	}
	elem, ok := value.(*element)
	if !ok {
		return false, ErrInvalidElement
	}
	if !elem.expireAt.IsZero() && elem.expireAt.Before(time.Now()) {
		s.maps.Delete(key)
		s.elems.Put(elem)
		return false, ErrExpired
	}
	return true, nil
}

func (s *mapCache) Close(ctx context.Context) error {
	close(s.stopCleanup) // Signal cleanup goroutine to stop
	return nil
}

func (s *mapCache) Get(ctx context.Context, key string) (string, error) {
	value, ok := s.maps.Load(key)
	if !ok {
		return "", ErrNotFound
	}
	elem, ok := value.(*element)
	if !ok {
		return "", ErrInvalidElement
	}
	if !elem.expireAt.IsZero() && elem.expireAt.Before(time.Now()) {
		s.maps.Delete(key)
		s.elems.Put(elem)
		return "", ErrExpired
	}
	return elem.value, nil
}

func (s *mapCache) GetAndDelete(ctx context.Context, key string) (string, error) {
	value, ok := s.maps.LoadAndDelete(key)
	if !ok {
		return "", ErrNotFound
	}
	elem, ok := value.(*element)
	if !ok {
		return "", ErrInvalidElement
	}
	if !elem.expireAt.IsZero() && elem.expireAt.Before(time.Now()) {
		s.elems.Put(elem)
		return "", ErrExpired
	}
	v := elem.value
	s.elems.Put(elem)
	return v, nil
}

func (s *mapCache) Set(ctx context.Context, key string, value string, expiration ...time.Duration) error {
	var expireAt time.Time
	if len(expiration) > 0 && expiration[0].Seconds() > 0 {
		expireAt = time.Now().Add(expiration[0])
	}
	elem := s.getElement()
	elem.value = value
	elem.expireAt = expireAt
	s.maps.Store(key, elem)
	return nil
}

func (s *mapCache) Delete(ctx context.Context, key string) error {
	if v, ok := s.maps.LoadAndDelete(key); ok {
		s.elems.Put(v)
		return nil
	}
	return ErrNotFound
}

func (s *mapCache) getElement() *element {
	return s.elems.Get().(*element)
}

func (s *mapCache) putElement(elem *element) {
	s.elems.Put(elem)
}

func NewMemoryCache(cfg *configv1.Memory) storageiface.Cache {
	interval := DefaultCleanupInterval // Default cleanup interval
	if cfg != nil && cfg.CleanupInterval > 0 {
		interval = time.Duration(cfg.CleanupInterval) * time.Second
	}

	mc := &mapCache{
		elems: &sync.Pool{
			New: func() any {
				return &element{}
			},
		},
		maps:            &sync.Map{},
		cleanupInterval: interval,
		stopCleanup:     make(chan struct{}),
	}
	go mc.cleanupLoop()
	return mc
}

// cleanupLoop runs periodically to remove expired entries.
func (s *mapCache) cleanupLoop() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpiredEntries()
		case <-s.stopCleanup:
			return // Stop the goroutine
		}
	}
}

// cleanupExpiredEntries iterates through the map and removes expired items.
func (s *mapCache) cleanupExpiredEntries() {
	now := time.Now()
	s.maps.Range(func(key, value interface{}) bool {
		elem, ok := value.(*element)
		if !ok {
			// Log error or handle unexpected type if necessary
			return true // Continue iteration
		}

		if !elem.expireAt.IsZero() && elem.expireAt.Before(now) {
			// Item is expired, delete it and put back to pool
			if actual, loaded := s.maps.LoadAndDelete(key); loaded {
				if actualElem, ok := actual.(*element); ok {
					s.elems.Put(actualElem)
				}
			}
		}
		return true // Continue iteration
	})
}

var _ storageiface.Cache = (*mapCache)(nil)
