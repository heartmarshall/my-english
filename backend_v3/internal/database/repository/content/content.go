// Package content содержит репозитории для работы с контентом словаря:
// смыслы (senses), примеры (examples), переводы (translations),
// изображения (images) и произношения (pronunciations).
package content

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository/base"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// ============================================================================
// SENSES REPOSITORY
// ============================================================================

// SenseRepository предоставляет методы для работы со смыслами слов.
type SenseRepository struct {
	*base.Base[model.Sense]
}

// NewSenseRepository создаёт новый репозиторий смыслов.
func NewSenseRepository(q database.Querier) *SenseRepository {
	return &SenseRepository{
		Base: base.MustNewBase[model.Sense](q, base.Config{
			Table:   schema.Senses.Name.String(),
			Columns: schema.Senses.Columns(),
		}),
	}
}

// GetByID получает смысл по ID.
func (r *SenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Sense, error) {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return nil, err
	}
	return r.Base.GetByID(ctx, schema.Senses.ID.Bare(), id)
}

// ListByEntryIDs получает смыслы для списка записей словаря.
func (r *SenseRepository) ListByEntryIDs(ctx context.Context, entryIDs []uuid.UUID) ([]model.Sense, error) {
	if len(entryIDs) == 0 {
		return []model.Sense{}, nil
	}
	return r.ListByUUIDs(ctx, schema.Senses.EntryID.Bare(), entryIDs)
}

// Create создает новый смысл.
func (r *SenseRepository) Create(ctx context.Context, sense *model.Sense) (*model.Sense, error) {
	if sense == nil {
		return nil, fmt.Errorf("%w: sense is required", database.ErrInvalidInput)
	}
	if err := base.ValidateUUID(sense.EntryID, "entry_id"); err != nil {
		return nil, err
	}
	if err := base.ValidateString(sense.SourceSlug, "source_slug"); err != nil {
		return nil, err
	}

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
//
// Производительность:
//   - Использует batch insert для минимизации round-trips к БД
//   - Автоматически разбивает большие батчи на чанки
//   - Рекомендуется использовать в транзакциях для атомарности
func (r *SenseRepository) BatchCreate(ctx context.Context, senses []model.Sense) ([]model.Sense, error) {
	// Проверяем контекст перед выполнением
	if err := ctx.Err(); err != nil {
		return nil, database.WrapDBError(err)
	}

	if len(senses) == 0 {
		return []model.Sense{}, nil
	}

	// Валидация всех элементов перед вставкой
	// Это позволяет вернуть ошибку до начала транзакции
	for i := range senses {
		if err := base.ValidateUUID(senses[i].EntryID, "entry_id"); err != nil {
			return nil, fmt.Errorf("sense[%d]: %w", i, err)
		}
		if err := base.ValidateString(senses[i].SourceSlug, "source_slug"); err != nil {
			return nil, fmt.Errorf("sense[%d]: %w", i, err)
		}
	}

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

// Delete удаляет смысл.
func (r *SenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return err
	}
	return r.Base.Delete(ctx, schema.Senses.ID.Bare(), id)
}

// ============================================================================
// EXAMPLES REPOSITORY
// ============================================================================

// ExampleRepository предоставляет методы для работы с примерами.
type ExampleRepository struct {
	*base.Base[model.Example]
}

// NewExampleRepository создаёт новый репозиторий примеров.
func NewExampleRepository(q database.Querier) *ExampleRepository {
	return &ExampleRepository{
		Base: base.MustNewBase[model.Example](q, base.Config{
			Table:   schema.Examples.Name.String(),
			Columns: schema.Examples.Columns(),
		}),
	}
}

// GetByID получает пример по ID.
func (r *ExampleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Example, error) {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return nil, err
	}
	return r.Base.GetByID(ctx, schema.Examples.ID.Bare(), id)
}

