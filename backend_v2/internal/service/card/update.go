package card

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateCardInput) (*model.Card, error) {
	// Валидация входных данных
	if err := s.validateUpdateInput(input); err != nil {
		return nil, err
	}

	var updatedCard *model.Card

	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx database.Querier) error {
		cardRepo := s.repos.Card(tx)

		// 1. Получаем текущую карточку
		current, err := cardRepo.GetByID(ctx, id)
		if err != nil {
			if database.IsNotFoundError(err) {
				return service.ErrCardNotFound
			}
			return err
		}

		// 2. Обновляем поля
		current.CustomTranscription = input.CustomTranscription
		current.CustomTranslations = input.CustomTranslations
		current.CustomNote = input.CustomNote
		current.CustomImageURL = input.CustomImageURL

		// Сохраняем изменения в БД
		updatedCard, err = cardRepo.Update(ctx, current)
		if err != nil {
			return err
		}

		// 3. Синхронизируем теги
		// Если input.Tags == nil, теги не обновляем
		// Если input.Tags != nil (даже пустой массив), синхронизируем теги
		if input.Tags != nil {
			tagRepo := s.repos.Tag(tx)
			cardTagRepo := s.repos.CardTag(tx)

			// Дедуплицируем и нормализуем теги
			uniqueTags := normalizeAndDeduplicateTags(input.Tags)

			// Собираем ID новых тегов с использованием атомарного GetOrCreate
			var newTagIDs []int
			for _, name := range uniqueTags {
				// GetOrCreate использует ON CONFLICT, безопасно для параллельных запросов
				tag, err := tagRepo.GetOrCreate(ctx, name)
				if err != nil {
					return err
				}
				newTagIDs = append(newTagIDs, tag.ID)
			}

			// Удаляем старые связи
			if err := cardTagRepo.DetachAll(ctx, id); err != nil {
				return err
			}
			// Добавляем новые (если есть)
			for _, tid := range newTagIDs {
				if err := cardTagRepo.Attach(ctx, id, tid); err != nil {
					return err
				}
			}
		}

		// Перезагружаем карточку с обновленными данными
		updatedCard, err = cardRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})

	return updatedCard, err
}

// validateUpdateInput проверяет корректность входных данных для обновления карточки.
func (s *Service) validateUpdateInput(input UpdateCardInput) error {
	// Проверка количества тегов
	if input.Tags != nil && len(input.Tags) > maxTagsPerCard {
		return service.ErrInvalidInput
	}

	// Проверка длины каждого тега
	for _, tag := range input.Tags {
		if len(strings.TrimSpace(tag)) > maxTagNameLength {
			return service.ErrInvalidInput
		}
	}

	// Проверка переводов
	for _, tr := range input.CustomTranslations {
		if len(tr) > maxTranslationsLength {
			return service.ErrInvalidInput
		}
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft Delete не требует транзакции, если мы не чистим связи
	return s.repos.Card(s.txManager.Q()).SoftDelete(ctx, id)
}
