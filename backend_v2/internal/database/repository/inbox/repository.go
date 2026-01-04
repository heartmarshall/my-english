package inbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// Repository работает с таблицей inbox_items.
type Repository struct {
	*repository.Base[model.InboxItem]
}

// New создаёт новый репозиторий inbox.
func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.InboxItem](q, schema.InboxItems.Name.String(), schema.InboxItems.Columns()),
	}
}

// GetByID находит элемент inbox по UUID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*model.InboxItem, error) {
	return r.Base.GetByID(ctx, schema.InboxItems.ID.String(), id)
}

// ListRecent возвращает последние элементы inbox, отсортированные по дате создания.
func (r *Repository) ListRecent(ctx context.Context, limit int) ([]model.InboxItem, error) {
	query := r.SelectBuilder().
		OrderBy(schema.InboxItems.CreatedAt.Desc())

	// Применяем лимит через функциональную опцию из родительского пакета repository
	query = repository.ApplyOptions(query, repository.WithLimit(limit))

	return r.List(ctx, query)
}

// Create создаёт новый элемент inbox.
// ID и CreatedAt генерируются базой данных (gen_random_uuid(), now()).
func (r *Repository) Create(ctx context.Context, item *model.InboxItem) (*model.InboxItem, error) {
	insert := r.InsertBuilder().
		Columns(schema.InboxItems.InsertColumns()...).
		Values(
			item.RawText,
			item.ContextNote,
		)

	return r.InsertReturning(ctx, insert)
}

// Delete удаляет элемент inbox по UUID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteByID(ctx, schema.InboxItems.ID.String(), id)
}
