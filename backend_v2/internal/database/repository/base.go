package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/heartmarshall/my-english/internal/database"
)

// Builder — глобальный squirrel builder с PostgreSQL placeholder format.
var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// ============================================================================
// BASE REPOSITORY
// ============================================================================

// Base предоставляет общие CRUD операции для репозиториев.
// Используй embedding для создания специализированных репозиториев.
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

// Q возвращает Querier для использования в специализированных методах.
func (r *Base[T]) Q() database.Querier {
	return r.querier
}

// Table возвращает имя таблицы.
func (r *Base[T]) Table() string {
	return r.table
}

// Columns возвращает список колонок для SELECT.
func (r *Base[T]) Columns() []string {
	return r.columns
}

// SelectBuilder возвращает базовый SELECT builder.
func (r *Base[T]) SelectBuilder() squirrel.SelectBuilder {
	return Builder.Select(r.columns...).From(r.table)
}

// ============================================================================
// READ OPERATIONS
// ============================================================================

// GetByID находит сущность по ID (первичный ключ).
// idColumn — имя колонки ID (например "id" или "cards.id").
func (r *Base[T]) GetByID(ctx context.Context, idColumn string, id any) (*T, error) {
	query := r.SelectBuilder().Where(squirrel.Eq{idColumn: id})
	return r.GetOne(ctx, query)
}

// GetOne выполняет запрос и возвращает одну сущность.
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

// List выполняет запрос и возвращает список сущностей.
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

// ListAll возвращает все сущности с опциональными модификаторами.
func (r *Base[T]) ListAll(ctx context.Context, opts ...QueryOption) ([]T, error) {
	query := ApplyOptions(r.SelectBuilder(), opts...)
	return r.List(ctx, query)
}

// FindBy находит сущности по условию.
func (r *Base[T]) FindBy(ctx context.Context, column string, value any, opts ...QueryOption) ([]T, error) {
	query := r.SelectBuilder().Where(squirrel.Eq{column: value})
	query = ApplyOptions(query, opts...)
	return r.List(ctx, query)
}

// ListByIDs находит сущности по списку ID (WHERE id IN (...)).
// idColumn — имя колонки ID (например "id" или "cards.id").
// ids — слайс ID для поиска. Если пустой, возвращает пустой слайс без запроса к БД.
func (r *Base[T]) ListByIDs(ctx context.Context, idColumn string, ids []any, opts ...QueryOption) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	query := r.SelectBuilder().Where(squirrel.Eq{idColumn: ids})
	query = ApplyOptions(query, opts...)
	return r.List(ctx, query)
}

// FindOneBy находит одну сущность по условию.
func (r *Base[T]) FindOneBy(ctx context.Context, column string, value any) (*T, error) {
	query := r.SelectBuilder().Where(squirrel.Eq{column: value}).Limit(1)
	return r.GetOne(ctx, query)
}

// Exists проверяет существование записи.
func (r *Base[T]) Exists(ctx context.Context, column string, value any) (bool, error) {
	query := Builder.Select("1").From(r.table).Where(squirrel.Eq{column: value}).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, database.WrapDBError(err)
	}

	var dummy int
	if err := pgxscan.Get(ctx, r.querier, &dummy, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return false, nil
		}
		return false, database.WrapDBError(err)
	}
	return true, nil
}

// Count возвращает количество записей.
func (r *Base[T]) Count(ctx context.Context, opts ...QueryOption) (int64, error) {
	query := Builder.Select("COUNT(*)").From(r.table)
	query = ApplyOptions(query, opts...)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, database.WrapDBError(err)
	}

	var count int64
	if err := pgxscan.Get(ctx, r.querier, &count, sql, args...); err != nil {
		return 0, database.WrapDBError(err)
	}
	return count, nil
}

// ============================================================================
// WRITE OPERATIONS
// ============================================================================

// InsertBuilder возвращает INSERT builder для таблицы.
func (r *Base[T]) InsertBuilder() squirrel.InsertBuilder {
	return Builder.Insert(r.table)
}

// UpdateBuilder возвращает UPDATE builder для таблицы.
func (r *Base[T]) UpdateBuilder() squirrel.UpdateBuilder {
	return Builder.Update(r.table)
}

// DeleteBuilder возвращает DELETE builder для таблицы.
func (r *Base[T]) DeleteBuilder() squirrel.DeleteBuilder {
	return Builder.Delete(r.table)
}

