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
// Все операции выполняются в транзакции.
func (s *Service) Create(ctx context.Context, input CreateWordInput) (*WordWithRelations, error) {
	// Валидация - только text обязателен
	text := strings.TrimSpace(strings.ToLower(input.Text))
	if text == "" {
		return nil, service.ErrInvalidInput
	}

	// Проверяем, не существует ли слово (вне транзакции для быстрого отклонения)
	_, err := s.words.GetByText(ctx, text)
	if err == nil {
		return nil, service.ErrWordAlreadyExists
	}
	if !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}

	var result *WordWithRelations

	// Выполняем создание в транзакции
	err = s.txRunner.RunInTx(ctx, func(ctx context.Context, tx database.Querier) error {
		txRepos := s.withTx(tx)

		// Создаём слово
		word := &model.Word{
			Text:          text,
			Transcription: input.Transcription,
			AudioURL:      input.AudioURL,
		}

		if err := txRepos.words.Create(ctx, word); err != nil {
			if errors.Is(err, database.ErrDuplicate) {
				return service.ErrWordAlreadyExists
			}
			return err
		}

		// Создаём meanings с примерами и тегами
		result = &WordWithRelations{
			Word:     *word,
			Meanings: make([]MeaningWithRelations, 0, len(input.Meanings)),
		}

		for _, meaningInput := range input.Meanings {
			mr, err := s.createMeaningTx(ctx, txRepos, word.ID, meaningInput)
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

// createMeaningTx создаёт meaning с примерами и тегами в рамках транзакции.
func (s *Service) createMeaningTx(ctx context.Context, r *repos, wordID int64, input CreateMeaningInput) (*MeaningWithRelations, error) {
	if len(input.Translations) == 0 {
		return nil, service.ErrInvalidInput
	}

	// Создаём meaning (translation_ru оставляем пустым или используем первый перевод для обратной совместимости)
	firstTranslation := input.Translations[0]
	meaning := &model.Meaning{
		WordID:        wordID,
		PartOfSpeech:  input.PartOfSpeech,
		DefinitionEn:  input.DefinitionEn,
		TranslationRu: firstTranslation, // Для обратной совместимости, но основное хранилище - translations
		CefrLevel:     input.CefrLevel,
		ImageURL:      input.ImageURL,
	}

	if err := r.meanings.Create(ctx, meaning); err != nil {
		return nil, err
	}

	// Создаём translations
	translations := make([]*model.Translation, 0, len(input.Translations))
	now := database.DefaultClock.Now()
	for _, tr := range input.Translations {
		if tr == "" {
			continue
		}
		translations = append(translations, &model.Translation{
			MeaningID:     meaning.ID,
			TranslationRu: tr,
			CreatedAt:     now,
		})
	}

	if len(translations) > 0 {
		if err := r.translations.CreateBatch(ctx, translations); err != nil {
			return nil, err
		}
	}

	result := MeaningWithRelations{
		Meaning:      *meaning,
		Translations: make([]model.Translation, len(translations)),
		Examples:     make([]model.Example, 0),
		Tags:         make([]model.Tag, 0),
	}

	// Конвертируем translations
	for i, t := range translations {
		result.Translations[i] = *t
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
			if err := r.examples.CreateBatch(ctx, examples); err != nil {
				return nil, err
			}
			result.Examples = make([]model.Example, len(examples))
			for i, ex := range examples {
				result.Examples[i] = *ex
			}
		}
	}

	// Создаём/получаем теги и привязываем
	if len(input.Tags) > 0 {
		tagIDs := make([]int64, 0, len(input.Tags))
		tags := make([]model.Tag, 0, len(input.Tags))

		for _, tagName := range input.Tags {
			tag, err := r.tags.GetOrCreate(ctx, tagName)
			if err != nil {
				return nil, err
			}
			tagIDs = append(tagIDs, tag.ID)
			tags = append(tags, tag)
		}

		if err := r.meaningTag.AttachTags(ctx, meaning.ID, tagIDs); err != nil {
			return nil, err
		}
		result.Tags = tags
	}

	return &result, nil
}
