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

	query, args, err := database.Builder.
		Insert(schema.Translations.Name.String()).
		Columns(schema.Translations.InsertColumns()...).
		Values(
			translation.MeaningID,
			translation.TranslationRu,
			translation.CreatedAt,
		).
		Suffix("ON CONFLICT (meaning_id, translation_ru) DO NOTHING RETURNING " + schema.Translations.ID.Bare()).
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

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	return nil
}

// DeleteByMeaningID удаляет все переводы для meaning.
func (r *Repo) DeleteByMeaningID(ctx context.Context, meaningID int64) error {
	query, args, err := database.Builder.
		Delete(schema.Translations.Name.String()).
		Where(schema.Translations.MeaningID.Eq(meaningID)).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	return nil
}

// Delete удаляет перевод по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.
		Delete(schema.Translations.Name.String()).
		Where(schema.Translations.ID.Eq(id)).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	return nil
}
