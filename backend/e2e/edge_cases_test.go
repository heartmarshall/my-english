package e2e

import (
	"context"
	"encoding/json"
	"testing"
)

// TestEdgeCases тестирует различные edge-cases и граничные случаи системы.
// Этот тест покрывает валидацию, обработку ошибок, несуществующие сущности,
// дубликаты и другие граничные случаи.
func TestEdgeCases(t *testing.T) {
	ctx := context.Background()

	// Настраиваем тестовую БД
	testDB := SetupTestDB(ctx, t)
	defer testDB.Cleanup(ctx)

	// Настраиваем тестовое приложение
	testApp := SetupTestApp(ctx, t, testDB.Pool)
	defer testApp.Cleanup(ctx)

	// ============================================
	// 1. ВАЛИДАЦИЯ ВХОДНЫХ ДАННЫХ
	// ============================================

	t.Run("create word with empty text", func(t *testing.T) {
		query := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
					}
				}
			}
		`

		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"text":     "",
				"meanings": []interface{}{},
			},
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for empty text")
		}

		// Проверяем код ошибки
		if len(resp.Errors) > 0 {
			if resp.Errors[0].Extensions != nil {
				if code, ok := resp.Errors[0].Extensions["code"].(string); ok {
					if code != "INVALID_INPUT" {
						t.Errorf("Expected error code INVALID_INPUT, got %s", code)
					}
				}
			}
		}
	})

	t.Run("create word with whitespace only text", func(t *testing.T) {
		query := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
					}
				}
			}
		`

		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"text":     "   \t\n  ",
				"meanings": []interface{}{},
			},
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for whitespace-only text")
		}
	})

	t.Run("create word with empty meanings list", func(t *testing.T) {
		query := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
						text
						meanings {
							id
						}
					}
				}
			}
		`

		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"text":     "testword",
				"meanings": []interface{}{},
			},
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		// Слово может быть создано без meanings, но это edge-case
		if len(resp.Errors) > 0 {
			// Если это ошибка - это нормально для edge-case
			t.Logf("Word creation with empty meanings returned error: %+v", resp.Errors)
		} else {
			var result struct {
				CreateWord struct {
					Word struct {
						ID       string        `json:"id"`
						Text     string        `json:"text"`
						Meanings []interface{} `json:"meanings"`
					} `json:"word"`
				} `json:"createWord"`
			}

			if err := json.Unmarshal(resp.Data, &result); err == nil {
				if result.CreateWord.Word.Text != "testword" {
					t.Errorf("Expected text 'testword', got '%s'", result.CreateWord.Word.Text)
				}
				if len(result.CreateWord.Word.Meanings) != 0 {
					t.Errorf("Expected empty meanings, got %d", len(result.CreateWord.Word.Meanings))
				}
			}
		}
	})

	t.Run("create word with meaning without translation", func(t *testing.T) {
		query := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
					}
				}
			}
		`

		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "testword2",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "", // Пустой перевод
					},
				},
			},
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		// Пустой перевод может быть невалидным
		if len(resp.Errors) > 0 {
			t.Logf("Expected error for empty translation: %+v", resp.Errors)
		}
	})

	// ============================================
	// 2. НЕСУЩЕСТВУЮЩИЕ СУЩНОСТИ
	// ============================================

	t.Run("update non-existent word", func(t *testing.T) {
		query := `
			mutation UpdateWord($id: ID!, $input: CreateWordInput!) {
				updateWord(id: $id, input: $input) {
					word {
						id
					}
				}
			}
		`

		// Используем невалидный ID (например, закодированный ID для несуществующего слова)
		variables := map[string]interface{}{
			"id": "999999", // Несуществующий ID
			"input": map[string]interface{}{
				"text": "updated",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "обновлено",
					},
				},
			},
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for non-existent word")
		}

		if len(resp.Errors) > 0 {
			if resp.Errors[0].Extensions != nil {
				if code, ok := resp.Errors[0].Extensions["code"].(string); ok {
					if code != "NOT_FOUND" {
						t.Errorf("Expected error code NOT_FOUND, got %s", code)
					}
				}
			}
		}
	})

	t.Run("delete non-existent word", func(t *testing.T) {
		query := `
			mutation DeleteWord($id: ID!) {
				deleteWord(id: $id)
			}
		`

		variables := map[string]interface{}{
			"id": "999999", // Несуществующий ID
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		// Удаление несуществующего слова может вернуть false или ошибку
		if len(resp.Errors) > 0 {
			// Ошибка - это нормально
			t.Logf("Delete non-existent word returned error: %+v", resp.Errors)
		} else {
			var result struct {
				DeleteWord bool `json:"deleteWord"`
			}

			if err := json.Unmarshal(resp.Data, &result); err == nil {
				if result.DeleteWord {
					t.Error("Expected deleteWord to return false for non-existent word")
				}
			}
		}
	})

	t.Run("review non-existent meaning", func(t *testing.T) {
		query := `
			mutation ReviewMeaning($meaningId: ID!, $grade: Int!) {
				reviewMeaning(meaningId: $meaningId, grade: $grade) {
					id
				}
			}
		`

		variables := map[string]interface{}{
			"meaningId": "999999", // Несуществующий ID
			"grade":     5,
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for non-existent meaning")
		}

		if len(resp.Errors) > 0 {
			if resp.Errors[0].Extensions != nil {
				if code, ok := resp.Errors[0].Extensions["code"].(string); ok {
					if code != "NOT_FOUND" {
						t.Errorf("Expected error code NOT_FOUND, got %s", code)
					}
				}
			}
		}
	})

	t.Run("convert non-existent inbox item", func(t *testing.T) {
		query := `
			mutation ConvertInboxItem($inboxId: ID!, $input: CreateWordInput!) {
				convertInboxItem(inboxId: $inboxId, input: $input) {
					word {
						id
					}
				}
			}
		`

		variables := map[string]interface{}{
			"inboxId": "999999", // Несуществующий ID
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

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for non-existent inbox item")
		}
	})

	t.Run("get non-existent word", func(t *testing.T) {
		query := `
			query GetWord($id: ID!) {
				word(id: $id) {
					id
					text
				}
			}
		`

		variables := map[string]interface{}{
			"id": "999999", // Несуществующий ID
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		// Запрос несуществующего слова должен вернуть null, а не ошибку
		var result struct {
			Word *struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"word"`
		}

		if err := json.Unmarshal(resp.Data, &result); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if result.Word != nil {
			t.Error("Expected word to be null for non-existent ID")
		}
	})

	// ============================================
	// 3. ДУБЛИКАТЫ
	// ============================================

	t.Run("create duplicate word", func(t *testing.T) {
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
				"text": "duplicate",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "дубликат",
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

		// Пытаемся создать то же слово снова
		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for duplicate word")
		}

		if len(resp.Errors) > 0 {
			if resp.Errors[0].Extensions != nil {
				if code, ok := resp.Errors[0].Extensions["code"].(string); ok {
					if code != "ALREADY_EXISTS" {
						t.Errorf("Expected error code ALREADY_EXISTS, got %s", code)
					}
				}
			}
		}
	})

	t.Run("update word to duplicate text", func(t *testing.T) {
		// Создаём первое слово
		createQuery1 := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
						text
					}
				}
			}
		`

		createVars1 := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "firstword",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "первое",
					},
				},
			},
		}

		createResp1, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery1,
			Variables: createVars1,
		})
		if err != nil {
			t.Fatalf("Failed to create first word: %v", err)
		}

		var createResult1 struct {
			CreateWord struct {
				Word struct {
					ID string `json:"id"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp1.Data, &createResult1); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		// Создаём второе слово
		createVars2 := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "secondword",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "второе",
					},
				},
			},
		}

		createResp2, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery1,
			Variables: createVars2,
		})
		if err != nil {
			t.Fatalf("Failed to create second word: %v", err)
		}

		var createResult2 struct {
			CreateWord struct {
				Word struct {
					ID string `json:"id"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp2.Data, &createResult2); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		// Пытаемся обновить второе слово на текст первого
		updateQuery := `
			mutation UpdateWord($id: ID!, $input: CreateWordInput!) {
				updateWord(id: $id, input: $input) {
					word {
						id
						text
					}
				}
			}
		`

		updateVars := map[string]interface{}{
			"id": createResult2.CreateWord.Word.ID,
			"input": map[string]interface{}{
				"text": "firstword", // Дубликат
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "второе",
					},
				},
			},
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     updateQuery,
			Variables: updateVars,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for duplicate word text")
		}

		if len(resp.Errors) > 0 {
			if resp.Errors[0].Extensions != nil {
				if code, ok := resp.Errors[0].Extensions["code"].(string); ok {
					if code != "ALREADY_EXISTS" {
						t.Errorf("Expected error code ALREADY_EXISTS, got %s", code)
					}
				}
			}
		}
	})

	// ============================================
	// 4. ГРАНИЧНЫЕ ЗНАЧЕНИЯ ДЛЯ REVIEW
	// ============================================

	t.Run("review with invalid grade too low", func(t *testing.T) {
		// Создаём слово для review
		createQuery := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
						meanings {
							id
						}
					}
				}
			}
		`

		createVars := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "reviewtest",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "тест",
					},
				},
			},
		}

		createResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars,
		})
		if err != nil {
			t.Fatalf("Failed to create word: %v", err)
		}

		var createResult struct {
			CreateWord struct {
				Word struct {
					Meanings []struct {
						ID string `json:"id"`
					} `json:"meanings"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp.Data, &createResult); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(createResult.CreateWord.Word.Meanings) == 0 {
			t.Fatal("Expected at least one meaning")
		}

		meaningID := createResult.CreateWord.Word.Meanings[0].ID

		// Пытаемся сделать review с невалидной оценкой (0)
		reviewQuery := `
			mutation ReviewMeaning($meaningId: ID!, $grade: Int!) {
				reviewMeaning(meaningId: $meaningId, grade: $grade) {
					id
				}
			}
		`

		reviewVars := map[string]interface{}{
			"meaningId": meaningID,
			"grade":     0, // Невалидная оценка
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     reviewQuery,
			Variables: reviewVars,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for invalid grade (0)")
		}

		if len(resp.Errors) > 0 {
			if resp.Errors[0].Extensions != nil {
				if code, ok := resp.Errors[0].Extensions["code"].(string); ok {
					if code != "INVALID_INPUT" {
						t.Errorf("Expected error code INVALID_INPUT, got %s", code)
					}
				}
			}
		}
	})

	t.Run("review with invalid grade too high", func(t *testing.T) {
		// Создаём слово для review
		createQuery := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
						meanings {
							id
						}
					}
				}
			}
		`

		createVars := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "reviewtest2",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "тест",
					},
				},
			},
		}

		createResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars,
		})
		if err != nil {
			t.Fatalf("Failed to create word: %v", err)
		}

		var createResult struct {
			CreateWord struct {
				Word struct {
					Meanings []struct {
						ID string `json:"id"`
					} `json:"meanings"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp.Data, &createResult); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(createResult.CreateWord.Word.Meanings) == 0 {
			t.Fatal("Expected at least one meaning")
		}

		meaningID := createResult.CreateWord.Word.Meanings[0].ID

		// Пытаемся сделать review с невалидной оценкой (6)
		reviewQuery := `
			mutation ReviewMeaning($meaningId: ID!, $grade: Int!) {
				reviewMeaning(meaningId: $meaningId, grade: $grade) {
					id
				}
			}
		`

		reviewVars := map[string]interface{}{
			"meaningId": meaningID,
			"grade":     6, // Невалидная оценка
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     reviewQuery,
			Variables: reviewVars,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for invalid grade (6)")
		}
	})

	t.Run("review with negative grade", func(t *testing.T) {
		// Создаём слово для review
		createQuery := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
						meanings {
							id
						}
					}
				}
			}
		`

		createVars := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "reviewtest3",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "тест",
					},
				},
			},
		}

		createResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars,
		})
		if err != nil {
			t.Fatalf("Failed to create word: %v", err)
		}

		var createResult struct {
			CreateWord struct {
				Word struct {
					Meanings []struct {
						ID string `json:"id"`
					} `json:"meanings"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp.Data, &createResult); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(createResult.CreateWord.Word.Meanings) == 0 {
			t.Fatal("Expected at least one meaning")
		}

		meaningID := createResult.CreateWord.Word.Meanings[0].ID

		// Пытаемся сделать review с отрицательной оценкой
		reviewQuery := `
			mutation ReviewMeaning($meaningId: ID!, $grade: Int!) {
				reviewMeaning(meaningId: $meaningId, grade: $grade) {
					id
				}
			}
		`

		reviewVars := map[string]interface{}{
			"meaningId": meaningID,
			"grade":     -1, // Отрицательная оценка
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     reviewQuery,
			Variables: reviewVars,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for negative grade")
		}
	})

	// ============================================
	// 5. ВАЛИДАЦИЯ INBOX
	// ============================================

	t.Run("add to inbox with empty text", func(t *testing.T) {
		query := `
			mutation AddToInbox($text: String!) {
				addToInbox(text: $text) {
					id
				}
			}
		`

		variables := map[string]interface{}{
			"text": "",
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) == 0 {
			t.Error("Expected GraphQL error for empty inbox text")
		}
	})

	t.Run("delete non-existent inbox item", func(t *testing.T) {
		query := `
			mutation DeleteInboxItem($id: ID!) {
				deleteInboxItem(id: $id)
			}
		`

		variables := map[string]interface{}{
			"id": "999999", // Несуществующий ID
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     query,
			Variables: variables,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		// Удаление несуществующего inbox item может вернуть false или ошибку
		if len(resp.Errors) > 0 {
			t.Logf("Delete non-existent inbox item returned error: %+v", resp.Errors)
		} else {
			var result struct {
				DeleteInboxItem bool `json:"deleteInboxItem"`
			}

			if err := json.Unmarshal(resp.Data, &result); err == nil {
				if result.DeleteInboxItem {
					t.Error("Expected deleteInboxItem to return false for non-existent item")
				}
			}
		}
	})

	// ============================================
	// 6. ОБНОВЛЕНИЕ С ПУСТЫМИ ДАННЫМИ
	// ============================================

	t.Run("update word to same text", func(t *testing.T) {
		// Создаём слово
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
				"text": "sametext",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "то же",
					},
				},
			},
		}

		createResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars,
		})
		if err != nil {
			t.Fatalf("Failed to create word: %v", err)
		}

		var createResult struct {
			CreateWord struct {
				Word struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp.Data, &createResult); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		wordID := createResult.CreateWord.Word.ID

		// Обновляем слово на тот же текст (должно работать)
		updateQuery := `
			mutation UpdateWord($id: ID!, $input: CreateWordInput!) {
				updateWord(id: $id, input: $input) {
					word {
						id
						text
					}
				}
			}
		`

		updateVars := map[string]interface{}{
			"id": wordID,
			"input": map[string]interface{}{
				"text": "sametext", // Тот же текст
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "то же",
					},
				},
			},
		}

		resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     updateQuery,
			Variables: updateVars,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		if len(resp.Errors) > 0 {
			t.Errorf("Unexpected error when updating to same text: %+v", resp.Errors)
		}

		var updateResult struct {
			UpdateWord struct {
				Word struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"word"`
			} `json:"updateWord"`
		}

		if err := json.Unmarshal(resp.Data, &updateResult); err == nil {
			if updateResult.UpdateWord.Word.ID != wordID {
				t.Errorf("Expected word ID to remain the same, got %s", updateResult.UpdateWord.Word.ID)
			}
			if updateResult.UpdateWord.Word.Text != "sametext" {
				t.Errorf("Expected text 'sametext', got '%s'", updateResult.UpdateWord.Word.Text)
			}
		}
	})

	// ============================================
	// 7. НЕВАЛИДНЫЕ ID
	// ============================================

	t.Run("operations with invalid ID format", func(t *testing.T) {
		invalidIDs := []string{
			"",
			"not-a-valid-id",
			"123abc",
			"   ",
		}

		for _, invalidID := range invalidIDs {
			// Тестируем updateWord с невалидным ID
			updateQuery := `
				mutation UpdateWord($id: ID!, $input: CreateWordInput!) {
					updateWord(id: $id, input: $input) {
						word {
							id
						}
					}
				}
			`

			variables := map[string]interface{}{
				"id": invalidID,
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

			resp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
				Query:     updateQuery,
				Variables: variables,
			})
			if err != nil {
				t.Logf("Request with invalid ID '%s' failed as expected: %v", invalidID, err)
				continue
			}

			// Должна быть ошибка валидации или NOT_FOUND
			if len(resp.Errors) == 0 {
				t.Errorf("Expected error for invalid ID '%s'", invalidID)
			}
		}
	})

	// ============================================
	// 8. КАСКАДНОЕ УДАЛЕНИЕ
	// ============================================

	t.Run("cascade delete when deleting word", func(t *testing.T) {
		// Создаём слово с meanings, examples и tags
		createQuery := `
			mutation CreateWord($input: CreateWordInput!) {
				createWord(input: $input) {
					word {
						id
						text
						meanings {
							id
							examples {
								id
							}
						}
					}
				}
			}
		`

		createVars := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "cascadetest",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "каскад",
						"examples": []map[string]interface{}{
							{
								"sentenceEn": "Test example",
								"sentenceRu": "Тестовый пример",
							},
						},
					},
				},
			},
		}

		createResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars,
		})
		if err != nil {
			t.Fatalf("Failed to create word: %v", err)
		}

		var createResult struct {
			CreateWord struct {
				Word struct {
					ID       string `json:"id"`
					Text     string `json:"text"`
					Meanings []struct {
						ID       string `json:"id"`
						Examples []struct {
							ID string `json:"id"`
						} `json:"examples"`
					} `json:"meanings"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp.Data, &createResult); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		wordID := createResult.CreateWord.Word.ID
		meaningID := createResult.CreateWord.Word.Meanings[0].ID

		// Удаляем слово
		deleteQuery := `
			mutation DeleteWord($id: ID!) {
				deleteWord(id: $id)
			}
		`

		deleteVars := map[string]interface{}{
			"id": wordID,
		}

		deleteResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     deleteQuery,
			Variables: deleteVars,
		})
		if err != nil {
			t.Fatalf("Failed to delete word: %v", err)
		}

		if len(deleteResp.Errors) > 0 {
			t.Fatalf("Unexpected error when deleting word: %+v", deleteResp.Errors)
		}

		// Проверяем, что слово удалено
		getQuery := `
			query GetWord($id: ID!) {
				word(id: $id) {
					id
				}
			}
		`

		getVars := map[string]interface{}{
			"id": wordID,
		}

		getResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     getQuery,
			Variables: getVars,
		})
		if err != nil {
			t.Fatalf("Failed to get word: %v", err)
		}

		var getResult struct {
			Word *struct {
				ID string `json:"id"`
			} `json:"word"`
		}

		if err := json.Unmarshal(getResp.Data, &getResult); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if getResult.Word != nil {
			t.Error("Expected word to be deleted")
		}

		// Проверяем, что meaning тоже удалён (через попытку review)
		reviewQuery := `
			mutation ReviewMeaning($meaningId: ID!, $grade: Int!) {
				reviewMeaning(meaningId: $meaningId, grade: $grade) {
					id
				}
			}
		`

		reviewVars := map[string]interface{}{
			"meaningId": meaningID,
			"grade":     5,
		}

		reviewResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     reviewQuery,
			Variables: reviewVars,
		})
		if err != nil {
			t.Fatalf("Failed to execute review: %v", err)
		}

		if len(reviewResp.Errors) == 0 {
			t.Error("Expected error when reviewing deleted meaning")
		}
	})

	// ============================================
	// 9. РЕГИСТР И НОРМАЛИЗАЦИЯ
	// ============================================

	t.Run("case insensitive word creation", func(t *testing.T) {
		// Создаём слово с заглавными буквами
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

		createVars1 := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "CASEtest",
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "тест",
					},
				},
			},
		}

		createResp1, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars1,
		})
		if err != nil {
			t.Fatalf("Failed to create word: %v", err)
		}

		var createResult1 struct {
			CreateWord struct {
				Word struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"word"`
			} `json:"createWord"`
		}

		if err := json.Unmarshal(createResp1.Data, &createResult1); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		// Проверяем, что текст нормализован (в нижнем регистре)
		if createResult1.CreateWord.Word.Text != "casetest" {
			t.Errorf("Expected normalized text 'casetest', got '%s'", createResult1.CreateWord.Word.Text)
		}

		// Пытаемся создать слово с тем же текстом в другом регистре
		createVars2 := map[string]interface{}{
			"input": map[string]interface{}{
				"text": "CASETEST", // Другой регистр
				"meanings": []map[string]interface{}{
					{
						"partOfSpeech":  "NOUN",
						"translationRu": "тест",
					},
				},
			},
		}

		createResp2, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
			Query:     createQuery,
			Variables: createVars2,
		})
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}

		// Должна быть ошибка дубликата
		if len(createResp2.Errors) == 0 {
			t.Error("Expected error for duplicate word (case insensitive)")
		}
	})
}
