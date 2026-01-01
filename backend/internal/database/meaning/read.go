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
	builder := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.ID.Eq(id))

	return database.NewQuery[model.Meaning](r.q, builder).One(ctx)
}

// GetByWordID возвращает все meanings для указанного слова.
func (r *Repo) GetByWordID(ctx context.Context, wordID int64) ([]model.Meaning, error) {
	builder := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.WordID.Eq(wordID)).
		OrderBy(schema.Meanings.CreatedAt.Asc())

	return database.NewQuery[model.Meaning](r.q, builder).List(ctx)
}

// GetByWordIDs возвращает все meanings для нескольких слов (batch loading).
func (r *Repo) GetByWordIDs(ctx context.Context, wordIDs []int64) ([]model.Meaning, error) {
	if len(wordIDs) == 0 {
		return make([]model.Meaning, 0), nil
	}

	builder := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.WordID.In(wordIDs)).
		OrderBy(schema.Meanings.WordID.Asc(), schema.Meanings.CreatedAt.Asc())

	return database.NewQuery[model.Meaning](r.q, builder).List(ctx)
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

	return database.NewQuery[model.Meaning](r.q, qb).List(ctx)
}

// Count возвращает количество meanings, соответствующих фильтру.
func (r *Repo) Count(ctx context.Context, filter *Filter) (int, error) {
	qb := database.Builder.Select("COUNT(*)").From(schema.Meanings.Name.String())
	qb = applyFilter(qb, filter)

	return database.NewQuery[int](r.q, qb).Scalar(ctx)
}

// Exists проверяет существование meaning по ID.
func (r *Repo) Exists(ctx context.Context, id int64) (bool, error) {
	builder := database.Builder.
		Select("1").
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.ID.Eq(id)).
		Limit(1)

	val, err := database.NewQuery[int](r.q, builder).Scalar(ctx)
	if err != nil {
		if err == database.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return val > 0, nil
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
