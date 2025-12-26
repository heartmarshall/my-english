package e2e

import (
	"context"
	"encoding/json"
	"testing"
)

// TestCreateWord тестирует создание слова через GraphQL API.
func TestCreateWord(t *testing.T) {
	ctx := context.Background()

	// Настраиваем тестовую БД
	testDB := SetupTestDB(ctx, t)
	defer testDB.Cleanup(ctx)

	// Настраиваем тестовое приложение
	testApp := SetupTestApp(ctx, t, testDB.Pool)
	defer testApp.Cleanup(ctx)

	// GraphQL запрос для создания слова
	query := `
		mutation CreateWord($input: CreateWordInput!) {
			createWord(input: $input) {
				word {
					id
					text
					transcription
					meanings {
						id
						partOfSpeech
						definitionEn
						translationRu
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"text":          "hello",
			"transcription": "həˈloʊ",
			"meanings": []map[string]interface{}{
				{
					"partOfSpeech":  "NOUN",
					"definitionEn":  "a greeting",
					"translationRu": "привет",
					"examples":      []interface{}{},
					"tags":          []interface{}{},
				},
			},
		},
	}

	req := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	resp, err := testApp.DoGraphQLRequest(ctx, req)
	if err != nil {
		t.Fatalf("Failed to execute GraphQL request: %v", err)
	}

	if len(resp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", resp.Errors)
	}

	// Проверяем ответ
	var result struct {
		CreateWord struct {
			Word struct {
				ID            string  `json:"id"`
				Text          string  `json:"text"`
				Transcription *string `json:"transcription"`
				Meanings      []struct {
					ID            string   `json:"id"`
					PartOfSpeech  string   `json:"partOfSpeech"`
					DefinitionEn  *string  `json:"definitionEn"`
					TranslationRu []string `json:"translationRu"`
				} `json:"meanings"`
			} `json:"word"`
		} `json:"createWord"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.CreateWord.Word.Text != "hello" {
		t.Errorf("Expected text 'hello', got '%s'", result.CreateWord.Word.Text)
	}

	if len(result.CreateWord.Word.Meanings) == 0 {
		t.Error("Expected at least one meaning")
	}

	if result.CreateWord.Word.Meanings[0].PartOfSpeech != "NOUN" {
		t.Errorf("Expected partOfSpeech 'NOUN', got '%s'", result.CreateWord.Word.Meanings[0].PartOfSpeech)
	}

	if len(result.CreateWord.Word.Meanings[0].TranslationRu) == 0 {
		t.Error("Expected at least one translation")
	} else if result.CreateWord.Word.Meanings[0].TranslationRu[0] != "привет" {
		t.Errorf("Expected translation 'привет', got '%s'", result.CreateWord.Word.Meanings[0].TranslationRu[0])
	}
}

