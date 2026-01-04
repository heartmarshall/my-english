package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	factory "github.com/heartmarshall/my-english/internal/database/repository/factory"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
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

// GetByID возвращает карточку по ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*model.Card, error) {
	// Используем обычный пул (чтение)
	card, err := s.repos.Card(s.txManager.Q()).GetByID(ctx, id)
	if err != nil {
		if database.IsNotFoundError(err) {
			return nil, service.ErrCardNotFound
		}
		return nil, err
	}
	return card, nil
}

// List возвращает список активных карточек с пагинацией и фильтрацией.
func (s *Service) List(ctx context.Context, filter *Filter, limit, offset int) ([]model.Card, error) {
	opts := []repository.QueryOption{
		repository.WithPagination(limit, offset),
		repository.WithOrderByDesc(schema.Cards.CreatedAt.Qualified()), // Сортировка по дате создания (новые сначала)
	}

	// Если есть фильтры, используем метод с фильтрацией
	if filter != nil && (len(filter.Tags) > 0 || len(filter.Statuses) > 0 || filter.Search != nil) {
		return s.repos.Card(s.txManager.Q()).ListWithFilters(
			ctx,
			filter.Tags,
			filter.Statuses,
			filter.Search,
			opts...,
		)
	}

	// Иначе используем простой метод без фильтров
	return s.repos.Card(s.txManager.Q()).ListActive(ctx, opts...)
}
