package database

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

// ============================================================================
// LEGACY QUERY EXECUTOR
// Сохранено для обратной совместимости.
// Для новых репозиториев используйте repository.Base[T].
// ============================================================================

// SQLBuilder интерфейс для squirrel.SelectBuilder, InsertBuilder и т.д.
type SQLBuilder interface {
	ToSql() (string, []interface{}, error)
}

// Exec выполняет generic запросы через squirrel builder.
type Exec[T any] struct {
	q   Querier
	b   SQLBuilder
	err error
}

// NewQuery создаёт новый экзекьютор.
func NewQuery[T any](q Querier, b SQLBuilder) *Exec[T] {
	return &Exec[T]{q: q, b: b}
}

// One возвращает одну сущность или ошибку ErrNotFound.
func (e *Exec[T]) One(ctx context.Context) (T, error) {
	var dest T
	if e.err != nil {
		return dest, e.err
	}

	sql, args, err := e.b.ToSql()
	if err != nil {
		return dest, WrapDBError(err)
	}

	if err := pgxscan.Get(ctx, e.q, &dest, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return dest, ErrNotFound
		}
		return dest, WrapDBError(err)
	}

	return dest, nil
}

// List возвращает слайс сущностей.
func (e *Exec[T]) List(ctx context.Context) ([]T, error) {
	if e.err != nil {
		return nil, e.err
	}

	sql, args, err := e.b.ToSql()
	if err != nil {
		return nil, WrapDBError(err)
	}

	var dest []T
	if err := pgxscan.Select(ctx, e.q, &dest, sql, args...); err != nil {
		return nil, WrapDBError(err)
	}

	return dest, nil
}

// Scalar возвращает примитив (int, string, bool).
func (e *Exec[T]) Scalar(ctx context.Context) (T, error) {
	return e.One(ctx)
}

// ============================================================================
// STANDALONE HELPERS
// ============================================================================

// ExecOnly выполняет запрос (INSERT/UPDATE/DELETE) и возвращает кол-во затронутых строк.
func ExecOnly(ctx context.Context, q Querier, b SQLBuilder) (int64, error) {
	sql, args, err := b.ToSql()
	if err != nil {
		return 0, WrapDBError(err)
	}

	tag, err := q.Exec(ctx, sql, args...)
	if err != nil {
		return 0, WrapDBError(err)
	}
	return tag.RowsAffected(), nil
}

// ExecInsertWithReturn выполняет INSERT ... RETURNING id и сканирует результат.
func ExecInsertWithReturn[T any](ctx context.Context, q Querier, b SQLBuilder) (T, error) {
	var id T
	sql, args, err := b.ToSql()
	if err != nil {
		return id, WrapDBError(err)
	}

	err = q.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return id, WrapDBError(err)
	}
	return id, nil
}

// ============================================================================
// RAW SQL HELPERS
// ============================================================================

// GetOne сканирует одну структуру из raw SQL.
func GetOne[T any](ctx context.Context, q Querier, sql string, args ...any) (*T, error) {
	var dest T
	err := pgxscan.Get(ctx, q, &dest, sql, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &dest, nil
}

// Select сканирует список структур из raw SQL.
func Select[T any](ctx context.Context, q Querier, sql string, args ...any) ([]T, error) {
	var dest []T
	err := pgxscan.Select(ctx, q, &dest, sql, args...)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// GetScalar сканирует одно скалярное значение.
func GetScalar[T any](ctx context.Context, q Querier, sql string, args ...any) (T, error) {
	var dest T
	err := pgxscan.Get(ctx, q, &dest, sql, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return dest, nil
		}
		return dest, err
	}
	return dest, nil
}

// SelectScalars сканирует список скалярных значений.
func SelectScalars[T any](ctx context.Context, q Querier, sql string, args ...any) ([]T, error) {
	var dest []T
	err := pgxscan.Select(ctx, q, &dest, sql, args...)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// CheckExists проверяет наличие строки.
func CheckExists(ctx context.Context, q Querier, sql string, args ...any) (bool, error) {
	var dummy int
	err := pgxscan.Get(ctx, q, &dummy, sql, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ============================================================================
// BUILDER
// ============================================================================

// Builder — глобальный squirrel builder с PostgreSQL placeholder format.
var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
