package sense

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.Sense]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Sense](q, schema.Senses.Name.String(), schema.Senses.Columns()),
	}
}

// GetByID возвращает смысл по UUID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Sense, error) {
	return r.Base.GetByID(ctx, schema.Senses.ID.String(), id)
}

// GetByIDs возвращает список смыслов по списку ID (для DataLoader'а card.sense).
func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Sense, error) {
	// Base.FindBy отлично работает со слайсами, превращая их в IN (...)
	return r.FindBy(ctx, schema.Senses.ID.String(), ids)
}

// ListByLexemeID возвращает список смыслов для конкретной лексемы.
func (r *Repository) ListByLexemeID(ctx context.Context, lexemeID uuid.UUID) ([]model.Sense, error) {
	return r.FindBy(ctx, schema.Senses.LexemeID.String(), lexemeID)
}

// ListByLexemeIDs загружает смыслы для списка лексем (Batch).
func (r *Repository) ListByLexemeIDs(ctx context.Context, lexemeIDs []uuid.UUID) ([]model.Sense, error) {
	return r.FindBy(ctx, schema.Senses.LexemeID.String(), lexemeIDs)
}

// Create создаёт новый смысл.
func (r *Repository) Create(ctx context.Context, sense *model.Sense) (*model.Sense, error) {
	insert := r.InsertBuilder().
		Columns(schema.Senses.InsertColumns()...).
		Values(
			sense.LexemeID,
			sense.PartOfSpeech,
			sense.Definition,
			sense.CefrLevel,
			sense.SourceID,
			sense.ExternalRefID,
		)

	return r.InsertReturning(ctx, insert)
}
