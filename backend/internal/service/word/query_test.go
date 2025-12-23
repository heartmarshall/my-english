package word_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
	"github.com/heartmarshall/my-english/internal/service/word"
)

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("found with relations", func(t *testing.T) {
		now := time.Now()
		testWord := &model.Word{ID: 1, Text: "hello", CreatedAt: now}
		testMeaning := &model.Meaning{ID: 1, WordID: 1, TranslationRu: "привет"}
		testExample := &model.Example{ID: 1, MeaningID: 1, SentenceEn: "Hello!"}
		testTag := &model.Tag{ID: 1, Name: "greetings"}

		wordRepo := &mockWordRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (*model.Word, error) {
				if id == 1 {
					return testWord, nil
				}
				return nil, database.ErrNotFound
			},
		}

		meaningRepo := &mockMeaningRepository{
			GetByWordIDFunc: func(ctx context.Context, wordID int64) ([]*model.Meaning, error) {
				return []*model.Meaning{testMeaning}, nil
			},
		}

		exampleRepo := &mockExampleRepository{
			GetByMeaningIDsFunc: func(ctx context.Context, meaningIDs []int64) ([]*model.Example, error) {
				return []*model.Example{testExample}, nil
			},
		}

		tagRepo := &mockTagRepository{
			GetByIDsFunc: func(ctx context.Context, ids []int64) ([]*model.Tag, error) {
				return []*model.Tag{testTag}, nil
			},
		}

		meaningTagRepo := &mockMeaningTagRepository{
			GetByMeaningIDsFunc: func(ctx context.Context, meaningIDs []int64) ([]*model.MeaningTag, error) {
				return []*model.MeaningTag{{MeaningID: 1, TagID: 1}}, nil
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   meaningRepo,
			Examples:   exampleRepo,
			Tags:       tagRepo,
			MeaningTag: meaningTagRepo,
		})

		result, err := svc.GetByID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Word.Text != "hello" {
			t.Errorf("expected text='hello', got %q", result.Word.Text)
		}
		if len(result.Meanings) != 1 {
			t.Fatalf("expected 1 meaning, got %d", len(result.Meanings))
		}
		if len(result.Meanings[0].Examples) != 1 {
			t.Errorf("expected 1 example, got %d", len(result.Meanings[0].Examples))
		}
		if len(result.Meanings[0].Tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(result.Meanings[0].Tags))
		}
	})

	t.Run("not found", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (*model.Word, error) {
				return nil, database.ErrNotFound
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		_, err := svc.GetByID(ctx, 999)

		if !errors.Is(err, service.ErrWordNotFound) {
			t.Errorf("expected ErrWordNotFound, got %v", err)
		}
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("returns words", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			ListFunc: func(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]*model.Word, error) {
				return []*model.Word{
					{ID: 1, Text: "hello"},
					{ID: 2, Text: "world"},
				}, nil
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		words, err := svc.List(ctx, nil, 20, 0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(words) != 2 {
			t.Errorf("expected 2 words, got %d", len(words))
		}
	})

	t.Run("with filter", func(t *testing.T) {
		var capturedFilter *model.WordFilter

		wordRepo := &mockWordRepository{
			ListFunc: func(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]*model.Word, error) {
				capturedFilter = filter
				return []*model.Word{}, nil
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		search := "hel"
		filter := &word.WordFilter{Search: &search}

		_, err := svc.List(ctx, filter, 10, 0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedFilter == nil {
			t.Fatal("filter was not passed to repository")
		}
		if capturedFilter.Search == nil || *capturedFilter.Search != "hel" {
			t.Error("search filter was not passed correctly")
		}
	})
}

func TestService_Count(t *testing.T) {
	ctx := context.Background()

	t.Run("returns count", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			CountFunc: func(ctx context.Context, filter *model.WordFilter) (int, error) {
				return 42, nil
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		count, err := svc.Count(ctx, nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 42 {
			t.Errorf("expected count=42, got %d", count)
		}
	})
}
