package cardtag

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// Repository для таблицы card_tags.
// Не наследуем Base полностью, так как тут составной PK и специфичная логика.
type Repository struct {
	q database.Querier
}

func New(q database.Querier) *Repository {
	return &Repository{q: q}
}

func (r *Repository) Attach(ctx context.Context, cardID uuid.UUID, tagID int) error {
	query := repository.Builder.
		Insert(schema.CardTags.Name.String()).
		Columns(schema.CardTags.InsertColumns()...).
		Values(cardID, tagID).
		Suffix("ON CONFLICT DO NOTHING")

	_, err := database.ExecOnly(ctx, r.q, query)
	return err
}

func (r *Repository) Detach(ctx context.Context, cardID uuid.UUID, tagID int) error {
	query := repository.Builder.
		Delete(schema.CardTags.Name.String()).
		Where(schema.CardTags.CardID.Eq(cardID)).
		Where(schema.CardTags.TagID.Eq(tagID))

	_, err := database.ExecOnly(ctx, r.q, query)
	return err
}

func (r *Repository) DetachAll(ctx context.Context, cardID uuid.UUID) error {
	query := repository.Builder.
		Delete(schema.CardTags.Name.String()).
		Where(schema.CardTags.CardID.Eq(cardID))

	_, err := database.ExecOnly(ctx, r.q, query)
	return err
}

// GetTagIDsByCardID возвращает ID тегов для карточки.
func (r *Repository) GetTagIDsByCardID(ctx context.Context, cardID uuid.UUID) ([]int, error) {
	query := repository.Builder.
		Select(schema.CardTags.TagID.String()).
		From(schema.CardTags.Name.String()).
		Where(schema.CardTags.CardID.Eq(cardID))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	return database.SelectScalars[int](ctx, r.q, sql, args...)
}

// GetTagsByCardID загружает полные модели тегов для карточки (JOIN).
func (r *Repository) GetTagsByCardID(ctx context.Context, cardID uuid.UUID) ([]model.Tag, error) {
	query := repository.Builder.
		Select(schema.Tags.Columns()...).
		From(schema.Tags.Name.String()).
		Join(schema.CardTags.Name.String() + " ON " + schema.Tags.ID.Qualified() + " = " + schema.CardTags.TagID.Qualified()).
		Where(schema.CardTags.CardID.Eq(cardID))

	// Используем helper из database package
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	return database.Select[model.Tag](ctx, r.q, sql, args...)
}

func (r *Repository) ListByCardIDs(ctx context.Context, cardIDs []uuid.UUID) ([]model.CardTag, error) {
	// squirrel.Eq с слайсом генерирует IN (...)
	query := repository.Builder.
		Select(schema.CardTags.Columns()...).
		From(schema.CardTags.Name.String()).
		Where(schema.CardTags.CardID.In(cardIDs))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	return database.Select[model.CardTag](ctx, r.q, sql, args...)
}
