package word

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
)

// GetWordForms возвращает формы слова из словаря по тексту слова пользователя.
func (s *Service) GetWordForms(ctx context.Context, wordText string) ([]model.DictionaryWordForm, error) {
	// Ищем слово в словаре по тексту
	dictWords, err := s.dictionary.GetByText(ctx, wordText)
	if err != nil {
		// Если слово не найдено в словаре, возвращаем пустой список
		return []model.DictionaryWordForm{}, nil
	}

	if len(dictWords) == 0 {
		return []model.DictionaryWordForm{}, nil
	}

	// Берем первое найденное слово
	dictWord := dictWords[0]

	// Получаем формы слова
	forms, err := s.dictionary.GetFormsByWordID(ctx, dictWord.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

// GetWordByFormText возвращает слово пользователя по форме слова из словаря.
// Полезно для поиска основного слова по любой его форме (например, найти "go" по "went").
func (s *Service) GetWordByFormText(ctx context.Context, formText string) (*WordWithRelations, error) {
	// Ищем слово в словаре по форме
	dictWord, err := s.dictionary.GetWordByFormText(ctx, formText)
	if err != nil {
		// Если форма не найдена, возвращаем nil
		return nil, nil
	}

	// Ищем слово пользователя по тексту из словаря
	return s.GetByText(ctx, dictWord.Text)
}
