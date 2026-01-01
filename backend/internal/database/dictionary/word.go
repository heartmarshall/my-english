package dictionary

import (
	"context"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetByText возвращает слово из словаря по тексту.
func (r *Repo) GetByText(ctx context.Context, text string) (model.DictionaryWord, error) {
	builder := database.Builder.
		Select(schema.DictionaryWords.All()...).
		From(schema.DictionaryWords.Name.String()).
		Where(schema.DictionaryWords.Text.Eq(text))

	return database.NewQuery[model.DictionaryWord](r.q, builder).One(ctx)
}

// SearchSimilar использует триграммный поиск для поиска похожих слов в словаре.
func (r *Repo) SearchSimilar(ctx context.Context, query string, limit int, similarityThreshold float64) ([]model.DictionaryWord, error) {
	if limit <= 0 {
		limit = 10
	}
	if similarityThreshold <= 0 {
		similarityThreshold = 0.3
	}

	selectCols := make([]string, 0, len(schema.DictionaryWords.All()))
	for _, col := range schema.DictionaryWords.All() {
		selectCols = append(selectCols, string(col))
	}

	sqlQuery := `
		SELECT ` + strings.Join(selectCols, ", ") + `
		FROM ` + schema.DictionaryWords.Name.String() + `
		WHERE word_similarity($1, ` + string(schema.DictionaryWords.Text) + `) > $2
		   OR ` + string(schema.DictionaryWords.Text) + ` % $1
		ORDER BY word_similarity($1, ` + string(schema.DictionaryWords.Text) + `) DESC
		LIMIT $3
	`

	args := []interface{}{query, similarityThreshold, limit}

	// Для прямого SQL нужно создать обертку, реализующую SQLBuilder
	// Пока используем старый метод, так как это прямой SQL
	words, err := database.Select[model.DictionaryWord](ctx, r.q, sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	return words, nil
}

// Create создаёт новое слово в словаре.
func (r *Repo) Create(ctx context.Context, word *model.DictionaryWord) error {
	if word == nil {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	builder := database.Builder.
		Insert(schema.DictionaryWords.Name.String()).
		Columns(schema.DictionaryWords.InsertColumns()...).
		Values(
			word.Text,
			word.Transcription,
			word.AudioURL,
			word.FrequencyRank,
			word.Source,
			word.SourceID,
			word.CreatedAt,
			word.UpdatedAt,
		).
		Suffix("ON CONFLICT (text) DO UPDATE SET updated_at = EXCLUDED.updated_at RETURNING " + schema.DictionaryWords.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		return err
	}

	word.ID = id

	return nil
}
