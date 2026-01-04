package word

import (
	"context"
	"strings"
)

// Suggest возвращает подсказки для автокомплита.
func (s *Service) Suggest(ctx context.Context, query string) ([]Suggestion, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return []Suggestion{}, nil
	}

	// 1. Сначала ищем в локальном словаре пользователя
	localSuggestions, err := s.suggestFromLocal(ctx, query)
	if err != nil {
		return []Suggestion{}, err
	}

	// Если найдено точное совпадение у пользователя, это приоритет
	if len(localSuggestions) > 0 && localSuggestions[0].ExistingWordID != nil {
		return localSuggestions, nil
	}

	// 2. Ищем во внешнем/системном словаре
	dictionarySuggestions, err := s.suggestFromDictionary(ctx, query)
	if err != nil {
		// Логируем ошибку, но не прерываем работу, если есть локальные подсказки
		// slog.Error("failed to fetch from dictionary", "error", err)
	}

	// Объединяем результаты
	suggestions := make([]Suggestion, 0, len(localSuggestions)+len(dictionarySuggestions))
	suggestions = append(suggestions, localSuggestions...)
	suggestions = append(suggestions, dictionarySuggestions...)

	return suggestions, nil
}

// suggestFromDictionary реализует логику Read-Through Cache.
func (s *Service) suggestFromDictionary(ctx context.Context, query string) ([]Suggestion, error) {
	// А. Ищем в кэше (БД)
	cachedWords, err := s.dictionary.GetByText(ctx, query)
	if err != nil {
		return nil, err
	}

	// Б. Если в кэше пусто, идем во внешний мир
	if len(cachedWords) == 0 {
		if s.fetcher != nil {
			externalData, err := s.fetcher.FetchWord(ctx, query)
			if err != nil {
				// Ошибка внешнего API не должна ломать приложение, просто идем дальше
				// (можно добавить логирование)
			} else if externalData != nil {
				// Сохраняем в БД (Write-through)
				if err := s.dictionary.SaveWordData(ctx, externalData); err == nil {
					// Добавляем сохраненное слово в список для отображения
					cachedWords = append(cachedWords, externalData.Word)
				}
			}
		}
	}

	// В. Если все еще пусто — пробуем нечеткий поиск по кэшу (на случай опечаток)
	if len(cachedWords) == 0 {
		cachedWords, err = s.dictionary.SearchSimilar(ctx, query, 5, 0.3)
		if err != nil {
			return nil, err
		}
	}

	// Г. Формируем DTO для ответа
	suggestions := make([]Suggestion, 0, len(cachedWords))
	for _, dw := range cachedWords {
		// 1. Получаем переводы (как и раньше)
		translations, err := s.getDictionaryTranslations(ctx, dw.ID)
		if err != nil {
			continue
		}

		// 2. Получаем определение (НОВОЕ)
		var definition *string
		meanings, err := s.dictionary.GetMeaningsByWordID(ctx, dw.ID)
		if err == nil && len(meanings) > 0 {
			// Берем определение из первого значения
			definition = meanings[0].DefinitionEn
		}

		suggestions = append(suggestions, Suggestion{
			Text:           dw.Text,
			Transcription:  dw.Transcription,
			Translations:   translations,
			Definition:     definition, // <-- Заполняем
			Origin:         "DICTIONARY",
			ExistingWordID: nil,
		})
	}

	return suggestions, nil
}

// suggestFromLocal (остается без изменений из предыдущей версии)
func (s *Service) suggestFromLocal(ctx context.Context, query string) ([]Suggestion, error) {
	// ... см. предыдущую реализацию ...
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

		// Получаем определение из первого значения
		var definition *string
		if len(wordWithRelations.Meanings) > 0 && wordWithRelations.Meanings[0].Meaning.DefinitionEn != nil {
			definition = wordWithRelations.Meanings[0].Meaning.DefinitionEn
		}

		existingWordID := wordWithRelations.Word.ID
		return []Suggestion{
			{
				Text:           wordWithRelations.Word.Text,
				Transcription:  wordWithRelations.Word.Transcription,
				Translations:   translations,
				Definition:     definition,
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

		// Получаем определение из первого значения
		var definition *string
		if len(meanings) > 0 && meanings[0].DefinitionEn != nil {
			definition = meanings[0].DefinitionEn
		}

		existingWordID := w.ID
		suggestions = append(suggestions, Suggestion{
			Text:           w.Text,
			Transcription:  w.Transcription,
			Translations:   translations,
			Definition:     definition,
			Origin:         "LOCAL",
			ExistingWordID: &existingWordID,
		})
	}

	return suggestions, nil
}

// getDictionaryTranslations (вспомогательный метод для загрузки переводов)
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
