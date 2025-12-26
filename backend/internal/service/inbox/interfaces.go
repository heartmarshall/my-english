package inbox

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
)

// InboxRepository определяет интерфейс для работы с inbox_items.
type InboxRepository interface {
	Create(ctx context.Context, item *model.InboxItem) error
	GetByID(ctx context.Context, id int64) (model.InboxItem, error)
	List(ctx context.Context) ([]model.InboxItem, error)
	Delete(ctx context.Context, id int64) error
}

