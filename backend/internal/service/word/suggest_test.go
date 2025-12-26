package word_test

import (
	"context"
	"testing"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/word"
)

func TestService_Suggest(t *testing.T) {
	ctx := context.Background()

	t.Run("exact match found", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				if text != "hello" {
					t.Errorf("expected text='hello', got %q", text)
				}
				transcription := "həˈloʊ"
				return model.Word{
					ID:            1,
					Text:          "hello",
					Transcription: &transcription,
				}, nil
			},
		}

		meaningRepo := &mockMeaningRepository{
			GetByWordIDFunc: func(ctx context.Context, wordID int64) ([]model.Meaning, error) {
				if wordID != 1 {
					t.Errorf("expected wordID=1, got %d", wordID)
				}
				return []model.Meaning{
					{
						ID:            1,
						WordID:        1,
						TranslationRu: "привет",
					},
					{
						ID:            2,
						WordID:        1,
						TranslationRu: "здравствуй",
					},
				}, nil
			},
		}

		exampleRepo := &mockExampleRepository{}
		tagRepo := &mockTagRepository{}
		meaningTagRepo := &mockMeaningTagRepository{}

		txRunner := &mockTxRunner{}
		repoFactory := &mockRepositoryFactory{}

		svc := word.New(word.Deps{
			Words:       wordRepo,
			Meanings:    meaningRepo,
			Examples:    exampleRepo,
			Tags:        tagRepo,
			MeaningTag:  meaningTagRepo,
			TxRunner:    txRunner,
			RepoFactory: repoFactory,
		})

		suggestions, err := svc.Suggest(ctx, "hello")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(suggestions) != 1 {
			t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
		}

		s := suggestions[0]
		if s.Text != "hello" {
			t.Errorf("expected Text='hello', got %q", s.Text)
		}
		if s.Origin != "LOCAL" {
			t.Errorf("expected Origin='LOCAL', got %q", s.Origin)
		}
		if s.ExistingWordID == nil || *s.ExistingWordID != 1 {
			t.Errorf("expected ExistingWordID=1, got %v", s.ExistingWordID)
		}
		if len(s.Translations) != 2 {
			t.Errorf("expected 2 translations, got %d", len(s.Translations))
		}
		if s.Translations[0] != "привет" {
			t.Errorf("expected first translation='привет', got %q", s.Translations[0])
		}
		if s.Translations[1] != "здравствуй" {
			t.Errorf("expected second translation='здравствуй', got %q", s.Translations[1])
		}
	})

	t.Run("exact match not found, trigram match found", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				return model.Word{}, database.ErrNotFound
			},
			SearchSimilarFunc: func(ctx context.Context, query string, limit int, similarityThreshold float64) ([]model.Word, error) {
				if query != "hel" {
					t.Errorf("expected query='hel', got %q", query)
				}
				if limit != 5 {
					t.Errorf("expected limit=5, got %d", limit)
				}
				if similarityThreshold != 0.3 {
					t.Errorf("expected similarityThreshold=0.3, got %f", similarityThreshold)
				}
				return []model.Word{
					{ID: 2, Text: "hello"},
					{ID: 3, Text: "help"},
				}, nil
			},
		}

		meaningRepo := &mockMeaningRepository{
			GetByWordIDFunc: func(ctx context.Context, wordID int64) ([]model.Meaning, error) {
				if wordID == 2 {
					return []model.Meaning{
						{ID: 3, WordID: 2, TranslationRu: "привет"},
					}, nil
				}
				if wordID == 3 {
					return []model.Meaning{
						{ID: 4, WordID: 3, TranslationRu: "помощь"},
					}, nil
				}
				return nil, database.ErrNotFound
			},
		}

		exampleRepo := &mockExampleRepository{}
		tagRepo := &mockTagRepository{}
		meaningTagRepo := &mockMeaningTagRepository{}

		txRunner := &mockTxRunner{}
		repoFactory := &mockRepositoryFactory{}

		svc := word.New(word.Deps{
			Words:       wordRepo,
			Meanings:    meaningRepo,
			Examples:    exampleRepo,
			Tags:        tagRepo,
			MeaningTag:  meaningTagRepo,
			TxRunner:    txRunner,
			RepoFactory: repoFactory,
		})

		suggestions, err := svc.Suggest(ctx, "hel")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(suggestions) != 2 {
			t.Fatalf("expected 2 suggestions, got %d", len(suggestions))
		}

		if suggestions[0].Text != "hello" {
			t.Errorf("expected first suggestion Text='hello', got %q", suggestions[0].Text)
		}
		if suggestions[1].Text != "help" {
			t.Errorf("expected second suggestion Text='help', got %q", suggestions[1].Text)
		}
	})

	t.Run("not found", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				return model.Word{}, database.ErrNotFound
			},
			ListFunc: func(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]model.Word, error) {
				return []model.Word{}, nil
			},
		}

		meaningRepo := &mockMeaningRepository{}
		exampleRepo := &mockExampleRepository{}
		tagRepo := &mockTagRepository{}
		meaningTagRepo := &mockMeaningTagRepository{}

		txRunner := &mockTxRunner{}
		repoFactory := &mockRepositoryFactory{}

		svc := word.New(word.Deps{
			Words:       wordRepo,
			Meanings:    meaningRepo,
			Examples:    exampleRepo,
			Tags:        tagRepo,
			MeaningTag:  meaningTagRepo,
			TxRunner:    txRunner,
			RepoFactory: repoFactory,
		})

		suggestions, err := svc.Suggest(ctx, "nonexistent")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(suggestions) != 0 {
			t.Errorf("expected 0 suggestions, got %d", len(suggestions))
		}
	})

	t.Run("empty query", func(t *testing.T) {
		wordRepo := &mockWordRepository{}
		meaningRepo := &mockMeaningRepository{}
		exampleRepo := &mockExampleRepository{}
		tagRepo := &mockTagRepository{}
		meaningTagRepo := &mockMeaningTagRepository{}

		txRunner := &mockTxRunner{}
		repoFactory := &mockRepositoryFactory{}

		svc := word.New(word.Deps{
			Words:       wordRepo,
			Meanings:    meaningRepo,
			Examples:    exampleRepo,
			Tags:        tagRepo,
			MeaningTag:  meaningTagRepo,
			TxRunner:    txRunner,
			RepoFactory: repoFactory,
		})

		suggestions, err := svc.Suggest(ctx, "   ")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(suggestions) != 0 {
			t.Errorf("expected 0 suggestions, got %d", len(suggestions))
		}
	})

	t.Run("word found but no translations", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				return model.Word{
					ID:   1,
					Text: "hello",
				}, nil
			},
		}

		meaningRepo := &mockMeaningRepository{
			GetByWordIDFunc: func(ctx context.Context, wordID int64) ([]model.Meaning, error) {
				return []model.Meaning{
					{
						ID:            1,
						WordID:        1,
						TranslationRu: "", // Пустой перевод
					},
				}, nil
			},
		}

		exampleRepo := &mockExampleRepository{}
		tagRepo := &mockTagRepository{}
		meaningTagRepo := &mockMeaningTagRepository{}

		txRunner := &mockTxRunner{}
		repoFactory := &mockRepositoryFactory{}

		svc := word.New(word.Deps{
			Words:       wordRepo,
			Meanings:    meaningRepo,
			Examples:    exampleRepo,
			Tags:        tagRepo,
			MeaningTag:  meaningTagRepo,
			TxRunner:    txRunner,
			RepoFactory: repoFactory,
		})

		suggestions, err := svc.Suggest(ctx, "hello")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Слово без переводов не должно возвращаться
		if len(suggestions) != 0 {
			t.Errorf("expected 0 suggestions (no translations), got %d", len(suggestions))
		}
	})
}
