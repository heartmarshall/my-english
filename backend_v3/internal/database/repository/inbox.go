package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type InboxRepository struct {
	*Base[model.InboxItem]
}

func NewInboxRepository(q database.Querier) *InboxRepository {
	return &InboxRepository{
		Base: NewBase[model.InboxItem](q, schema.InboxItems.Name.String(), schema.InboxItems.Columns()),
	}
}

func (r *InboxRepository) Create(ctx context.Context, item *model.InboxItem) (*model.InboxItem, error) {
	insert := r.InsertBuilder().
		Columns(schema.InboxItems.InsertColumns()...).
		Values(
			item.Text,
			item.Context,
		)
	return r.InsertReturning(ctx, insert)
}

func (r *InboxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.Base.Delete(ctx, schema.InboxItems.ID.String(), id)
}

func (r *InboxRepository) List(ctx context.Context, query squirrel.SelectBuilder) ([]model.InboxItem, error) {
	return r.Base.List(ctx, query)
}

func (r *InboxRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.InboxItem, error) {
	return r.Base.GetByID(ctx, schema.InboxItems.ID.String(), id)
}

func (r *InboxRepository) FindOneBy(ctx context.Context, column string, value any) (*model.InboxItem, error) {
	return r.Base.FindOneBy(ctx, column, value)
}
