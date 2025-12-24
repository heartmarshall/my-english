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

	query, args, err := database.Builder.
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
		Suffix("RETURNING " + schema.Meanings.ID.Bare()).
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
		Update(schema.Meanings.Name.String()).
		Set(schema.Meanings.WordID.Bare(), meaning.WordID).
		Set(schema.Meanings.PartOfSpeech.Bare(), meaning.PartOfSpeech).
		Set(schema.Meanings.DefinitionEn.Bare(), meaning.DefinitionEn). // Прямая передача
		Set(schema.Meanings.TranslationRu.Bare(), meaning.TranslationRu).
		Set(schema.Meanings.CefrLevel.Bare(), meaning.CefrLevel). // Прямая передача
		Set(schema.Meanings.ImageURL.Bare(), meaning.ImageURL).   // Прямая передача
		Set(schema.Meanings.UpdatedAt.Bare(), now).
		Where(schema.Meanings.ID.Eq(meaning.ID)).
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
		Delete(schema.Meanings.Name.String()).
		Where(schema.Meanings.ID.Eq(id)).
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
		Delete(schema.Meanings.Name.String()).
		Where(schema.Meanings.WordID.Eq(wordID)).
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
