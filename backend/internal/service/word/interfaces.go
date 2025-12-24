package word

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// WordRepository определяет интерфейс для работы со словами.
// Интерфейс определён здесь (у потребителя) согласно Go idiom.
type WordRepository interface {
	Create(ctx context.Context, word *model.Word) error
	GetByID(ctx context.Context, id int64) (model.Word, error)
	GetByText(ctx context.Context, text string) (model.Word, error)
	List(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]model.Word, error)
	Count(ctx context.Context, filter *model.WordFilter) (int, error)
	Update(ctx context.Context, word *model.Word) error
	Delete(ctx context.Context, id int64) error
}

// MeaningRepository определяет интерфейс для работы со значениями.
type MeaningRepository interface {
	Create(ctx context.Context, meaning *model.Meaning) error
	GetByID(ctx context.Context, id int64) (model.Meaning, error)
	GetByWordID(ctx context.Context, wordID int64) ([]model.Meaning, error)
	Update(ctx context.Context, meaning *model.Meaning) error
	Delete(ctx context.Context, id int64) error
	DeleteByWordID(ctx context.Context, wordID int64) (int64, error)
}

// ExampleRepository определяет интерфейс для работы с примерами.
type ExampleRepository interface {
	Create(ctx context.Context, example *model.Example) error
	CreateBatch(ctx context.Context, examples []*model.Example) error
	GetByMeaningID(ctx context.Context, meaningID int64) ([]model.Example, error)
	GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Example, error)
	DeleteByMeaningID(ctx context.Context, meaningID int64) (int64, error)
}

// TagRepository определяет интерфейс для работы с тегами.
type TagRepository interface {
	GetByName(ctx context.Context, name string) (model.Tag, error)
	GetByNames(ctx context.Context, names []string) ([]model.Tag, error)
	GetByIDs(ctx context.Context, ids []int64) ([]model.Tag, error)
	GetOrCreate(ctx context.Context, name string) (model.Tag, error)
}

// MeaningTagRepository определяет интерфейс для связи meaning-tag.
type MeaningTagRepository interface {
	AttachTags(ctx context.Context, meaningID int64, tagIDs []int64) error
	GetTagIDsByMeaningID(ctx context.Context, meaningID int64) ([]int64, error)
	GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.MeaningTag, error)
	SyncTags(ctx context.Context, meaningID int64, tagIDs []int64) error
	DetachAllFromMeaning(ctx context.Context, meaningID int64) error
}

// RepositoryFactory создаёт репозитории с указанным Querier.
// Используется для работы в транзакциях.
type RepositoryFactory interface {
	Words(q database.Querier) WordRepository
	Meanings(q database.Querier) MeaningRepository
	Examples(q database.Querier) ExampleRepository
	Tags(q database.Querier) TagRepository
	MeaningTags(q database.Querier) MeaningTagRepository
}

// TxRunner определяет интерфейс для запуска транзакций.
type TxRunner interface {
	// RunInTx выполняет функцию в рамках транзакции.
	// Querier передаётся в функцию для создания репозиториев.
	RunInTx(ctx context.Context, fn func(ctx context.Context, tx database.Querier) error) error
}
