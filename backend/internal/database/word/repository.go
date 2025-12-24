// Package word содержит репозиторий для работы с таблицей words.
package word

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
)

// columns — список колонок таблицы words.
var columns = []string{
	schema.WordColumns.ID.String(),
	schema.WordColumns.Text.String(),
	schema.WordColumns.Transcription.String(),
	schema.WordColumns.AudioURL.String(),
	schema.WordColumns.FrequencyRank.String(),
	schema.WordColumns.CreatedAt.String(),
}

// Repo — реализация репозитория для PostgreSQL.
type Repo struct {
	q     database.Querier
	clock database.Clock
}

// Option — функциональная опция для конфигурации Repo.
type Option func(*Repo)

// WithClock устанавливает кастомный clock (полезно для тестов).
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
