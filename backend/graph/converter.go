package graph

import (
	"encoding/base64"
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

// FromGraphQLLearningStatus конвертирует GraphQL enum в domain enum.
func FromGraphQLLearningStatus(status LearningStatus) model.LearningStatus {
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
	createdAt := Time(w.CreatedAt)
	return &Word{
		ID:            ToGraphQLID(w.ID),
		Text:          w.Text,
		Transcription: w.Transcription,
		AudioURL:      w.AudioURL,
		FrequencyRank: w.FrequencyRank,
		CreatedAt:     &createdAt,
		Meanings:      meanings,
		// Forms загрузятся через field resolver
	}
}

// ToGraphQLMeaning конвертирует domain Meaning в GraphQL Meaning.
func ToGraphQLMeaning(m *model.Meaning, translations []model.Translation, examples []*Example, tags []*Tag) *Meaning {
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

	// Преобразуем translations в массив строк
	translationRuArray := make([]string, 0, len(translations))
	for _, tr := range translations {
		if tr.TranslationRu != "" {
			translationRuArray = append(translationRuArray, tr.TranslationRu)
		}
	}

	// Fallback на старое поле для обратной совместимости
	if len(translationRuArray) == 0 && m.TranslationRu != "" {
		translationRuArray = []string{m.TranslationRu}
	}

	return &Meaning{
		ID:            ToGraphQLID(m.ID),
		WordID:        ToGraphQLID(m.WordID),
		PartOfSpeech:  ToGraphQLPartOfSpeech(m.PartOfSpeech),
		DefinitionEn:  m.DefinitionEn,
		TranslationRu: translationRuArray,
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
func ToCreateWordInput(input CreateWordInput) word.CreateWordInput {
	meanings := make([]word.CreateMeaningInput, 0, len(input.Meanings))
	for _, m := range input.Meanings {
		meanings = append(meanings, ToCreateMeaningInput(m))
	}

	return word.CreateWordInput{
		Text:          input.Text,
		Transcription: input.Transcription,
		AudioURL:      input.AudioURL,
		SourceContext: input.SourceContext,
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

	// Конвертируем TranslationRu (string) в Translations ([]string)
	translations := []string{}
	if input.TranslationRu != "" {
		translations = []string{input.TranslationRu}
	}

	return word.CreateMeaningInput{
		PartOfSpeech: FromGraphQLPartOfSpeech(input.PartOfSpeech),
		DefinitionEn: input.DefinitionEn,
		Translations: translations,
		ImageURL:     input.ImageURL,
		Examples:     examples,
		Tags:         input.Tags,
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
func ToUpdateWordInput(input CreateWordInput) word.UpdateWordInput {
	meanings := make([]word.UpdateMeaningInput, 0, len(input.Meanings))
	for _, m := range input.Meanings {
		meanings = append(meanings, ToUpdateMeaningInput(m))
	}

	return word.UpdateWordInput{
		Text:          input.Text,
		Transcription: input.Transcription,
		AudioURL:      input.AudioURL,
		SourceContext: input.SourceContext,
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

	// Конвертируем TranslationRu (string) в Translations ([]string)
	translations := []string{}
	if input.TranslationRu != "" {
		translations = []string{input.TranslationRu}
	}

	return word.UpdateMeaningInput{
		PartOfSpeech: FromGraphQLPartOfSpeech(input.PartOfSpeech),
		DefinitionEn: input.DefinitionEn,
		Translations: translations,
		ImageURL:     input.ImageURL,
		Examples:     examples,
		Tags:         input.Tags,
	}
}

// ToWordFilter конвертирует GraphQL WordFilter в service filter.
func ToWordFilter(filter *WordFilter) *word.WordFilter {
	if filter == nil {
		return nil
	}

	var status *model.LearningStatus
	if filter.Status != nil {
		statusVal := FromGraphQLLearningStatus(*filter.Status)
		status = &statusVal
	}

	return &word.WordFilter{
		Search: filter.Search,
		Status: status,
		Tags:   filter.Tags,
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
		meanings = append(meanings, MeaningWithRelationsToGraphQL(&mr))
	}

	return ToGraphQLWord(&wr.Word, meanings)
}

// MeaningWithRelationsToGraphQL конвертирует MeaningWithRelations в GraphQL Meaning.
func MeaningWithRelationsToGraphQL(mr *word.MeaningWithRelations) *Meaning {
	if mr == nil {
		return nil
	}

	examples := make([]*model.Example, len(mr.Examples))
	for i := range mr.Examples {
		examples[i] = &mr.Examples[i]
	}
	tags := make([]*model.Tag, len(mr.Tags))
	for i := range mr.Tags {
		tags[i] = &mr.Tags[i]
	}

	return ToGraphQLMeaning(
		&mr.Meaning,
		mr.Translations,
		ToGraphQLExamples(examples),
		ToGraphQLTags(tags),
	)
}

// ToGraphQLMeaningBasic конвертирует domain Meaning в GraphQL Meaning без examples и tags.
// Используется с field resolvers, которые загружают relations через DataLoader.
// translations могут быть переданы, если уже загружены, иначе будет использован fallback на TranslationRu
func ToGraphQLMeaningBasic(m *model.Meaning, translations ...[]model.Translation) *Meaning {
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

	// Преобразуем translations в массив строк
	translationRuArray := []string{}
	if len(translations) > 0 && len(translations[0]) > 0 {
		for _, tr := range translations[0] {
			if tr.TranslationRu != "" {
				translationRuArray = append(translationRuArray, tr.TranslationRu)
			}
		}
	}

	// Fallback на старое поле для обратной совместимости
	if len(translationRuArray) == 0 && m.TranslationRu != "" {
		translationRuArray = []string{m.TranslationRu}
	}

	return &Meaning{
		ID:            ToGraphQLID(m.ID),
		WordID:        ToGraphQLID(m.WordID),
		PartOfSpeech:  ToGraphQLPartOfSpeech(m.PartOfSpeech),
		DefinitionEn:  m.DefinitionEn,
		TranslationRu: translationRuArray,
		CefrLevel:     m.CefrLevel,
		ImageURL:      m.ImageURL,
		Status:        ToGraphQLLearningStatus(m.LearningStatus),
		NextReviewAt:  nextReviewAt,
		ReviewCount:   reviewCount,
		// Examples и Tags загрузятся через field resolvers
		// TranslationRu также может загружаться через field resolver
	}
}

// ToGraphQLWordBasic конвертирует domain Word в GraphQL Word без meanings.
// Используется с field resolvers, которые загружают meanings через DataLoader.
func ToGraphQLWordBasic(w *model.Word) *Word {
	if w == nil {
		return nil
	}
	createdAt := Time(w.CreatedAt)
	return &Word{
		ID:            ToGraphQLID(w.ID),
		Text:          w.Text,
		Transcription: w.Transcription,
		AudioURL:      w.AudioURL,
		FrequencyRank: w.FrequencyRank,
		CreatedAt:     &createdAt,
		// Meanings и Forms загрузятся через field resolvers
	}
}

// --- Cursor Pagination ---

// EncodeCursor кодирует offset в base64 cursor.
func EncodeCursor(offset int) string {
	data := strconv.FormatInt(int64(offset), 10)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// DecodeCursor декодирует base64 cursor в offset.
func DecodeCursor(cursor string) (int, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}
	offset, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(offset), nil
}

// --- InboxItem Conversion ---

// ToGraphQLInboxItem конвертирует domain InboxItem в GraphQL InboxItem.
func ToGraphQLInboxItem(item *model.InboxItem) *InboxItem {
	if item == nil {
		return nil
	}
	t := Time(item.CreatedAt)
	return &InboxItem{
		ID:            ToGraphQLID(item.ID),
		Text:          item.Text,
		SourceContext: item.SourceContext,
		CreatedAt:     t,
	}
}

// ToGraphQLInboxItems конвертирует slice domain InboxItems в GraphQL InboxItems.
func ToGraphQLInboxItems(items []model.InboxItem) []*InboxItem {
	result := make([]*InboxItem, 0, len(items))
	for i := range items {
		result = append(result, ToGraphQLInboxItem(&items[i]))
	}
	return result
}

// --- Suggestion Conversion ---

// ToGraphQLSuggestion конвертирует service Suggestion в GraphQL Suggestion.
func ToGraphQLSuggestion(s *word.Suggestion) *Suggestion {
	if s == nil {
		return nil
	}

	var existingWordID *string
	if s.ExistingWordID != nil {
		id := ToGraphQLID(*s.ExistingWordID)
		existingWordID = &id
	}

	var origin SuggestionOrigin
	switch s.Origin {
	case "LOCAL":
		origin = SuggestionOriginLocal
	case "DICTIONARY":
		origin = SuggestionOriginDictionary
	default:
		origin = SuggestionOriginLocal
	}

	return &Suggestion{
		Text:           s.Text,
		Transcription:  s.Transcription,
		Translations:   s.Translations,
		Definition:     s.Definition, // <-- Маппим новое поле
		Origin:         origin,
		ExistingWordID: existingWordID,
	}
}

// ToGraphQLSuggestions конвертирует slice service Suggestions в GraphQL Suggestions.
func ToGraphQLSuggestions(suggestions []word.Suggestion) []*Suggestion {
	result := make([]*Suggestion, 0, len(suggestions))
	for i := range suggestions {
		result = append(result, ToGraphQLSuggestion(&suggestions[i]))
	}
	return result
}

// --- WordForm Conversion ---

// ToGraphQLWordForm конвертирует domain DictionaryWordForm в GraphQL WordForm.
func ToGraphQLWordForm(f *model.DictionaryWordForm) *WordForm {
	if f == nil {
		return nil
	}
	return &WordForm{
		ID:       ToGraphQLID(f.ID),
		FormText: f.FormText,
		FormType: f.FormType,
	}
}

// ToGraphQLWordForms конвертирует slice DictionaryWordForm в GraphQL WordForms.
func ToGraphQLWordForms(forms []model.DictionaryWordForm) []*WordForm {
	result := make([]*WordForm, 0, len(forms))
	for i := range forms {
		result = append(result, ToGraphQLWordForm(&forms[i]))
	}
	return result
}

// --- SynonymAntonym Conversion ---

// ToGraphQLRelationType конвертирует domain RelationType в GraphQL RelationType.
func ToGraphQLRelationType(rt model.RelationType) RelationType {
	switch rt {
	case model.RelationTypeSynonym:
		return RelationTypeSynonym
	case model.RelationTypeAntonym:
		return RelationTypeAntonym
	default:
		return RelationTypeSynonym
	}
}

// ToGraphQLSynonymAntonym конвертирует domain DictionarySynonymAntonym в GraphQL SynonymAntonym.
// meaningID - ID текущего значения, чтобы определить relatedMeaningId
func ToGraphQLSynonymAntonym(sa *model.DictionarySynonymAntonym, currentMeaningID int64) *SynonymAntonym {
	if sa == nil {
		return nil
	}

	// Определяем relatedMeaningId - это другой meaning (не текущий)
	var relatedMeaningID int64
	if sa.MeaningID1 == currentMeaningID {
		relatedMeaningID = sa.MeaningID2
	} else {
		relatedMeaningID = sa.MeaningID1
	}

	return &SynonymAntonym{
		ID:               ToGraphQLID(sa.ID),
		RelatedMeaningID: ToGraphQLID(relatedMeaningID),
		RelationType:     ToGraphQLRelationType(sa.RelationType),
	}
}

// ToGraphQLSynonymAntonyms конвертирует slice DictionarySynonymAntonym в GraphQL SynonymAntonyms.
func ToGraphQLSynonymAntonyms(relations []model.DictionarySynonymAntonym, currentMeaningID int64) []*SynonymAntonym {
	result := make([]*SynonymAntonym, 0, len(relations))
	for i := range relations {
		result = append(result, ToGraphQLSynonymAntonym(&relations[i], currentMeaningID))
	}
	return result
}
