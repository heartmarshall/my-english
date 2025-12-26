package e2e

import (
	"context"
	"encoding/json"
	"testing"
)

// TestStudyOperations тестирует операции изучения слов.
func TestStudyOperations(t *testing.T) {
	ctx := context.Background()

	// Настраиваем тестовую БД
	testDB := SetupTestDB(ctx, t)
	defer testDB.Cleanup(ctx)

	// Настраиваем тестовое приложение
	testApp := SetupTestApp(ctx, t, testDB.Pool)
	defer testApp.Cleanup(ctx)

	// Создаём слово для изучения
	createQuery := `
		mutation CreateWord($input: CreateWordInput!) {
			createWord(input: $input) {
				word {
					id
					meanings {
						id
						status
					}
				}
			}
		}
	`

	createVars := map[string]interface{}{
		"input": map[string]interface{}{
			"text": "study",
			"meanings": []map[string]interface{}{
				{
					"partOfSpeech":  "NOUN",
					"translationRu": "изучение",
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

	if len(createResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", createResp.Errors)
	}

	var createResult struct {
		CreateWord struct {
			Word struct {
				ID      string `json:"id"`
				Meanings []struct {
					ID     string `json:"id"`
					Status string `json:"status"`
				} `json:"meanings"`
			} `json:"word"`
		} `json:"createWord"`
	}

	if err := json.Unmarshal(createResp.Data, &createResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(createResult.CreateWord.Word.Meanings) == 0 {
		t.Fatal("Expected at least one meaning")
	}

	meaningID := createResult.CreateWord.Word.Meanings[0].ID

	// Проверяем начальный статус
	if createResult.CreateWord.Word.Meanings[0].Status != "NEW" {
		t.Errorf("Expected initial status 'NEW', got '%s'", createResult.CreateWord.Word.Meanings[0].Status)
	}

	// Выполняем review
	reviewQuery := `
		mutation ReviewMeaning($meaningId: ID!, $grade: Int!) {
			reviewMeaning(meaningId: $meaningId, grade: $grade) {
				id
				status
				reviewCount
				nextReviewAt
			}
		}
	`

	reviewVars := map[string]interface{}{
		"meaningId": meaningID,
		"grade":     5, // Отличная оценка
	}

	reviewResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     reviewQuery,
		Variables: reviewVars,
	})
	if err != nil {
		t.Fatalf("Failed to review meaning: %v", err)
	}

	if len(reviewResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", reviewResp.Errors)
	}

	var reviewResult struct {
		ReviewMeaning struct {
			ID          string  `json:"id"`
			Status      string  `json:"status"`
			ReviewCount int     `json:"reviewCount"`
			NextReviewAt *string `json:"nextReviewAt"`
		} `json:"reviewMeaning"`
	}

	if err := json.Unmarshal(reviewResp.Data, &reviewResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if reviewResult.ReviewMeaning.ReviewCount != 1 {
		t.Errorf("Expected reviewCount 1, got %d", reviewResult.ReviewMeaning.ReviewCount)
	}

	if reviewResult.ReviewMeaning.Status != "LEARNING" {
		t.Errorf("Expected status 'LEARNING' after review, got '%s'", reviewResult.ReviewMeaning.Status)
	}

	// Проверяем study queue
	queueQuery := `
		query GetStudyQueue {
			studyQueue(limit: 10) {
				id
				status
			}
		}
	`

	queueResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query: queueQuery,
	})
	if err != nil {
		t.Fatalf("Failed to get study queue: %v", err)
	}

	if len(queueResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", queueResp.Errors)
	}

	var queueResult struct {
		StudyQueue []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"studyQueue"`
	}

	if err := json.Unmarshal(queueResp.Data, &queueResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Проверяем статистику
	statsQuery := `
		query GetStats {
			stats {
				totalWords
				masteredCount
				learningCount
				dueForReviewCount
			}
		}
	`

	statsResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query: statsQuery,
	})
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if len(statsResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", statsResp.Errors)
	}

	var statsResult struct {
		Stats struct {
			TotalWords        int `json:"totalWords"`
			MasteredCount     int `json:"masteredCount"`
			LearningCount     int `json:"learningCount"`
			DueForReviewCount int `json:"dueForReviewCount"`
		} `json:"stats"`
	}

	if err := json.Unmarshal(statsResp.Data, &statsResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if statsResult.Stats.TotalWords < 1 {
		t.Errorf("Expected at least 1 total word, got %d", statsResult.Stats.TotalWords)
	}
}