// TestGetWords тестирует получение списка слов.
func TestGetWords(t *testing.T) {
	ctx := context.Background()

	// Настраиваем тестовую БД
	testDB := SetupTestDB(ctx, t)
	defer testDB.Cleanup(ctx)

	// Настраиваем тестовое приложение
	testApp := SetupTestApp(ctx, t, testDB.Pool)
	defer testApp.Cleanup(ctx)

	// Сначала создаём слово
	createQuery := `
		mutation CreateWord($input: CreateWordInput!) {
			createWord(input: $input) {
				word {
					id
					text
				}
			}
		}
	`

	createVars := map[string]interface{}{
		"input": map[string]interface{}{
			"text": "test",
			"meanings": []map[string]interface{}{
				{
					"partOfSpeech":  "NOUN",
					"translationRu": "тест",
				},
			},
		},
	}

	_, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     createQuery,
		Variables: createVars,
	})
	if err != nil {
		t.Fatalf("Failed to create word: %v", err)
	}

	// Теперь получаем список слов
	query := `
		query GetWords {
			words(first: 10) {
				edges {
					node {
						id
						text
					}
				}
				totalCount
			}
		}
	`

	req := GraphQLRequest{
		Query: query,
	}

	resp, err := testApp.DoGraphQLRequest(ctx, req)
	if err != nil {
		t.Fatalf("Failed to execute GraphQL request: %v", err)
	}

	if len(resp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", resp.Errors)
	}

	// Проверяем ответ
	var result struct {
		Words struct {
			Edges []struct {
				Node struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"node"`
			} `json:"edges"`
			TotalCount int `json:"totalCount"`
		} `json:"words"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.Words.TotalCount < 1 {
		t.Errorf("Expected at least 1 word, got %d", result.Words.TotalCount)
	}

	found := false
	for _, edge := range result.Words.Edges {
		if edge.Node.Text == "test" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find word 'test' in the list")
	}
}

// TestSuggest тестирует операцию suggest.
func TestSuggest(t *testing.T) {
	ctx := context.Background()

	// Настраиваем тестовую БД
	testDB := SetupTestDB(ctx, t)
	defer testDB.Cleanup(ctx)

	// Настраиваем тестовое приложение
	testApp := SetupTestApp(ctx, t, testDB.Pool)
	defer testApp.Cleanup(ctx)

	// Создаём слово для поиска
	createQuery := `
		mutation CreateWord($input: CreateWordInput!) {
			createWord(input: $input) {
				word {
					id
					text
				}
			}
		}
	`

	createVars := map[string]interface{}{
		"input": map[string]interface{}{
			"text": "hello",
			"meanings": []map[string]interface{}{
				{
					"partOfSpeech":  "NOUN",
					"translationRu": "привет",
				},
			},
		},
	}

	_, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     createQuery,
		Variables: createVars,
	})
	if err != nil {
		t.Fatalf("Failed to create word: %v", err)
	}

	// Тестируем suggest
	query := `
		query Suggest($query: String!) {
			suggest(query: $query) {
				text
				transcription
				translations
				origin
				existingWordId
			}
		}
	`

	req := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"query": "hel",
		},
	}

	resp, err := testApp.DoGraphQLRequest(ctx, req)
	if err != nil {
		t.Fatalf("Failed to execute GraphQL request: %v", err)
	}

	if len(resp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", resp.Errors)
	}

	// Проверяем ответ
	var result struct {
		Suggest []struct {
			Text           string   `json:"text"`
			Transcription  *string  `json:"transcription"`
			Translations   []string `json:"translations"`
			Origin         string   `json:"origin"`
			ExistingWordID *string  `json:"existingWordId"`
		} `json:"suggest"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result.Suggest) == 0 {
		t.Error("Expected at least one suggestion")
	}

	found := false
	for _, suggestion := range result.Suggest {
		if suggestion.Text == "hello" {
			found = true
			if suggestion.Origin != "LOCAL" {
				t.Errorf("Expected origin 'LOCAL', got '%s'", suggestion.Origin)
			}
			if suggestion.ExistingWordID == nil {
				t.Error("Expected existingWordId to be set for LOCAL suggestion")
			}
			break
		}
	}

	if !found {
		t.Error("Expected to find 'hello' in suggestions")
	}
}

// TestInboxToWordUpdate тестирует полный цикл: слово попадает в inbox,
// затем создаётся с неполной информацией, и потом обновляется с полной информацией.
func TestInboxToWordUpdate(t *testing.T) {
	ctx := context.Background()

	// Настраиваем тестовую БД
	testDB := SetupTestDB(ctx, t)
	defer testDB.Cleanup(ctx)

	// Настраиваем тестовое приложение
	testApp := SetupTestApp(ctx, t, testDB.Pool)
	defer testApp.Cleanup(ctx)

	// Шаг 1: Добавляем слово в inbox
	addToInboxQuery := `
		mutation AddToInbox($text: String!, $sourceContext: String) {
			addToInbox(text: $text, sourceContext: $sourceContext) {
				id
				text
				sourceContext
			}
		}
	`

	addToInboxVars := map[string]interface{}{
		"text":          "beautiful",
		"sourceContext": "Harry Potter, page 42",
	}

	addToInboxResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     addToInboxQuery,
		Variables: addToInboxVars,
	})
	if err != nil {
		t.Fatalf("Failed to add to inbox: %v", err)
	}

	if len(addToInboxResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", addToInboxResp.Errors)
	}

	var addToInboxResult struct {
		AddToInbox struct {
			ID            string  `json:"id"`
			Text          string  `json:"text"`
			SourceContext *string `json:"sourceContext"`
		} `json:"addToInbox"`
	}

	if err := json.Unmarshal(addToInboxResp.Data, &addToInboxResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if addToInboxResult.AddToInbox.Text != "beautiful" {
		t.Errorf("Expected text 'beautiful', got '%s'", addToInboxResult.AddToInbox.Text)
	}

	inboxID := addToInboxResult.AddToInbox.ID

	// Шаг 2: Превращаем inbox item в слово с неполной информацией
	// Создаём слово только с текстом и одним переводом, без транскрипции, определений и примеров
	convertQuery := `
		mutation ConvertInboxItem($inboxId: ID!, $input: CreateWordInput!) {
			convertInboxItem(inboxId: $inboxId, input: $input) {
				word {
					id
					text
					transcription
					meanings {
						id
						partOfSpeech
						definitionEn
						translationRu
						examples {
							id
							sentenceEn
						}
					}
				}
			}
		}
	`

	convertVars := map[string]interface{}{
		"inboxId": inboxID,
		"input": map[string]interface{}{
			"text":          "beautiful",
			"transcription": nil, // Неполная информация - нет транскрипции
			"meanings": []map[string]interface{}{
				{
					"partOfSpeech":  "ADJECTIVE",
					"translationRu": "красивый",
					// Нет definitionEn, examples, tags - неполная информация
				},
			},
		},
	}

	convertResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     convertQuery,
		Variables: convertVars,
	})
	if err != nil {
		t.Fatalf("Failed to convert inbox item: %v", err)
	}

	if len(convertResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", convertResp.Errors)
	}

	var convertResult struct {
		ConvertInboxItem struct {
			Word struct {
				ID            string  `json:"id"`
				Text          string  `json:"text"`
				Transcription *string `json:"transcription"`
				Meanings      []struct {
					ID            string   `json:"id"`
					PartOfSpeech  string   `json:"partOfSpeech"`
					DefinitionEn  *string  `json:"definitionEn"`
					TranslationRu []string `json:"translationRu"`
					Examples      []struct {
						ID         string `json:"id"`
						SentenceEn string `json:"sentenceEn"`
					} `json:"examples"`
				} `json:"meanings"`
			} `json:"word"`
		} `json:"convertInboxItem"`
	}

	if err := json.Unmarshal(convertResp.Data, &convertResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if convertResult.ConvertInboxItem.Word.Text != "beautiful" {
		t.Errorf("Expected text 'beautiful', got '%s'", convertResult.ConvertInboxItem.Word.Text)
	}

	// Проверяем, что информация действительно неполная
	if convertResult.ConvertInboxItem.Word.Transcription != nil {
		t.Error("Expected transcription to be nil (incomplete information)")
	}

	if len(convertResult.ConvertInboxItem.Word.Meanings) == 0 {
		t.Fatal("Expected at least one meaning")
	}

	meaning := convertResult.ConvertInboxItem.Word.Meanings[0]
	if meaning.DefinitionEn != nil {
		t.Error("Expected definitionEn to be nil (incomplete information)")
	}

	if len(meaning.Examples) > 0 {
		t.Error("Expected no examples (incomplete information)")
	}

	wordID := convertResult.ConvertInboxItem.Word.ID

	// Шаг 3: Обновляем слово с полной информацией
	updateQuery := `
		mutation UpdateWord($id: ID!, $input: CreateWordInput!) {
			updateWord(id: $id, input: $input) {
				word {
					id
					text
					transcription
					meanings {
						id
						partOfSpeech
						definitionEn
						translationRu
						examples {
							id
							sentenceEn
							sentenceRu
							sourceName
						}
					}
				}
			}
		}
	`

	transcription := "ˈbjuːtɪfəl"
	definitionEn := "pleasing the senses or mind aesthetically"
	updateVars := map[string]interface{}{
		"id": wordID,
		"input": map[string]interface{}{
			"text":          "beautiful",
			"transcription": transcription, // Добавляем транскрипцию
			"meanings": []map[string]interface{}{
				{
					"partOfSpeech":  "ADJECTIVE",
					"definitionEn":  definitionEn, // Добавляем определение
					"translationRu": "красивый",
					"examples": []map[string]interface{}{ // Добавляем примеры
						{
							"sentenceEn": "She has a beautiful smile",
							"sentenceRu": "У неё красивая улыбка",
							"sourceName": "BOOK",
						},
						{
							"sentenceEn": "What a beautiful day!",
							"sentenceRu": "Какой прекрасный день!",
							"sourceName": nil,
						},
					},
				},
			},
		},
	}

	updateResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     updateQuery,
		Variables: updateVars,
	})
	if err != nil {
		t.Fatalf("Failed to update word: %v", err)
	}

	if len(updateResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", updateResp.Errors)
	}

	var updateResult struct {
		UpdateWord struct {
			Word struct {
				ID            string  `json:"id"`
				Text          string  `json:"text"`
				Transcription *string `json:"transcription"`
				Meanings      []struct {
					ID            string   `json:"id"`
					PartOfSpeech  string   `json:"partOfSpeech"`
					DefinitionEn  *string  `json:"definitionEn"`
					TranslationRu []string `json:"translationRu"`
					Examples      []struct {
						ID         string  `json:"id"`
						SentenceEn string  `json:"sentenceEn"`
						SentenceRu *string `json:"sentenceRu"`
						SourceName *string `json:"sourceName"`
					} `json:"examples"`
				} `json:"meanings"`
			} `json:"word"`
		} `json:"updateWord"`
	}

	if err := json.Unmarshal(updateResp.Data, &updateResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Проверяем, что слово обновлено с полной информацией
	if updateResult.UpdateWord.Word.ID != wordID {
		t.Errorf("Expected word ID to remain the same, got '%s'", updateResult.UpdateWord.Word.ID)
	}

	if updateResult.UpdateWord.Word.Text != "beautiful" {
		t.Errorf("Expected text 'beautiful', got '%s'", updateResult.UpdateWord.Word.Text)
	}

	// Проверяем транскрипцию
	if updateResult.UpdateWord.Word.Transcription == nil {
		t.Error("Expected transcription to be set after update")
	} else if *updateResult.UpdateWord.Word.Transcription != transcription {
		t.Errorf("Expected transcription '%s', got '%s'", transcription, *updateResult.UpdateWord.Word.Transcription)
	}

	// Проверяем meanings
	if len(updateResult.UpdateWord.Word.Meanings) == 0 {
		t.Fatal("Expected at least one meaning after update")
	}

	updatedMeaning := updateResult.UpdateWord.Word.Meanings[0]

	// Проверяем определение
	if updatedMeaning.DefinitionEn == nil {
		t.Error("Expected definitionEn to be set after update")
	} else if *updatedMeaning.DefinitionEn != definitionEn {
		t.Errorf("Expected definitionEn '%s', got '%s'", definitionEn, *updatedMeaning.DefinitionEn)
	}

	// Проверяем перевод
	if len(updatedMeaning.TranslationRu) == 0 {
		t.Error("Expected at least one translation")
	} else if updatedMeaning.TranslationRu[0] != "красивый" {
		t.Errorf("Expected translation 'красивый', got '%s'", updatedMeaning.TranslationRu[0])
	}

	// Проверяем примеры
	if len(updatedMeaning.Examples) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(updatedMeaning.Examples))
	}

	// Проверяем первый пример
	if updatedMeaning.Examples[0].SentenceEn != "She has a beautiful smile" {
		t.Errorf("Expected first example 'She has a beautiful smile', got '%s'", updatedMeaning.Examples[0].SentenceEn)
	}

	if updatedMeaning.Examples[0].SentenceRu == nil || *updatedMeaning.Examples[0].SentenceRu != "У неё красивая улыбка" {
		t.Error("Expected first example to have Russian translation")
	}

	if updatedMeaning.Examples[0].SourceName == nil || *updatedMeaning.Examples[0].SourceName != "BOOK" {
		t.Error("Expected first example to have sourceName 'BOOK'")
	}

	// Проверяем второй пример
	if updatedMeaning.Examples[1].SentenceEn != "What a beautiful day!" {
		t.Errorf("Expected second example 'What a beautiful day!', got '%s'", updatedMeaning.Examples[1].SentenceEn)
	}

	// Проверяем, что inbox item был удалён после конвертации
	inboxItemsQuery := `
		query GetInboxItems {
			inboxItems {
				id
				text
			}
		}
	`

	inboxItemsResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query: inboxItemsQuery,
	})
	if err != nil {
		t.Fatalf("Failed to get inbox items: %v", err)
	}

	if len(inboxItemsResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", inboxItemsResp.Errors)
	}

	var inboxItemsResult struct {
		InboxItems []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"inboxItems"`
	}

	if err := json.Unmarshal(inboxItemsResp.Data, &inboxItemsResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Проверяем, что inbox item больше нет в списке
	for _, item := range inboxItemsResult.InboxItems {
		if item.ID == inboxID {
			t.Error("Expected inbox item to be deleted after conversion")
		}
	}
}
