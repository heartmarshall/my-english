package stats

import (
	"context"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	factory "github.com/heartmarshall/my-english/internal/database/repository/factory"
	"github.com/heartmarshall/my-english/internal/database/schema"
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

// GetDashboardStats возвращает статистику для дашборда.
func (s *Service) GetDashboardStats(ctx context.Context) (totalCards, masteredCount, learningCount, dueCount int, err error) {
	cardRepo := s.repos.Card(s.txManager.Q())
	srsRepo := s.repos.SRS(s.txManager.Q())

	// Общее количество активных карточек
	totalCards64, err := cardRepo.Count(ctx, repository.WithWhere(schema.Cards.IsDeleted.Eq(false)))
	if err != nil {
		return
	}
	totalCards = int(totalCards64)

	// Количество карточек со статусом MASTERED
	masteredCount64, err := srsRepo.Count(ctx, repository.WithWhere(schema.SRSStates.Status.Eq(model.LearningStatusMastered)))
	if err != nil {
		return
	}
	masteredCount = int(masteredCount64)

	// Количество карточек со статусом LEARNING
	learningCount64, err := srsRepo.Count(ctx, repository.WithWhere(schema.SRSStates.Status.Eq(model.LearningStatusLearning)))
	if err != nil {
		return
	}
	learningCount = int(learningCount64)

	// Количество карточек, которые нужно повторить (due_date <= now)
	now := time.Now()
	dueCount64, err := srsRepo.Count(ctx,
		repository.WithWhere(schema.SRSStates.DueDate.LtOrEq(now)),
		repository.WithWhere(schema.SRSStates.Status.NotEq(model.LearningStatusMastered)),
	)
	if err != nil {
		return
	}
	dueCount = int(dueCount64)

	return
}
