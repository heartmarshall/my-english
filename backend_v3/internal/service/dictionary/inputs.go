package dictionary

import "github.com/heartmarshall/my-english/internal/model"

// CreateWordInput — полный набор данных для создания слова.
type CreateWordInput struct {
	Text           string
	Senses         []SenseInput
	Images         []ImageInput
	Pronunciations []PronunciationInput
	CreateCard     bool
}

type SenseInput struct {
	Definition   *string
	PartOfSpeech *model.PartOfSpeech
	SourceSlug   string
	Translations []TranslationInput
	Examples     []ExampleInput
}

type TranslationInput struct {
	Text       string
	SourceSlug string
}

type ExampleInput struct {
	Sentence    string
	Translation *string
	SourceSlug  string
}

type ImageInput struct {
	URL        string
	Caption    *string
	SourceSlug string
}

type PronunciationInput struct {
	AudioURL      string
	Transcription *string
	Region        *string
	SourceSlug    string
}

// UpdateWordInput — полный набор данных для обновления слова.
type UpdateWordInput struct {
	ID             string // UUID слова
	Text           *string
	Senses         []SenseInput
	Images         []ImageInput
	Pronunciations []PronunciationInput
}

// DeleteWordInput — входные данные для удаления слова.
type DeleteWordInput struct {
	ID string // UUID слова
}

// AddSenseInput — входные данные для добавления нового смысла к записи.
type AddSenseInput struct {
	EntryID      string // UUID записи словаря
	Definition   *string
	PartOfSpeech *model.PartOfSpeech
	SourceSlug   string
	Translations []TranslationInput
	Examples     []ExampleInput
}

// AddExamplesInput — входные данные для добавления примеров к смыслу.
type AddExamplesInput struct {
	SenseID  string // UUID смысла
	Examples []ExampleInput
}

// AddTranslationsInput — входные данные для добавления переводов к смыслу.
type AddTranslationsInput struct {
	SenseID      string // UUID смысла
	Translations []TranslationInput
}

// AddImagesInput — входные данные для добавления изображений к записи.
type AddImagesInput struct {
	EntryID string // UUID записи словаря
	Images  []ImageInput
}

// AddPronunciationsInput — входные данные для добавления произношений к записи.
type AddPronunciationsInput struct {
	EntryID        string // UUID записи словаря
	Pronunciations []PronunciationInput
}

// DeleteSenseInput — входные данные для удаления смысла.
type DeleteSenseInput struct {
	ID string // UUID смысла
}

// DeleteExampleInput — входные данные для удаления примера.
type DeleteExampleInput struct {
	ID string // UUID примера
}

// DeleteTranslationInput — входные данные для удаления перевода.
type DeleteTranslationInput struct {
	ID string // UUID перевода
}

// DeleteImageInput — входные данные для удаления изображения.
type DeleteImageInput struct {
	ID string // UUID изображения
}

// DeletePronunciationInput — входные данные для удаления произношения.
type DeletePronunciationInput struct {
	ID string // UUID произношения
}
