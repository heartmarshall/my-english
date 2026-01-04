package graph

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/card"
	"github.com/heartmarshall/my-english/internal/service/study"
)

// toGraphQLCard конвертирует model.Card в GraphQL Card.
func toGraphQLCard(m *model.Card) *Card {
	if m == nil {
		return nil
	}

	return &Card{
		ID:                  m.ID.String(),
		CustomText:          m.CustomText,
		CustomTranscription: m.CustomTranscription,
		CustomTranslations:  m.CustomTranslations,
		CustomNote:          m.CustomNote,
		CustomImageURL:      m.CustomImageURL,
		CreatedAt:           m.CreatedAt,
		// Sense, Progress, Tags будут загружены через field resolvers
	}
}

// toGraphQLLexeme конвертирует model.Lexeme в GraphQL Lexeme.
func toGraphQLLexeme(m *model.Lexeme) *Lexeme {
	if m == nil {
		return nil
	}

	return &Lexeme{
		ID:   m.ID.String(),
		Text: m.TextDisplay,
		// Pronunciations, Senses, Inflections, Lemma будут загружены через field resolvers
	}
}

// toGraphQLSense конвертирует model.Sense в GraphQL Sense.
func toGraphQLSense(m *model.Sense) *Sense {
	if m == nil {
		return nil
	}

	pos := toGraphQLPartOfSpeech(model.PartOfSpeech(m.PartOfSpeech))

	return &Sense{
		ID:           m.ID.String(),
		PartOfSpeech: pos,
		Definition:   m.Definition,
		CefrLevel:    m.CefrLevel,
		// Lexeme, Source, Translations, Examples, Relations будут загружены через field resolvers
	}
}

// toGraphQLInboxItem конвертирует model.InboxItem в GraphQL InboxItem.
func toGraphQLInboxItem(m *model.InboxItem) *InboxItem {
	if m == nil {
		return nil
	}

	return &InboxItem{
		ID:        m.ID.String(),
		Text:      m.RawText,
		Context:   m.ContextNote,
		CreatedAt: m.CreatedAt,
	}
}

// toGraphQLTag конвертирует model.Tag в GraphQL Tag.
func toGraphQLTag(m *model.Tag) *Tag {
	if m == nil {
		return nil
	}

	return &Tag{
		ID:    strconv.Itoa(m.ID), // Tag.ID в модели - int, в GraphQL - string
		Name:  m.Name,
		Color: m.ColorHex,
	}
}

// toGraphQLPartOfSpeech конвертирует model.PartOfSpeech в GraphQL PartOfSpeech.
func toGraphQLPartOfSpeech(pos model.PartOfSpeech) PartOfSpeech {
	switch pos {
	case model.PartOfSpeechNoun:
		return PartOfSpeechNoun
	case model.PartOfSpeechVerb:
		return PartOfSpeechVerb
	case model.PartOfSpeechAdjective:
		return PartOfSpeechAdjective
	case model.PartOfSpeechAdverb:
		return PartOfSpeechAdverb
	case model.PartOfSpeechPronoun:
		return PartOfSpeechPronoun
	case model.PartOfSpeechPreposition:
		return PartOfSpeechPreposition
	case model.PartOfSpeechConjunction:
		return PartOfSpeechConjunction
	case model.PartOfSpeechInterjection:
		return PartOfSpeechInterjection
	case model.PartOfSpeechPhrase:
		return PartOfSpeechPhrase
	case model.PartOfSpeechIdiom:
		return PartOfSpeechIdiom
	case model.PartOfSpeechOther:
		return PartOfSpeechOther
	default:
		return PartOfSpeechOther
	}
}

// toGraphQLLearningStatus конвертирует model.LearningStatus в GraphQL LearningStatus.
func toGraphQLLearningStatus(status model.LearningStatus) LearningStatus {
	switch status {
	case model.LearningStatusNew:
		return LearningStatusNew
	case model.LearningStatusLearning:
		return LearningStatusLearning
	case model.LearningStatusReview:
		return LearningStatusReview
	case model.LearningStatusMastered:
		return LearningStatusMastered
	default:
		return LearningStatusNew
	}
}

// toModelLearningStatus конвертирует GraphQL LearningStatus в model.LearningStatus.
func toModelLearningStatus(status LearningStatus) model.LearningStatus {
	switch status {
	case LearningStatusNew:
		return model.LearningStatusNew
	case LearningStatusLearning:
		return model.LearningStatusLearning
	case LearningStatusReview:
		return model.LearningStatusReview
	case LearningStatusMastered:
		return model.LearningStatusMastered
	default:
		return model.LearningStatusNew
	}
}

