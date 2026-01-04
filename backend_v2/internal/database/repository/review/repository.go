package review

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.ReviewLog]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.ReviewLog](q, schema.ReviewLogs.Name.String(), schema.ReviewLogs.Columns()),
	}
}

func (r *Repository) ListByCardID(ctx context.Context, cardID uuid.UUID, limit int) ([]model.ReviewLog, error) {
	query := r.SelectBuilder().
		Where(schema.ReviewLogs.CardID.Eq(cardID)).
		OrderBy(schema.ReviewLogs.ReviewedAt.Desc())

	query = repository.ApplyOptions(query, repository.WithLimit(limit))
	return r.List(ctx, query)
}

func (r *Repository) Create(ctx context.Context, log *model.ReviewLog) error {
	// ID генерируется БД (serial/identity), поэтому не возвращаем структуру целиком, если не надо
	insert := r.InsertBuilder().
		Columns(schema.ReviewLogs.InsertColumns()...).
		Values(
			log.CardID,
			log.Grade,
			log.DurationMs,
			log.StateBefore,
			log.StateAfter,
		)

	// Используем ExecOnly, так как нам обычно не нужен ID лога сразу
	_, err := database.ExecOnly(ctx, r.Q(), insert)
	return err
}
