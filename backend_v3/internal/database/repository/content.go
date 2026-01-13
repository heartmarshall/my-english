package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// --- SENSES ---

type SenseRepository struct {
	*Base[model.Sense]
}

func NewSenseRepository(q database.Querier) *SenseRepository {
	return &SenseRepository{
		Base: NewBase[model.Sense](q, schema.Senses.Name.String(), schema.Senses.Columns()),
	}
}

func (r *SenseRepository) ListByEntryIDs(ctx context.Context, entryIDs []uuid.UUID) ([]model.Sense, error) {
	ids := make([]any, len(entryIDs))
	for i, id := range entryIDs {
		ids[i] = id
	}
	return r.ListByIDs(ctx, schema.Senses.EntryID.String(), ids)
}

// Create создает новый смысл
func (r *SenseRepository) Create(ctx context.Context, sense *model.Sense) (*model.Sense, error) {
	insert := r.InsertBuilder().
		Columns(schema.Senses.InsertColumns()...).
		Values(
			sense.EntryID,
			sense.Definition,
			sense.PartOfSpeech,
			sense.SourceSlug,
			sense.CefrLevel,
		)

	return r.InsertReturning(ctx, insert)
}

// BatchCreate создает несколько смыслов за один запрос.
func (r *SenseRepository) BatchCreate(ctx context.Context, senses []model.Sense) ([]model.Sense, error) {
	columns := schema.Senses.InsertColumns()
	valuesFunc := func(s model.Sense) []any {
		return []any{
			s.EntryID,
			s.Definition,
			s.PartOfSpeech,
			s.SourceSlug,
			s.CefrLevel,
		}
	}
	return r.BatchInsertReturning(ctx, columns, senses, valuesFunc)
}

// --- EXAMPLES ---

type ExampleRepository struct {
	*Base[model.Example]
}

func NewExampleRepository(q database.Querier) *ExampleRepository {
	return &ExampleRepository{
		Base: NewBase[model.Example](q, schema.Examples.Name.String(), schema.Examples.Columns()),
	}
}

func (r *ExampleRepository) ListBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.Example, error) {
	ids := make([]any, len(senseIDs))
	for i, id := range senseIDs {
		ids[i] = id
	}
	return r.ListByIDs(ctx, schema.Examples.SenseID.String(), ids)
}

// BatchCreate создает несколько примеров за один запрос.
func (r *ExampleRepository) BatchCreate(ctx context.Context, examples []model.Example) ([]model.Example, error) {
	columns := schema.Examples.InsertColumns()
	valuesFunc := func(e model.Example) []any {
		return []any{
			e.SenseID,
			e.Sentence,
			e.Translation,
			e.SourceSlug,
		}
	}
	return r.BatchInsertReturning(ctx, columns, examples, valuesFunc)
}

// --- TRANSLATIONS ---

type TranslationRepository struct {
	*Base[model.Translation]
}

func NewTranslationRepository(q database.Querier) *TranslationRepository {
	return &TranslationRepository{
		Base: NewBase[model.Translation](q, schema.Translations.Name.String(), schema.Translations.Columns()),
	}
}

func (r *TranslationRepository) ListBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.Translation, error) {
	ids := make([]any, len(senseIDs))
	for i, id := range senseIDs {
		ids[i] = id
	}
	return r.ListByIDs(ctx, schema.Translations.SenseID.String(), ids)
}

// BatchCreate создает несколько переводов за один запрос.
func (r *TranslationRepository) BatchCreate(ctx context.Context, translations []model.Translation) ([]model.Translation, error) {
	columns := schema.Translations.InsertColumns()
	valuesFunc := func(t model.Translation) []any {
		return []any{
			t.SenseID,
			t.Text,
			t.SourceSlug,
		}
	}
	return r.BatchInsertReturning(ctx, columns, translations, valuesFunc)
}

// --- IMAGES ---

type ImageRepository struct {
	*Base[model.Image]
}

func NewImageRepository(q database.Querier) *ImageRepository {
	return &ImageRepository{
		Base: NewBase[model.Image](q, schema.Images.Name.String(), schema.Images.Columns()),
	}
}

func (r *ImageRepository) ListByEntryIDs(ctx context.Context, entryIDs []uuid.UUID) ([]model.Image, error) {
	ids := make([]any, len(entryIDs))
	for i, id := range entryIDs {
		ids[i] = id
	}
	return r.ListByIDs(ctx, schema.Images.EntryID.String(), ids)
}

// BatchCreate создает несколько изображений за один запрос.
func (r *ImageRepository) BatchCreate(ctx context.Context, images []model.Image) ([]model.Image, error) {
	columns := schema.Images.InsertColumns()
	valuesFunc := func(img model.Image) []any {
		return []any{
			img.EntryID,
			img.URL,
			img.Caption,
			img.SourceSlug,
		}
	}
	return r.BatchInsertReturning(ctx, columns, images, valuesFunc)
}

// --- PRONUNCIATIONS ---

type PronunciationRepository struct {
	*Base[model.Pronunciation]
}

func NewPronunciationRepository(q database.Querier) *PronunciationRepository {
	return &PronunciationRepository{
		Base: NewBase[model.Pronunciation](q, schema.Pronunciations.Name.String(), schema.Pronunciations.Columns()),
	}
}

func (r *PronunciationRepository) ListByEntryIDs(ctx context.Context, entryIDs []uuid.UUID) ([]model.Pronunciation, error) {
	ids := make([]any, len(entryIDs))
	for i, id := range entryIDs {
		ids[i] = id
	}
	return r.ListByIDs(ctx, schema.Pronunciations.EntryID.String(), ids)
}

// BatchCreate создает несколько произношений за один запрос.
func (r *PronunciationRepository) BatchCreate(ctx context.Context, pronunciations []model.Pronunciation) ([]model.Pronunciation, error) {
	columns := schema.Pronunciations.InsertColumns()
	valuesFunc := func(p model.Pronunciation) []any {
		return []any{
			p.EntryID,
			p.AudioURL,
			p.Transcription,
			p.Region,
			p.SourceSlug,
		}
	}
	return r.BatchInsertReturning(ctx, columns, pronunciations, valuesFunc)
}
