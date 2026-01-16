package graph

import (
	"github.com/heartmarshall/my-english/graph/model"
	"github.com/heartmarshall/my-english/internal/service/dictionary"
)

// mapCreateWordInput конвертирует GraphQL input в сервисный input
func mapCreateWordInput(input model.CreateWordInput) dictionary.CreateWordInput {
	return dictionary.CreateWordInput{
		Text:           input.Text,
		Senses:         mapSensesInput(input.Senses),
		Images:         mapImagesInput(input.Images),
		Pronunciations: mapPronunciationsInput(input.Pronunciations),
		CreateCard:     input.CreateCard,
	}
}

// mapUpdateWordInput конвертирует GraphQL input в сервисный input
func mapUpdateWordInput(id string, input model.UpdateWordInput) dictionary.UpdateWordInput {
	return dictionary.UpdateWordInput{
		ID:             id,
		Text:           input.Text,
		Senses:         mapSensesInput(input.Senses),
		Images:         nil,
		Pronunciations: nil,
	}
}

func mapSensesInput(inputs []*model.SenseInput) []dictionary.SenseInput {
	if len(inputs) == 0 {
		return nil
	}
	res := make([]dictionary.SenseInput, len(inputs))
	for i, in := range inputs {
		if in == nil {
			continue
		}
		res[i] = dictionary.SenseInput{
			Definition:   in.Definition, // ИСПРАВЛЕНО: передаем указатель как есть
			PartOfSpeech: in.PartOfSpeech,
			SourceSlug:   getString(in.SourceSlug),
			Translations: mapTranslationsInput(in.Translations),
			Examples:     mapExamplesInput(in.Examples),
		}
	}
	return res
}

func mapTranslationsInput(inputs []*model.TranslationInput) []dictionary.TranslationInput {
	if len(inputs) == 0 {
		return nil
	}
	res := make([]dictionary.TranslationInput, len(inputs))
	for i, in := range inputs {
		if in == nil {
			continue
		}
		res[i] = dictionary.TranslationInput{
			Text:       in.Text,
			SourceSlug: getString(in.SourceSlug),
		}
	}
	return res
}

func mapExamplesInput(inputs []*model.ExampleInput) []dictionary.ExampleInput {
	if len(inputs) == 0 {
		return nil
	}
	res := make([]dictionary.ExampleInput, len(inputs))
	for i, in := range inputs {
		if in == nil {
			continue
		}
		res[i] = dictionary.ExampleInput{
			Sentence:    in.Sentence,
			Translation: in.Translation,
			SourceSlug:  getString(in.SourceSlug),
		}
	}
	return res
}

func mapImagesInput(inputs []*model.ImageInput) []dictionary.ImageInput {
	if len(inputs) == 0 {
		return nil
	}
	res := make([]dictionary.ImageInput, len(inputs))
	for i, in := range inputs {
		if in == nil {
			continue
		}
		res[i] = dictionary.ImageInput{
			URL:        in.URL,
			Caption:    in.Caption,
			SourceSlug: getString(in.SourceSlug),
		}
	}
	return res
}

func mapPronunciationsInput(inputs []*model.PronunciationInput) []dictionary.PronunciationInput {
	if len(inputs) == 0 {
		return nil
	}
	res := make([]dictionary.PronunciationInput, len(inputs))
	for i, in := range inputs {
		if in == nil {
			continue
		}
		res[i] = dictionary.PronunciationInput{
			AudioURL:      in.AudioURL,
			Transcription: in.Transcription,
			Region:        in.Region,
			SourceSlug:    getString(in.SourceSlug),
		}
	}
	return res
}

// mapDictionaryFilter мапит фильтр для поиска
func mapDictionaryFilter(f *model.WordFilter) dictionary.DictionaryFilter {
	if f == nil {
		return dictionary.DictionaryFilter{}
	}

	return dictionary.DictionaryFilter{
		Search:       getString(f.Search),
		PartOfSpeech: f.PartOfSpeech,
		HasCard:      f.HasCard,
		Limit:        getInt(f.Limit, 20),
		Offset:       getInt(f.Offset, 0),
		SortBy:       f.SortBy,
		SortDir:      f.SortDir,
	}
}

// Helpers

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getInt(i *int, def int) int {
	if i == nil {
		return def
	}
	return *i
}
