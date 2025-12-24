package example

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
)

var columns = []string{
	schema.ExampleColumns.ID.String(),
	schema.ExampleColumns.MeaningID.String(),
	schema.ExampleColumns.SentenceEn.String(),
	schema.ExampleColumns.SentenceRu.String(),
	schema.ExampleColumns.SourceName.String(),
}

// Repo — реализация репозитория для работы с examples.
type Repo struct {
	q database.Querier
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repo {
	return &Repo{q: q}
}
