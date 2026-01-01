package inbox

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает inbox item по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (model.InboxItem, error) {
	builder := database.Builder.
		Select(schema.InboxItems.All()...).
		From(schema.InboxItems.Name.String()).
		Where(schema.InboxItems.ID.Eq(id))

	return database.NewQuery[model.InboxItem](r.q, builder).One(ctx)
}

// List возвращает список всех inbox items, отсортированных по дате создания (новые первыми).
func (r *Repo) List(ctx context.Context) ([]model.InboxItem, error) {
	builder := database.Builder.
		Select(schema.InboxItems.All()...).
		From(schema.InboxItems.Name.String()).
		OrderBy(schema.InboxItems.CreatedAt.Desc())

	return database.NewQuery[model.InboxItem](r.q, builder).List(ctx)
}

