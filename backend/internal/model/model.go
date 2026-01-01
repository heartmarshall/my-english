package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Word представляет модель слова
type Word struct {
	ID            int64     `db:"id"`
	Text          string    `db:"text"`
	Transcription *string   `db:"transcription"`
	AudioURL      *string   `db:"audio_url"`
	FrequencyRank *int      `db:"frequency_rank"`
	CreatedAt     time.Time `db:"created_at"`
}

// SortOrder представляет тип сортировки
type SortOrder string

const (
	SortOrderAlphabetical SortOrder = "alphabetical"
	SortOrderCreatedAt    SortOrder = "created_at"
)

// WordFilter содержит параметры фильтрации при поиске слов
type WordFilter struct {
	Search *string
	Status *LearningStatus
	Tags   []string
	SortBy *SortOrder
}

// LearningStatus представляет статус изучения слова
type LearningStatus string

const (
	LearningStatusNew      LearningStatus = "new"
	LearningStatusLearning LearningStatus = "learning"
	LearningStatusReview   LearningStatus = "review"
	LearningStatusMastered LearningStatus = "mastered"
)

// Value реализует driver.Valuer для LearningStatus
func (ls LearningStatus) Value() (driver.Value, error) {
	return string(ls), nil
}

// Scan реализует sql.Scanner для LearningStatus
func (ls *LearningStatus) Scan(value interface{}) error {
	if value == nil {
		*ls = LearningStatusNew
		return nil
	}
	switch v := value.(type) {
	case string:
		*ls = LearningStatus(v)
	case []byte:
		*ls = LearningStatus(v)
	default:
		return fmt.Errorf("cannot scan %T into LearningStatus", value)
	}
	return nil
}

// IsValid проверяет, является ли статус валидным значением enum.
func (ls LearningStatus) IsValid() bool {
	switch ls {
	case LearningStatusNew, LearningStatusLearning, LearningStatusReview, LearningStatusMastered:
		return true
	}
	return false
}

// PartOfSpeech представляет часть речи
type PartOfSpeech string

const (
	PartOfSpeechNoun      PartOfSpeech = "noun"
	PartOfSpeechVerb      PartOfSpeech = "verb"
	PartOfSpeechAdjective PartOfSpeech = "adjective"
	PartOfSpeechAdverb    PartOfSpeech = "adverb"
	PartOfSpeechOther     PartOfSpeech = "other"
)

// Value реализует driver.Valuer для PartOfSpeech
func (pos PartOfSpeech) Value() (driver.Value, error) {
	return string(pos), nil
}

// Scan реализует sql.Scanner для PartOfSpeech
func (pos *PartOfSpeech) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("PartOfSpeech cannot be nil")
	}
	switch v := value.(type) {
	case string:
		*pos = PartOfSpeech(v)
	case []byte:
		*pos = PartOfSpeech(v)
	default:
		return fmt.Errorf("cannot scan %T into PartOfSpeech", value)
	}
	return nil
}

// IsValid проверяет, является ли часть речи валидным значением enum.
func (pos PartOfSpeech) IsValid() bool {
	switch pos {
	case PartOfSpeechNoun, PartOfSpeechVerb, PartOfSpeechAdjective,
		PartOfSpeechAdverb, PartOfSpeechOther:
		return true
	}
	return false
}

