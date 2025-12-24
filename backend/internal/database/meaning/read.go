package meaning

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает meaning по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (model.Meaning, error) {
	query, args, err := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.ID.Eq(id)).
		ToSql()
	if err != nil {
		return model.Meaning{}, err
	}

	meaning, err := database.GetOne[model.Meaning](ctx, r.q, query, args...)
	if err != nil {
		return model.Meaning{}, err
	}
	return *meaning, nil
}

// GetByWordID возвращает все meanings для указанного слова.
func (r *Repo) GetByWordID(ctx context.Context, wordID int64) ([]model.Meaning, error) {
	query, args, err := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.WordID.Eq(wordID)).
		OrderBy(schema.Meanings.CreatedAt.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Meaning](ctx, r.q, query, args...)
}

// GetByWordIDs возвращает все meanings для нескольких слов (batch loading).
func (r *Repo) GetByWordIDs(ctx context.Context, wordIDs []int64) ([]model.Meaning, error) {
	if len(wordIDs) == 0 {
		return make([]model.Meaning, 0), nil
	}

	query, args, err := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.WordID.In(wordIDs)).
		OrderBy(schema.Meanings.WordID.Asc(), schema.Meanings.CreatedAt.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Meaning](ctx, r.q, query, args...)
}

// List возвращает список meanings с фильтрацией и пагинацией.
func (r *Repo) List(ctx context.Context, filter *Filter, limit, offset int) ([]model.Meaning, error) {
	limit, offset = database.NormalizePagination(limit, offset)

	qb := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String())

	qb = applyFilter(qb, filter)

	qb = qb.
		OrderBy(schema.Meanings.CreatedAt.Desc()).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Meaning](ctx, r.q, query, args...)
}

// Count возвращает количество meanings, соответствующих фильтру.
func (r *Repo) Count(ctx context.Context, filter *Filter) (int, error) {
	qb := database.Builder.Select("COUNT(*)").From(schema.Meanings.Name.String())
	qb = applyFilter(qb, filter)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	return database.GetScalar[int](ctx, r.q, query, args...)
}

// Exists проверяет существование meaning по ID.
func (r *Repo) Exists(ctx context.Context, id int64) (bool, error) {
	query, args, err := database.Builder.
		Select("1").
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.ID.Eq(id)).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	return database.CheckExists(ctx, r.q, query, args...)
}

// applyFilter применяет фильтры к query builder.
func applyFilter(qb squirrel.SelectBuilder, filter *Filter) squirrel.SelectBuilder {
	if filter == nil {
		return qb
	}
	if filter.WordID != nil {
		qb = qb.Where(schema.Meanings.WordID.Eq(*filter.WordID))
	}
	if filter.PartOfSpeech != nil {
		qb = qb.Where(schema.Meanings.PartOfSpeech.Eq(*filter.PartOfSpeech))
	}
	if filter.LearningStatus != nil {
		qb = qb.Where(schema.Meanings.LearningStatus.Eq(*filter.LearningStatus))
	}
	return qb
}
