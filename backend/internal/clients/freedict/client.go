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
)

const (
	apiURL     = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	sourceName = "free_dictionary"
)

// Client реализует взаимодействие с Free Dictionary API.
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

// FetchWord получает информацию о слове из API.
// Обрати внимание: мы возвращаем конкретный тип, а не интерфейс.
func (c *Client) FetchWord(ctx context.Context, word string) (*model.DictionaryWordData, error) {
	// Экранируем слово
	safeWord := url.PathEscape(strings.TrimSpace(word))
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
		// Возвращаем nil, nil, если слово не найдено (это штатная ситуация)
		return nil, nil
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

	return c.mapToModel(apiResponse[0]), nil
}

// --- Internal Structures & Mapping ---

type apiEntry struct {
	Word      string        `json:"word"`
	Phonetic  string        `json:"phonetic"`
	Phonetics []apiPhonetic `json:"phonetics"`
	Meanings  []apiMeaning  `json:"meanings"`
}

type apiPhonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type apiMeaning struct {
	PartOfSpeech string          `json:"partOfSpeech"`
	Definitions  []apiDefinition `json:"definitions"`
}

type apiDefinition struct {
	Definition string `json:"definition"`
}

func (c *Client) mapToModel(entry apiEntry) *model.DictionaryWordData {
	// Логика выбора лучшей фонетики
	audio := ""
	transcription := entry.Phonetic
	for _, p := range entry.Phonetics {
		if p.Audio != "" && audio == "" {
			audio = p.Audio
		}
		if p.Text != "" && transcription == "" {
			transcription = p.Text
		}
	}

	word := model.DictionaryWord{
		Text:          entry.Word,
		Transcription: &transcription,
		AudioURL:      &audio,
		Source:        sourceName,
	}

	var meaningsData []model.DictionaryMeaningData

	for i, m := range entry.Meanings {
		pos := mapPartOfSpeech(m.PartOfSpeech)

		// Берем определения
		for _, def := range m.Definitions {
			definition := def.Definition

			dm := model.DictionaryMeaning{
				PartOfSpeech: pos,
				DefinitionEn: &definition,
				OrderIndex:   i,
			}

			meaningsData = append(meaningsData, model.DictionaryMeaningData{
				Meaning:      dm,
				Translations: []model.DictionaryTranslation{}, // Переводов нет
			})
		}
	}

	return &model.DictionaryWordData{
		Word:     word,
		Meanings: meaningsData,
	}
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
	default:
		return model.PartOfSpeechOther
	}
}
