package card

import (
	"context"
	"errors"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

func (s *Service) Create(ctx context.Context, input CreateCardInput) (*model.Card, error) {
	// 1. Валидация
	if input.SenseID == nil && (input.CustomText == nil || strings.TrimSpace(*input.CustomText) == "") {
		return nil, service.ErrInvalidInput // Должен быть либо SenseID, либо CustomText
	}

	var createdCard *model.Card

	// 2. Транзакция
	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx database.Querier) error {
		// --- A. Валидация SenseID (если указан) ---
		if input.SenseID != nil {
			sense, err := s.repos.Sense(tx).GetByID(ctx, *input.SenseID)
			if err != nil {
				if database.IsNotFoundError(err) {
					return service.ErrSenseNotFound
				}
				return err
			}
			if sense == nil {
				return service.ErrSenseNotFound
			}
		}

		// --- B. Создаем Card ---
		cardRepo := s.repos.Card(tx)

		cardModel := &model.Card{
			SenseID:             input.SenseID,
			CustomText:          input.CustomText,
			CustomTranscription: input.CustomTranscription,
			CustomTranslations:  input.CustomTranslations,
			CustomNote:          input.CustomNote,
			CustomImageURL:      input.CustomImageURL,
		}

		var err error
		createdCard, err = cardRepo.Create(ctx, cardModel)
		if err != nil {
			return err
		}

		// --- C. Инициализируем SRS State (Status: New) ---
		srsRepo := s.repos.SRS(tx)
		initialState := &model.SRSState{
			CardID:        createdCard.ID,
			Status:        model.LearningStatusNew,
			DueDate:       nil,              // Для новых карт даты нет, пока не выучим
			AlgorithmData: map[string]any{}, // Пустой JSON (или дефолты алгоритма)
		}

		if _, err := srsRepo.Upsert(ctx, initialState); err != nil {
			return err
		}

		// --- D. Обрабатываем Теги ---
		if len(input.Tags) > 0 {
			tagRepo := s.repos.Tag(tx)
			cardTagRepo := s.repos.CardTag(tx)

			for _, tagName := range input.Tags {
				tagName = strings.TrimSpace(tagName)
				if tagName == "" {
					continue
				}

				// GetOrCreate Tag
				// Пытаемся найти
				tag, err := tagRepo.GetByName(ctx, tagName)
				if err != nil && !errors.Is(err, database.ErrNotFound) {
					return err
				}

				// Если не нашли - создаем
				if tag == nil {
					tag, err = tagRepo.Create(ctx, &model.Tag{Name: tagName})
					if err != nil {
						// Если параллельно создали такой же тег (race condition), пробуем найти снова
						if database.IsDuplicateError(err) {
							tag, err = tagRepo.GetByName(ctx, tagName)
							if err != nil {
								return err
							}
						} else {
							return err
						}
					}
				}

				// Привязываем к карточке
				if err := cardTagRepo.Attach(ctx, createdCard.ID, tag.ID); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdCard, nil
}