// ListBySenseIDs получает примеры для списка смыслов.
func (r *ExampleRepository) ListBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.Example, error) {
	if len(senseIDs) == 0 {
		return []model.Example{}, nil
	}
	return r.ListByUUIDs(ctx, schema.Examples.SenseID.Bare(), senseIDs)
}

// BatchCreate создает несколько примеров за один запрос.
func (r *ExampleRepository) BatchCreate(ctx context.Context, examples []model.Example) ([]model.Example, error) {
	// Проверяем контекст перед выполнением
	if err := ctx.Err(); err != nil {
		return nil, database.WrapDBError(err)
	}

	if len(examples) == 0 {
		return []model.Example{}, nil
	}

	// Валидация всех элементов перед вставкой
	for i := range examples {
		if err := base.ValidateUUID(examples[i].SenseID, "sense_id"); err != nil {
			return nil, fmt.Errorf("example[%d]: %w", i, err)
		}
		if err := base.ValidateString(examples[i].Sentence, "sentence"); err != nil {
			return nil, fmt.Errorf("example[%d]: %w", i, err)
		}
		if err := base.ValidateString(examples[i].SourceSlug, "source_slug"); err != nil {
			return nil, fmt.Errorf("example[%d]: %w", i, err)
		}
	}

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

// Delete удаляет пример.
func (r *ExampleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return err
	}
	return r.Base.Delete(ctx, schema.Examples.ID.Bare(), id)
}

// ============================================================================
// TRANSLATIONS REPOSITORY
// ============================================================================

// TranslationRepository предоставляет методы для работы с переводами.
type TranslationRepository struct {
	*base.Base[model.Translation]
}

// NewTranslationRepository создаёт новый репозиторий переводов.
func NewTranslationRepository(q database.Querier) *TranslationRepository {
	return &TranslationRepository{
		Base: base.MustNewBase[model.Translation](q, base.Config{
			Table:   schema.Translations.Name.String(),
			Columns: schema.Translations.Columns(),
		}),
	}
}

// GetByID получает перевод по ID.
func (r *TranslationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Translation, error) {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return nil, err
	}
	return r.Base.GetByID(ctx, schema.Translations.ID.Bare(), id)
}

// ListBySenseIDs получает переводы для списка смыслов.
func (r *TranslationRepository) ListBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.Translation, error) {
	if len(senseIDs) == 0 {
		return []model.Translation{}, nil
	}
	return r.ListByUUIDs(ctx, schema.Translations.SenseID.Bare(), senseIDs)
}

// BatchCreate создает несколько переводов за один запрос.
func (r *TranslationRepository) BatchCreate(ctx context.Context, translations []model.Translation) ([]model.Translation, error) {
	// Проверяем контекст перед выполнением
	if err := ctx.Err(); err != nil {
		return nil, database.WrapDBError(err)
	}

	if len(translations) == 0 {
		return []model.Translation{}, nil
	}

	// Валидация всех элементов перед вставкой
	for i := range translations {
		if err := base.ValidateUUID(translations[i].SenseID, "sense_id"); err != nil {
			return nil, fmt.Errorf("translation[%d]: %w", i, err)
		}
		if err := base.ValidateString(translations[i].Text, "text"); err != nil {
			return nil, fmt.Errorf("translation[%d]: %w", i, err)
		}
		if err := base.ValidateString(translations[i].SourceSlug, "source_slug"); err != nil {
			return nil, fmt.Errorf("translation[%d]: %w", i, err)
		}
	}

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

// Delete удаляет перевод.
func (r *TranslationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return err
	}
	return r.Base.Delete(ctx, schema.Translations.ID.Bare(), id)
}

// ============================================================================
// IMAGES REPOSITORY
// ============================================================================

// ImageRepository предоставляет методы для работы с изображениями.
type ImageRepository struct {
	*base.Base[model.Image]
}

