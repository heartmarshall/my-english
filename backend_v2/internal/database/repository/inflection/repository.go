package inflection

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.Inflection]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Inflection](q, schema.Inflections.Name.String(), schema.Inflections.Columns()),
	}
}

// GetFormsByLemmaID возвращает все формы для базового слова (Go -> Went, Gone).
func (r *Repository) GetFormsByLemmaID(ctx context.Context, lemmaID uuid.UUID) ([]model.Inflection, error) {
	return r.FindBy(ctx, schema.Inflections.LemmaLexemeID.String(), lemmaID)
}

// GetLemmaByInflectedID возвращает базовое слово для формы (Went -> Go).
// Обычно форма ссылается только на одну лемму, но возвращаем список для универсальности.
func (r *Repository) GetLemmaByInflectedID(ctx context.Context, inflectedID uuid.UUID) ([]model.Inflection, error) {
	return r.FindBy(ctx, schema.Inflections.InflectedLexemeID.String(), inflectedID)
}

// Create создаёт связь между леммой и формой.
func (r *Repository) Create(ctx context.Context, i *model.Inflection) (*model.Inflection, error) {
	insert := r.InsertBuilder().
		Columns(schema.Inflections.InsertColumns()...).
		Values(i.InflectedLexemeID, i.LemmaLexemeID, i.Type).
		Suffix("ON CONFLICT DO NOTHING RETURNING *") // Игнорируем дубли

	return r.InsertReturning(ctx, insert)
}
