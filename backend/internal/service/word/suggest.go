package word

import (
	"context"
	"strings"
)

// Suggest возвращает подсказки для автокомплита.
// Ищет в локальном словаре пользователя и во внутреннем словаре (dictionary_words).
// Приоритет: локальный словарь > внутренний словарь.
func (s *Service) Suggest(ctx context.Context, query string) ([]Suggestion, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return []Suggestion{}, nil
	}

	// Сначала ищем в локальном словаре
	localSuggestions, err := s.suggestFromLocal(ctx, query)
	if err != nil {
		return []Suggestion{}, err
	}

	// Если найдено точное совпадение в локальном словаре, возвращаем только его
	if len(localSuggestions) > 0 && localSuggestions[0].ExistingWordID != nil {
		return localSuggestions, nil
	}

	// Ищем во внутреннем словаре
	dictionarySuggestions, err := s.suggestFromDictionary(ctx, query)
	if err != nil {
		// Логируем ошибку, но продолжаем с результатами из локального словаря
		// TODO: добавить логирование
	}

	// Объединяем результаты: сначала локальные, потом из словаря
	suggestions := make([]Suggestion, 0, len(localSuggestions)+len(dictionarySuggestions))
	suggestions = append(suggestions, localSuggestions...)
	suggestions = append(suggestions, dictionarySuggestions...)

	return suggestions, nil
}

// suggestFromLocal ищет подсказки в локальном словаре пользователя.
func (s *Service) suggestFromLocal(ctx context.Context, query string) ([]Suggestion, error) {
	// Пытаемся найти точное совпадение
	wordWithRelations, err := s.GetByText(ctx, query)
	if err == nil {
		// Найдено в локальном словаре
		translations := make([]string, 0)
		for _, m := range wordWithRelations.Meanings {
			// Используем translations из таблицы translations
			for _, tr := range m.Translations {
				if tr.TranslationRu != "" {
					translations = append(translations, tr.TranslationRu)
				}
			}
			// Fallback на старое поле для обратной совместимости
			if len(m.Translations) == 0 && m.Meaning.TranslationRu != "" {
				translations = append(translations, m.Meaning.TranslationRu)
			}
		}

		// Если нет переводов, возвращаем пустой массив
		if len(translations) == 0 {
			return []Suggestion{}, nil
		}

		existingWordID := wordWithRelations.Word.ID
		return []Suggestion{
			{
				Text:           wordWithRelations.Word.Text,
				Transcription:  wordWithRelations.Word.Transcription,
				Translations:   translations,
				Origin:         "LOCAL",
				ExistingWordID: &existingWordID,
			},
		}, nil
	}

	// Если не найдено точное совпадение, используем триграммный поиск
	words, err := s.words.SearchSimilar(ctx, query, 5, 0.3)
	if err != nil {
		return []Suggestion{}, err
	}

	if len(words) == 0 {
		return []Suggestion{}, nil
	}

	// Для найденных слов загружаем meanings и формируем подсказки
	suggestions := make([]Suggestion, 0, len(words))
	for _, w := range words {
		// Загружаем meanings для получения переводов
		meanings, err := s.meanings.GetByWordID(ctx, w.ID)
		if err != nil {
			continue
		}

		// Загружаем translations для meanings
		meaningIDs := make([]int64, 0, len(meanings))
		for _, m := range meanings {
			meaningIDs = append(meaningIDs, m.ID)
		}

		allTranslations, err := s.translations.GetByMeaningIDs(ctx, meaningIDs)
		if err != nil {
			continue
		}

		// Группируем translations по meaningID
		translationsByMeaning := make(map[int64][]string)
		for _, tr := range allTranslations {
			translationsByMeaning[tr.MeaningID] = append(translationsByMeaning[tr.MeaningID], tr.TranslationRu)
		}

		translations := make([]string, 0)
		for _, m := range meanings {
			if trs, ok := translationsByMeaning[m.ID]; ok {
				translations = append(translations, trs...)
			} else if m.TranslationRu != "" {
				// Fallback на старое поле для обратной совместимости
				translations = append(translations, m.TranslationRu)
			}
		}

		// Пропускаем слова без переводов
		if len(translations) == 0 {
			continue
		}

		existingWordID := w.ID
		suggestions = append(suggestions, Suggestion{
			Text:           w.Text,
			Transcription:  w.Transcription,
			Translations:   translations,
			Origin:         "LOCAL",
			ExistingWordID: &existingWordID,
		})
	}

	return suggestions, nil
}

// suggestFromDictionary ищет подсказки во внутреннем словаре.
func (s *Service) suggestFromDictionary(ctx context.Context, query string) ([]Suggestion, error) {
	// Пытаемся найти точное совпадение во внутреннем словаре
	dictWord, err := s.dictionary.GetByText(ctx, query)
	if err == nil {
		// Найдено во внутреннем словаре
		translations, err := s.getDictionaryTranslations(ctx, dictWord.ID)
		if err != nil {
			return []Suggestion{}, err
		}

		if len(translations) == 0 {
			return []Suggestion{}, nil
		}

		return []Suggestion{
			{
				Text:          dictWord.Text,
				Transcription: dictWord.Transcription,
				Translations:  translations,
				Origin:        "DICTIONARY",
				// ExistingWordID = nil, так как это не слово пользователя
			},
		}, nil
	}

	// Если не найдено точное совпадение, используем триграммный поиск
	dictWords, err := s.dictionary.SearchSimilar(ctx, query, 5, 0.3)
	if err != nil {
		return []Suggestion{}, err
	}

	if len(dictWords) == 0 {
		return []Suggestion{}, nil
	}

	// Для найденных слов загружаем meanings и формируем подсказки
	suggestions := make([]Suggestion, 0, len(dictWords))
	for _, dw := range dictWords {
		translations, err := s.getDictionaryTranslations(ctx, dw.ID)
		if err != nil {
			continue
		}

		// Пропускаем слова без переводов
		if len(translations) == 0 {
			continue
		}

		suggestions = append(suggestions, Suggestion{
			Text:          dw.Text,
			Transcription: dw.Transcription,
			Translations:  translations,
			Origin:        "DICTIONARY",
			// ExistingWordID = nil, так как это не слово пользователя
		})
	}

	return suggestions, nil
}

// getDictionaryTranslations загружает переводы для слова из внутреннего словаря.
func (s *Service) getDictionaryTranslations(ctx context.Context, wordID int64) ([]string, error) {
	// Загружаем meanings для слова
	meanings, err := s.dictionary.GetMeaningsByWordID(ctx, wordID)
	if err != nil {
		return nil, err
	}

	if len(meanings) == 0 {
		return []string{}, nil
	}

	// Собираем meaningIDs
	meaningIDs := make([]int64, 0, len(meanings))
	for _, m := range meanings {
		meaningIDs = append(meaningIDs, m.ID)
	}

	// Загружаем translations
	allTranslations, err := s.dictionary.GetTranslationsByMeaningIDs(ctx, meaningIDs)
	if err != nil {
		return nil, err
	}

	// Преобразуем в массив строк
	translations := make([]string, 0, len(allTranslations))
	for _, tr := range allTranslations {
		if tr.TranslationRu != "" {
			translations = append(translations, tr.TranslationRu)
		}
	}

	return translations, nil
}

