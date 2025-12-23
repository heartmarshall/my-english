package word

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает слово по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Word, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	word, err := r.scanRow(r.q.QueryRowContext(ctx, query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return word, nil
}

// GetByText возвращает слово по тексту (case-insensitive).
func (r *Repo) GetByText(ctx context.Context, text string) (*model.Word, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"text": strings.ToLower(strings.TrimSpace(text))}).
		ToSql()
	if err != nil {
		return nil, err
	}

	word, err := r.scanRow(r.q.QueryRowContext(ctx, query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return word, nil
}

// List возвращает список слов с фильтрацией и пагинацией.
// limit ограничивается до [1, MaxLimit], offset не может быть отрицательным.
func (r *Repo) List(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]*model.Word, error) {
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

	rows, err := r.q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// Count возвращает количество слов, соответствующих фильтру.
func (r *Repo) Count(ctx context.Context, filter *model.WordFilter) (int, error) {
	qb := database.Builder.
		Select("COUNT(*)").
		From(tableName)

	qb = applyFilter(qb, filter)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = r.q.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

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
	err = r.q.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
