package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/heartmarshall/my-english/internal/database"
)

// Builder — глобальный squirrel builder с PostgreSQL placeholder format.
var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// Base предоставляет общие CRUD операции для репозиториев.
type Base[T any] struct {
	querier database.Querier
	table   string
	columns []string
}

// NewBase создаёт базовый репозиторий.
func NewBase[T any](q database.Querier, table string, columns []string) *Base[T] {
	return &Base[T]{
		querier: q,
		table:   table,
		columns: columns,
	}
}

func (r *Base[T]) Q() database.Querier { return r.querier }

func (r *Base[T]) SelectBuilder() squirrel.SelectBuilder {
	return Builder.Select(r.columns...).From(r.table)
}

func (r *Base[T]) InsertBuilder() squirrel.InsertBuilder {
	return Builder.Insert(r.table)
}

func (r *Base[T]) UpdateBuilder() squirrel.UpdateBuilder {
	return Builder.Update(r.table)
}

func (r *Base[T]) DeleteBuilder() squirrel.DeleteBuilder {
	return Builder.Delete(r.table)
}

// --- READ OPERATIONS ---

func (r *Base[T]) GetByID(ctx context.Context, idColumn string, id any) (*T, error) {
	query := r.SelectBuilder().Where(squirrel.Eq{idColumn: id})
	return r.GetOne(ctx, query)
}

func (r *Base[T]) GetOne(ctx context.Context, query squirrel.SelectBuilder) (*T, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var dest T
	if err := pgxscan.Get(ctx, r.querier, &dest, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, database.ErrNotFound
		}
		return nil, database.WrapDBError(err)
	}
	return &dest, nil
}

func (r *Base[T]) List(ctx context.Context, query squirrel.SelectBuilder) ([]T, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var dest []T
	if err := pgxscan.Select(ctx, r.querier, &dest, sql, args...); err != nil {
		return nil, database.WrapDBError(err)
	}
	return dest, nil
}

func (r *Base[T]) ListByIDs(ctx context.Context, idColumn string, ids []any) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}
	query := r.SelectBuilder().Where(squirrel.Eq{idColumn: ids})
	return r.List(ctx, query)
}

func (r *Base[T]) FindOneBy(ctx context.Context, column string, value any) (*T, error) {
	query := r.SelectBuilder().Where(squirrel.Eq{column: value}).Limit(1)
	return r.GetOne(ctx, query)
}

// --- WRITE OPERATIONS ---

func (r *Base[T]) Create(ctx context.Context, insert squirrel.InsertBuilder) (*T, error) {
	// Обертка для InsertReturning для удобства
	return r.InsertReturning(ctx, insert)
}

func (r *Base[T]) InsertReturning(ctx context.Context, insert squirrel.InsertBuilder) (*T, error) {
	query := insert.Suffix("RETURNING *")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var dest T
	if err := pgxscan.Get(ctx, r.querier, &dest, sql, args...); err != nil {
		return nil, database.WrapDBError(err)
	}
	return &dest, nil
}

// BatchInsertReturning выполняет множественную вставку и возвращает созданные сущности.
func (r *Base[T]) BatchInsertReturning(ctx context.Context, columns []string, items []T, valuesFunc func(T) []any) ([]T, error) {
	if len(items) == 0 {
		return []T{}, nil
	}

	insert := r.InsertBuilder().Columns(columns...)
	for _, item := range items {
		insert = insert.Values(valuesFunc(item)...)
	}

	query := insert.Suffix("RETURNING *")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var dest []T
	if err := pgxscan.Select(ctx, r.querier, &dest, sql, args...); err != nil {
		return nil, database.WrapDBError(err)
	}
	return dest, nil
}

func (r *Base[T]) Update(ctx context.Context, update squirrel.UpdateBuilder) (*T, error) {
	query := update.Suffix("RETURNING *")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var dest T
	if err := pgxscan.Get(ctx, r.querier, &dest, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, database.ErrNotFound
		}
		return nil, database.WrapDBError(err)
	}
	return &dest, nil
}

func (r *Base[T]) Delete(ctx context.Context, idColumn string, id any) error {
	del := r.DeleteBuilder().Where(squirrel.Eq{idColumn: id})
	sql, args, err := del.ToSql()
	if err != nil {
		return database.WrapDBError(err)
	}
	tag, err := r.querier.Exec(ctx, sql, args...)
	if err != nil {
		return database.WrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return database.ErrNotFound
	}
	return nil
}
