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
// Все операции выполняются в транзакции.
func (s *Service) Update(ctx context.Context, id int64, input UpdateWordInput) (*WordWithRelations, error) {
	// Валидация
	text := strings.TrimSpace(strings.ToLower(input.Text))
	if text == "" {
		return nil, service.ErrInvalidInput
	}

	// Получаем существующее слово (вне транзакции для быстрого отклонения)
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

	var result *WordWithRelations

	// Выполняем обновление в транзакции
	err = s.txRunner.RunInTx(ctx, func(ctx context.Context, tx database.Querier) error {
		txRepos := s.withTx(tx)

		// Обновляем слово
		word.Text = text
		word.Transcription = input.Transcription
		word.AudioURL = input.AudioURL

		if err := txRepos.words.Update(ctx, &word); err != nil {
			if errors.Is(err, database.ErrDuplicate) {
				return service.ErrWordAlreadyExists
			}
			return err
		}

		// Удаляем старые meanings (CASCADE удалит examples и meaning_tags)
		if _, err := txRepos.meanings.DeleteByWordID(ctx, id); err != nil {
			return err
		}

		// Создаём новые meanings
		result = &WordWithRelations{
			Word:     word,
			Meanings: make([]MeaningWithRelations, 0, len(input.Meanings)),
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

			mr, err := s.createMeaningTx(ctx, txRepos, word.ID, createInput)
			if err != nil {
				return err
			}
			result.Meanings = append(result.Meanings, *mr)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
