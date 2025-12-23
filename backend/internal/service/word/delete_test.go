package word_test

import (
	"context"
	"errors"
	"testing"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/service"
	"github.com/heartmarshall/my-english/internal/service/word"
)

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		deleted := false
		wordRepo := &mockWordRepository{
			DeleteFunc: func(ctx context.Context, id int64) error {
				deleted = true
				return nil
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		err := svc.Delete(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !deleted {
			t.Error("word was not deleted")
		}
	})

	t.Run("not found", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			DeleteFunc: func(ctx context.Context, id int64) error {
				return database.ErrNotFound
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		err := svc.Delete(ctx, 999)

		if !errors.Is(err, service.ErrWordNotFound) {
			t.Errorf("expected ErrWordNotFound, got %v", err)
		}
	})
}
