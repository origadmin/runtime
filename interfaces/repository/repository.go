package repository

import "context"

// Repository is a generic interface for basic CRUD operations on an entity.
// T is the type of the entity, K is the type of the entity's ID.
type Repository[T any, K comparable] interface {
	// Save creates or updates an entity.
	Save(ctx context.Context, entity T) (T, error)
	// FindByID retrieves an entity by its ID.
	FindByID(ctx context.Context, id K) (T, error)
	// FindAll retrieves all entities.
	FindAll(ctx context.Context) ([]T, error)
	// DeleteByID deletes an entity by its ID.
	DeleteByID(ctx context.Context, id K) error
	// ExistsByID checks if an entity with the given ID exists.
	ExistsByID(ctx context.Context, id K) (bool, error)
}