// Insert выполняет INSERT и возвращает ID (RETURNING).
func (r *Base[T]) Insert(ctx context.Context, insert squirrel.InsertBuilder, returningCol string) (any, error) {
	query := insert.Suffix("RETURNING " + returningCol)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var id any
	if err := r.querier.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return nil, database.WrapDBError(err)
	}
	return id, nil
}

// InsertReturning выполняет INSERT RETURNING * и сканирует результат в структуру.
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

// InsertReturningMany выполняет batch INSERT и возвращает все вставленные записи.
// Используется для INSERT с ON CONFLICT, который может вернуть несколько строк.
func (r *Base[T]) InsertReturningMany(ctx context.Context, insert squirrel.InsertBuilder) ([]T, error) {
	sql, args, err := insert.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var dest []T
	if err := pgxscan.Select(ctx, r.querier, &dest, sql, args...); err != nil {
		return nil, database.WrapDBError(err)
	}
	return dest, nil
}

// Update выполняет UPDATE и возвращает количество затронутых строк.
func (r *Base[T]) Update(ctx context.Context, update squirrel.UpdateBuilder) (int64, error) {
	sql, args, err := update.ToSql()
	if err != nil {
		return 0, database.WrapDBError(err)
	}

	tag, err := r.querier.Exec(ctx, sql, args...)
	if err != nil {
		return 0, database.WrapDBError(err)
	}
	return tag.RowsAffected(), nil
}

// UpdateReturning выполняет UPDATE RETURNING * и сканирует результат.
func (r *Base[T]) UpdateReturning(ctx context.Context, update squirrel.UpdateBuilder) (*T, error) {
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

// Delete выполняет DELETE и возвращает количество удалённых строк.
func (r *Base[T]) Delete(ctx context.Context, del squirrel.DeleteBuilder) (int64, error) {
	sql, args, err := del.ToSql()
	if err != nil {
		return 0, database.WrapDBError(err)
	}

	tag, err := r.querier.Exec(ctx, sql, args...)
	if err != nil {
		return 0, database.WrapDBError(err)
	}
	return tag.RowsAffected(), nil
}

// DeleteByID удаляет запись по ID.
func (r *Base[T]) DeleteByID(ctx context.Context, idColumn string, id any) error {
	del := r.DeleteBuilder().Where(squirrel.Eq{idColumn: id})

	affected, err := r.Delete(ctx, del)
	if err != nil {
		return err
	}
	if affected == 0 {
		return database.ErrNotFound
	}
	return nil
}

// ============================================================================
// BATCH OPERATIONS
// ============================================================================

// BatchInsert выполняет вставку нескольких записей.
// columns — колонки для INSERT, valuesFunc — функция, возвращающая значения для каждой записи.
func (r *Base[T]) BatchInsert(ctx context.Context, columns []string, items []T, valuesFunc func(T) []any) (int64, error) {
	if len(items) == 0 {
		return 0, nil
	}

	insert := r.InsertBuilder().Columns(columns...)
	for _, item := range items {
		insert = insert.Values(valuesFunc(item)...)
	}

	sql, args, err := insert.ToSql()
	if err != nil {
		return 0, database.WrapDBError(err)
	}

	tag, err := r.querier.Exec(ctx, sql, args...)
	if err != nil {
		return 0, database.WrapDBError(err)
	}
	return tag.RowsAffected(), nil
}

// ============================================================================
// SCALAR QUERIES
// ============================================================================

// GetScalar выполняет запрос и возвращает скалярное значение.
func GetScalar[R any](ctx context.Context, q database.Querier, query squirrel.SelectBuilder) (R, error) {
	var result R

	sql, args, err := query.ToSql()
	if err != nil {
		return result, database.WrapDBError(err)
	}

	if err := pgxscan.Get(ctx, q, &result, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return result, nil
		}
		return result, database.WrapDBError(err)
	}
	return result, nil
}

// SelectScalars выполняет запрос и возвращает список скалярных значений.
func SelectScalars[R any](ctx context.Context, q database.Querier, query squirrel.SelectBuilder) ([]R, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var dest []R
	if err := pgxscan.Select(ctx, q, &dest, sql, args...); err != nil {
		return nil, database.WrapDBError(err)
	}
	return dest, nil
}
