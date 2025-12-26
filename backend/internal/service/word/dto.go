package word

import "github.com/heartmarshall/my-english/internal/model"

// CreateWordInput — входные данные для создания слова.
type CreateWordInput struct {
	Text          string
	Transcription *string
	AudioURL      *string
	SourceContext *string
	Meanings      []CreateMeaningInput
}

// CreateMeaningInput — входные данные для создания значения.
type CreateMeaningInput struct {
	PartOfSpeech model.PartOfSpeech
	DefinitionEn *string
	Translations []string // Множественные переводы (вместо одного TranslationRu)
	CefrLevel    *string
	ImageURL     *string
	Examples     []CreateExampleInput
	Tags         []string // Имена тегов, не ID
}

// CreateExampleInput — входные данные для создания примера.
type CreateExampleInput struct {
	SentenceEn string
	SentenceRu *string
	SourceName *model.ExampleSource
}

// UpdateWordInput — входные данные для обновления слова.
type UpdateWordInput struct {
	Text          string
	Transcription *string
	AudioURL      *string
	SourceContext *string
	Meanings      []UpdateMeaningInput
}

// UpdateMeaningInput — входные данные для обновления значения.
type UpdateMeaningInput struct {
	ID           *int64 // Если nil — создаём новое, иначе обновляем
	PartOfSpeech model.PartOfSpeech
	DefinitionEn *string
	Translations []string // Множественные переводы (вместо одного TranslationRu)
	CefrLevel    *string
	ImageURL     *string
	Examples     []CreateExampleInput // Примеры пересоздаются полностью
	Tags         []string
}

// WordFilter — параметры фильтрации для списка слов.
type WordFilter struct {
	Search *string
	Status *model.LearningStatus
	Tags   []string
}

// WordWithRelations — слово со всеми связанными данными.
type WordWithRelations struct {
	Word     model.Word
	Meanings []MeaningWithRelations
}

// MeaningWithRelations — значение со всеми связанными данными.
type MeaningWithRelations struct {
	Meaning      model.Meaning
	Translations []model.Translation // Переводы из таблицы translations
	Examples     []model.Example
	Tags         []model.Tag
}

// Suggestion представляет подсказку для автокомплита.
type Suggestion struct {
	Text           string
	Transcription  *string
	Translations   []string
	Origin         string // "LOCAL" или "DICTIONARY"
	ExistingWordID *int64
}
