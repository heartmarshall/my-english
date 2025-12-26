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

	query, args, err := database.Builder.
		Insert(schema.InboxItems.Name.String()).
		Columns(schema.InboxItems.InsertColumns()...).
		Values(
			text,
			item.SourceContext,
			now,
		).
		Suffix("RETURNING " + schema.InboxItems.ID.Bare()).
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(&item.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	item.CreatedAt = now
	return nil
}

// Delete удаляет inbox item по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.
		Delete(schema.InboxItems.Name.String()).
		Where(schema.InboxItems.ID.Eq(id)).
		ToSql()
	if err != nil {
		return err
	}

	cmd, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	if cmd.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}

