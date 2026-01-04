package lexeme

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database" // Импорт родителя
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// Repository работает с таблицей lexemes.
// Мы называем структуру просто Repository, так как она внутри пакета lexeme.
// Использование: lexeme.Repository
type Repository struct {
	*repository.Base[model.Lexeme]
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Lexeme](q, schema.Lexemes.Name.String(), schema.Lexemes.Columns()),
	}
}

// FindByTextNormalized ищет точное совпадение.
func (r *Repository) FindByTextNormalized(ctx context.Context, text string) (*model.Lexeme, error) {
	return r.FindOneBy(ctx, schema.Lexemes.TextNormalized.String(), text)
}

// SearchFuzzy выполняет нечеткий поиск.
func (r *Repository) SearchFuzzy(ctx context.Context, query string, limit int) ([]model.Lexeme, error) {
	similarityExpr := fmt.Sprintf("similarity(%s, ?)", schema.Lexemes.TextNormalized.String())

	builder := r.SelectBuilder().
		Where(squirrel.Expr(fmt.Sprintf("%s %% ?", schema.Lexemes.TextNormalized.String()), query)).
		OrderBy(similarityExpr+" DESC", schema.Lexemes.TextDisplay.String()).
		Limit(uint64(limit))

	return r.List(ctx, builder)
}

// CreateWithConflictIgnore создает лексему или возвращает существующую.
func (r *Repository) CreateWithConflictIgnore(ctx context.Context, lexeme *model.Lexeme) (*model.Lexeme, error) {
	existing, err := r.FindByTextNormalized(ctx, lexeme.TextNormalized)
	if err == nil {
		return existing, nil
	}
	if err != database.ErrNotFound {
		return nil, err
	}

	insert := r.InsertBuilder().
		Columns(schema.Lexemes.InsertColumns()...).
		Values(lexeme.TextNormalized, lexeme.TextDisplay).
		Suffix("ON CONFLICT (text_normalized) DO UPDATE SET text_display = EXCLUDED.text_display RETURNING *")

	return r.InsertReturning(ctx, insert)
}
