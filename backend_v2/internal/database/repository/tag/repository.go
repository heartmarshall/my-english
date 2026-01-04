package tag

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.Tag]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Tag](q, schema.Tags.Name.String(), schema.Tags.Columns()),
	}
}

func (r *Repository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	return r.FindOneBy(ctx, schema.Tags.NameCol.String(), name)
}

func (r *Repository) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	insert := r.InsertBuilder().
		Columns(schema.Tags.InsertColumns()...).
		Values(tag.Name, tag.ColorHex)

	return r.InsertReturning(ctx, insert)
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int) ([]model.Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return r.FindBy(ctx, schema.Tags.ID.String(), ids)
}
