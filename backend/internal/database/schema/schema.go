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

// InboxItemTable определяет схему таблицы inbox_items
type InboxItemTable struct {
	Name          Table
	ID            Column
	Text          Column
	SourceContext Column
	CreatedAt     Column
}

var InboxItems = InboxItemTable{
	Name:          "inbox_items",
	ID:            "inbox_items.id",
	Text:          "inbox_items.text",
	SourceContext: "inbox_items.source_context",
	CreatedAt:     "inbox_items.created_at",
}

func (t InboxItemTable) All() []string {
	return []string{
		string(t.ID),
		string(t.Text),
		string(t.SourceContext),
		string(t.CreatedAt),
	}
}

func (t InboxItemTable) InsertColumns() []string {
	return []string{"text", "source_context", "created_at"}
}

// TranslationTable определяет схему таблицы translations
type TranslationTable struct {
	Name          Table
	ID            Column
	MeaningID     Column
	TranslationRu Column
	CreatedAt     Column
}

var Translations = TranslationTable{
	Name:          "translations",
	ID:            "translations.id",
	MeaningID:     "translations.meaning_id",
	TranslationRu: "translations.translation_ru",
	CreatedAt:     "translations.created_at",
}

func (t TranslationTable) All() []string {
	return []string{
		string(t.ID),
		string(t.MeaningID),
		string(t.TranslationRu),
		string(t.CreatedAt),
	}
}

func (t TranslationTable) InsertColumns() []string {
	return []string{"meaning_id", "translation_ru", "created_at"}
}

// DictionaryWordTable определяет схему таблицы dictionary_words
type DictionaryWordTable struct {
	Name          Table
	ID            Column
	Text          Column
	Transcription Column
	AudioURL      Column
	FrequencyRank Column
	Source        Column
	SourceID      Column
	CreatedAt     Column
	UpdatedAt     Column
}

var DictionaryWords = DictionaryWordTable{
	Name:          "dictionary_words",
	ID:            "dictionary_words.id",
	Text:          "dictionary_words.text",
	Transcription: "dictionary_words.transcription",
	AudioURL:      "dictionary_words.audio_url",
	FrequencyRank: "dictionary_words.frequency_rank",
	Source:        "dictionary_words.source",
	SourceID:      "dictionary_words.source_id",
	CreatedAt:     "dictionary_words.created_at",
	UpdatedAt:     "dictionary_words.updated_at",
}

func (t DictionaryWordTable) All() []string {
	return []string{
		string(t.ID),
		string(t.Text),
		string(t.Transcription),
		string(t.AudioURL),
		string(t.FrequencyRank),
		string(t.Source),
		string(t.SourceID),
		string(t.CreatedAt),
		string(t.UpdatedAt),
	}
}

func (t DictionaryWordTable) InsertColumns() []string {
	return []string{"text", "transcription", "audio_url", "frequency_rank", "source", "source_id", "created_at", "updated_at"}
}

// DictionaryMeaningTable определяет схему таблицы dictionary_meanings
type DictionaryMeaningTable struct {
	Name            Table
	ID              Column
	DictionaryWordID Column
	PartOfSpeech    Column
	DefinitionEn    Column
	CefrLevel       Column
	ImageURL        Column
	OrderIndex      Column
	CreatedAt       Column
	UpdatedAt       Column
}

var DictionaryMeanings = DictionaryMeaningTable{
	Name:            "dictionary_meanings",
	ID:              "dictionary_meanings.id",
	DictionaryWordID: "dictionary_meanings.dictionary_word_id",
	PartOfSpeech:    "dictionary_meanings.part_of_speech",
	DefinitionEn:    "dictionary_meanings.definition_en",
	CefrLevel:       "dictionary_meanings.cefr_level",
	ImageURL:        "dictionary_meanings.image_url",
	OrderIndex:      "dictionary_meanings.order_index",
	CreatedAt:       "dictionary_meanings.created_at",
	UpdatedAt:       "dictionary_meanings.updated_at",
}

