// Package word содержит репозиторий для работы с таблицей words.
package word

import (
	"github.com/heartmarshall/my-english/internal/database"
)

// Константы таблицы.
const (
	tableName = "words"
)

// columns — список колонок таблицы words.
var columns = []string{
	"id",
	"text",
	"transcription",
	"audio_url",
	"frequency_rank",
	"created_at",
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
