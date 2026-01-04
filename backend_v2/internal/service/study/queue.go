package study

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetQueue возвращает список карточек для изучения с опциональной фильтрацией.
// Фильтрация по тегам выполняется на уровне SQL для оптимальной производительности.
func (s *Service) GetQueue(ctx context.Context, filter *Filter, limit int) ([]model.Card, error) {
	srsRepo := s.repos.SRS(s.txManager.Q())

	var dueStates []model.SRSState
	var err error

	// Определяем параметры фильтрации
	var statuses []model.LearningStatus
	var tagNames []string

	if filter != nil {
		statuses = filter.Statuses
		tagNames = filter.Tags
	}

	// Используем оптимизированный метод с фильтрацией по тегам на уровне SQL
	if len(tagNames) > 0 {
		dueStates, err = srsRepo.ListDueForReviewWithTags(ctx, statuses, tagNames, limit)
	} else if len(statuses) > 0 {
		dueStates, err = srsRepo.ListDueForReviewWithFilter(ctx, statuses, limit)
	} else {
		dueStates, err = srsRepo.ListDueForReview(ctx, limit)
	}

	if err != nil {
		return nil, err
	}

	if len(dueStates) == 0 {
		return []model.Card{}, nil
	}

	// Собираем ID карточек
	cardIDs := make([]uuid.UUID, len(dueStates))
	for i, state := range dueStates {
		cardIDs[i] = state.CardID
	}

	// Загружаем сами карточки batch-запросом
	cardRepo := s.repos.Card(s.txManager.Q())
	cards, err := cardRepo.ListByIDs(ctx, cardIDs)
	if err != nil {
		return nil, err
	}

	// Сохраняем порядок из dueStates (отсортированный по due_date)
	cards = s.sortCardsByOrder(cards, cardIDs)

	return cards, nil
}

// sortCardsByOrder сортирует карточки в соответствии с заданным порядком ID.
// Это нужно, потому что ListByIDs не гарантирует порядок результатов.
func (s *Service) sortCardsByOrder(cards []model.Card, orderedIDs []uuid.UUID) []model.Card {
	if len(cards) == 0 {
		return cards
	}

	// Создаём мапу для быстрого доступа
	cardMap := make(map[uuid.UUID]model.Card, len(cards))
	for _, card := range cards {
		cardMap[card.ID] = card
	}

	// Собираем результат в нужном порядке
	result := make([]model.Card, 0, len(cards))
	for _, id := range orderedIDs {
		if card, ok := cardMap[id]; ok {
			result = append(result, card)
		}
	}

	return result
}
