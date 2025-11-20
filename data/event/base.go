package event

import (
	"time"

	"github.com/google/uuid"

	"github.com/origadmin/runtime/interfaces/event"
)

// BaseEvent provides common fields for events.
type BaseEvent struct {
	name      string
	id        string
	timestamp time.Time
}

// NewBaseEvent creates a new BaseEvent with a unique ID and current timestamp.
func NewBaseEvent(name string) BaseEvent {
	return BaseEvent{
		name:      name,
		id:        uuid.NewString(),
		timestamp: time.Now(),
	}
}

func (b *BaseEvent) EventName() string    { return b.name }
func (b *BaseEvent) EventID() string      { return b.id }
func (b *BaseEvent) Timestamp() time.Time { return b.timestamp }

// Ensure BaseEvent implements Event interface.
var _ event.Event = (*BaseEvent)(nil)
