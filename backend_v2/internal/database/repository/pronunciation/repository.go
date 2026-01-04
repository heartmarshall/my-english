package pronunciation

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.Pronunciation]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Pronunciation](q, schema.Pronunciations.Name.String(), schema.Pronunciations.Columns()),
	}
}

func (r *Repository) ListByLexemeID(ctx context.Context, lexemeID uuid.UUID) ([]model.Pronunciation, error) {
	return r.FindBy(ctx, schema.Pronunciations.LexemeID.String(), lexemeID)
}

// ListByLexemeIDs batch-загрузка произношений.
func (r *Repository) ListByLexemeIDs(ctx context.Context, lexemeIDs []uuid.UUID) ([]model.Pronunciation, error) {
	return r.FindBy(ctx, schema.Pronunciations.LexemeID.String(), lexemeIDs)
}

func (r *Repository) Create(ctx context.Context, p *model.Pronunciation) (*model.Pronunciation, error) {
	insert := r.InsertBuilder().
		Columns(schema.Pronunciations.InsertColumns()...).
		Values(p.LexemeID, p.AudioURL, p.Transcription, p.Region, p.SourceID)

	return r.InsertReturning(ctx, insert)
}
