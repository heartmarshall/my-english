package tag

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает tag по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	query, args, err := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.ID.Eq(id)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.GetOne[model.Tag](ctx, r.q, query, args...)
}

// GetByName возвращает tag по имени.
func (r *Repo) GetByName(ctx context.Context, name string) (model.Tag, error) {
	query, args, err := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.NameCol.Eq(name)).
		ToSql()
	if err != nil {
		return model.Tag{}, err
	}

	tag, err := database.GetOne[model.Tag](ctx, r.q, query, args...)
	if err != nil {
		return model.Tag{}, err
	}
	return *tag, nil
}

// GetByNames возвращает tags по списку имён.
func (r *Repo) GetByNames(ctx context.Context, names []string) ([]model.Tag, error) {
	if len(names) == 0 {
		return make([]model.Tag, 0), nil
	}

	query, args, err := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.NameCol.In(names)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Tag](ctx, r.q, query, args...)
}

// GetByIDs возвращает tags по списку ID.
func (r *Repo) GetByIDs(ctx context.Context, ids []int64) ([]model.Tag, error) {
	if len(ids) == 0 {
		return make([]model.Tag, 0), nil
	}

	query, args, err := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		Where(schema.Tags.ID.In(ids)).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Tag](ctx, r.q, query, args...)
}

// List возвращает все tags.
func (r *Repo) List(ctx context.Context) ([]model.Tag, error) {
	query, args, err := database.Builder.
		Select(schema.Tags.All()...).
		From(schema.Tags.Name.String()).
		OrderBy(schema.Tags.NameCol.Asc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Tag](ctx, r.q, query, args...)
}
