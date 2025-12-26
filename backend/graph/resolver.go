package graph

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/word"
)

//go:generate go run github.com/99designs/gqlgen generate

// WordService определяет интерфейс для работы со словами.
type WordService interface {
	Create(ctx context.Context, input word.CreateWordInput) (*word.WordWithRelations, error)
	GetByID(ctx context.Context, id int64) (*word.WordWithRelations, error)
	GetByText(ctx context.Context, text string) (*word.WordWithRelations, error)
	List(ctx context.Context, filter *word.WordFilter, limit, offset int) ([]model.Word, error)
	Count(ctx context.Context, filter *word.WordFilter) (int, error)
	Update(ctx context.Context, id int64, input word.UpdateWordInput) (*word.WordWithRelations, error)
	Delete(ctx context.Context, id int64) error
	Suggest(ctx context.Context, query string) ([]word.Suggestion, error)
}

// StudyService определяет интерфейс для системы изучения.
type StudyService interface {
	GetStudyQueue(ctx context.Context, limit int) ([]model.Meaning, error)
	GetStats(ctx context.Context) (*model.Stats, error)
	ReviewMeaning(ctx context.Context, meaningID int64, grade int) (model.Meaning, error)
}

// InboxService определяет интерфейс для работы с inbox.
type InboxService interface {
	Create(ctx context.Context, text string, sourceContext *string) (*model.InboxItem, error)
	GetByID(ctx context.Context, id int64) (*model.InboxItem, error)
	List(ctx context.Context) ([]model.InboxItem, error)
	Delete(ctx context.Context, id int64) error
}

// Deps — зависимости для резолвера.
// Resolver использует только сервисы, не репозитории.
type Deps struct {
	Words WordService
	Study StudyService
	Inbox InboxService
}

// Resolver — корневой резолвер GraphQL.
// Использует только сервисы. DataLoaders инжектируются через middleware.
type Resolver struct {
	words WordService
	study StudyService
	inbox InboxService
}

// NewResolver создаёт новый резолвер с зависимостями.
func NewResolver(deps Deps) *Resolver {
	return &Resolver{
		words: deps.Words,
		study: deps.Study,
		inbox: deps.Inbox,
	}
}
