package example

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.Example]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Example](q, schema.Examples.Name.String(), schema.Examples.Columns()),
	}
}

func (r *Repository) ListBySenseID(ctx context.Context, senseID uuid.UUID) ([]model.Example, error) {
	return r.FindBy(ctx, schema.Examples.SenseID.String(), senseID)
}

// ListBySenseIDs batch-загрузка примеров.
func (r *Repository) ListBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.Example, error) {
	return r.FindBy(ctx, schema.Examples.SenseID.String(), senseIDs)
}

func (r *Repository) Create(ctx context.Context, ex *model.Example) (*model.Example, error) {
	insert := r.InsertBuilder().
		Columns(schema.Examples.InsertColumns()...).
		Values(ex.SenseID, ex.SentenceEn, ex.SentenceRu, ex.TargetWordRange, ex.SourceName)

	return r.InsertReturning(ctx, insert)
}
