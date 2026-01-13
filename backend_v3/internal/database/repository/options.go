// Package repository содержит базовые компоненты для работы с репозиториями.
package repository

import (
	"github.com/Masterminds/squirrel"
)

// ============================================================================
// PAGINATION
// ============================================================================

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// Pagination содержит параметры пагинации.
type Pagination struct {
	Limit  int
	Offset int
}

// NewPagination создаёт нормализованную пагинацию.
func NewPagination(limit, offset int) Pagination {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return Pagination{Limit: limit, Offset: offset}
}

// Apply применяет пагинацию к SelectBuilder.
func (p Pagination) Apply(b squirrel.SelectBuilder) squirrel.SelectBuilder {
	return b.Limit(uint64(p.Limit)).Offset(uint64(p.Offset))
}

// ============================================================================
// ORDERING
// ============================================================================

// OrderDirection направление сортировки.
type OrderDirection string

const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

// OrderBy описывает один элемент сортировки.
type OrderBy struct {
	Column    string
	Direction OrderDirection
}

// String возвращает строку для ORDER BY.
func (o OrderBy) String() string {
	if o.Direction == "" {
		o.Direction = OrderAsc
	}
	return o.Column + " " + string(o.Direction)
}

// Ordering содержит параметры сортировки.
type Ordering []OrderBy

// Apply применяет сортировку к SelectBuilder.
func (o Ordering) Apply(b squirrel.SelectBuilder) squirrel.SelectBuilder {
	if len(o) == 0 {
		return b
	}
	orderClauses := make([]string, len(o))
	for i, ob := range o {
		orderClauses[i] = ob.String()
	}
	return b.OrderBy(orderClauses...)
}

// ============================================================================
// QUERY OPTIONS
// ============================================================================

// QueryOption функция для модификации SelectBuilder.
type QueryOption func(squirrel.SelectBuilder) squirrel.SelectBuilder

// WithPagination добавляет пагинацию.
func WithPagination(limit, offset int) QueryOption {
	p := NewPagination(limit, offset)
	return func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
		return p.Apply(b)
	}
}

// WithOrderBy добавляет сортировку.
func WithOrderBy(column string, direction OrderDirection) QueryOption {
	return func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
		return b.OrderBy(column + " " + string(direction))
	}
}

// WithOrderByDesc добавляет сортировку по убыванию.
func WithOrderByDesc(column string) QueryOption {
	return WithOrderBy(column, OrderDesc)
}

// WithOrderByAsc добавляет сортировку по возрастанию.
func WithOrderByAsc(column string) QueryOption {
	return WithOrderBy(column, OrderAsc)
}

// WithWhere добавляет условие WHERE.
func WithWhere(pred any, args ...any) QueryOption {
	return func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
		return b.Where(pred, args...)
	}
}

// WithLimit добавляет LIMIT.
func WithLimit(limit int) QueryOption {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
		return b.Limit(uint64(limit))
	}
}

// WithOffset добавляет OFFSET.
func WithOffset(offset int) QueryOption {
	if offset < 0 {
		offset = 0
	}
	return func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
		return b.Offset(uint64(offset))
	}
}

// ApplyOptions применяет все опции к SelectBuilder.
func ApplyOptions(b squirrel.SelectBuilder, opts ...QueryOption) squirrel.SelectBuilder {
	for _, opt := range opts {
		b = opt(b)
	}
	return b
}

// ============================================================================
// FILTER HELPERS
// ============================================================================

// Filter представляет условие фильтрации.
type Filter struct {
	Column string
	Op     FilterOp
	Value  any
}

// FilterOp операция фильтрации.
type FilterOp string

const (
	FilterOpEq    FilterOp = "eq"
	FilterOpNotEq FilterOp = "neq"
	FilterOpLt    FilterOp = "lt"
	FilterOpLtEq  FilterOp = "lte"
	FilterOpGt    FilterOp = "gt"
	FilterOpGtEq  FilterOp = "gte"
	FilterOpIn    FilterOp = "in"
	FilterOpLike  FilterOp = "like"
	FilterOpILike FilterOp = "ilike"
)

// ToSquirrel конвертирует Filter в squirrel condition.
func (f Filter) ToSquirrel() squirrel.Sqlizer {
	switch f.Op {
	case FilterOpEq:
		return squirrel.Eq{f.Column: f.Value}
	case FilterOpNotEq:
		return squirrel.NotEq{f.Column: f.Value}
	case FilterOpLt:
		return squirrel.Lt{f.Column: f.Value}
	case FilterOpLtEq:
		return squirrel.LtOrEq{f.Column: f.Value}
	case FilterOpGt:
		return squirrel.Gt{f.Column: f.Value}
	case FilterOpGtEq:
		return squirrel.GtOrEq{f.Column: f.Value}
	case FilterOpIn:
		return squirrel.Eq{f.Column: f.Value}
	case FilterOpLike:
		return squirrel.Like{f.Column: f.Value}
	case FilterOpILike:
		return squirrel.ILike{f.Column: f.Value}
	default:
		return squirrel.Eq{f.Column: f.Value}
	}
}

// WithFilters добавляет несколько фильтров.
func WithFilters(filters ...Filter) QueryOption {
	return func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
		for _, f := range filters {
			b = b.Where(f.ToSquirrel())
		}
		return b
	}
}

