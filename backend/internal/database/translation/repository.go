package translation

import (
	"github.com/heartmarshall/my-english/internal/database"
)

// Repo — реализация репозитория для работы с translations.
type Repo struct {
	q database.Querier
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repo {
	return &Repo{q: q}
}
