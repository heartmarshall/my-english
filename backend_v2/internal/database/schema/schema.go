// Package schema содержит определения таблиц и колонок БД.
// Используется для type-safe построения SQL запросов через squirrel.
package schema

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

// Table описывает таблицу базы данных.
type Table string

func (t Table) String() string { return string(t) }

// Column описывает колонку с опциональным префиксом таблицы.
type Column string

func (c Column) String() string { return string(c) }

// Qualified возвращает полное имя колонки "table.column".
func (c Column) Qualified() string { return string(c) }

// Bare возвращает "голое" имя колонки без таблицы.
func (c Column) Bare() string {
	parts := strings.Split(string(c), ".")
	if len(parts) > 1 {
		return parts[1]
	}
	return string(c)
}

// --- Fluent API для условий WHERE ---

func (c Column) Eq(val any) squirrel.Eq          { return squirrel.Eq{string(c): val} }
func (c Column) NotEq(val any) squirrel.NotEq    { return squirrel.NotEq{string(c): val} }
func (c Column) Lt(val any) squirrel.Lt          { return squirrel.Lt{string(c): val} }
func (c Column) LtOrEq(val any) squirrel.LtOrEq  { return squirrel.LtOrEq{string(c): val} }
func (c Column) Gt(val any) squirrel.Gt          { return squirrel.Gt{string(c): val} }
func (c Column) GtOrEq(val any) squirrel.GtOrEq  { return squirrel.GtOrEq{string(c): val} }
func (c Column) In(vals any) squirrel.Eq         { return squirrel.Eq{string(c): vals} }
func (c Column) NotIn(vals any) squirrel.NotEq   { return squirrel.NotEq{string(c): vals} }
func (c Column) Like(pat string) squirrel.Like   { return squirrel.Like{string(c): pat} }
func (c Column) ILike(pat string) squirrel.ILike { return squirrel.ILike{string(c): pat} }
func (c Column) IsNull() squirrel.Eq             { return squirrel.Eq{string(c): nil} }
func (c Column) IsNotNull() squirrel.NotEq       { return squirrel.NotEq{string(c): nil} }

// --- Fluent API для ORDER BY ---

func (c Column) Desc() string { return string(c) + " DESC" }
func (c Column) Asc() string  { return string(c) + " ASC" }

// ============================================================================
// SYSTEM LAYER
// ============================================================================

type DataSourcesTable struct {
	Name        Table
	ID          Column
	Slug        Column
	DisplayName Column
	TrustLevel  Column
	WebsiteURL  Column
	CreatedAt   Column
}

var DataSources = DataSourcesTable{
	Name:        "data_sources",
	ID:          "data_sources.id",
	Slug:        "data_sources.slug",
	DisplayName: "data_sources.display_name",
	TrustLevel:  "data_sources.trust_level",
	WebsiteURL:  "data_sources.website_url",
	CreatedAt:   "data_sources.created_at",
}

func (t DataSourcesTable) Columns() []string {
	return []string{
		string(t.ID), string(t.Slug), string(t.DisplayName),
		string(t.TrustLevel), string(t.WebsiteURL), string(t.CreatedAt),
	}
}

func (t DataSourcesTable) InsertColumns() []string {
	return []string{"slug", "display_name", "trust_level", "website_url"}
}

// ============================================================================
// LINGUISTIC LAYER
// ============================================================================

// --- Lexemes ---

type LexemesTable struct {
	Name           Table
	ID             Column
	TextNormalized Column
	TextDisplay    Column
	CreatedAt      Column
}

var Lexemes = LexemesTable{
	Name:           "lexemes",
	ID:             "lexemes.id",
	TextNormalized: "lexemes.text_normalized",
	TextDisplay:    "lexemes.text_display",
	CreatedAt:      "lexemes.created_at",
}

func (t LexemesTable) Columns() []string {
	return []string{
		string(t.ID), string(t.TextNormalized),
		string(t.TextDisplay), string(t.CreatedAt),
	}
}

func (t LexemesTable) InsertColumns() []string {
	return []string{"text_normalized", "text_display"}
}

type PronunciationsTable struct {
	Name          Table
	ID            Column
	LexemeID      Column
	AudioURL      Column
	Transcription Column
	Region        Column
	SourceID      Column
}

var Pronunciations = PronunciationsTable{
	Name:          "pronunciations",
	ID:            "pronunciations.id",
	LexemeID:      "pronunciations.lexeme_id",
	AudioURL:      "pronunciations.audio_url",
	Transcription: "pronunciations.transcription",
	Region:        "pronunciations.region",
	SourceID:      "pronunciations.source_id",
}

func (t PronunciationsTable) Columns() []string {
	return []string{
		string(t.ID), string(t.LexemeID), string(t.AudioURL),
		string(t.Transcription), string(t.Region), string(t.SourceID),
	}
}

