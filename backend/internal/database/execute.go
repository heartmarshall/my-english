package database

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
)

// SQLBuilder интерфейс для squirrel.SelectBuilder, InsertBuilder и т.д.
type SQLBuilder interface {
	ToSql() (string, []interface{}, error)
}

// Exec executes generic queries and maps them.
type Exec[T any] struct {
	q   Querier
	b   SQLBuilder
	err error
}

// NewQuery создает новый экзекьютор
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
	// Re-use logic from One, works for scalars too in scany
	return e.One(ctx)
}

// ExecOnly выполняет запрос (INSERT/UPDATE/DELETE) и возвращает кол-во затронутых строк.
// Для этого не нужен Generic type, сделаем функцию-хелпер.
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
