package example

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает example по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Example, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(schema.Examples.String()).
		Where(schema.ExampleColumns.ID.Eq(id)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.GetOne[model.Example](ctx, r.q, query, args...)
}

// GetByMeaningID возвращает все examples для указанного meaning.
func (r *Repo) GetByMeaningID(ctx context.Context, meaningID int64) ([]model.Example, error) {
	query, args, err := database.Builder.
		Select(columns...).
		From(schema.Examples.String()).
		Where(schema.ExampleColumns.MeaningID.Eq(meaningID)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Example](ctx, r.q, query, args...)
}

// GetByMeaningIDs возвращает examples для нескольких meanings (для batch loading).
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Example, error) {
	if len(meaningIDs) == 0 {
		return make([]model.Example, 0), nil
	}

	query, args, err := database.Builder.
		Select(columns...).
		From(schema.Examples.String()).
		Where(schema.ExampleColumns.MeaningID.In(meaningIDs)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Example](ctx, r.q, query, args...)
}
