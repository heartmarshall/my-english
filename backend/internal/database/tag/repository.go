package tag

import (
	"github.com/heartmarshall/my-english/internal/database"
)

const (
	tableName = "tags"
)

var columns = []string{
	"id",
	"name",
}

// Repo — реализация репозитория для работы с tags.
type Repo struct {
	q database.Querier
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repo {
	return &Repo{q: q}
}
