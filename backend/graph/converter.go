package graph

import (
	"strconv"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/word"
)

// --- ID Conversion ---

// ToGraphQLID конвертирует int64 в GraphQL ID.
func ToGraphQLID(id int64) string {
	return strconv.FormatInt(id, 10)
}

// FromGraphQLID конвертирует GraphQL ID в int64.
func FromGraphQLID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

// --- Enum Conversion ---

// ToGraphQLPartOfSpeech конвертирует domain enum в GraphQL enum.
func ToGraphQLPartOfSpeech(pos model.PartOfSpeech) PartOfSpeech {
	switch pos {
	case model.PartOfSpeechNoun:
		return PartOfSpeechNoun
	case model.PartOfSpeechVerb:
		return PartOfSpeechVerb
	case model.PartOfSpeechAdjective:
		return PartOfSpeechAdjective
	case model.PartOfSpeechAdverb:
		return PartOfSpeechAdverb
	default:
		return PartOfSpeechOther
	}
}

// FromGraphQLPartOfSpeech конвертирует GraphQL enum в domain enum.
func FromGraphQLPartOfSpeech(pos PartOfSpeech) model.PartOfSpeech {
	switch pos {
	case PartOfSpeechNoun:
		return model.PartOfSpeechNoun
	case PartOfSpeechVerb:
		return model.PartOfSpeechVerb
	case PartOfSpeechAdjective:
		return model.PartOfSpeechAdjective
	case PartOfSpeechAdverb:
		return model.PartOfSpeechAdverb
	default:
		return model.PartOfSpeechOther
	}
}

