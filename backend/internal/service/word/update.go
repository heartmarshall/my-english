package word

import (
	"context"
	"errors"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/service"
)

// Update обновляет слово и все связанные данные.
// Стратегия: полная замена meanings (удаляем старые, создаём новые).
func (s *Service) Update(ctx context.Context, id int64, input UpdateWordInput) (*WordWithRelations, error) {
	// Валидация
	text := strings.TrimSpace(strings.ToLower(input.Text))
	if text == "" {
		return nil, service.ErrInvalidInput
	}

	// Получаем существующее слово
	word, err := s.words.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, service.ErrWordNotFound
		}
		return nil, err
	}

	// Проверяем, не занят ли новый текст другим словом
	if text != word.Text {
		existing, err := s.words.GetByText(ctx, text)
		if err == nil && existing.ID != id {
			return nil, service.ErrWordAlreadyExists
		}
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return nil, err
		}
	}

	// Обновляем слово
	word.Text = text
	word.Transcription = input.Transcription
	word.AudioURL = input.AudioURL

	if err := s.words.Update(ctx, word); err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			return nil, service.ErrWordAlreadyExists
		}
		return nil, err
	}

	// Удаляем старые meanings (CASCADE удалит examples и meaning_tags)
	if _, err := s.meanings.DeleteByWordID(ctx, id); err != nil {
		return nil, err
	}

	// Создаём новые meanings
	result := &WordWithRelations{
		Word:     word,
		Meanings: make([]*MeaningWithRelations, 0, len(input.Meanings)),
	}

	for _, meaningInput := range input.Meanings {
		createInput := CreateMeaningInput{
			PartOfSpeech:  meaningInput.PartOfSpeech,
			DefinitionEn:  meaningInput.DefinitionEn,
			TranslationRu: meaningInput.TranslationRu,
			CefrLevel:     meaningInput.CefrLevel,
			ImageURL:      meaningInput.ImageURL,
			Examples:      meaningInput.Examples,
			Tags:          meaningInput.Tags,
		}

		mr, err := s.createMeaning(ctx, word.ID, createInput)
		if err != nil {
			return nil, err
		}
		result.Meanings = append(result.Meanings, mr)
	}

	return result, nil
}
