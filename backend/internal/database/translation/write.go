package translation

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// Create создаёт новый перевод.
func (r *Repo) Create(ctx context.Context, translation *model.Translation) error {
	if translation == nil {
		return database.ErrInvalidInput
	}

	builder := database.Builder.
		Insert(schema.Translations.Name.String()).
		Columns(schema.Translations.InsertColumns()...).
		Values(
			translation.MeaningID,
			translation.TranslationRu,
			translation.CreatedAt,
		).
		Suffix("ON CONFLICT (meaning_id, translation_ru) DO NOTHING RETURNING " + schema.Translations.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		// Если конфликт, возвращаем nil (ON CONFLICT DO NOTHING)
		if err.Error() == "no rows in result set" {
			return nil
		}
		return err
	}

	translation.ID = id

	return nil
}

// CreateBatch создаёт несколько переводов за раз.
func (r *Repo) CreateBatch(ctx context.Context, translations []*model.Translation) error {
	if len(translations) == 0 {
		return nil
	}

	qb := database.Builder.
		Insert(schema.Translations.Name.String()).
		Columns(schema.Translations.InsertColumns()...)

	for _, t := range translations {
		qb = qb.Values(t.MeaningID, t.TranslationRu, t.CreatedAt)
	}

	qb = qb.Suffix("ON CONFLICT (meaning_id, translation_ru) DO NOTHING")

	_, err := database.ExecOnly(ctx, r.q, qb)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByMeaningID удаляет все переводы для meaning.
func (r *Repo) DeleteByMeaningID(ctx context.Context, meaningID int64) error {
	builder := database.Builder.
		Delete(schema.Translations.Name.String()).
		Where(schema.Translations.MeaningID.Eq(meaningID))

	_, err := database.ExecOnly(ctx, r.q, builder)
	return err
}

// Delete удаляет перевод по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	builder := database.Builder.
		Delete(schema.Translations.Name.String()).
		Where(schema.Translations.ID.Eq(id))

	_, err := database.ExecOnly(ctx, r.q, builder)
	return err
}
