package e2e

import (
	"context"
	"encoding/json"
	"testing"
)

// TestInboxOperations тестирует операции с inbox.
func TestInboxOperations(t *testing.T) {
	ctx := context.Background()

	// Настраиваем тестовую БД
	testDB := SetupTestDB(ctx, t)
	defer testDB.Cleanup(ctx)

	// Настраиваем тестовое приложение
	testApp := SetupTestApp(ctx, t, testDB.Pool)
	defer testApp.Cleanup(ctx)

	// Добавляем элемент в inbox
	addQuery := `
		mutation AddToInbox($text: String!, $sourceContext: String) {
			addToInbox(text: $text, sourceContext: $sourceContext) {
				id
				text
				sourceContext
				createdAt
			}
		}
	`

	addVars := map[string]interface{}{
		"text":          "new word",
		"sourceContext": "test context",
	}

	addResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     addQuery,
		Variables: addVars,
	})
	if err != nil {
		t.Fatalf("Failed to add to inbox: %v", err)
	}

	if len(addResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", addResp.Errors)
	}

	var addResult struct {
		AddToInbox struct {
			ID            string  `json:"id"`
			Text          string  `json:"text"`
			SourceContext *string `json:"sourceContext"`
		} `json:"addToInbox"`
	}

	if err := json.Unmarshal(addResp.Data, &addResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if addResult.AddToInbox.Text != "new word" {
		t.Errorf("Expected text 'new word', got '%s'", addResult.AddToInbox.Text)
	}

	// Получаем список inbox items
	listQuery := `
		query GetInboxItems {
			inboxItems {
				id
				text
				sourceContext
			}
		}
	`

	listResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query: listQuery,
	})
	if err != nil {
		t.Fatalf("Failed to get inbox items: %v", err)
	}

	if len(listResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", listResp.Errors)
	}

	var listResult struct {
		InboxItems []struct {
			ID            string  `json:"id"`
			Text          string  `json:"text"`
			SourceContext *string `json:"sourceContext"`
		} `json:"inboxItems"`
	}

	if err := json.Unmarshal(listResp.Data, &listResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(listResult.InboxItems) == 0 {
		t.Error("Expected at least one inbox item")
	}

	found := false
	for _, item := range listResult.InboxItems {
		if item.Text == "new word" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'new word' in inbox items")
	}

	// Удаляем элемент из inbox
	deleteQuery := `
		mutation DeleteInboxItem($id: ID!) {
			deleteInboxItem(id: $id)
		}
	`

	deleteVars := map[string]interface{}{
		"id": addResult.AddToInbox.ID,
	}

	deleteResp, err := testApp.DoGraphQLRequest(ctx, GraphQLRequest{
		Query:     deleteQuery,
		Variables: deleteVars,
	})
	if err != nil {
		t.Fatalf("Failed to delete inbox item: %v", err)
	}

	if len(deleteResp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %+v", deleteResp.Errors)
	}

	var deleteResult struct {
		DeleteInboxItem bool `json:"deleteInboxItem"`
	}

	if err := json.Unmarshal(deleteResp.Data, &deleteResult); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !deleteResult.DeleteInboxItem {
		t.Error("Expected deleteInboxItem to return true")
	}
}

