package inbox

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByID возвращает inbox item по ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (model.InboxItem, error) {
	query, args, err := database.Builder.
		Select(schema.InboxItems.All()...).
		From(schema.InboxItems.Name.String()).
		Where(schema.InboxItems.ID.Eq(id)).
		ToSql()
	if err != nil {
		return model.InboxItem{}, err
	}

	item, err := database.GetOne[model.InboxItem](ctx, r.q, query, args...)
	if err != nil {
		return model.InboxItem{}, err
	}
	return *item, nil
}

// List возвращает список всех inbox items, отсортированных по дате создания (новые первыми).
func (r *Repo) List(ctx context.Context) ([]model.InboxItem, error) {
	query, args, err := database.Builder.
		Select(schema.InboxItems.All()...).
		From(schema.InboxItems.Name.String()).
		OrderBy(schema.InboxItems.CreatedAt.Desc()).
		ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.InboxItem](ctx, r.q, query, args...)
}

