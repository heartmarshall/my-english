package meaning

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetDueForReview возвращает meanings, которые нужно повторить (next_review_at < NOW()).
func (r *Repo) GetDueForReview(ctx context.Context, limit int) ([]*model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Lt{"next_review_at": r.clock.Now()}).
		OrderBy("next_review_at ASC").
		Limit(uint64(limit)).
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

// GetByStatus возвращает meanings с указанным статусом обучения.
func (r *Repo) GetByStatus(ctx context.Context, status model.LearningStatus, limit int) ([]*model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"learning_status": status}).
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
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

// GetStudyQueue возвращает очередь для изучения:
// meanings со статусом NEW или next_review_at < NOW().
func (r *Repo) GetStudyQueue(ctx context.Context, limit int) ([]*model.Meaning, error) {
	limit = database.NormalizeLimit(limit, database.DefaultSRSLimit)

	now := r.clock.Now()

	query, args, err := database.Builder.
		Select(columns...).
		From(tableName).
		Where(squirrel.Or{
			squirrel.Eq{"learning_status": model.LearningStatusNew},
			squirrel.Lt{"next_review_at": now},
		}).
		OrderBy("COALESCE(next_review_at, created_at) ASC").
		Limit(uint64(limit)).
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

// GetStats возвращает статистику по изучению слов.
// Использует один эффективный запрос с FILTER.
func (r *Repo) GetStats(ctx context.Context) (*model.Stats, error) {
	now := r.clock.Now()

	// Используем raw SQL, так как squirrel не поддерживает FILTER синтаксис.
	const query = `
		SELECT 
			COUNT(DISTINCT word_id),
			COUNT(*) FILTER (WHERE learning_status = $1),
			COUNT(*) FILTER (WHERE learning_status = $2),
			COUNT(*) FILTER (WHERE next_review_at < $3 OR learning_status = $4)
		FROM meanings
	`

	var stats model.Stats
	err := r.q.QueryRow(ctx, query,
		model.LearningStatusMastered,
		model.LearningStatusLearning,
		now,
		model.LearningStatusNew,
	).Scan(
		&stats.TotalWords,
		&stats.MasteredCount,
		&stats.LearningCount,
		&stats.DueForReviewCount,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// UpdateSRS обновляет только SRS-поля meaning.
// Бизнес-логика расчёта новых значений должна быть в сервисном слое.
// Возвращает database.ErrInvalidInput если srs == nil.
// Возвращает database.ErrNotFound, если meaning не найден.
func (r *Repo) UpdateSRS(ctx context.Context, id int64, srs *SRSUpdate) error {
	if srs == nil {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	qb := database.Builder.
		Update(tableName).
		Set("learning_status", srs.LearningStatus).
		Set("updated_at", now).
		Where(squirrel.Eq{"id": id})

	// Опционально обновляем поля, если они заданы
	if srs.NextReviewAt != nil {
		qb = qb.Set("next_review_at", database.NullTime(srs.NextReviewAt))
	}
	if srs.Interval != nil {
		qb = qb.Set("interval", *srs.Interval)
	}
	if srs.EaseFactor != nil {
		qb = qb.Set("ease_factor", *srs.EaseFactor)
	}
	if srs.ReviewCount != nil {
		qb = qb.Set("review_count", *srs.ReviewCount)
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
