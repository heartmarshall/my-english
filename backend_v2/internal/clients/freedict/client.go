package freedict

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/dictionary"
)

const (
	apiURL     = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	SourceSlug = "freedict"
)

// Client реализует интерфейс dictionary.Provider.
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SourceSlug возвращает идентификатор источника.
func (c *Client) SourceSlug() string {
	return SourceSlug
}

// Fetch получает данные о слове и преобразует их в ImportedWord.
func (c *Client) Fetch(ctx context.Context, query string) (*dictionary.ImportedWord, error) {
	safeWord := url.PathEscape(strings.TrimSpace(query))
	reqURL := apiURL + safeWord

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Штатная ситуация: слово не найдено
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResponse []apiEntry
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse) == 0 {
		return nil, nil
	}

	return c.mapToImportedWord(apiResponse), nil
}

// mapToImportedWord преобразует ответ API во внутреннюю структуру.
// FreeDict может вернуть несколько entries для одного слова (например, как сущ. и как глагол),
// мы объединяем их в один ImportedWord.
func (c *Client) mapToImportedWord(entries []apiEntry) *dictionary.ImportedWord {
	result := &dictionary.ImportedWord{
		Text:           entries[0].Word, // Берем написание из первого вхождения
		Pronunciations: make([]dictionary.ImportedPronunciation, 0),
		Senses:         make([]dictionary.ImportedSense, 0),
	}

	// Используем мапы для дедупликации, так как API часто дублирует данные
	seenAudios := make(map[string]bool)

	for _, entry := range entries {
		// 1. Собираем произношения
		for _, p := range entry.Phonetics {
			if p.Audio == "" {
				continue
			}
			if seenAudios[p.Audio] {
				continue
			}

			seenAudios[p.Audio] = true

			// Пытаемся определить регион по URL (API часто пишет -us.mp3 / -uk.mp3)
			region := model.AccentRegionGeneral
			if strings.Contains(p.Audio, "-us.mp3") {
				region = model.AccentRegionUS
			} else if strings.Contains(p.Audio, "-uk.mp3") {
				region = model.AccentRegionUK
			} else if strings.Contains(p.Audio, "-au.mp3") {
				region = model.AccentRegionAU
			}

			// Текст транскрипции берем либо из фонетики, либо из корня entry
			transcription := p.Text
			if transcription == "" {
				transcription = entry.Phonetic
			}

			result.Pronunciations = append(result.Pronunciations, dictionary.ImportedPronunciation{
				AudioURL:      p.Audio,
				Transcription: transcription,
				Region:        region,
			})
		}

		// 2. Собираем смыслы (Senses)
		for _, m := range entry.Meanings {
			pos := mapPartOfSpeech(m.PartOfSpeech)

			for _, def := range m.Definitions {
				sense := dictionary.ImportedSense{
					PartOfSpeech: pos,
					Definition:   def.Definition,
					Translations: []string{}, // FreeDict не дает переводов
					Examples:     make([]dictionary.ImportedExample, 0),
				}

				if def.Example != "" {
					sense.Examples = append(sense.Examples, dictionary.ImportedExample{
						SentenceEn: def.Example,
						// SentenceRu пустой
					})
				}

				result.Senses = append(result.Senses, sense)
			}
		}
	}

	return result
}

func mapPartOfSpeech(pos string) model.PartOfSpeech {
	switch strings.ToLower(pos) {
	case "noun":
		return model.PartOfSpeechNoun
	case "verb":
		return model.PartOfSpeechVerb
	case "adjective":
		return model.PartOfSpeechAdjective
	case "adverb":
		return model.PartOfSpeechAdverb
	case "pronoun":
		return model.PartOfSpeechPronoun
	case "preposition":
		return model.PartOfSpeechPreposition
	case "conjunction":
		return model.PartOfSpeechConjunction
	case "interjection":
		return model.PartOfSpeechInterjection
	default:
		return model.PartOfSpeechOther
	}
}
