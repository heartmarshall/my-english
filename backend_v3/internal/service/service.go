package service

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/service/dictionary"
	"github.com/heartmarshall/my-english/internal/service/inbox"
	"github.com/heartmarshall/my-english/internal/service/study"
	"github.com/heartmarshall/my-english/internal/service/suggestion"
)

// Services объединяет все сервисы приложения.
type Services struct {
	Dictionary *dictionary.Service
	Inbox      *inbox.Service
	Study      *study.Service
	Suggestion *suggestion.Service
}

// Deps — зависимости для создания сервисов.
type Deps struct {
	Repos     *repository.Registry
	TxManager *database.TxManager
	Providers []suggestion.Provider
}

// NewServices инициализирует сервисный слой.
func NewServices(deps Deps) *Services {
	dictSvc := dictionary.NewService(deps.Repos, deps.TxManager)
	return &Services{
		Dictionary: dictSvc,
		Inbox:      inbox.NewService(deps.Repos, deps.TxManager, dictSvc),
		Study:      study.NewService(deps.Repos, deps.TxManager),
		Suggestion: suggestion.NewService(deps.Providers...),
	}
}
