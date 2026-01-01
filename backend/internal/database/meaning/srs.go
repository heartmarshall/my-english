package meaning

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetDueForReview возвращает meanings, которые нужно повторить (next_review_at < NOW()).
func (r *Repo) GetDueForReview(ctx context.Context, limit int) ([]model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	builder := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.NextReviewAt.Lt(r.clock.Now())).
		OrderBy(schema.Meanings.NextReviewAt.Asc()).
		Limit(uint64(limit))

	return database.NewQuery[model.Meaning](r.q, builder).List(ctx)
}

// GetByStatus возвращает meanings с указанным статусом обучения.
func (r *Repo) GetByStatus(ctx context.Context, status model.LearningStatus, limit int) ([]model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	builder := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(schema.Meanings.LearningStatus.Eq(status)).
		OrderBy(schema.Meanings.CreatedAt.Asc()).
		Limit(uint64(limit))

	return database.NewQuery[model.Meaning](r.q, builder).List(ctx)
}

// GetStudyQueue возвращает очередь для изучения.
func (r *Repo) GetStudyQueue(ctx context.Context, limit int) ([]model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	now := r.clock.Now()

	builder := database.Builder.
		Select(schema.Meanings.All()...).
		From(schema.Meanings.Name.String()).
		Where(squirrel.Or{
			schema.Meanings.LearningStatus.Eq(model.LearningStatusNew),
			schema.Meanings.NextReviewAt.Lt(now),
		}).
		OrderBy("COALESCE(" + schema.Meanings.NextReviewAt.String() + ", " + schema.Meanings.CreatedAt.String() + ") ASC").
		Limit(uint64(limit))

	return database.NewQuery[model.Meaning](r.q, builder).List(ctx)
}

// GetStats возвращает статистику.
// Используем SQL Aliases (as total_words), чтобы scany мог замапить колонки на поля структуры Stats.
func (r *Repo) GetStats(ctx context.Context) (*model.Stats, error) {
	now := r.clock.Now()

	// Для прямого SQL нужно создать обертку, реализующую SQLBuilder
	// Пока оставим как есть, так как это сложный случай с прямым SQL
	query := `
		SELECT 
			COUNT(DISTINCT ` + schema.Meanings.WordID.String() + `) as total_words,
			COUNT(*) FILTER (WHERE ` + schema.Meanings.LearningStatus.String() + ` = $1) as mastered_count,
			COUNT(*) FILTER (WHERE ` + schema.Meanings.LearningStatus.String() + ` = $2) as learning_count,
			COUNT(*) FILTER (WHERE ` + schema.Meanings.NextReviewAt.String() + ` < $3 OR ` + schema.Meanings.LearningStatus.String() + ` = $4) as due_for_review_count
		FROM ` + schema.Meanings.Name.String() + `
	`

	stats, err := database.GetOne[model.Stats](ctx, r.q, query,
		model.LearningStatusMastered,
		model.LearningStatusLearning,
		now,
		model.LearningStatusNew,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// UpdateSRS обновляет только SRS-поля meaning.
func (r *Repo) UpdateSRS(ctx context.Context, id int64, srs *SRSUpdate) error {
	if srs == nil {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	qb := database.Builder.
		Update(schema.Meanings.Name.String()).
		Set(schema.Meanings.LearningStatus.Bare(), srs.LearningStatus).
		Set(schema.Meanings.UpdatedAt.Bare(), now).
		Where(schema.Meanings.ID.Eq(id))

	if srs.NextReviewAt != nil {
		qb = qb.Set(schema.Meanings.NextReviewAt.Bare(), srs.NextReviewAt)
	}
	if srs.Interval != nil {
		qb = qb.Set(schema.Meanings.Interval.Bare(), srs.Interval)
	}
	if srs.EaseFactor != nil {
		qb = qb.Set(schema.Meanings.EaseFactor.Bare(), srs.EaseFactor)
	}
	if srs.ReviewCount != nil {
		qb = qb.Set(schema.Meanings.ReviewCount.Bare(), srs.ReviewCount)
	}

	rowsAffected, err := database.ExecOnly(ctx, r.q, qb)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}
