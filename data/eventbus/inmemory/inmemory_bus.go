package inmemory

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/origadmin/runtime/interfaces/eventbus"
)

// SimpleInMemoryEventBus is a basic, in-process event bus implementation.
// Note: This implementation is for simple use cases and does not persist events.
// Handlers are executed synchronously in the publisher's goroutine.
type SimpleInMemoryEventBus struct {
	handlers map[string][]eventbus.EventHandler
	mu       sync.RWMutex
}

// NewSimpleInMemoryEventBus creates a new SimpleInMemoryEventBus.
func NewSimpleInMemoryEventBus() *SimpleInMemoryEventBus {
	return &SimpleInMemoryEventBus{
		handlers: make(map[string][]eventbus.EventHandler),
	}
}

// Publish executes all registered handlers for a given event.
func (b *SimpleInMemoryEventBus) Publish(ctx context.Context, evt eventbus.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers, found := b.handlers[evt.EventName()]
	if !found {
		return nil // No handlers for this event, not an error.
	}

	// In a real-world system, handlers might run in separate goroutines
	// with a worker pool to avoid blocking the publisher and to add concurrency.
	for _, handler := range handlers {
		if err := handler(ctx, evt); err != nil {
			// In a real system, you'd use a structured logger.
			fmt.Printf("error handling event %s: %v\n", evt.EventName(), err)
			// Decide on error handling strategy: continue or stop?
			// For this simple bus, we continue.
		}
	}
	return nil
}

// Subscribe registers an event handler for a specific event name.
func (b *SimpleInMemoryEventBus) Subscribe(ctx context.Context, eventName string, handler eventbus.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
	return nil
}

// Unsubscribe removes a specific event handler for an event name.
// Note: Unsubscribing based on function equality can be unreliable in Go.
// A more robust implementation would require handlers to have unique IDs.
func (b *SimpleInMemoryEventBus) Unsubscribe(ctx context.Context, eventName string, handler eventbus.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers, found := b.handlers[eventName]
	if !found {
		return nil
	}

	// Find the handler to remove.
	// This is O(n) and relies on reflect.ValueOf, which is not ideal for performance.
	// For a production system, consider a map of handler IDs.
	handlerPtr := reflect.ValueOf(handler).Pointer()
	for i, h := range handlers {
		if reflect.ValueOf(h).Pointer() == handlerPtr {
			b.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
			return nil
		}
	}
	return nil
}

// Start does nothing for this simple in-memory bus.
func (b *SimpleInMemoryEventBus) Start(ctx context.Context) error {
	return nil
}

// Stop does nothing for this simple in-memory bus.
func (b *SimpleInMemoryEventBus) Stop(ctx context.Context) error {
	return nil
}

// Ensure SimpleInMemoryEventBus implements EventBus interface.
var _ eventbus.EventBus = (*SimpleInMemoryEventBus)(nil)
