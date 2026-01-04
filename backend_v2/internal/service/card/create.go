package card

import (
	"context"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

// Константы валидации
const (
	maxCustomTextLength   = 1000
	maxTagNameLength      = 50
	maxTagsPerCard        = 20
	maxTranslationsLength = 500
)

func (s *Service) Create(ctx context.Context, input CreateCardInput) (*model.Card, error) {
	// 1. Валидация входных данных
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Подготавливаем и дедуплицируем теги
	uniqueTags := normalizeAndDeduplicateTags(input.Tags)

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

		// --- D. Обрабатываем Теги (используем атомарный GetOrCreate) ---
		if len(uniqueTags) > 0 {
			tagRepo := s.repos.Tag(tx)
			cardTagRepo := s.repos.CardTag(tx)

			for _, tagName := range uniqueTags {
				// GetOrCreate использует ON CONFLICT, безопасно для параллельных запросов
				tag, err := tagRepo.GetOrCreate(ctx, tagName)
				if err != nil {
					return err
				}

				// Привязываем к карточке (Attach использует ON CONFLICT DO NOTHING)
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

// validateCreateInput проверяет корректность входных данных для создания карточки.
func (s *Service) validateCreateInput(input CreateCardInput) error {
	// Должен быть либо SenseID, либо CustomText
	hasCustomText := input.CustomText != nil && strings.TrimSpace(*input.CustomText) != ""
	if input.SenseID == nil && !hasCustomText {
		return service.ErrInvalidInput
	}

	// Проверка длины CustomText
	if input.CustomText != nil && len(*input.CustomText) > maxCustomTextLength {
		return service.ErrInvalidInput
	}

	// Проверка количества тегов
	if len(input.Tags) > maxTagsPerCard {
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

// normalizeAndDeduplicateTags нормализует и удаляет дубликаты тегов.
func normalizeAndDeduplicateTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		normalized := strings.TrimSpace(tag)
		if normalized == "" {
			continue
		}
		// Приводим к нижнему регистру для дедупликации
		key := strings.ToLower(normalized)
		if !seen[key] {
			seen[key] = true
			result = append(result, normalized) // Сохраняем оригинальный регистр
		}
	}

	return result
}
