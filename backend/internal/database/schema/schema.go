package schema

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

// Table описывает таблицу базы данных.
type Table string

func (t Table) String() string { return string(t) }

// Column описывает колонку.
type Column string

func (c Column) String() string { return string(c) }

// definition описывает метаданные колонки (имя и таблицу).
// Это нужно, чтобы генерировать "table.column".
type definition struct {
	table string
	name  string
}

// Qualified возвращает полное имя колонки "table.column".
// Используй это в WHERE и ORDER BY, чтобы избежать ambiguous errors.
func (c Column) Qualified() string {
	return string(c)
}

// Bare возвращает "голое" имя колонки без таблицы (например, для AS в SELECT).
// В текущей реализации c хранит уже полное имя, так что нужно парсить,
// либо хранить структуру сложнее. Для простоты будем считать, что c — это всегда table.col
func (c Column) Bare() string {
	parts := strings.Split(string(c), ".")
	if len(parts) > 1 {
		return parts[1]
	}
	return string(c)
}

// --- Fluent API Helpers ---

func (c Column) Eq(val any) squirrel.Eq  { return squirrel.Eq{string(c): val} }
func (c Column) Neq(val any) squirrel.Eq { return squirrel.Eq{string(c): val} } // Not Equal (map syntax handles this differently usually, strictly Eq is map)
// squirrel.NotEq{"col": val}
func (c Column) Lt(val any) squirrel.Lt          { return squirrel.Lt{string(c): val} }
func (c Column) Gt(val any) squirrel.Gt          { return squirrel.Gt{string(c): val} }
func (c Column) In(vals any) squirrel.Eq         { return squirrel.Eq{string(c): vals} }
func (c Column) ILike(pat string) squirrel.ILike { return squirrel.ILike{string(c): pat} }

func (c Column) Desc() string { return string(c) + " DESC" }
func (c Column) Asc() string  { return string(c) + " ASC" }

// --- Tables Definitions ---

// WordTable определяет схему таблицы words
type WordTable struct {
	Name          Table
	ID            Column
	Text          Column
	Transcription Column
	AudioURL      Column
	FrequencyRank Column
	CreatedAt     Column
}

// Метод All возвращает все колонки таблицы.
// Больше не нужно руками собирать var columns = []string{...} в репозиториях!
func (t WordTable) All() []string {
	return []string{
		string(t.ID),
		string(t.Text),
		string(t.Transcription),
		string(t.AudioURL),
		string(t.FrequencyRank),
		string(t.CreatedAt),
	}
}

// Инициализация с префиксами таблиц
var Words = WordTable{
	Name:          "words",
	ID:            "words.id",
	Text:          "words.text",
	Transcription: "words.transcription",
	AudioURL:      "words.audio_url",
	FrequencyRank: "words.frequency_rank",
	CreatedAt:     "words.created_at",
}

// MeaningTable определяет схему таблицы meanings
type MeaningTable struct {
	Name           Table
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
}

var Meanings = MeaningTable{
	Name:           "meanings",
	ID:             "meanings.id",
	WordID:         "meanings.word_id",
	TranslationRu:  "meanings.translation_ru",
	PartOfSpeech:   "meanings.part_of_speech",
	DefinitionEn:   "meanings.definition_en",
	CefrLevel:      "meanings.cefr_level",
	ImageURL:       "meanings.image_url",
	LearningStatus: "meanings.learning_status",
	NextReviewAt:   "meanings.next_review_at",
	Interval:       "meanings.interval",
	EaseFactor:     "meanings.ease_factor",
	ReviewCount:    "meanings.review_count",
	CreatedAt:      "meanings.created_at",
	UpdatedAt:      "meanings.updated_at",
}

func (t MeaningTable) All() []string {
	return []string{
		string(t.ID),
		string(t.WordID),
		string(t.PartOfSpeech),
		string(t.DefinitionEn),
		string(t.TranslationRu),
		string(t.CefrLevel),
		string(t.ImageURL),
		string(t.LearningStatus),
		string(t.NextReviewAt),
		string(t.Interval),
		string(t.EaseFactor),
		string(t.ReviewCount),
		string(t.CreatedAt),
		string(t.UpdatedAt),
	}
}

// Helper для получения списка колонок без префикса (если нужно для INSERT)
// INSERT INTO table (col) VALUES (...) - тут нельзя писать table.col
func (t WordTable) InsertColumns() []string {
	return []string{"text", "transcription", "audio_url", "frequency_rank", "created_at"}
}

func (t MeaningTable) InsertColumns() []string {
	return []string{
		"word_id", "part_of_speech", "definition_en", "translation_ru",
		"cefr_level", "image_url", "learning_status", "next_review_at",
		"interval", "ease_factor", "review_count", "created_at", "updated_at",
	}
}

// TagTable определяет схему таблицы tags
type TagTable struct {
	Name    Table
	ID      Column
	NameCol Column
}

var Tags = TagTable{
	Name:    "tags",
	ID:      "tags.id",
	NameCol: "tags.name",
}

func (t TagTable) All() []string {
	return []string{
		string(t.ID),
		string(t.NameCol),
	}
}

func (t TagTable) InsertColumns() []string {
	return []string{"name"}
}

// ExampleTable определяет схему таблицы examples
type ExampleTable struct {
	Name       Table
	ID         Column
	MeaningID  Column
	SentenceEn Column
	SentenceRu Column
	SourceName Column
}

var Examples = ExampleTable{
	Name:       "examples",
	ID:         "examples.id",
	MeaningID:  "examples.meaning_id",
	SentenceEn: "examples.sentence_en",
	SentenceRu: "examples.sentence_ru",
	SourceName: "examples.source_name",
}

func (t ExampleTable) All() []string {
	return []string{
		string(t.ID),
		string(t.MeaningID),
		string(t.SentenceEn),
		string(t.SentenceRu),
		string(t.SourceName),
	}
}

func (t ExampleTable) InsertColumns() []string {
	return []string{"meaning_id", "sentence_en", "sentence_ru", "source_name"}
}

// MeaningTagTable определяет схему таблицы meanings_tags
type MeaningTagTable struct {
	Name      Table
	MeaningID Column
	TagID     Column
}

var MeaningTags = MeaningTagTable{
	Name:      "meanings_tags",
	MeaningID: "meanings_tags.meaning_id",
	TagID:     "meanings_tags.tag_id",
}

func (t MeaningTagTable) All() []string {
	return []string{
		string(t.MeaningID),
		string(t.TagID),
	}
}

func (t MeaningTagTable) InsertColumns() []string {
	return []string{"meaning_id", "tag_id"}
}