func (t PronunciationsTable) InsertColumns() []string {
	return []string{"lexeme_id", "audio_url", "transcription", "region", "source_id"}
}

// --- Inflections ---

type InflectionsTable struct {
	Name              Table
	InflectedLexemeID Column
	LemmaLexemeID     Column
	Type              Column
}

var Inflections = InflectionsTable{
	Name:              "inflections",
	InflectedLexemeID: "inflections.inflected_lexeme_id",
	LemmaLexemeID:     "inflections.lemma_lexeme_id",
	Type:              "inflections.type",
}

func (t InflectionsTable) Columns() []string {
	return []string{
		string(t.InflectedLexemeID), string(t.LemmaLexemeID), string(t.Type),
	}
}

func (t InflectionsTable) InsertColumns() []string {
	return []string{"inflected_lexeme_id", "lemma_lexeme_id", "type"}
}

// --- Senses ---

type SensesTable struct {
	Name          Table
	ID            Column
	LexemeID      Column
	PartOfSpeech  Column
	Definition    Column
	CefrLevel     Column
	SourceID      Column
	ExternalRefID Column
	CreatedAt     Column
}

var Senses = SensesTable{
	Name:          "senses",
	ID:            "senses.id",
	LexemeID:      "senses.lexeme_id",
	PartOfSpeech:  "senses.part_of_speech",
	Definition:    "senses.definition",
	CefrLevel:     "senses.cefr_level",
	SourceID:      "senses.source_id",
	ExternalRefID: "senses.external_ref_id",
	CreatedAt:     "senses.created_at",
}

func (t SensesTable) Columns() []string {
	return []string{
		string(t.ID), string(t.LexemeID), string(t.PartOfSpeech),
		string(t.Definition), string(t.CefrLevel), string(t.SourceID),
		string(t.ExternalRefID), string(t.CreatedAt),
	}
}

func (t SensesTable) InsertColumns() []string {
	return []string{"lexeme_id", "part_of_speech", "definition", "cefr_level", "source_id", "external_ref_id"}
}

// --- SenseTranslations ---

type SenseTranslationsTable struct {
	Name        Table
	ID          Column
	SenseID     Column
	Translation Column
	SourceID    Column
}

var SenseTranslations = SenseTranslationsTable{
	Name:        "sense_translations",
	ID:          "sense_translations.id",
	SenseID:     "sense_translations.sense_id",
	Translation: "sense_translations.translation",
	SourceID:    "sense_translations.source_id",
}

func (t SenseTranslationsTable) Columns() []string {
	return []string{
		string(t.ID), string(t.SenseID), string(t.Translation), string(t.SourceID),
	}
}

func (t SenseTranslationsTable) InsertColumns() []string {
	return []string{"sense_id", "translation", "source_id"}
}

// --- SenseRelations ---

type SenseRelationsTable struct {
	Name            Table
	SourceSenseID   Column
	TargetSenseID   Column
	Type            Column
	IsBidirectional Column
	SourceID        Column
}

var SenseRelations = SenseRelationsTable{
	Name:            "sense_relations",
	SourceSenseID:   "sense_relations.source_sense_id",
	TargetSenseID:   "sense_relations.target_sense_id",
	Type:            "sense_relations.type",
	IsBidirectional: "sense_relations.is_bidirectional",
	SourceID:        "sense_relations.source_id",
}

func (t SenseRelationsTable) Columns() []string {
	return []string{
		string(t.SourceSenseID), string(t.TargetSenseID),
		string(t.Type), string(t.IsBidirectional), string(t.SourceID),
	}
}

func (t SenseRelationsTable) InsertColumns() []string {
	return []string{"source_sense_id", "target_sense_id", "type", "is_bidirectional", "source_id"}
}

// --- Examples ---

type ExamplesTable struct {
	Name            Table
	ID              Column
	SenseID         Column
	SentenceEn      Column
	SentenceRu      Column
	TargetWordRange Column
	SourceName      Column
}

var Examples = ExamplesTable{
	Name:            "examples",
	ID:              "examples.id",
	SenseID:         "examples.sense_id",
	SentenceEn:      "examples.sentence_en",
	SentenceRu:      "examples.sentence_ru",
	TargetWordRange: "examples.target_word_range",
	SourceName:      "examples.source_name",
}

func (t ExamplesTable) Columns() []string {
	return []string{
		string(t.ID), string(t.SenseID), string(t.SentenceEn),
		string(t.SentenceRu), string(t.TargetWordRange), string(t.SourceName),
	}
}

func (t ExamplesTable) InsertColumns() []string {
	return []string{"sense_id", "sentence_en", "sentence_ru", "target_word_range", "source_name"}
}

// ============================================================================
// USER LAYER
// ============================================================================

// --- InboxItems ---

type InboxItemsTable struct {
	Name        Table
	ID          Column
	RawText     Column
	ContextNote Column
	CreatedAt   Column
}