func (t DictionaryMeaningTable) All() []string {
	return []string{
		string(t.ID),
		string(t.DictionaryWordID),
		string(t.PartOfSpeech),
		string(t.DefinitionEn),
		string(t.CefrLevel),
		string(t.ImageURL),
		string(t.OrderIndex),
		string(t.CreatedAt),
		string(t.UpdatedAt),
	}
}

func (t DictionaryMeaningTable) InsertColumns() []string {
	return []string{"dictionary_word_id", "part_of_speech", "definition_en", "cefr_level", "image_url", "order_index", "created_at", "updated_at"}
}

// DictionaryTranslationTable определяет схему таблицы dictionary_translations
type DictionaryTranslationTable struct {
	Name                Table
	ID                  Column
	DictionaryMeaningID Column
	TranslationRu       Column
	CreatedAt           Column
}

var DictionaryTranslations = DictionaryTranslationTable{
	Name:                "dictionary_translations",
	ID:                  "dictionary_translations.id",
	DictionaryMeaningID: "dictionary_translations.dictionary_meaning_id",
	TranslationRu:       "dictionary_translations.translation_ru",
	CreatedAt:           "dictionary_translations.created_at",
}

func (t DictionaryTranslationTable) All() []string {
	return []string{
		string(t.ID),
		string(t.DictionaryMeaningID),
		string(t.TranslationRu),
		string(t.CreatedAt),
	}
}

func (t DictionaryTranslationTable) InsertColumns() []string {
	return []string{"dictionary_meaning_id", "translation_ru", "created_at"}
}

// DictionaryWordFormTable определяет схему таблицы dictionary_word_forms
type DictionaryWordFormTable struct {
	Name            Table
	ID              Column
	DictionaryWordID Column
	FormText        Column
	FormType        Column
	CreatedAt       Column
	UpdatedAt       Column
}

var DictionaryWordForms = DictionaryWordFormTable{
	Name:            "dictionary_word_forms",
	ID:              "dictionary_word_forms.id",
	DictionaryWordID: "dictionary_word_forms.dictionary_word_id",
	FormText:        "dictionary_word_forms.form_text",
	FormType:        "dictionary_word_forms.form_type",
	CreatedAt:       "dictionary_word_forms.created_at",
	UpdatedAt:       "dictionary_word_forms.updated_at",
}

func (t DictionaryWordFormTable) All() []string {
	return []string{
		string(t.ID),
		string(t.DictionaryWordID),
		string(t.FormText),
		string(t.FormType),
		string(t.CreatedAt),
		string(t.UpdatedAt),
	}
}

func (t DictionaryWordFormTable) InsertColumns() []string {
	return []string{"dictionary_word_id", "form_text", "form_type", "created_at", "updated_at"}
}

// DictionarySynonymAntonymTable определяет схему таблицы dictionary_synonyms_antonyms
type DictionarySynonymAntonymTable struct {
	Name         Table
	ID           Column
	MeaningID1   Column
	MeaningID2   Column
	RelationType Column
	CreatedAt    Column
	UpdatedAt    Column
}

var DictionarySynonymsAntonyms = DictionarySynonymAntonymTable{
	Name:         "dictionary_synonyms_antonyms",
	ID:           "dictionary_synonyms_antonyms.id",
	MeaningID1:   "dictionary_synonyms_antonyms.meaning_id_1",
	MeaningID2:   "dictionary_synonyms_antonyms.meaning_id_2",
	RelationType: "dictionary_synonyms_antonyms.relation_type",
	CreatedAt:    "dictionary_synonyms_antonyms.created_at",
	UpdatedAt:    "dictionary_synonyms_antonyms.updated_at",
}

func (t DictionarySynonymAntonymTable) All() []string {
	return []string{
		string(t.ID),
		string(t.MeaningID1),
		string(t.MeaningID2),
		string(t.RelationType),
		string(t.CreatedAt),
		string(t.UpdatedAt),
	}
}

func (t DictionarySynonymAntonymTable) InsertColumns() []string {
	return []string{"meaning_id_1", "meaning_id_2", "relation_type", "created_at", "updated_at"}
}
