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