var InboxItems = InboxItemsTable{
	Name:        "inbox_items",
	ID:          "inbox_items.id",
	RawText:     "inbox_items.raw_text",
	ContextNote: "inbox_items.context_note",
	CreatedAt:   "inbox_items.created_at",
}

func (t InboxItemsTable) Columns() []string {
	return []string{
		string(t.ID), string(t.RawText), string(t.ContextNote), string(t.CreatedAt),
	}
}

func (t InboxItemsTable) InsertColumns() []string {
	return []string{"raw_text", "context_note"}
}

// --- Tags ---

type TagsTable struct {
	Name     Table
	ID       Column
	NameCol  Column
	ColorHex Column
}

var Tags = TagsTable{
	Name:     "tags",
	ID:       "tags.id",
	NameCol:  "tags.name",
	ColorHex: "tags.color_hex",
}

func (t TagsTable) Columns() []string {
	return []string{string(t.ID), string(t.NameCol), string(t.ColorHex)}
}

func (t TagsTable) InsertColumns() []string {
	return []string{"name", "color_hex"}
}

// --- Cards ---

type CardsTable struct {
	Name                Table
	ID                  Column
	SenseID             Column
	CustomText          Column
	CustomTranscription Column
	CustomTranslations  Column
	CustomNote          Column
	CustomImageURL      Column
	CreatedAt           Column
	UpdatedAt           Column
	IsDeleted           Column
}

var Cards = CardsTable{
	Name:                "cards",
	ID:                  "cards.id",
	SenseID:             "cards.sense_id",
	CustomText:          "cards.custom_text",
	CustomTranscription: "cards.custom_transcription",
	CustomTranslations:  "cards.custom_translations",
	CustomNote:          "cards.custom_note",
	CustomImageURL:      "cards.custom_image_url",
	CreatedAt:           "cards.created_at",
	UpdatedAt:           "cards.updated_at",
	IsDeleted:           "cards.is_deleted",
}

func (t CardsTable) Columns() []string {
	return []string{
		string(t.ID), string(t.SenseID), string(t.CustomText),
		string(t.CustomTranscription), string(t.CustomTranslations),
		string(t.CustomNote), string(t.CustomImageURL),
		string(t.CreatedAt), string(t.UpdatedAt), string(t.IsDeleted),
	}
}

func (t CardsTable) InsertColumns() []string {
	return []string{
		"sense_id", "custom_text", "custom_transcription",
		"custom_translations", "custom_note", "custom_image_url",
	}
}

// --- CardTags ---

type CardTagsTable struct {
	Name   Table
	CardID Column
	TagID  Column
}

var CardTags = CardTagsTable{
	Name:   "card_tags",
	CardID: "card_tags.card_id",
	TagID:  "card_tags.tag_id",
}

func (t CardTagsTable) Columns() []string {
	return []string{string(t.CardID), string(t.TagID)}
}

func (t CardTagsTable) InsertColumns() []string {
	return []string{"card_id", "tag_id"}
}

// --- SRSStates ---

type SRSStatesTable struct {
	Name          Table
	CardID        Column
	Status        Column
	DueDate       Column
	AlgorithmData Column
	LastReviewAt  Column
}

var SRSStates = SRSStatesTable{
	Name:          "srs_states",
	CardID:        "srs_states.card_id",
	Status:        "srs_states.status",
	DueDate:       "srs_states.due_date",
	AlgorithmData: "srs_states.algorithm_data",
	LastReviewAt:  "srs_states.last_review_at",
}

func (t SRSStatesTable) Columns() []string {
	return []string{
		string(t.CardID), string(t.Status), string(t.DueDate),
		string(t.AlgorithmData), string(t.LastReviewAt),
	}
}

func (t SRSStatesTable) InsertColumns() []string {
	return []string{"card_id", "status", "due_date", "algorithm_data", "last_review_at"}
}

// --- ReviewLogs ---

type ReviewLogsTable struct {
	Name        Table
	ID          Column
	CardID      Column
	Grade       Column
	DurationMs  Column
	ReviewedAt  Column
	StateBefore Column
	StateAfter  Column
}

var ReviewLogs = ReviewLogsTable{
	Name:        "review_logs",
	ID:          "review_logs.id",
	CardID:      "review_logs.card_id",
	Grade:       "review_logs.grade",
	DurationMs:  "review_logs.duration_ms",
	ReviewedAt:  "review_logs.reviewed_at",
	StateBefore: "review_logs.state_before",
	StateAfter:  "review_logs.state_after",
}

func (t ReviewLogsTable) Columns() []string {
	return []string{
		string(t.ID), string(t.CardID), string(t.Grade),
		string(t.DurationMs), string(t.ReviewedAt),
		string(t.StateBefore), string(t.StateAfter),
	}
}

func (t ReviewLogsTable) InsertColumns() []string {
	return []string{"card_id", "grade", "duration_ms", "state_before", "state_after"}
}
