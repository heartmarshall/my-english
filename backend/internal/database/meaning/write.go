package meaning

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// Create создаёт новое meaning в базе данных.
func (r *Repo) Create(ctx context.Context, meaning *model.Meaning) error {
	if meaning == nil {
		return database.ErrInvalidInput
	}

	if meaning.WordID == 0 || meaning.TranslationRu == "" {
		return database.ErrInvalidInput
	}

	// Значения по умолчанию
	if meaning.LearningStatus == "" {
		meaning.LearningStatus = model.LearningStatusNew
	}
	if meaning.PartOfSpeech == "" {
		meaning.PartOfSpeech = model.PartOfSpeechOther
	}

	now := r.clock.Now()

	query, args, err := database.Builder.
		Insert(schema.Meanings.String()).
		Columns(
			schema.MeaningColumns.WordID.String(),
			schema.MeaningColumns.PartOfSpeech.String(),
			schema.MeaningColumns.DefinitionEn.String(),
			schema.MeaningColumns.TranslationRu.String(),
			schema.MeaningColumns.CefrLevel.String(),
			schema.MeaningColumns.ImageURL.String(),
			schema.MeaningColumns.LearningStatus.String(),
			schema.MeaningColumns.NextReviewAt.String(),
			schema.MeaningColumns.Interval.String(),
			schema.MeaningColumns.EaseFactor.String(),
			schema.MeaningColumns.ReviewCount.String(),
			schema.MeaningColumns.CreatedAt.String(),
			schema.MeaningColumns.UpdatedAt.String(),
		).
		Values(
			meaning.WordID,
			meaning.PartOfSpeech,
			meaning.DefinitionEn, // Прямая передача *string
			meaning.TranslationRu,
			meaning.CefrLevel, // Прямая передача *string
			meaning.ImageURL,  // Прямая передача *string
			meaning.LearningStatus,
			meaning.NextReviewAt, // Прямая передача *time.Time
			meaning.Interval,     // Прямая передача *int
			meaning.EaseFactor,   // Прямая передача *float64
			meaning.ReviewCount,  // Прямая передача *int
			now,
			now,
		).
		Suffix(schema.MeaningColumns.ID.Returning()).
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(&meaning.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	meaning.CreatedAt = now
	meaning.UpdatedAt = now
	return nil
}

// Update обновляет лингвистические поля meaning (не SRS).
func (r *Repo) Update(ctx context.Context, meaning *model.Meaning) error {
	if meaning == nil || meaning.TranslationRu == "" {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	query, args, err := database.Builder.
		Update(schema.Meanings.String()).
		Set(schema.MeaningColumns.WordID.String(), meaning.WordID).
		Set(schema.MeaningColumns.PartOfSpeech.String(), meaning.PartOfSpeech).
		Set(schema.MeaningColumns.DefinitionEn.String(), meaning.DefinitionEn). // Прямая передача
		Set(schema.MeaningColumns.TranslationRu.String(), meaning.TranslationRu).
		Set(schema.MeaningColumns.CefrLevel.String(), meaning.CefrLevel). // Прямая передача
		Set(schema.MeaningColumns.ImageURL.String(), meaning.ImageURL).   // Прямая передача
		Set(schema.MeaningColumns.UpdatedAt.String(), now).
		Where(squirrel.Eq{schema.MeaningColumns.ID.String(): meaning.ID}).
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	if commandTag.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	meaning.UpdatedAt = now
	return nil
}

// Delete удаляет meaning по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.
		Delete(schema.Meanings.String()).
		Where(squirrel.Eq{schema.MeaningColumns.ID.String(): id}).
		ToSql()
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

// DeleteByWordID удаляет все meanings для указанного слова.
func (r *Repo) DeleteByWordID(ctx context.Context, wordID int64) (int64, error) {
	query, args, err := database.Builder.
		Delete(schema.Meanings.String()).
		Where(squirrel.Eq{schema.MeaningColumns.WordID.String(): wordID}).
		ToSql()
	if err != nil {
		return 0, err
	}

	commandTag, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return commandTag.RowsAffected(), nil
}
