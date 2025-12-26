package translation

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByMeaningID возвращает все переводы для meaning.
func (r *Repo) GetByMeaningID(ctx context.Context, meaningID int64) ([]model.Translation, error) {
	query, args, err := database.Builder.
		Select(schema.Translations.All()...).
		From(schema.Translations.Name.String()).
		Where(schema.Translations.MeaningID.Eq(meaningID)).
		OrderBy(schema.Translations.CreatedAt.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Translation](ctx, r.q, query, args...)
}

// GetByMeaningIDs возвращает переводы для нескольких meanings (batch loading).
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Translation, error) {
	if len(meaningIDs) == 0 {
		return []model.Translation{}, nil
	}

	query, args, err := database.Builder.
		Select(schema.Translations.All()...).
		From(schema.Translations.Name.String()).
		Where(schema.Translations.MeaningID.In(meaningIDs)).
		OrderBy(schema.Translations.MeaningID.Asc(), schema.Translations.CreatedAt.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Translation](ctx, r.q, query, args...)
}

// GetByID возвращает перевод по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (model.Translation, error) {
	query, args, err := database.Builder.
		Select(schema.Translations.All()...).
		From(schema.Translations.Name.String()).
		Where(schema.Translations.ID.Eq(id)).
		ToSql()
	if err != nil {
		return model.Translation{}, err
	}

	translation, err := database.GetOne[model.Translation](ctx, r.q, query, args...)
	if err != nil {
		return model.Translation{}, err
	}
	return *translation, nil
}
