package word

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
)

// GetDictionaryMeaningsByWordID возвращает значения слова из словаря по ID слова.
func (s *Service) GetDictionaryMeaningsByWordID(ctx context.Context, dictionaryWordID int64) ([]model.DictionaryMeaning, error) {
	meanings, err := s.dictionary.GetMeaningsByWordID(ctx, dictionaryWordID)
	if err != nil {
		return nil, err
	}
	return meanings, nil
}

// GetDictionarySynonymsAntonymsByUserMeaning возвращает синонимы и антонимы для пользовательского значения.
// Ищет соответствующее значение в словаре по тексту слова и части речи.
func (s *Service) GetDictionarySynonymsAntonymsByUserMeaning(ctx context.Context, wordText string, partOfSpeech model.PartOfSpeech) ([]model.DictionarySynonymAntonym, int64, error) {
	// Ищем слово в словаре
	dictWords, err := s.dictionary.GetByText(ctx, wordText)
	if err != nil {
		// Если слово не найдено в словаре, возвращаем пустой список
		return []model.DictionarySynonymAntonym{}, 0, nil
	}

	if len(dictWords) == 0 {
		return []model.DictionarySynonymAntonym{}, 0, nil
	}

	// Берем первое найденное слово
	dictWord := dictWords[0]

	// Получаем значения для слова из словаря
	dictMeanings, err := s.dictionary.GetMeaningsByWordID(ctx, dictWord.ID)
	if err != nil {
		return nil, 0, err
	}

	// Ищем значение с такой же частью речи
	var dictMeaningID int64
	for _, dm := range dictMeanings {
		if dm.PartOfSpeech == partOfSpeech {
			dictMeaningID = dm.ID
			break
		}
	}

	// Если не найдено точное соответствие, берем первое значение
	if dictMeaningID == 0 && len(dictMeanings) > 0 {
		dictMeaningID = dictMeanings[0].ID
	}

	if dictMeaningID == 0 {
		return []model.DictionarySynonymAntonym{}, 0, nil
	}

	// Получаем связи для найденного значения
	relations, err := s.dictionary.GetRelationsByMeaningID(ctx, dictMeaningID)
	if err != nil {
		return nil, 0, err
	}

	return relations, dictMeaningID, nil
}

// GetDictionarySynonymsAntonyms возвращает синонимы и антонимы для значения из словаря.
func (s *Service) GetDictionarySynonymsAntonyms(ctx context.Context, dictionaryMeaningID int64) ([]model.DictionarySynonymAntonym, error) {
	relations, err := s.dictionary.GetRelationsByMeaningID(ctx, dictionaryMeaningID)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// GetDictionarySynonyms возвращает только синонимы для значения из словаря.
func (s *Service) GetDictionarySynonyms(ctx context.Context, dictionaryMeaningID int64) ([]model.DictionarySynonymAntonym, error) {
	synonyms, err := s.dictionary.GetSynonymsByMeaningID(ctx, dictionaryMeaningID)
	if err != nil {
		return nil, err
	}
	return synonyms, nil
}

// GetDictionaryAntonyms возвращает только антонимы для значения из словаря.
func (s *Service) GetDictionaryAntonyms(ctx context.Context, dictionaryMeaningID int64) ([]model.DictionarySynonymAntonym, error) {
	antonyms, err := s.dictionary.GetAntonymsByMeaningID(ctx, dictionaryMeaningID)
	if err != nil {
		return nil, err
	}
	return antonyms, nil
}
