package inbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/dictionary"
)

type Service struct {
	repos      *repository.Registry
	tx         *database.TxManager
	dictionary *dictionary.Service // Зависимость от другого сервиса
}

// NewService создает сервис входящих.
// Обрати внимание: мы внедряем DictionaryService, чтобы переиспользовать сложную логику создания слова.
func NewService(repos *repository.Registry, tx *database.TxManager, dict *dictionary.Service) *Service {
	return &Service{
		repos:      repos,
		tx:         tx,
		dictionary: dict,
	}
}

// AddToInbox создает новую заметку.
func (s *Service) AddToInbox(ctx context.Context, text string, contextStr *string) (*model.InboxItem, error) {
	item := &model.InboxItem{
		Text:    text,
		Context: contextStr,
	}

	// Простая операция, можно без явной транзакции (авто-коммит),
	// но для единообразия и метрик лучше оставить как есть.
	return s.repos.Inbox.Create(ctx, item)
}

// Delete удаляет заметку.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repos.Inbox.Delete(ctx, id)
}

// List возвращает все входящие.
func (s *Service) List(ctx context.Context) ([]model.InboxItem, error) {
	// Для MVP возвращаем всё. В будущем добавим пагинацию.
	return s.repos.Inbox.List(ctx, s.repos.Inbox.SelectBuilder().OrderBy("created_at DESC"))
}

// ConvertToWord превращает InboxItem в DictionaryEntry.
func (s *Service) ConvertToWord(ctx context.Context, inboxID uuid.UUID, input dictionary.CreateWordInput) (*model.DictionaryEntry, error) {
	// 1. Проверяем, что inbox item существует
	_, err := s.repos.Inbox.GetByID(ctx, inboxID)
	if err != nil {
		return nil, err
	}

	// 2. Создаем слово через DictionaryService
	entry, err := s.dictionary.CreateWord(ctx, input)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
