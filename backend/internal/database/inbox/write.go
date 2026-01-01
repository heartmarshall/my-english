package inbox

import (
	"context"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// Create создаёт новый inbox item.
func (r *Repo) Create(ctx context.Context, item *model.InboxItem) error {
	if item == nil {
		return database.ErrInvalidInput
	}

	text := strings.TrimSpace(item.Text)
	if text == "" {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	builder := database.Builder.
		Insert(schema.InboxItems.Name.String()).
		Columns(schema.InboxItems.InsertColumns()...).
		Values(
			text,
			item.SourceContext,
			now,
		).
		Suffix("RETURNING " + schema.InboxItems.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		return err
	}

	item.ID = id

	item.CreatedAt = now
	return nil
}

// Delete удаляет inbox item по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	builder := database.Builder.
		Delete(schema.InboxItems.Name.String()).
		Where(schema.InboxItems.ID.Eq(id))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

