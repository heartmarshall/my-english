package example

import (
	"github.com/heartmarshall/my-english/internal/database"
)

const (
	tableName = "examples"
)

var columns = []string{
	"id",
	"meaning_id",
	"sentence_en",
	"sentence_ru",
	"source_name",
}

// Repo — реализация репозитория для работы с examples.
type Repo struct {
	q database.Querier
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repo {
	return &Repo{q: q}
}
