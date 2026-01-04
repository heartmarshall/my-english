package study

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetQueue возвращает список карточек для изучения с опциональной фильтрацией.
func (s *Service) GetQueue(ctx context.Context, filter *Filter, limit int) ([]model.Card, error) {
	// 1. Получаем состояния SRS, которые пора учить
	// Используем обычный пул (чтение)
	srsRepo := s.repos.SRS(s.txManager.Q())

	var dueStates []model.SRSState
	var err error

	// Применяем фильтр по статусам, если указан
	if filter != nil && len(filter.Statuses) > 0 {
		dueStates, err = srsRepo.ListDueForReviewWithFilter(ctx, filter.Statuses, limit)
	} else {
		dueStates, err = srsRepo.ListDueForReview(ctx, limit)
	}

	if err != nil {
		return nil, err
	}

	if len(dueStates) == 0 {
		return []model.Card{}, nil
	}

	// 2. Собираем ID карточек
	cardIDs := make([]uuid.UUID, len(dueStates))
	for i, state := range dueStates {
		cardIDs[i] = state.CardID
	}

	// 3. Загружаем сами карточки batch-запросом
	cardRepo := s.repos.Card(s.txManager.Q())
	cards, err := cardRepo.ListByIDs(ctx, cardIDs)
	if err != nil {
		return nil, err
	}

	// 4. Применяем фильтр по тегам, если указан
	if filter != nil && len(filter.Tags) > 0 {
		cards = s.filterCardsByTags(ctx, cards, filter.Tags)
	}

	return cards, nil
}

// filterCardsByTags фильтрует карточки по тегам.
func (s *Service) filterCardsByTags(ctx context.Context, cards []model.Card, tagNames []string) []model.Card {
	if len(tagNames) == 0 {
		return cards
	}

	// Получаем ID тегов по именам
	tagRepo := s.repos.Tag(s.txManager.Q())
	tagIDsMap := make(map[int]bool)

	for _, tagName := range tagNames {
		tag, err := tagRepo.GetByName(ctx, tagName)
		if err != nil || tag == nil {
			continue // Пропускаем несуществующие теги
		}
		tagIDsMap[tag.ID] = true
	}

	if len(tagIDsMap) == 0 {
		return []model.Card{} // Если ни один тег не найден, возвращаем пустой список
	}

	// Получаем связи карточек с тегами
	cardTagRepo := s.repos.CardTag(s.txManager.Q())
	cardIDs := make([]uuid.UUID, len(cards))
	for i := range cards {
		cardIDs[i] = cards[i].ID
	}

	cardTags, err := cardTagRepo.ListByCardIDs(ctx, cardIDs)
	if err != nil {
		// В случае ошибки возвращаем исходный список
		return cards
	}

	// Создаем мапу: cardID -> множество tagID
	cardTagMap := make(map[uuid.UUID]map[int]bool)
	for _, ct := range cardTags {
		if cardTagMap[ct.CardID] == nil {
			cardTagMap[ct.CardID] = make(map[int]bool)
		}
		cardTagMap[ct.CardID][ct.TagID] = true
	}

	// Фильтруем карточки: карточка должна иметь все указанные теги
	filtered := make([]model.Card, 0)
	for _, card := range cards {
		cardTags := cardTagMap[card.ID]
		hasAllTags := true
		for tagID := range tagIDsMap {
			if !cardTags[tagID] {
				hasAllTags = false
				break
			}
		}
		if hasAllTags {
			filtered = append(filtered, card)
		}
	}

	return filtered
}