// ToGraphQLLearningStatus конвертирует domain enum в GraphQL enum.
func ToGraphQLLearningStatus(status model.LearningStatus) LearningStatus {
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

// ToGraphQLExampleSource конвертирует domain enum в GraphQL enum.
func ToGraphQLExampleSource(src *model.ExampleSource) *ExampleSource {
	if src == nil || *src == "" {
		return nil
	}
	var result ExampleSource
	switch *src {
	case model.ExampleSourceFilm:
		result = ExampleSourceFilm
	case model.ExampleSourceBook:
		result = ExampleSourceBook
	case model.ExampleSourceChat:
		result = ExampleSourceChat
	case model.ExampleSourceVideo:
		result = ExampleSourceVideo
	case model.ExampleSourcePodcast:
		result = ExampleSourcePodcast
	default:
		return nil
	}
	return &result
}

// FromGraphQLExampleSource конвертирует GraphQL string в domain enum.
func FromGraphQLExampleSource(src *string) *model.ExampleSource {
	if src == nil {
		return nil
	}
	var result model.ExampleSource
	switch *src {
	case "FILM", "film":
		result = model.ExampleSourceFilm
	case "BOOK", "book":
		result = model.ExampleSourceBook
	case "CHAT", "chat":
		result = model.ExampleSourceChat
	case "VIDEO", "video":
		result = model.ExampleSourceVideo
	case "PODCAST", "podcast":
		result = model.ExampleSourcePodcast
	default:
		return nil
	}
	return &result
}

// --- Model Conversion ---

// ToGraphQLWord конвертирует domain Word в GraphQL Word.
func ToGraphQLWord(w *model.Word, meanings []*Meaning) *Word {
	if w == nil {
		return nil
	}
	return &Word{
		ID:            ToGraphQLID(w.ID),
		Text:          w.Text,
		Transcription: w.Transcription,
		AudioURL:      w.AudioURL,
		FrequencyRank: w.FrequencyRank,
		Meanings:      meanings,
	}
}

// ToGraphQLMeaning конвертирует domain Meaning в GraphQL Meaning.
func ToGraphQLMeaning(m *model.Meaning, examples []*Example, tags []*Tag) *Meaning {
	if m == nil {
		return nil
	}

	reviewCount := 0
	if m.ReviewCount != nil {
		reviewCount = *m.ReviewCount
	}

	var nextReviewAt *Time
	if m.NextReviewAt != nil {
		t := Time(*m.NextReviewAt)
		nextReviewAt = &t
	}

	return &Meaning{
		ID:            ToGraphQLID(m.ID),
		WordID:        ToGraphQLID(m.WordID),
		PartOfSpeech:  ToGraphQLPartOfSpeech(m.PartOfSpeech),
		DefinitionEn:  m.DefinitionEn,
		TranslationRu: m.TranslationRu,
		CefrLevel:     m.CefrLevel,
		ImageURL:      m.ImageURL,
		Status:        ToGraphQLLearningStatus(m.LearningStatus),
		NextReviewAt:  nextReviewAt,
		ReviewCount:   reviewCount,
		Examples:      examples,
		Tags:          tags,
	}
}

// ToGraphQLExample конвертирует domain Example в GraphQL Example.
func ToGraphQLExample(e *model.Example) *Example {
	if e == nil {
		return nil
	}
	return &Example{
		ID:         ToGraphQLID(e.ID),
		SentenceEn: e.SentenceEn,
		SentenceRu: e.SentenceRu,
		SourceName: ToGraphQLExampleSource(e.SourceName),
	}
}

// ToGraphQLExamples конвертирует slice domain Examples в GraphQL Examples.
func ToGraphQLExamples(examples []*model.Example) []*Example {
	result := make([]*Example, 0, len(examples))
	for _, e := range examples {
		result = append(result, ToGraphQLExample(e))
	}
	return result
}

// ToGraphQLTag конвертирует domain Tag в GraphQL Tag.
func ToGraphQLTag(t *model.Tag) *Tag {
	if t == nil {
		return nil
	}
	return &Tag{
		ID:   ToGraphQLID(t.ID),
		Name: t.Name,
	}
}

// ToGraphQLTags конвертирует slice domain Tags в GraphQL Tags.
func ToGraphQLTags(tags []*model.Tag) []*Tag {
	result := make([]*Tag, 0, len(tags))
	for _, t := range tags {
		result = append(result, ToGraphQLTag(t))
	}
	return result
}

// ToGraphQLStats конвертирует domain Stats в GraphQL DashboardStats.
func ToGraphQLStats(s *model.Stats) *DashboardStats {
	if s == nil {
		return nil
	}
	return &DashboardStats{
		TotalWords:        s.TotalWords,
		MasteredCount:     s.MasteredCount,
		LearningCount:     s.LearningCount,
		DueForReviewCount: s.DueForReviewCount,
	}
}

// --- Input Conversion ---

// ToCreateWordInput конвертирует GraphQL input в service input.
func ToCreateWordInput(input AddWordInput) word.CreateWordInput {
	meanings := make([]word.CreateMeaningInput, 0, len(input.Meanings))
	for _, m := range input.Meanings {
		meanings = append(meanings, ToCreateMeaningInput(m))
	}

	return word.CreateWordInput{
		Text:          input.Text,
		Transcription: input.Transcription,
		AudioURL:      input.AudioURL,
		Meanings:      meanings,
	}
}

// ToCreateMeaningInput конвертирует GraphQL MeaningInput в service input.
func ToCreateMeaningInput(input *MeaningInput) word.CreateMeaningInput {
	if input == nil {
		return word.CreateMeaningInput{}
	}

	examples := make([]word.CreateExampleInput, 0, len(input.Examples))
	for _, e := range input.Examples {
		examples = append(examples, ToCreateExampleInput(e))
	}

	return word.CreateMeaningInput{
		PartOfSpeech:  FromGraphQLPartOfSpeech(input.PartOfSpeech),
		DefinitionEn:  input.DefinitionEn,
		TranslationRu: input.TranslationRu,
		ImageURL:      input.ImageURL,
		Examples:      examples,
		Tags:          input.Tags,
	}
}

// ToCreateExampleInput конвертирует GraphQL ExampleInput в service input.
func ToCreateExampleInput(input *ExampleInput) word.CreateExampleInput {
	if input == nil {
		return word.CreateExampleInput{}
	}

	return word.CreateExampleInput{
		SentenceEn: input.SentenceEn,
		SentenceRu: input.SentenceRu,
		SourceName: FromGraphQLExampleSource(input.SourceName),
	}
}

// ToUpdateWordInput конвертирует GraphQL input в service update input.
func ToUpdateWordInput(input AddWordInput) word.UpdateWordInput {
	meanings := make([]word.UpdateMeaningInput, 0, len(input.Meanings))
	for _, m := range input.Meanings {
		meanings = append(meanings, ToUpdateMeaningInput(m))
	}

	return word.UpdateWordInput{
		Text:          input.Text,
		Transcription: input.Transcription,
		AudioURL:      input.AudioURL,
		Meanings:      meanings,
	}
}

// ToUpdateMeaningInput конвертирует GraphQL MeaningInput в service update input.
func ToUpdateMeaningInput(input *MeaningInput) word.UpdateMeaningInput {
	if input == nil {
		return word.UpdateMeaningInput{}
	}

	examples := make([]word.CreateExampleInput, 0, len(input.Examples))
	for _, e := range input.Examples {
		examples = append(examples, ToCreateExampleInput(e))
	}

	return word.UpdateMeaningInput{
		PartOfSpeech:  FromGraphQLPartOfSpeech(input.PartOfSpeech),
		DefinitionEn:  input.DefinitionEn,
		TranslationRu: input.TranslationRu,
		ImageURL:      input.ImageURL,
		Examples:      examples,
		Tags:          input.Tags,
	}
}

// ToWordFilter конвертирует GraphQL WordFilter в service filter.
func ToWordFilter(filter *WordFilter) *word.WordFilter {
	if filter == nil {
		return nil
	}
	return &word.WordFilter{
		Search: filter.Search,
	}
}

// --- WordWithRelations Conversion ---

// WordWithRelationsToGraphQL конвертирует WordWithRelations в GraphQL Word.
func WordWithRelationsToGraphQL(wr *word.WordWithRelations) *Word {
	if wr == nil {
		return nil
	}

	meanings := make([]*Meaning, 0, len(wr.Meanings))
	for _, mr := range wr.Meanings {
		meanings = append(meanings, MeaningWithRelationsToGraphQL(mr))
	}

	return ToGraphQLWord(wr.Word, meanings)
}

// MeaningWithRelationsToGraphQL конвертирует MeaningWithRelations в GraphQL Meaning.
func MeaningWithRelationsToGraphQL(mr *word.MeaningWithRelations) *Meaning {
	if mr == nil {
		return nil
	}

	return ToGraphQLMeaning(
		mr.Meaning,
		ToGraphQLExamples(mr.Examples),
		ToGraphQLTags(mr.Tags),
	)
}
