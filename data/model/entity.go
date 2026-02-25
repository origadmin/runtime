package model

import (
	"time"

	"github.com/origadmin/runtime/contracts/model"
)

// Entity provides common fields for all entities.
// It can be embedded into concrete business entities.
type Entity struct {
	ID        string    `json:"id"`         // Unique identifier for the entity
	CreatedAt time.Time `json:"created_at"` // Timestamp when the entity was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the entity was last updated
}

// GetID returns the ID of the Entity.
func (b *Entity) GetID() string {
	return b.ID
}

// SetID sets the ID of the Entity.
func (b *Entity) SetID(id string) {
	b.ID = id
}

// Ensure Entity implements Identifiable interface.
var _ model.Identifiable = (*Entity)(nil)
