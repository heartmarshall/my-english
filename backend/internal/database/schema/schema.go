package schema

import (
	"github.com/Masterminds/squirrel"
)

// Table — описание таблицы для использования в Squirrel
type Table string

// String реализует интерфейс fmt.Stringer (и подходит для squirrel)
func (t Table) String() string { return string(t) }

// Column — описание колонки
type Column string

func (c Column) String() string { return string(c) }

// --- Вспомогательные функции для работы с Squirrel ---

// Eq создает условие равенства для Where с использованием Column
func (c Column) Eq(value interface{}) squirrel.Eq {
	return squirrel.Eq{c.String(): value}
}

// In создает условие IN для Where с использованием Column
func (c Column) In(values interface{}) squirrel.Eq {
	return squirrel.Eq{c.String(): values}
}

// ILike создает условие ILIKE для Where с использованием Column
func (c Column) ILike(pattern string) squirrel.ILike {
	return squirrel.ILike{c.String(): pattern}
}

// Lt создает условие < для Where с использованием Column
func (c Column) Lt(value interface{}) squirrel.Lt {
	return squirrel.Lt{c.String(): value}
}

// OrderByASC создает строку для OrderBy с ASC
func (c Column) OrderByASC() string {
	return c.String() + " ASC"
}

// OrderByDESC создает строку для OrderBy с DESC
func (c Column) OrderByDESC() string {
	return c.String() + " DESC"
}

// OrderBy создает строку для OrderBy с указанным направлением
func (c Column) OrderBy(direction string) string {
	return c.String() + " " + direction
}

// Returning создает строку для RETURNING
func (c Column) Returning() string {
	return "RETURNING " + c.String()
}

// ColumnsToStrings преобразует список Column в []string
func ColumnsToStrings(cols ...Column) []string {
	result := make([]string, len(cols))
	for i, col := range cols {
		result[i] = col.String()
	}
	return result
}

// Определяем таблицы
const (
	Words       Table = "words"
	Meanings    Table = "meanings"
	Examples    Table = "examples"
	Tags        Table = "tags"
	MeaningTags Table = "meanings_tags"
)

var (
	WordColumns = struct {
		ID            Column
		Text          Column
		Transcription Column
		AudioURL      Column
		FrequencyRank Column
		CreatedAt     Column
	}{
		ID:            "id",
		Text:          "text",
		Transcription: "transcription",
		AudioURL:      "audio_url",
		FrequencyRank: "frequency_rank",
		CreatedAt:     "created_at",
	}

	MeaningColumns = struct {
		ID             Column
		WordID         Column
		TranslationRu  Column
		PartOfSpeech   Column
		DefinitionEn   Column
		CefrLevel      Column
		ImageURL       Column
		LearningStatus Column
		NextReviewAt   Column
		Interval       Column
		EaseFactor     Column
		ReviewCount    Column
		CreatedAt      Column
		UpdatedAt      Column
	}{
		ID:             "id",
		WordID:         "word_id",
		TranslationRu:  "translation_ru",
		PartOfSpeech:   "part_of_speech",
		DefinitionEn:   "definition_en",
		CefrLevel:      "cefr_level",
		ImageURL:       "image_url",
		LearningStatus: "learning_status",
		NextReviewAt:   "next_review_at",
		Interval:       "interval",
		EaseFactor:     "ease_factor",
		ReviewCount:    "review_count",
		CreatedAt:      "created_at",
		UpdatedAt:      "updated_at",
	}

	ExampleColumns = struct {
		ID         Column
		MeaningID  Column
		SentenceEn Column
		SentenceRu Column
		SourceName Column
		CreatedAt  Column
		UpdatedAt  Column
	}{
		ID:         "id",
		MeaningID:  "meaning_id",
		SentenceEn: "sentence_en",
		SentenceRu: "sentence_ru",
		SourceName: "source_name",
		CreatedAt:  "created_at",
		UpdatedAt:  "updated_at",
	}

	TagColumns = struct {
		ID        Column
		Name      Column
		CreatedAt Column
		UpdatedAt Column
	}{
		ID:        "id",
		Name:      "name",
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
	}

	MeaningTagColumns = struct {
		MeaningID Column
		TagID     Column
	}{
		MeaningID: "meaning_id",
		TagID:     "tag_id",
	}
)