// Meaning представляет модель значения слова из таблицы meanings
type Meaning struct {
	ID             int64          `db:"id"`
	WordID         int64          `db:"word_id"`
	PartOfSpeech   PartOfSpeech   `db:"part_of_speech"`
	DefinitionEn   *string        `db:"definition_en"`
	TranslationRu  string         `db:"translation_ru"`
	CefrLevel      *string        `db:"cefr_level"`
	ImageURL       *string        `db:"image_url"`
	LearningStatus LearningStatus `db:"learning_status"`
	NextReviewAt   *time.Time     `db:"next_review_at"`
	Interval       *int           `db:"interval"`
	EaseFactor     *float64       `db:"ease_factor"`
	ReviewCount    *int           `db:"review_count"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

// Stats представляет статистику по изучению слов
type Stats struct {
	TotalWords        int
	MasteredCount     int
	LearningCount     int
	DueForReviewCount int
}

// ExampleSource представляет источник примера
type ExampleSource string

const (
	ExampleSourceFilm    ExampleSource = "film"
	ExampleSourceBook    ExampleSource = "book"
	ExampleSourceChat    ExampleSource = "chat"
	ExampleSourceVideo   ExampleSource = "video"
	ExampleSourcePodcast ExampleSource = "podcast"
)

// Value реализует driver.Valuer для ExampleSource
func (es ExampleSource) Value() (driver.Value, error) {
	if es == "" {
		return nil, nil
	}
	return string(es), nil
}

// Scan реализует sql.Scanner для ExampleSource
func (es *ExampleSource) Scan(value interface{}) error {
	if value == nil {
		*es = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*es = ExampleSource(v)
	case []byte:
		*es = ExampleSource(v)
	default:
		return fmt.Errorf("cannot scan %T into ExampleSource", value)
	}
	return nil
}

// IsValid проверяет, является ли источник валидным значением enum.
func (es ExampleSource) IsValid() bool {
	switch es {
	case ExampleSourceFilm, ExampleSourceBook, ExampleSourceChat,
		ExampleSourceVideo, ExampleSourcePodcast, "":
		return true
	}
	return false
}

// Example представляет пример использования слова
type Example struct {
	ID         int64          `db:"id"`
	MeaningID  int64          `db:"meaning_id"`
	SentenceEn string         `db:"sentence_en"`
	SentenceRu *string        `db:"sentence_ru"`
	SourceName *ExampleSource `db:"source_name"`
}

// Tag представляет тег для категоризации значений
type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// MeaningTag представляет связь meaning-tag (many-to-many)
type MeaningTag struct {
	MeaningID int64 `db:"meaning_id"`
	TagID     int64 `db:"tag_id"`
}

// InboxItem представляет элемент корзины входящих слов
type InboxItem struct {
	ID            int64     `db:"id"`
	Text          string    `db:"text"`
	SourceContext *string   `db:"source_context"`
	CreatedAt     time.Time `db:"created_at"`
}

// Translation представляет перевод значения слова
type Translation struct {
	ID            int64     `db:"id"`
	MeaningID     int64     `db:"meaning_id"`
	TranslationRu string    `db:"translation_ru"`
	CreatedAt     time.Time `db:"created_at"`
}

// DictionaryWord представляет слово из внутреннего словаря (не пользовательского)
type DictionaryWord struct {
	ID            int64     `db:"id"`
	Text          string    `db:"text"`
	Transcription *string   `db:"transcription"`
	AudioURL      *string   `db:"audio_url"`
	FrequencyRank *int      `db:"frequency_rank"`
	Source        string    `db:"source"`    // Источник: 'free_dictionary', 'oxford', 'custom' и т.д.
	SourceID      *string   `db:"source_id"` // ID слова в источнике
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// DictionaryMeaning представляет значение слова из внутреннего словаря
type DictionaryMeaning struct {
	ID               int64        `db:"id"`
	DictionaryWordID int64        `db:"dictionary_word_id"`
	PartOfSpeech     PartOfSpeech `db:"part_of_speech"`
	DefinitionEn     *string      `db:"definition_en"`
	CefrLevel        *string      `db:"cefr_level"`
	ImageURL         *string      `db:"image_url"`
	OrderIndex       int          `db:"order_index"`
	CreatedAt        time.Time    `db:"created_at"`
	UpdatedAt        time.Time    `db:"updated_at"`
}

// DictionaryTranslation представляет перевод значения из внутреннего словаря
type DictionaryTranslation struct {
	ID                  int64     `db:"id"`
	DictionaryMeaningID int64     `db:"dictionary_meaning_id"`
	TranslationRu       string    `db:"translation_ru"`
	CreatedAt           time.Time `db:"created_at"`
}

// DictionaryWordForm представляет форму слова из внутреннего словаря
// Формы слов: времена глаголов (go, went, gone), множественное число существительных (mouse, mice),
// степени сравнения прилагательных (big, bigger, biggest) и т.д.
type DictionaryWordForm struct {
	ID               int64     `db:"id"`
	DictionaryWordID int64     `db:"dictionary_word_id"`
	FormText         string    `db:"form_text"`
	FormType         *string   `db:"form_type"` // Тип формы: 'past_tense', 'past_participle', 'plural', 'comparative', 'superlative', 'third_person_singular', 'present_participle', 'gerund' и т.д.
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// RelationType представляет тип связи между значениями словаря
type RelationType string

const (
	RelationTypeSynonym RelationType = "synonym"
	RelationTypeAntonym RelationType = "antonym"
)

// Value реализует driver.Valuer для RelationType
func (rt RelationType) Value() (driver.Value, error) {
	return string(rt), nil
}

// Scan реализует sql.Scanner для RelationType
func (rt *RelationType) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("RelationType cannot be nil")
	}
	switch v := value.(type) {
	case string:
		*rt = RelationType(v)
	case []byte:
		*rt = RelationType(v)
	default:
		return fmt.Errorf("cannot scan %T into RelationType", value)
	}
	return nil
}

// IsValid проверяет, является ли тип связи валидным значением.
func (rt RelationType) IsValid() bool {
	switch rt {
	case RelationTypeSynonym, RelationTypeAntonym:
		return true
	}
	return false
}

// DictionarySynonymAntonym представляет связь между значениями словаря (синонимы/антонимы)
type DictionarySynonymAntonym struct {
	ID           int64        `db:"id"`
	MeaningID1   int64        `db:"meaning_id_1"`
	MeaningID2   int64        `db:"meaning_id_2"`
	RelationType RelationType `db:"relation_type"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}
