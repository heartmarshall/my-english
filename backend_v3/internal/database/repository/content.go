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