// toModelCreateCardInput конвертирует GraphQL CreateCardInput в card.CreateCardInput.
func toModelCreateCardInput(input CreateCardInput) card.CreateCardInput {
	var senseID *uuid.UUID
	if input.SenseID != nil {
		id, _ := uuid.Parse(*input.SenseID)
		senseID = &id
	}

	return card.CreateCardInput{
		SenseID:             senseID,
		CustomText:          input.CustomText,
		CustomTranscription: input.Transcription,
		CustomTranslations:  input.Translations,
		CustomNote:          input.Note,
		CustomImageURL:      nil, // Не поддерживается в GraphQL схеме
		Tags:                input.Tags,
	}
}

// toModelUpdateCardInput конвертирует GraphQL UpdateCardInput в card.UpdateCardInput.
func toModelUpdateCardInput(input UpdateCardInput) card.UpdateCardInput {
	return card.UpdateCardInput{
		CustomTranscription: input.CustomTranscription,
		CustomTranslations:  input.CustomTranslations,
		CustomNote:          input.CustomNote,
		CustomImageURL:      nil, // Не поддерживается в GraphQL схеме
		Tags:                input.Tags,
	}
}

// parseUUID парсит строку в UUID.
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// calculateNextReviewInDays вычисляет количество дней до следующего повторения.
func calculateNextReviewInDays(dueDate *time.Time) float64 {
	if dueDate == nil {
		return 0
	}
	now := time.Now()
	diff := dueDate.Sub(now)
	return diff.Hours() / 24.0
}

// toGraphQLDataSource конвертирует model.DataSource в GraphQL DataSource.
func toGraphQLDataSource(m *model.DataSource) *DataSource {
	if m == nil {
		return nil
	}

	return &DataSource{
		ID:          strconv.Itoa(m.ID), // DataSource.ID в модели - int, в GraphQL - string
		Slug:        m.Slug,
		DisplayName: m.DisplayName,
		TrustLevel:  m.TrustLevel,
		WebsiteURL:  m.WebsiteURL,
	}
}

// toGraphQLMorphologicalType конвертирует model.MorphologicalType в GraphQL MorphologicalType.
func toGraphQLMorphologicalType(mt model.MorphologicalType) MorphologicalType {
	switch mt {
	case model.MorphTypePlural:
		return MorphologicalTypePlural
	case model.MorphTypePastTense:
		return MorphologicalTypePastTense
	case model.MorphTypePastParticiple:
		return MorphologicalTypePastParticiple
	case model.MorphTypePresentParticiple:
		return MorphologicalTypePresentParticiple
	case model.MorphTypeComparative:
		return MorphologicalTypeComparative
	case model.MorphTypeSuperlative:
		return MorphologicalTypeSuperlative
	default:
		return MorphologicalTypePlural
	}
}

// toGraphQLRelationType конвертирует model.RelationType в GraphQL RelationType.
func toGraphQLRelationType(rt model.RelationType) RelationType {
	switch rt {
	case model.RelationTypeSynonym:
		return RelationTypeSynonym
	case model.RelationTypeAntonym:
		return RelationTypeAntonym
	case model.RelationTypeRelated:
		return RelationTypeRelated
	case model.RelationTypeCollocation:
		return RelationTypeCollocation
	default:
		return RelationTypeRelated
	}
}

// toGraphQLSenseRelation конвертирует model.SenseRelation в GraphQL SenseRelation.
// targetSense должен быть загружен отдельно через DataLoader.
func toGraphQLSenseRelation(m *model.SenseRelation, targetSense *Sense) *SenseRelation {
	if m == nil {
		return nil
	}

	return &SenseRelation{
		Sense:           targetSense,
		Type:            toGraphQLRelationType(model.RelationType(m.Type)),
		IsBidirectional: m.IsBidirectional,
	}
}

// toModelCardFilter конвертирует GraphQL CardsFilter в card.Filter.
func toModelCardFilter(filter *CardsFilter) *card.Filter {
	if filter == nil {
		return nil
	}

	var statuses []model.LearningStatus
	if len(filter.Statuses) > 0 {
		statuses = make([]model.LearningStatus, len(filter.Statuses))
		for i, s := range filter.Statuses {
			statuses[i] = toModelLearningStatus(s)
		}
	}

	return &card.Filter{
		Search:   filter.Search,
		Tags:     filter.Tags,
		Statuses: statuses,
	}
}

// toModelStudyFilter конвертирует GraphQL StudyFilter в study.Filter.
func toModelStudyFilter(filter *StudyFilter) *study.Filter {
	if filter == nil {
		return nil
	}

	var statuses []model.LearningStatus
	if len(filter.Statuses) > 0 {
		statuses = make([]model.LearningStatus, len(filter.Statuses))
		for i, s := range filter.Statuses {
			statuses[i] = toModelLearningStatus(s)
		}
	}

	return &study.Filter{
		Tags:     filter.Tags,
		Statuses: statuses,
	}
}
