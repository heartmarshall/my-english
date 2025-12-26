package word

import (
	"context"
	"errors"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

// GetByID возвращает слово по ID со всеми связанными данными.
func (s *Service) GetByID(ctx context.Context, id int64) (*WordWithRelations, error) {
	word, err := s.words.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, service.ErrWordNotFound
		}
		return nil, err
	}

	return s.loadRelations(ctx, word)
}

// GetByText возвращает слово по тексту со всеми связанными данными.
func (s *Service) GetByText(ctx context.Context, text string) (*WordWithRelations, error) {
	word, err := s.words.GetByText(ctx, text)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, service.ErrWordNotFound
		}
		return nil, err
	}

	return s.loadRelations(ctx, word)
}

// List возвращает список слов с пагинацией.
// Возвращает слова без связанных данных (для списка).
func (s *Service) List(ctx context.Context, filter *WordFilter, limit, offset int) ([]model.Word, error) {
	var modelFilter *model.WordFilter
	if filter != nil {
		modelFilter = &model.WordFilter{
			Search: filter.Search,
			Status: filter.Status,
			Tags:   filter.Tags,
		}
	}

	return s.words.List(ctx, modelFilter, limit, offset)
}

// Count возвращает количество слов, соответствующих фильтру.
func (s *Service) Count(ctx context.Context, filter *WordFilter) (int, error) {
	var modelFilter *model.WordFilter
	if filter != nil {
		modelFilter = &model.WordFilter{
			Search: filter.Search,
			Status: filter.Status,
			Tags:   filter.Tags,
		}
	}

	return s.words.Count(ctx, modelFilter)
}

// loadRelations загружает связанные данные для слова.
func (s *Service) loadRelations(ctx context.Context, word model.Word) (*WordWithRelations, error) {
	meanings, err := s.meanings.GetByWordID(ctx, word.ID)
	if err != nil {
		return nil, err
	}

	result := &WordWithRelations{
		Word:     word,
		Meanings: make([]MeaningWithRelations, 0, len(meanings)),
	}

	if len(meanings) == 0 {
		return result, nil
	}

	// Собираем ID meanings для batch loading
	meaningIDs := make([]int64, len(meanings))
	for i, m := range meanings {
		meaningIDs[i] = m.ID
	}

	// Загружаем примеры для всех meanings
	examples, err := s.examples.GetByMeaningIDs(ctx, meaningIDs)
	if err != nil {
		return nil, err
	}

	// Загружаем translations для всех meanings
	translations, err := s.translations.GetByMeaningIDs(ctx, meaningIDs)
	if err != nil {
		return nil, err
	}

	// Загружаем связи meaning-tag
	meaningTags, err := s.meaningTag.GetByMeaningIDs(ctx, meaningIDs)
	if err != nil {
		return nil, err
	}

	// Собираем уникальные tag IDs
	tagIDSet := make(map[int64]struct{})
	for _, mt := range meaningTags {
		tagIDSet[mt.TagID] = struct{}{}
	}

	tagIDs := make([]int64, 0, len(tagIDSet))
	for id := range tagIDSet {
		tagIDs = append(tagIDs, id)
	}

	// Загружаем теги
	var tags []model.Tag
	if len(tagIDs) > 0 {
		tags, err = s.tags.GetByIDs(ctx, tagIDs)
		if err != nil {
			return nil, err
		}
	}

	// Создаём мапы для быстрого доступа
	examplesByMeaning := make(map[int64][]model.Example)
	for _, ex := range examples {
		examplesByMeaning[ex.MeaningID] = append(examplesByMeaning[ex.MeaningID], ex)
	}

	translationsByMeaning := make(map[int64][]model.Translation)
	for _, tr := range translations {
		translationsByMeaning[tr.MeaningID] = append(translationsByMeaning[tr.MeaningID], tr)
	}

	tagsByID := make(map[int64]model.Tag)
	for _, t := range tags {
		tagsByID[t.ID] = t
	}

	tagIDsByMeaning := make(map[int64][]int64)
	for _, mt := range meaningTags {
		tagIDsByMeaning[mt.MeaningID] = append(tagIDsByMeaning[mt.MeaningID], mt.TagID)
	}

	// Собираем результат
	for _, m := range meanings {
		mr := MeaningWithRelations{
			Meaning:     m,
			Translations: translationsByMeaning[m.ID],
			Examples:    examplesByMeaning[m.ID],
			Tags:        make([]model.Tag, 0),
		}

		if mr.Examples == nil {
			mr.Examples = make([]model.Example, 0)
		}
		if mr.Translations == nil {
			mr.Translations = make([]model.Translation, 0)
		}

		for _, tagID := range tagIDsByMeaning[m.ID] {
			if tag, ok := tagsByID[tagID]; ok {
				mr.Tags = append(mr.Tags, tag)
			}
		}

		result.Meanings = append(result.Meanings, mr)
	}

	return result, nil
}

// LoadMeaningsForWords загружает meanings для списка слов (для batch loading в GraphQL).
func (s *Service) LoadMeaningsForWords(ctx context.Context, wordIDs []int64) (map[int64][]model.Meaning, error) {
	result := make(map[int64][]model.Meaning)

	for _, wordID := range wordIDs {
		meanings, err := s.meanings.GetByWordID(ctx, wordID)
		if err != nil {
			return nil, err
		}
		result[wordID] = meanings
	}

	return result, nil
}
