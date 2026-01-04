package dictionary

import "github.com/heartmarshall/my-english/internal/model"

// ImportedWord — структура, представляющая слово, полученное из внешнего источника.
// Используется для импорта в базу данных.
type ImportedWord struct {
	Text           string
	Pronunciations []ImportedPronunciation
	Senses         []ImportedSense
}

type ImportedPronunciation struct {
	AudioURL      string
	Transcription string
	Region        model.AccentRegion
}

type ImportedSense struct {
	PartOfSpeech model.PartOfSpeech
	Definition   string
	Translations []string
	Examples     []ImportedExample
}

type ImportedExample struct {
	SentenceEn string
	SentenceRu string
}
