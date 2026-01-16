package dictionary

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/model"
)

// createSenses создает смыслы и связанные с ними сущности (переводы и примеры).
func (s *Service) createSenses(ctx context.Context, entryID uuid.UUID, senses []SenseInput) error {
	for i, senseIn := range senses {
		sense := buildSense(entryID, senseIn)
		createdSense, err := s.repos.Senses.Create(ctx, sense)
		if err != nil {
			return fmt.Errorf("create sense[%d]: %w", i, err)
		}

		if err := s.createTranslations(ctx, createdSense.ID, senseIn.Translations); err != nil {
			return fmt.Errorf("create translations for sense[%d]: %w", i, err)
		}

		if err := s.createExamples(ctx, createdSense.ID, senseIn.Examples); err != nil {
			return fmt.Errorf("create examples for sense[%d]: %w", i, err)
		}
	}
	return nil
}

// createTranslations создает переводы для смысла.
func (s *Service) createTranslations(ctx context.Context, senseID uuid.UUID, translations []TranslationInput) error {
	if len(translations) == 0 {
		return nil
	}

	models := buildTranslations(senseID, translations)
	_, err := s.repos.Translations.BatchCreate(ctx, models)
	if err != nil {
		return fmt.Errorf("batch create translations: %w", err)
	}
	return nil
}

// createExamples создает примеры для смысла.
func (s *Service) createExamples(ctx context.Context, senseID uuid.UUID, examples []ExampleInput) error {
	if len(examples) == 0 {
		return nil
	}

	models := buildExamples(senseID, examples)
	_, err := s.repos.Examples.BatchCreate(ctx, models)
	if err != nil {
		return fmt.Errorf("batch create examples: %w", err)
	}
	return nil
}

// createImages создает изображения для записи.
func (s *Service) createImages(ctx context.Context, entryID uuid.UUID, images []ImageInput) error {
	if len(images) == 0 {
		return nil
	}

	models := buildImages(entryID, images)
	_, err := s.repos.Images.BatchCreate(ctx, models)
	if err != nil {
		return fmt.Errorf("batch create images: %w", err)
	}
	return nil
}

// createPronunciations создает произношения для записи.
func (s *Service) createPronunciations(ctx context.Context, entryID uuid.UUID, pronunciations []PronunciationInput) error {
	if len(pronunciations) == 0 {
		return nil
	}

	models := buildPronunciations(entryID, pronunciations)
	_, err := s.repos.Pronunciations.BatchCreate(ctx, models)
	if err != nil {
		return fmt.Errorf("batch create pronunciations: %w", err)
	}
	return nil
}

// createCardIfNeeded создает карточку для изучения, если требуется.
// Использует дефолтные значения для новой карточки согласно алгоритму SM-2.
func (s *Service) createCardIfNeeded(ctx context.Context, entryID uuid.UUID, createCard bool) error {
	if !createCard {
		return nil
	}

	card := &model.Card{
		EntryID:      entryID,
		Status:       model.StatusNew,
		IntervalDays: 0,
		EaseFactor:   DefaultEaseFactor,
	}
	_, err := s.repos.Cards.Create(ctx, card)
	if err != nil {
		return fmt.Errorf("create card: %w", err)
	}
	return nil
}
