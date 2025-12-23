package tag

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает tag по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	tag, err := r.scanRow(r.q.QueryRow(ctx, query, args...))
	if err != nil {
		if database.IsNotFoundError(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return tag, nil
}

// GetByName возвращает tag по имени.
func (r *Repo) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, err
	}

	tag, err := r.scanRow(r.q.QueryRow(ctx, query, args...))
	if err != nil {
		if database.IsNotFoundError(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return tag, nil
}

// GetByNames возвращает tags по списку имён.
func (r *Repo) GetByNames(ctx context.Context, names []string) ([]*model.Tag, error) {
	if len(names) == 0 {
		return make([]*model.Tag, 0), nil
	}

	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"name": names}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// GetByIDs возвращает tags по списку ID.
func (r *Repo) GetByIDs(ctx context.Context, ids []int64) ([]*model.Tag, error) {
	if len(ids) == 0 {
		return make([]*model.Tag, 0), nil
	}

	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// List возвращает все tags.
func (r *Repo) List(ctx context.Context) ([]*model.Tag, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}
