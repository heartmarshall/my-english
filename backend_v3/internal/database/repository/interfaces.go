package repository

import (
	"context"
)

// ============================================================================
// GENERIC REPOSITORY INTERFACES
// Используй эти интерфейсы для определения контрактов репозиториев.
// ============================================================================

// Reader интерфейс для операций чтения.
type Reader[T any, ID any] interface {
	GetByID(ctx context.Context, id ID) (*T, error)
	List(ctx context.Context, opts ...QueryOption) ([]T, error)
	ListByIDs(ctx context.Context, ids []ID) ([]T, error)
	Exists(ctx context.Context, id ID) (bool, error)
	Count(ctx context.Context, opts ...QueryOption) (int64, error)
}

// Writer интерфейс для операций записи.
type Writer[T any, ID any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id ID) error
}

// Repository объединяет Reader и Writer.
type Repository[T any, ID any] interface {
	Reader[T, ID]
	Writer[T, ID]
}

// ============================================================================
// PAGINATED RESULT
// ============================================================================

// PagedResult содержит результаты с пагинацией.
type PagedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	HasMore    bool  `json:"has_more"`
}

// NewPagedResult создаёт результат с пагинацией.
func NewPagedResult[T any](items []T, totalCount int64, limit, offset int) PagedResult[T] {
	return PagedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
		HasMore:    int64(offset+len(items)) < totalCount,
	}
}
