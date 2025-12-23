package example

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает example по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Example, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	example, err := r.scanRow(r.q.QueryRow(ctx, query, args...))
	if err != nil {
		if database.IsNotFoundError(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return example, nil
}

// GetByMeaningID возвращает все examples для указанного meaning.
func (r *Repo) GetByMeaningID(ctx context.Context, meaningID int64) ([]*model.Example, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"meaning_id": meaningID}).
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

// GetByMeaningIDs возвращает examples для нескольких meanings (для batch loading).
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]*model.Example, error) {
	if len(meaningIDs) == 0 {
		return make([]*model.Example, 0), nil
	}

	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"meaning_id": meaningIDs}).
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
