package inbox

import (
	"context"
	"strings"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

// Service содержит бизнес-логику для работы с inbox.
type Service struct {
	inbox InboxRepository
}

// Deps — зависимости для создания сервиса.
type Deps struct {
	Inbox InboxRepository
}

// New создаёт новый сервис.
func New(deps Deps) *Service {
	return &Service{
		inbox: deps.Inbox,
	}
}

// Create создаёт новый inbox item.
func (s *Service) Create(ctx context.Context, text string, sourceContext *string) (*model.InboxItem, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, service.ErrInvalidInput
	}

	item := &model.InboxItem{
		Text:         text,
		SourceContext: sourceContext,
	}

	if err := s.inbox.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetByID возвращает inbox item по ID.
func (s *Service) GetByID(ctx context.Context, id int64) (*model.InboxItem, error) {
	item, err := s.inbox.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// List возвращает список всех inbox items.
func (s *Service) List(ctx context.Context) ([]model.InboxItem, error) {
	return s.inbox.List(ctx)
}

// Delete удаляет inbox item по ID.
func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.inbox.Delete(ctx, id)
}

