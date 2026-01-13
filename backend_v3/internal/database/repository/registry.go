package repository

import (
	"github.com/heartmarshall/my-english/internal/database"
)

// Registry объединяет все репозитории.
type Registry struct {
	Dictionary     *DictionaryRepository
	Senses         *SenseRepository
	Translations   *TranslationRepository
	Examples       *ExampleRepository
	Images         *ImageRepository
	Pronunciations *PronunciationRepository
	Cards          *CardRepository
	ReviewLogs     *ReviewLogRepository
	Inbox          *InboxRepository
	Audit          *AuditRepository
}

// NewRegistry создает все репозитории, используя переданный Querier.
func NewRegistry(q database.Querier) *Registry {
	return &Registry{
		Dictionary:     NewDictionaryRepository(q),
		Senses:         NewSenseRepository(q),
		Translations:   NewTranslationRepository(q),
		Examples:       NewExampleRepository(q),
		Images:         NewImageRepository(q),
		Pronunciations: NewPronunciationRepository(q),
		Cards:          NewCardRepository(q),
		ReviewLogs:     NewReviewLogRepository(q),
		Inbox:          NewInboxRepository(q),
		Audit:          NewAuditRepository(q),
	}
}
