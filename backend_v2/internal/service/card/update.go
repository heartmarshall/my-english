package card

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateCardInput) (*model.Card, error) {
	var updatedCard *model.Card

	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx database.Querier) error {
		cardRepo := s.repos.Card(tx)

		// 1. Получаем текущую карточку (с блокировкой на обновление, если нужно, но пока просто Get)
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

			// Собираем ID новых тегов
			var newTagIDs []int
			for _, name := range input.Tags {
				name = strings.TrimSpace(name)
				if name == "" {
					continue
				}

				tag, err := tagRepo.GetByName(ctx, name)
				if err != nil && !errors.Is(err, database.ErrNotFound) {
					return err
				}

				if tag == nil {
					tag, err = tagRepo.Create(ctx, &model.Tag{Name: name}) // Упрощенно
					if err != nil {
						return err
					}
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

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft Delete не требует транзакции, если мы не чистим связи
	return s.repos.Card(s.txManager.Q()).SoftDelete(ctx, id)
}
