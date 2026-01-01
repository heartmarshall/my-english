package tag

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает tag по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	builder := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.ID.Eq(id))

	tag, err := database.NewQuery[model.Tag](r.q, builder).One(ctx)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetByName возвращает tag по имени.
func (r *Repo) GetByName(ctx context.Context, name string) (model.Tag, error) {
	builder := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.NameCol.Eq(name))

	return database.NewQuery[model.Tag](r.q, builder).One(ctx)
}

// GetByNames возвращает tags по списку имён.
func (r *Repo) GetByNames(ctx context.Context, names []string) ([]model.Tag, error) {
	if len(names) == 0 {
		return make([]model.Tag, 0), nil
	}

	builder := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.NameCol.In(names))

	return database.NewQuery[model.Tag](r.q, builder).List(ctx)
}

// GetByIDs возвращает tags по списку ID.
func (r *Repo) GetByIDs(ctx context.Context, ids []int64) ([]model.Tag, error) {
	if len(ids) == 0 {
		return make([]model.Tag, 0), nil
	}

	builder := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.ID.In(ids))

	return database.NewQuery[model.Tag](r.q, builder).List(ctx)
}

// List возвращает все tags.
func (r *Repo) List(ctx context.Context) ([]model.Tag, error) {
	builder := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		OrderBy(schema.Tags.NameCol.Asc())

	return database.NewQuery[model.Tag](r.q, builder).List(ctx)
}
