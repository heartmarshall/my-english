package dictionary

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetFormsByWordID возвращает все формы слова по ID слова из словаря.
func (r *Repo) GetFormsByWordID(ctx context.Context, wordID int64) ([]model.DictionaryWordForm, error) {
	builder := database.Builder.
		Select(schema.DictionaryWordForms.All()...).
		From(schema.DictionaryWordForms.Name.String()).
		Where(schema.DictionaryWordForms.DictionaryWordID.Eq(wordID)).
		OrderBy(schema.DictionaryWordForms.FormType.Asc(), schema.DictionaryWordForms.FormText.Asc())

	return database.NewQuery[model.DictionaryWordForm](r.q, builder).List(ctx)
}

// GetFormsByWordIDs возвращает формы для нескольких слов (batch loading).
func (r *Repo) GetFormsByWordIDs(ctx context.Context, wordIDs []int64) ([]model.DictionaryWordForm, error) {
	if len(wordIDs) == 0 {
		return []model.DictionaryWordForm{}, nil
	}

	builder := database.Builder.
		Select(schema.DictionaryWordForms.All()...).
		From(schema.DictionaryWordForms.Name.String()).
		Where(schema.DictionaryWordForms.DictionaryWordID.In(wordIDs)).
		OrderBy(schema.DictionaryWordForms.DictionaryWordID.Asc(), schema.DictionaryWordForms.FormType.Asc(), schema.DictionaryWordForms.FormText.Asc())

	return database.NewQuery[model.DictionaryWordForm](r.q, builder).List(ctx)
}

// GetWordByFormText возвращает слово из словаря по тексту формы.
// Полезно для поиска основного слова по любой его форме (например, найти "go" по "went").
func (r *Repo) GetWordByFormText(ctx context.Context, formText string) (*model.DictionaryWord, error) {
	builder := database.Builder.
		Select(schema.DictionaryWords.All()...).
		From(schema.DictionaryWords.Name.String()).
		Join(schema.DictionaryWordForms.Name.String() + " ON " + schema.DictionaryWordForms.DictionaryWordID.Qualified() + " = " + schema.DictionaryWords.ID.Qualified()).
		Where(schema.DictionaryWordForms.FormText.Eq(formText)).
		Limit(1)

	word, err := database.NewQuery[model.DictionaryWord](r.q, builder).One(ctx)
	if err != nil {
		return nil, err
	}
	return &word, nil
}

// SearchSimilarForms выполняет поиск похожих форм слов с использованием триграмм.
func (r *Repo) SearchSimilarForms(ctx context.Context, query string, limit int, similarityThreshold float64) ([]model.DictionaryWordForm, error) {
	trigramCond := squirrel.Or{
		squirrel.Expr("word_similarity(?, ?) > ?", query, schema.DictionaryWordForms.FormText, similarityThreshold),
		squirrel.Expr("? % ?", schema.DictionaryWordForms.FormText, query),
	}

	// 1. Внутренний запрос с расчетом similarity
	innerBuilder := database.Builder.
		Select(schema.DictionaryWordForms.All()...).
		Column(squirrel.Expr("word_similarity(?, ?) AS similarity", query, schema.DictionaryWordForms.FormText)).
		From(schema.DictionaryWordForms.Name.String()).
		Where(trigramCond)

	// 2. Внешний запрос для чистой проекции и сортировки
	finalCols := make([]string, 0)
	for _, col := range schema.DictionaryWordForms.All() {
		finalCols = append(finalCols, "sub."+schema.Column(col).Bare())
	}

	outerBuilder := database.Builder.
		Select(finalCols...).
		FromSelect(innerBuilder, "sub").
		OrderBy("sub.similarity DESC").
		Limit(uint64(limit))

	return database.NewQuery[model.DictionaryWordForm](r.q, outerBuilder).List(ctx)
}

// CreateForm создаёт новую форму слова в словаре.
func (r *Repo) CreateForm(ctx context.Context, form *model.DictionaryWordForm) error {
	if form == nil {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()
	form.CreatedAt = now
	form.UpdatedAt = now

	builder := database.Builder.
		Insert(schema.DictionaryWordForms.Name.String()).
		Columns(schema.DictionaryWordForms.InsertColumns()...).
		Values(
			form.DictionaryWordID,
			form.FormText,
			form.FormType,
			form.CreatedAt,
			form.UpdatedAt,
		).
		Suffix("ON CONFLICT (dictionary_word_id, form_text, form_type) DO UPDATE SET updated_at = EXCLUDED.updated_at RETURNING " + schema.DictionaryWordForms.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		return err
	}

	form.ID = id
	return nil
}

// CreateForms создаёт несколько форм слова за один раз (batch insert).
func (r *Repo) CreateForms(ctx context.Context, forms []model.DictionaryWordForm) error {
	if len(forms) == 0 {
		return nil
	}

	now := r.clock.Now()
	builder := database.Builder.
		Insert(schema.DictionaryWordForms.Name.String()).
		Columns(schema.DictionaryWordForms.InsertColumns()...)

	for _, form := range forms {
		form.CreatedAt = now
		form.UpdatedAt = now
		builder = builder.Values(
			form.DictionaryWordID,
			form.FormText,
			form.FormType,
			form.CreatedAt,
			form.UpdatedAt,
		)
	}

	builder = builder.Suffix("ON CONFLICT (dictionary_word_id, form_text, form_type) DO UPDATE SET updated_at = EXCLUDED.updated_at")

	_, err := database.ExecOnly(ctx, r.q, builder)
	return err
}

// DeleteForm удаляет форму слова по ID.
func (r *Repo) DeleteForm(ctx context.Context, formID int64) error {
	builder := database.Builder.
		Delete(schema.DictionaryWordForms.Name.String()).
		Where(schema.DictionaryWordForms.ID.Eq(formID))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return database.ErrNotFound
	}
	return nil
}

// DeleteFormsByWordID удаляет все формы слова по ID слова.
func (r *Repo) DeleteFormsByWordID(ctx context.Context, wordID int64) error {
	builder := database.Builder.
		Delete(schema.DictionaryWordForms.Name.String()).
		Where(schema.DictionaryWordForms.DictionaryWordID.Eq(wordID))

	_, err := database.ExecOnly(ctx, r.q, builder)
	return err
}

