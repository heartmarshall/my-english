package dictionary

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	repository "github.com/heartmarshall/my-english/internal/database/repository/factory"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

type Service struct {
	repos     *repository.Factory
	txManager *database.TxManager
	providers []Provider
}

type Deps struct {
	Repos     *repository.Factory
	TxManager *database.TxManager
	Providers []Provider
}

func New(deps Deps) *Service {
	return &Service{
		repos:     deps.Repos,
		txManager: deps.TxManager,
		providers: deps.Providers,
	}
}

// GetLexeme возвращает лексему по ID (без вложенных связей, для простоты).
// В реальном GraphQL резолвере связи будут грузиться через DataLoaders.
func (s *Service) GetLexeme(ctx context.Context, id uuid.UUID) (*model.Lexeme, error) {
	lexeme, err := s.repos.Lexeme(s.txManager.Q()).GetByID(ctx, "id", id)
	if err != nil {
		if database.IsNotFoundError(err) {
			return nil, service.ErrLexemeNotFound
		}
		return nil, err
	}
	return lexeme, nil
}
