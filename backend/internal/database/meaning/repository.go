package meaning

import (
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

var columns = []string{
	schema.MeaningColumns.ID.String(),
	schema.MeaningColumns.WordID.String(),
	schema.MeaningColumns.PartOfSpeech.String(),
	schema.MeaningColumns.DefinitionEn.String(),
	schema.MeaningColumns.TranslationRu.String(),
	schema.MeaningColumns.CefrLevel.String(),
	schema.MeaningColumns.ImageURL.String(),
	schema.MeaningColumns.LearningStatus.String(),
	schema.MeaningColumns.NextReviewAt.String(),
	schema.MeaningColumns.Interval.String(),
	schema.MeaningColumns.EaseFactor.String(),
	schema.MeaningColumns.ReviewCount.String(),
	schema.MeaningColumns.CreatedAt.String(),
	schema.MeaningColumns.UpdatedAt.String(),
}

// Filter содержит параметры фильтрации для запросов поиска.
type Filter struct {
	WordID         *int64
	PartOfSpeech   *model.PartOfSpeech
	LearningStatus *model.LearningStatus
}

// SRSUpdate содержит поля для обновления SRS данных.
type SRSUpdate struct {
	LearningStatus model.LearningStatus
	NextReviewAt   *time.Time
	Interval       *int
	EaseFactor     *float64
	ReviewCount    *int
}

// Repo — реализация репозитория для работы с meanings.
type Repo struct {
	q     database.Querier
	clock database.Clock
}

// Option — функциональная опция для конфигурации Repo.
type Option func(*Repo)

// WithClock устанавливает объект clock для работы с временем.
func WithClock(c database.Clock) Option {
	return func(r *Repo) {
		r.clock = c
	}
}

// New создаёт новый репозиторий.
func New(q database.Querier, opts ...Option) *Repo {
	r := &Repo{
		q:     q,
		clock: database.DefaultClock,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}
