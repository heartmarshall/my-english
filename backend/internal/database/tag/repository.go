package tag

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
)

var columns = []string{
	schema.TagColumns.ID.String(),
	schema.TagColumns.Name.String(),
}

// Repo — реализация репозитория для работы с tags.
type Repo struct {
	q database.Querier
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repo {
	return &Repo{q: q}
}
