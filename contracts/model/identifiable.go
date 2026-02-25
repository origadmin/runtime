package model

// Identifiable defines an interface for entities that have an ID.
type Identifiable interface {
	GetID() string
	SetID(id string)
}
