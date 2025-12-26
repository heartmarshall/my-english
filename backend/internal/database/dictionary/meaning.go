package dictionary

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByWordID возвращает все значения для слова из словаря.
func (r *Repo) GetMeaningsByWordID(ctx context.Context, wordID int64) ([]model.DictionaryMeaning, error) {
	query, args, err := database.Builder.
		Select(schema.DictionaryMeanings.All()...).
		From(schema.DictionaryMeanings.Name.String()).
		Where(schema.DictionaryMeanings.DictionaryWordID.Eq(wordID)).
		OrderBy(schema.DictionaryMeanings.OrderIndex.Asc(), schema.DictionaryMeanings.CreatedAt.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.DictionaryMeaning](ctx, r.q, query, args...)
}

// GetTranslationsByMeaningID возвращает все переводы для значения из словаря.
func (r *Repo) GetTranslationsByMeaningID(ctx context.Context, meaningID int64) ([]model.DictionaryTranslation, error) {
	query, args, err := database.Builder.
		Select(schema.DictionaryTranslations.All()...).
		From(schema.DictionaryTranslations.Name.String()).
		Where(schema.DictionaryTranslations.DictionaryMeaningID.Eq(meaningID)).
		OrderBy(schema.DictionaryTranslations.CreatedAt.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.DictionaryTranslation](ctx, r.q, query, args...)
}

// GetTranslationsByMeaningIDs возвращает переводы для нескольких значений (batch loading).
func (r *Repo) GetTranslationsByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.DictionaryTranslation, error) {
	if len(meaningIDs) == 0 {
		return []model.DictionaryTranslation{}, nil
	}

	query, args, err := database.Builder.
		Select(schema.DictionaryTranslations.All()...).
		From(schema.DictionaryTranslations.Name.String()).
		Where(schema.DictionaryTranslations.DictionaryMeaningID.In(meaningIDs)).
		OrderBy(schema.DictionaryTranslations.DictionaryMeaningID.Asc(), schema.DictionaryTranslations.CreatedAt.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.DictionaryTranslation](ctx, r.q, query, args...)
}

// CreateMeaning создаёт новое значение в словаре.
func (r *Repo) CreateMeaning(ctx context.Context, meaning *model.DictionaryMeaning) error {
	if meaning == nil {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()
	meaning.CreatedAt = now
	meaning.UpdatedAt = now

	query, args, err := database.Builder.
		Insert(schema.DictionaryMeanings.Name.String()).
		Columns(schema.DictionaryMeanings.InsertColumns()...).
		Values(
			meaning.DictionaryWordID,
			meaning.PartOfSpeech,
			meaning.DefinitionEn,
			meaning.CefrLevel,
			meaning.ImageURL,
			meaning.OrderIndex,
			meaning.CreatedAt,
			meaning.UpdatedAt,
		).
		Suffix("RETURNING " + schema.DictionaryMeanings.ID.Bare()).
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(&meaning.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	return nil
}

// CreateTranslation создаёт новый перевод для значения из словаря.
func (r *Repo) CreateTranslation(ctx context.Context, translation *model.DictionaryTranslation) error {
	if translation == nil {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()
	translation.CreatedAt = now

	query, args, err := database.Builder.
		Insert(schema.DictionaryTranslations.Name.String()).
		Columns(schema.DictionaryTranslations.InsertColumns()...).
		Values(
			translation.DictionaryMeaningID,
			translation.TranslationRu,
			translation.CreatedAt,
		).
		Suffix("ON CONFLICT (dictionary_meaning_id, translation_ru) DO NOTHING RETURNING " + schema.DictionaryTranslations.ID.Bare()).
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(&translation.ID)
	if err != nil {
		// Если конфликт, возвращаем nil (ON CONFLICT DO NOTHING)
		if err.Error() == "no rows in result set" {
			return nil
		}
		return database.WrapDBError(err)
	}

	return nil
}

