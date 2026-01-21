package eventbus

import (
	"context"
	"time"
)

// Event is the base interface for all domain events.
type Event interface {
	EventName() string    // Returns the name of the event (e.g., "UserCreated")
	EventID() string      // Returns a unique ID for this specific event instance
	Timestamp() time.Time // Returns the time when the event occurred
	// 可以添加其他通用的事件元数据，如 Source, CorrelationID 等
}

// Publisher is an interface for publishing events.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// Subscriber is an interface for subscribing to events.
type Subscriber interface {
	Subscribe(ctx context.Context, eventName string, handler EventHandler) error
	Unsubscribe(ctx context.Context, eventName string, handler EventHandler) error
}

// EventHandler is a function type that handles a specific event.
type EventHandler func(ctx context.Context, event Event) error

// EventBus combines publishing and subscribing capabilities.
type EventBus interface {
	Publisher
	Subscriber
	// 如果事件总线是后台进程，可以添加 Start 和 Stop 方法来管理其生命周期
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
