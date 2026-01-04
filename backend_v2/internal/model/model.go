// Package model содержит доменные модели приложения.
package model

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// SYSTEM LAYER
// ============================================================================

// DataSource представляет источник данных (freedict, user, system).
type DataSource struct {
	ID          int       `db:"id"`
	Slug        string    `db:"slug"`
	DisplayName string    `db:"display_name"`
	TrustLevel  int       `db:"trust_level"`
	WebsiteURL  *string   `db:"website_url"`
	CreatedAt   time.Time `db:"created_at"`
}

// ============================================================================
// LINGUISTIC LAYER
// ============================================================================

// Lexeme представляет слово/фразу в словаре.
type Lexeme struct {
	ID             uuid.UUID `db:"id"`
	TextNormalized string    `db:"text_normalized"`
	TextDisplay    string    `db:"text_display"`
	CreatedAt      time.Time `db:"created_at"`
}

// Pronunciation представляет произношение лексемы.
type Pronunciation struct {
	ID            uuid.UUID `db:"id"`
	LexemeID      uuid.UUID `db:"lexeme_id"`
	AudioURL      string    `db:"audio_url"`
	Transcription *string   `db:"transcription"`
	Region        string    `db:"region"` // us, uk, au, general
	SourceID      *int      `db:"source_id"`
}

// Inflection представляет морфологическую связь между формами слова.
type Inflection struct {
	InflectedLexemeID uuid.UUID `db:"inflected_lexeme_id"`
	LemmaLexemeID     uuid.UUID `db:"lemma_lexeme_id"`
	Type              string    `db:"type"` // plural, past_tense, etc.
}

// Sense представляет смысл/значение слова.
type Sense struct {
	ID            uuid.UUID `db:"id"`
	LexemeID      uuid.UUID `db:"lexeme_id"`
	PartOfSpeech  string    `db:"part_of_speech"`
	Definition    string    `db:"definition"`
	CefrLevel     *string   `db:"cefr_level"`
	SourceID      int       `db:"source_id"`
	ExternalRefID *string   `db:"external_ref_id"`
	CreatedAt     time.Time `db:"created_at"`
}

// SenseTranslation представляет перевод смысла.
type SenseTranslation struct {
	ID          uuid.UUID `db:"id"`
	SenseID     uuid.UUID `db:"sense_id"`
	Translation string    `db:"translation"`
	SourceID    *int      `db:"source_id"`
}

// SenseRelation представляет семантическую связь между смыслами.
type SenseRelation struct {
	SourceSenseID   uuid.UUID `db:"source_sense_id"`
	TargetSenseID   uuid.UUID `db:"target_sense_id"`
	Type            string    `db:"type"` // synonym, antonym, related, collocation
	IsBidirectional bool      `db:"is_bidirectional"`
	SourceID        *int      `db:"source_id"`
}

// Example представляет пример использования слова.
type Example struct {
	ID              uuid.UUID  `db:"id"`
	SenseID         *uuid.UUID `db:"sense_id"`
	SentenceEn      string     `db:"sentence_en"`
	SentenceRu      *string    `db:"sentence_ru"`
	TargetWordRange []int      `db:"target_word_range"` // [start, end]
	SourceName      *string    `db:"source_name"`
}

// ============================================================================
// USER LAYER
// ============================================================================

// InboxItem представляет элемент inbox (GTD).
type InboxItem struct {
	ID          uuid.UUID `db:"id"`
	RawText     string    `db:"raw_text"`
	ContextNote *string   `db:"context_note"`
	CreatedAt   time.Time `db:"created_at"`
}

// Tag представляет пользовательский тег.
type Tag struct {
	ID       int     `db:"id"`
	Name     string  `db:"name"`
	ColorHex *string `db:"color_hex"`
}

// Card представляет личную карточку пользователя.
type Card struct {
	ID                  uuid.UUID  `db:"id"`
	SenseID             *uuid.UUID `db:"sense_id"`
	CustomText          *string    `db:"custom_text"`
	CustomTranscription *string    `db:"custom_transcription"`
	CustomTranslations  []string   `db:"custom_translations"`
	CustomNote          *string    `db:"custom_note"`
	CustomImageURL      *string    `db:"custom_image_url"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
	IsDeleted           bool       `db:"is_deleted"`
}

// CardTag связывает карточку с тегом.
type CardTag struct {
	CardID uuid.UUID `db:"card_id"`
	TagID  int       `db:"tag_id"`
}

// LearningStatus представляет статус изучения.
type LearningStatus string

const (
	LearningStatusNew      LearningStatus = "new"
	LearningStatusLearning LearningStatus = "learning"
	LearningStatusReview   LearningStatus = "review"
	LearningStatusMastered LearningStatus = "mastered"
)

// SRSState представляет текущее состояние SRS для карточки.
type SRSState struct {
	CardID        uuid.UUID      `db:"card_id"`
	Status        LearningStatus `db:"status"`
	DueDate       *time.Time     `db:"due_date"`
	AlgorithmData map[string]any `db:"algorithm_data"`
	LastReviewAt  *time.Time     `db:"last_review_at"`
}

// ReviewLog представляет запись об одном повторении.
type ReviewLog struct {
	ID          int64          `db:"id"`
	CardID      uuid.UUID      `db:"card_id"`
	Grade       int            `db:"grade"` // 1-5
	DurationMs  *int           `db:"duration_ms"`
	ReviewedAt  time.Time      `db:"reviewed_at"`
	StateBefore map[string]any `db:"state_before"`
	StateAfter  map[string]any `db:"state_after"`
}
