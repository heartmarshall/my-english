package word

import "github.com/heartmarshall/my-english/internal/database"

// repos — набор репозиториев для работы в транзакции.
type repos struct {
	words      WordRepository
	meanings   MeaningRepository
	examples   ExampleRepository
	tags       TagRepository
	meaningTag MeaningTagRepository
}

// withTx создаёт набор репозиториев, работающих в транзакции.
func (s *Service) withTx(tx database.Querier) *repos {
	return &repos{
		words:      s.repoFactory.Words(tx),
		meanings:   s.repoFactory.Meanings(tx),
		examples:   s.repoFactory.Examples(tx),
		tags:       s.repoFactory.Tags(tx),
		meaningTag: s.repoFactory.MeaningTags(tx),
	}
}
