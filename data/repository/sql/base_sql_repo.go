package sql

import (
	"context"
	"database/sql"

	"github.com/origadmin/toolkits/errors" // 假设有一个通用的错误包
)

// BaseSQLRepository provides common functionality for SQL-based repositories.
// It's meant to be embedded in concrete repository implementations or used as a reference.
// T is the entity type, K is the ID type.
type BaseSQLRepository[T any, K comparable] struct {
	DB        *sql.DB
	TableName string
	// 可以添加其他通用字段，如 logger, metrics 等
}

// NewBaseSQLRepository creates a new BaseSQLRepository.
func NewBaseSQLRepository[T any, K comparable](db *sql.DB, tableName string) *BaseSQLRepository[T, K] {
	return &BaseSQLRepository[T, K]{
		DB:        db,
		TableName: tableName,
	}
}

// Save is a placeholder for a generic save operation.
// A truly generic Save method without an ORM is highly complex and usually implemented
// by concrete business repositories using reflection or a SQL builder.
func (r *BaseSQLRepository[T, K]) Save(ctx context.Context, entity T) (T, error) {
	// Implementations will typically use an ORM (ent-go, gorm) or a SQL builder here.
	// This method is a placeholder to show where the logic would go.
	return entity, errors.New("generic Save not implemented for BaseSQLRepository; implement in concrete repository")
}

// FindByID is a placeholder for a generic find by ID operation.
func (r *BaseSQLRepository[T, K]) FindByID(ctx context.Context, id K) (T, error) {
	var entity T
	// Implementations will typically use an ORM (ent-go, gorm) or a SQL builder here.
	// This method is a placeholder to show where the logic would go.
	return entity, errors.New("generic FindByID not implemented for BaseSQLRepository; implement in concrete repository")
}

// FindAll is a placeholder for a generic find all operation.
func (r *BaseSQLRepository[T, K]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	// Implementations will typically use an ORM (ent-go, gorm) or a SQL builder here.
	// This method is a placeholder to show where the logic would go.
	return entities, errors.New("generic FindAll not implemented for BaseSQLRepository; implement in concrete repository")
}

// DeleteByID is a placeholder for a generic delete by ID operation.
func (r *BaseSQLRepository[T, K]) DeleteByID(ctx context.Context, id K) error {
	// Implementations will typically use an ORM (ent-go, gorm) or a SQL builder here.
	// This method is a placeholder to show where the logic would go.
	return errors.New("generic DeleteByID not implemented for BaseSQLRepository; implement in concrete repository")
}

// ExistsByID is a placeholder for a generic exists by ID operation.
func (r *BaseSQLRepository[T, K]) ExistsByID(ctx context.Context, id K) (bool, error) {
	// Implementations will typically use an ORM (ent-go, gorm) or a SQL builder here.
	// This method is a placeholder to show where the logic would go.
	return false, errors.New("generic ExistsByID not implemented for BaseSQLRepository; implement in concrete repository")
}
