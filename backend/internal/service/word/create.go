package word

import (
	"context"
	"errors"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

// Create создаёт новое слово со всеми связанными данными.
func (s *Service) Create(ctx context.Context, input CreateWordInput) (*WordWithRelations, error) {
	// Валидация
	text := strings.TrimSpace(strings.ToLower(input.Text))
	if text == "" {
		return nil, service.ErrInvalidInput
	}

	if len(input.Meanings) == 0 {
		return nil, service.ErrInvalidInput
	}

	// Проверяем, не существует ли слово
	_, err := s.words.GetByText(ctx, text)
	if err == nil {
		return nil, service.ErrWordAlreadyExists
	}
	if !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}

	// Создаём слово
	word := &model.Word{
		Text:          text,
		Transcription: input.Transcription,
		AudioURL:      input.AudioURL,
	}

	if err := s.words.Create(ctx, word); err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			return nil, service.ErrWordAlreadyExists
		}
		return nil, err
	}

	// Создаём meanings с примерами и тегами
	result := &WordWithRelations{
		Word:     word,
		Meanings: make([]*MeaningWithRelations, 0, len(input.Meanings)),
	}

	for _, meaningInput := range input.Meanings {
		mr, err := s.createMeaning(ctx, word.ID, meaningInput)
		if err != nil {
			return nil, err
		}
		result.Meanings = append(result.Meanings, mr)
	}

	return result, nil
}

// createMeaning создаёт meaning с примерами и тегами.
func (s *Service) createMeaning(ctx context.Context, wordID int64, input CreateMeaningInput) (*MeaningWithRelations, error) {
	if input.TranslationRu == "" {
		return nil, service.ErrInvalidInput
	}

	// Создаём meaning
	meaning := &model.Meaning{
		WordID:        wordID,
		PartOfSpeech:  input.PartOfSpeech,
		DefinitionEn:  input.DefinitionEn,
		TranslationRu: input.TranslationRu,
		CefrLevel:     input.CefrLevel,
		ImageURL:      input.ImageURL,
	}

	if err := s.meanings.Create(ctx, meaning); err != nil {
		return nil, err
	}

	result := &MeaningWithRelations{
		Meaning:  meaning,
		Examples: make([]*model.Example, 0),
		Tags:     make([]*model.Tag, 0),
	}

	// Создаём примеры
	if len(input.Examples) > 0 {
		examples := make([]*model.Example, 0, len(input.Examples))
		for _, exInput := range input.Examples {
			if exInput.SentenceEn == "" {
				continue
			}
			examples = append(examples, &model.Example{
				MeaningID:  meaning.ID,
				SentenceEn: exInput.SentenceEn,
				SentenceRu: exInput.SentenceRu,
				SourceName: exInput.SourceName,
			})
		}

		if len(examples) > 0 {
			if err := s.examples.CreateBatch(ctx, examples); err != nil {
				return nil, err
			}
			result.Examples = examples
		}
	}

	// Создаём/получаем теги и привязываем
	if len(input.Tags) > 0 {
		tagIDs := make([]int64, 0, len(input.Tags))
		tags := make([]*model.Tag, 0, len(input.Tags))

		for _, tagName := range input.Tags {
			tag, err := s.tags.GetOrCreate(ctx, tagName)
			if err != nil {
				return nil, err
			}
			tagIDs = append(tagIDs, tag.ID)
			tags = append(tags, tag)
		}

		if err := s.meaningTag.AttachTags(ctx, meaning.ID, tagIDs); err != nil {
			return nil, err
		}
		result.Tags = tags
	}

	return result, nil
}
