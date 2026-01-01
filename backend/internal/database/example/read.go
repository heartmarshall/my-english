package example

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает example по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Example, error) {
	builder := database.Builder.
		Select(schema.Examples.All()...).
		From(schema.Examples.Name.String()).
		Where(schema.Examples.ID.Eq(id))

	example, err := database.NewQuery[model.Example](r.q, builder).One(ctx)
	if err != nil {
		return nil, err
	}
	return &example, nil
}

// GetByMeaningID возвращает все examples для указанного meaning.
func (r *Repo) GetByMeaningID(ctx context.Context, meaningID int64) ([]model.Example, error) {
	builder := database.Builder.
		Select(schema.Examples.All()...).
		From(schema.Examples.Name.String()).
		Where(schema.Examples.MeaningID.Eq(meaningID))

	return database.NewQuery[model.Example](r.q, builder).List(ctx)
}

// GetByMeaningIDs возвращает examples для нескольких meanings (для batch loading).
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Example, error) {
	if len(meaningIDs) == 0 {
		return make([]model.Example, 0), nil
	}

	builder := database.Builder.
		Select(schema.Examples.All()...).
		From(schema.Examples.Name.String()).
		Where(schema.Examples.MeaningID.In(meaningIDs))

	return database.NewQuery[model.Example](r.q, builder).List(ctx)
}
