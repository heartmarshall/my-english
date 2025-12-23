package meaning

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// Create создаёт новое meaning в базе данных.
// Возвращает database.ErrInvalidInput если meaning == nil или обязательные поля пустые.
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
		Insert(tableName).
		Columns(
			"word_id", "part_of_speech", "definition_en", "translation_ru",
			"cefr_level", "image_url", "learning_status", "next_review_at",
			"interval", "ease_factor", "review_count", "created_at", "updated_at",
		).
		Values(
			meaning.WordID,
			meaning.PartOfSpeech,
			database.NullString(meaning.DefinitionEn),
			meaning.TranslationRu,
			database.NullString(meaning.CefrLevel),
			database.NullString(meaning.ImageURL),
			meaning.LearningStatus,
			database.NullTime(meaning.NextReviewAt),
			database.NullInt(meaning.Interval),
			database.NullFloat(meaning.EaseFactor),
			database.NullInt(meaning.ReviewCount),
			now,
			now,
		).
		Suffix("RETURNING id").
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
// Возвращает database.ErrInvalidInput если meaning == nil.
// Возвращает database.ErrNotFound, если meaning не найден.
func (r *Repo) Update(ctx context.Context, meaning *model.Meaning) error {
	if meaning == nil {
		return database.ErrInvalidInput
	}

	if meaning.TranslationRu == "" {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	query, args, err := database.Builder.
		Update(tableName).
		Set("word_id", meaning.WordID).
		Set("part_of_speech", meaning.PartOfSpeech).
		Set("definition_en", database.NullString(meaning.DefinitionEn)).
		Set("translation_ru", meaning.TranslationRu).
		Set("cefr_level", database.NullString(meaning.CefrLevel)).
		Set("image_url", database.NullString(meaning.ImageURL)).
		Set("updated_at", now).
		Where(squirrel.Eq{"id": meaning.ID}).
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
// Возвращает database.ErrNotFound, если meaning не найден.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"id": id}).
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
// Возвращает количество удалённых записей.
func (r *Repo) DeleteByWordID(ctx context.Context, wordID int64) (int64, error) {
	query, args, err := database.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"word_id": wordID}).
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
