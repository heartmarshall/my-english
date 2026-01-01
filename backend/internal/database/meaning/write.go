package meaning

import (
	"context"

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

	builder := database.Builder.
		Insert(schema.Meanings.Name.String()).
		Columns(schema.Meanings.InsertColumns()...).
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
		Suffix("RETURNING " + schema.Meanings.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		return err
	}

	meaning.ID = id

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

	builder := database.Builder.
		Update(schema.Meanings.Name.String()).
		Set(schema.Meanings.WordID.Bare(), meaning.WordID).
		Set(schema.Meanings.PartOfSpeech.Bare(), meaning.PartOfSpeech).
		Set(schema.Meanings.DefinitionEn.Bare(), meaning.DefinitionEn). // Прямая передача
		Set(schema.Meanings.TranslationRu.Bare(), meaning.TranslationRu).
		Set(schema.Meanings.CefrLevel.Bare(), meaning.CefrLevel). // Прямая передача
		Set(schema.Meanings.ImageURL.Bare(), meaning.ImageURL).   // Прямая передача
		Set(schema.Meanings.UpdatedAt.Bare(), now).
		Where(schema.Meanings.ID.Eq(meaning.ID))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	meaning.UpdatedAt = now
	return nil
}

// Delete удаляет meaning по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	builder := database.Builder.
		Delete(schema.Meanings.Name.String()).
		Where(schema.Meanings.ID.Eq(id))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return database.ErrNotFound
	}
	return nil
}

// DeleteByWordID удаляет все meanings для указанного слова.
func (r *Repo) DeleteByWordID(ctx context.Context, wordID int64) (int64, error) {
	builder := database.Builder.
		Delete(schema.Meanings.Name.String()).
		Where(schema.Meanings.WordID.Eq(wordID))

	return database.ExecOnly(ctx, r.q, builder)
}
