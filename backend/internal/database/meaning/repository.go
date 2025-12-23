package meaning

import (
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

const (
	tableName = "meanings"
)

var columns = []string{
	"id",
	"word_id",
	"part_of_speech",
	"definition_en",
	"translation_ru",
	"cefr_level",
	"image_url",
	"learning_status",
	"next_review_at",
	"interval",
	"ease_factor",
	"review_count",
	"created_at",
	"updated_at",
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
