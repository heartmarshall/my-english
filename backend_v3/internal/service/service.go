package service

import (
	"fmt"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/service/dictionary"
	"github.com/heartmarshall/my-english/internal/service/inbox"
	"github.com/heartmarshall/my-english/internal/service/study"
	"github.com/heartmarshall/my-english/internal/service/suggestion"
)

// Services объединяет все сервисы приложения.
// Предоставляет единую точку доступа ко всем бизнес-сервисам.
type Services struct {
	Dictionary *dictionary.Service // Сервис для работы со словарем
	Inbox      *inbox.Service      // Сервис для работы с входящими заметками
	Study      *study.Service      // Сервис для работы с изучением карточек
	Suggestion *suggestion.Service // Сервис для получения подсказок из внешних источников
}

// Deps содержит зависимости, необходимые для создания сервисов.
type Deps struct {
	Repos     *repository.Registry  // Реестр репозиториев для доступа к данным
	TxManager *database.TxManager   // Менеджер транзакций для атомарных операций
	Providers []suggestion.Provider // Провайдеры подсказок из внешних источников
}

// NewServices инициализирует и возвращает все сервисы приложения.
// Сервисы создаются в правильном порядке с учетом зависимостей между ними.
// Возвращает ошибку, если не удалось создать сервисы.
func NewServices(deps Deps) (*Services, error) {
	if deps.Repos == nil {
		return nil, fmt.Errorf("repos cannot be nil")
	}
	if deps.TxManager == nil {
		return nil, fmt.Errorf("tx manager cannot be nil")
	}

	// Создаем сервисы в порядке зависимостей
	dictSvc, err := dictionary.NewService(deps.Repos, deps.TxManager)
	if err != nil {
		return nil, fmt.Errorf("create dictionary service: %w", err)
	}

	inboxSvc, err := inbox.NewService(deps.Repos, deps.TxManager, dictSvc)
	if err != nil {
		return nil, fmt.Errorf("create inbox service: %w", err)
	}

	studySvc, err := study.NewService(deps.Repos, deps.TxManager)
	if err != nil {
		return nil, fmt.Errorf("create study service: %w", err)
	}

	return &Services{
		Dictionary: dictSvc,
		Inbox:      inboxSvc,
		Study:      studySvc,
		Suggestion: suggestion.NewService(deps.Providers...),
	}, nil
}
