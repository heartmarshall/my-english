package sensetranslation

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.SenseTranslation]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.SenseTranslation](q, schema.SenseTranslations.Name.String(), schema.SenseTranslations.Columns()),
	}
}

// ListBySenseID возвращает список переводов для конкретного смысла.
func (r *Repository) ListBySenseID(ctx context.Context, senseID uuid.UUID) ([]model.SenseTranslation, error) {
	return r.FindBy(ctx, schema.SenseTranslations.SenseID.String(), senseID)
}

// Create создаёт новый перевод.
func (r *Repository) Create(ctx context.Context, t *model.SenseTranslation) (*model.SenseTranslation, error) {
	insert := r.InsertBuilder().
		Columns(schema.SenseTranslations.InsertColumns()...).
		Values(t.SenseID, t.Translation, t.SourceID)

	return r.InsertReturning(ctx, insert)
}

func (r *Repository) ListBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.SenseTranslation, error) {
	return r.FindBy(ctx, schema.SenseTranslations.SenseID.String(), senseIDs)
}
