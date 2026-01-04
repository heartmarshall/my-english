package inbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	factory "github.com/heartmarshall/my-english/internal/database/repository/factory"
	"github.com/heartmarshall/my-english/internal/model"
)

type Service struct {
	repos     *factory.Factory
	txManager *database.TxManager
}

type Deps struct {
	Repos     *factory.Factory
	TxManager *database.TxManager
}

func New(deps Deps) *Service {
	return &Service{
		repos:     deps.Repos,
		txManager: deps.TxManager,
	}
}

// List возвращает список элементов inbox.
func (s *Service) List(ctx context.Context, limit int) ([]model.InboxItem, error) {
	return s.repos.Inbox(s.txManager.Q()).ListRecent(ctx, limit)
}

// GetByID возвращает элемент inbox по ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*model.InboxItem, error) {
	return s.repos.Inbox(s.txManager.Q()).GetByID(ctx, id)
}

// Create создаёт новый элемент inbox.
func (s *Service) Create(ctx context.Context, text string, contextNote *string) (*model.InboxItem, error) {
	item := &model.InboxItem{
		RawText:     text,
		ContextNote: contextNote,
	}
	return s.repos.Inbox(s.txManager.Q()).Create(ctx, item)
}

// Delete удаляет элемент inbox.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repos.Inbox(s.txManager.Q()).Delete(ctx, id)
}
