package datasource

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.DataSource]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.DataSource](q, schema.DataSources.Name.String(), schema.DataSources.Columns()),
	}
}

// GetByID возвращает источник по ID (int).
func (r *Repository) GetByID(ctx context.Context, id int) (*model.DataSource, error) {
	return r.Base.GetByID(ctx, schema.DataSources.ID.String(), id)
}

// GetBySlug возвращает источник по уникальному коду (например "freedict").
func (r *Repository) GetBySlug(ctx context.Context, slug string) (*model.DataSource, error) {
	return r.FindOneBy(ctx, schema.DataSources.Slug.String(), slug)
}
