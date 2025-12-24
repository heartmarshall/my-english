// Package loader предоставляет сервис для batch-загрузки данных.
// Используется DataLoader'ами в транспортном слое.
package loader

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
)

// MeaningRepository определяет интерфейс для загрузки meanings.
type MeaningRepository interface {
	GetByWordIDs(ctx context.Context, wordIDs []int64) ([]model.Meaning, error)
}

// ExampleRepository определяет интерфейс для загрузки examples.
type ExampleRepository interface {
	GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Example, error)
}

// TagRepository определяет интерфейс для загрузки tags.
type TagRepository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]model.Tag, error)
}

// MeaningTagRepository определяет интерфейс для связей meaning-tag.
type MeaningTagRepository interface {
	GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.MeaningTag, error)
}

// Deps — зависимости сервиса.
type Deps struct {
	Meanings    MeaningRepository
	Examples    ExampleRepository
	Tags        TagRepository
	MeaningTags MeaningTagRepository
}

// Service предоставляет batch-операции для DataLoader.
type Service struct {
	meanings    MeaningRepository
	examples    ExampleRepository
	tags        TagRepository
	meaningTags MeaningTagRepository
}

// New создаёт новый сервис.
func New(deps Deps) *Service {
	return &Service{
		meanings:    deps.Meanings,
		examples:    deps.Examples,
		tags:        deps.Tags,
		meaningTags: deps.MeaningTags,
	}
}

// GetMeaningsByWordIDs загружает meanings для нескольких слов.
func (s *Service) GetMeaningsByWordIDs(ctx context.Context, wordIDs []int64) ([]model.Meaning, error) {
	return s.meanings.GetByWordIDs(ctx, wordIDs)
}

// GetExamplesByMeaningIDs загружает examples для нескольких meanings.
func (s *Service) GetExamplesByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Example, error) {
	return s.examples.GetByMeaningIDs(ctx, meaningIDs)
}

// GetTagsByMeaningIDs загружает теги для нескольких meanings.
// Возвращает связи MeaningTag для группировки.
func (s *Service) GetTagsByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.MeaningTag, error) {
	return s.meaningTags.GetByMeaningIDs(ctx, meaningIDs)
}

// GetTagsByIDs загружает теги по ID.
func (s *Service) GetTagsByIDs(ctx context.Context, ids []int64) ([]model.Tag, error) {
	return s.tags.GetByIDs(ctx, ids)
}