// NewImageRepository создаёт новый репозиторий изображений.
func NewImageRepository(q database.Querier) *ImageRepository {
	return &ImageRepository{
		Base: base.MustNewBase[model.Image](q, base.Config{
			Table:   schema.Images.Name.String(),
			Columns: schema.Images.Columns(),
		}),
	}
}

// GetByID получает изображение по ID.
func (r *ImageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Image, error) {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return nil, err
	}
	return r.Base.GetByID(ctx, schema.Images.ID.Bare(), id)
}

// ListByEntryIDs получает изображения для списка записей словаря.
func (r *ImageRepository) ListByEntryIDs(ctx context.Context, entryIDs []uuid.UUID) ([]model.Image, error) {
	if len(entryIDs) == 0 {
		return []model.Image{}, nil
	}
	return r.ListByUUIDs(ctx, schema.Images.EntryID.Bare(), entryIDs)
}

// BatchCreate создает несколько изображений за один запрос.
func (r *ImageRepository) BatchCreate(ctx context.Context, images []model.Image) ([]model.Image, error) {
	if len(images) == 0 {
		return []model.Image{}, nil
	}

	// Валидация всех элементов
	for i := range images {
		if err := base.ValidateUUID(images[i].EntryID, "entry_id"); err != nil {
			return nil, fmt.Errorf("image[%d]: %w", i, err)
		}
		if err := base.ValidateString(images[i].URL, "url"); err != nil {
			return nil, fmt.Errorf("image[%d]: %w", i, err)
		}
		if err := base.ValidateString(images[i].SourceSlug, "source_slug"); err != nil {
			return nil, fmt.Errorf("image[%d]: %w", i, err)
		}
	}

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

// Delete удаляет изображение.
func (r *ImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return err
	}
	return r.Base.Delete(ctx, schema.Images.ID.Bare(), id)
}

// ============================================================================
// PRONUNCIATIONS REPOSITORY
// ============================================================================

// PronunciationRepository предоставляет методы для работы с произношениями.
type PronunciationRepository struct {
	*base.Base[model.Pronunciation]
}

// NewPronunciationRepository создаёт новый репозиторий произношений.
func NewPronunciationRepository(q database.Querier) *PronunciationRepository {
	return &PronunciationRepository{
		Base: base.MustNewBase[model.Pronunciation](q, base.Config{
			Table:   schema.Pronunciations.Name.String(),
			Columns: schema.Pronunciations.Columns(),
		}),
	}
}

// GetByID получает произношение по ID.
func (r *PronunciationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Pronunciation, error) {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return nil, err
	}
	return r.Base.GetByID(ctx, schema.Pronunciations.ID.Bare(), id)
}

// ListByEntryIDs получает произношения для списка записей словаря.
func (r *PronunciationRepository) ListByEntryIDs(ctx context.Context, entryIDs []uuid.UUID) ([]model.Pronunciation, error) {
	if len(entryIDs) == 0 {
		return []model.Pronunciation{}, nil
	}
	return r.ListByUUIDs(ctx, schema.Pronunciations.EntryID.Bare(), entryIDs)
}

// BatchCreate создает несколько произношений за один запрос.
func (r *PronunciationRepository) BatchCreate(ctx context.Context, pronunciations []model.Pronunciation) ([]model.Pronunciation, error) {
	if len(pronunciations) == 0 {
		return []model.Pronunciation{}, nil
	}

	// Валидация всех элементов
	for i := range pronunciations {
		if err := base.ValidateUUID(pronunciations[i].EntryID, "entry_id"); err != nil {
			return nil, fmt.Errorf("pronunciation[%d]: %w", i, err)
		}
		if err := base.ValidateString(pronunciations[i].SourceSlug, "source_slug"); err != nil {
			return nil, fmt.Errorf("pronunciation[%d]: %w", i, err)
		}
	}

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

// Delete удаляет произношение.
func (r *PronunciationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := base.ValidateUUID(id, "id"); err != nil {
		return err
	}
	return r.Base.Delete(ctx, schema.Pronunciations.ID.Bare(), id)
}
