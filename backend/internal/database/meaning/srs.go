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

	query, args, err := database.Builder.
		Select(columns...).
		From(schema.Meanings.String()).
		Where(schema.MeaningColumns.NextReviewAt.Lt(r.clock.Now())).
		OrderBy(schema.MeaningColumns.NextReviewAt.OrderByASC()).
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Meaning](ctx, r.q, query, args...)
}

// GetByStatus возвращает meanings с указанным статусом обучения.
func (r *Repo) GetByStatus(ctx context.Context, status model.LearningStatus, limit int) ([]model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	query, args, err := database.Builder.
		Select(columns...).
		From(schema.Meanings.String()).
		Where(schema.MeaningColumns.LearningStatus.Eq(status)).
		OrderBy(schema.MeaningColumns.CreatedAt.OrderByASC()).
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Meaning](ctx, r.q, query, args...)
}

// GetStudyQueue возвращает очередь для изучения.
func (r *Repo) GetStudyQueue(ctx context.Context, limit int) ([]model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	now := r.clock.Now()

	query, args, err := database.Builder.
		Select(columns...).
		From(schema.Meanings.String()).
		Where(squirrel.Or{
			schema.MeaningColumns.LearningStatus.Eq(model.LearningStatusNew),
			schema.MeaningColumns.NextReviewAt.Lt(now),
		}).
		OrderBy("COALESCE(" + schema.MeaningColumns.NextReviewAt.String() + ", " + schema.MeaningColumns.CreatedAt.String() + ") ASC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Meaning](ctx, r.q, query, args...)
}

// GetStats возвращает статистику.
// Используем SQL Aliases (as total_words), чтобы scany мог замапить колонки на поля структуры Stats.
func (r *Repo) GetStats(ctx context.Context) (*model.Stats, error) {
	now := r.clock.Now()

	query := `
		SELECT 
			COUNT(DISTINCT ` + schema.MeaningColumns.WordID.String() + `) as total_words,
			COUNT(*) FILTER (WHERE ` + schema.MeaningColumns.LearningStatus.String() + ` = $1) as mastered_count,
			COUNT(*) FILTER (WHERE ` + schema.MeaningColumns.LearningStatus.String() + ` = $2) as learning_count,
			COUNT(*) FILTER (WHERE ` + schema.MeaningColumns.NextReviewAt.String() + ` < $3 OR ` + schema.MeaningColumns.LearningStatus.String() + ` = $4) as due_for_review_count
		FROM ` + schema.Meanings.String() + `
	`

	return database.GetOne[model.Stats](ctx, r.q, query,
		model.LearningStatusMastered,
		model.LearningStatusLearning,
		now,
		model.LearningStatusNew,
	)
}

// UpdateSRS обновляет только SRS-поля meaning.
func (r *Repo) UpdateSRS(ctx context.Context, id int64, srs *SRSUpdate) error {
	if srs == nil {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	qb := database.Builder.
		Update(schema.Meanings.String()).
		Set(schema.MeaningColumns.LearningStatus.String(), srs.LearningStatus).
		Set(schema.MeaningColumns.UpdatedAt.String(), now).
		Where(squirrel.Eq{schema.MeaningColumns.ID.String(): id})

	if srs.NextReviewAt != nil {
		qb = qb.Set(schema.MeaningColumns.NextReviewAt.String(), srs.NextReviewAt)
	}
	if srs.Interval != nil {
		qb = qb.Set(schema.MeaningColumns.Interval.String(), srs.Interval)
	}
	if srs.EaseFactor != nil {
		qb = qb.Set(schema.MeaningColumns.EaseFactor.String(), srs.EaseFactor)
	}
	if srs.ReviewCount != nil {
		qb = qb.Set(schema.MeaningColumns.ReviewCount.String(), srs.ReviewCount)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}
