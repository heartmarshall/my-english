package meaning

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает meaning по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Meaning, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	meaning, err := r.scanRow(r.q.QueryRow(ctx, query, args...))
	if err != nil {
		if database.IsNotFoundError(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return meaning, nil
}

// GetByWordID возвращает все meanings для указанного слова.
func (r *Repo) GetByWordID(ctx context.Context, wordID int64) ([]*model.Meaning, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"word_id": wordID}).
		OrderBy("created_at ASC").
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

// GetByWordIDs возвращает все meanings для нескольких слов (batch loading).
func (r *Repo) GetByWordIDs(ctx context.Context, wordIDs []int64) ([]*model.Meaning, error) {
	if len(wordIDs) == 0 {
		return make([]*model.Meaning, 0), nil
	}

	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"word_id": wordIDs}).
		OrderBy("word_id ASC", "created_at ASC").
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

// List возвращает список meanings с фильтрацией и пагинацией.
func (r *Repo) List(ctx context.Context, filter *Filter, limit, offset int) ([]*model.Meaning, error) {
	limit, offset = database.NormalizePagination(limit, offset)

	qb := database.Builder.
		Select(columns...).
		From(tableName)

	qb = applyFilter(qb, filter)

	qb = qb.
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := qb.ToSql()
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

// Count возвращает количество meanings, соответствующих фильтру.
func (r *Repo) Count(ctx context.Context, filter *Filter) (int, error) {
	qb := database.Builder.
		Select("COUNT(*)").
		From(tableName)

	qb = applyFilter(qb, filter)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = r.q.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Exists проверяет существование meaning по ID.
func (r *Repo) Exists(ctx context.Context, id int64) (bool, error) {
	query, args, err := database.Builder.
		Select("1").
		From(tableName).
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = r.q.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if database.IsNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
