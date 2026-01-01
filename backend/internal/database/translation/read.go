package translation

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByMeaningID возвращает все переводы для meaning.
func (r *Repo) GetByMeaningID(ctx context.Context, meaningID int64) ([]model.Translation, error) {
	builder := database.Builder.
		Select(schema.Translations.All()...).
		From(schema.Translations.Name.String()).
		Where(schema.Translations.MeaningID.Eq(meaningID)).
		OrderBy(schema.Translations.CreatedAt.Asc())

	return database.NewQuery[model.Translation](r.q, builder).List(ctx)
}

// GetByMeaningIDs возвращает переводы для нескольких meanings (batch loading).
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Translation, error) {
	if len(meaningIDs) == 0 {
		return []model.Translation{}, nil
	}

	builder := database.Builder.
		Select(schema.Translations.All()...).
		From(schema.Translations.Name.String()).
		Where(schema.Translations.MeaningID.In(meaningIDs)).
		OrderBy(schema.Translations.MeaningID.Asc(), schema.Translations.CreatedAt.Asc())

	return database.NewQuery[model.Translation](r.q, builder).List(ctx)
}

// GetByID возвращает перевод по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (model.Translation, error) {
	builder := database.Builder.
		Select(schema.Translations.All()...).
		From(schema.Translations.Name.String()).
		Where(schema.Translations.ID.Eq(id))

	return database.NewQuery[model.Translation](r.q, builder).One(ctx)
}
