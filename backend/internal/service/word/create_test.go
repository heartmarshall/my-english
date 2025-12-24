package word_test

import (
	"context"
	"errors"
	"testing"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
	"github.com/heartmarshall/my-english/internal/service/word"
)

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success with all data", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				return model.Word{}, database.ErrNotFound
			},
			CreateFunc: func(ctx context.Context, w *model.Word) error {
				w.ID = 1
				return nil
			},
		}

		meaningRepo := &mockMeaningRepository{
			CreateFunc: func(ctx context.Context, m *model.Meaning) error {
				m.ID = 1
				return nil
			},
		}

		exampleRepo := &mockExampleRepository{}
		tagRepo := &mockTagRepository{}
		meaningTagRepo := &mockMeaningTagRepository{}

		txRunner := &mockTxRunner{}
		repoFactory := &mockRepositoryFactory{
			wordRepo:       wordRepo,
			meaningRepo:    meaningRepo,
			exampleRepo:    exampleRepo,
			tagRepo:        tagRepo,
			meaningTagRepo: meaningTagRepo,
		}

		svc := word.New(word.Deps{
			Words:       wordRepo,
			Meanings:    meaningRepo,
			Examples:    exampleRepo,
			Tags:        tagRepo,
			MeaningTag:  meaningTagRepo,
			TxRunner:    txRunner,
			RepoFactory: repoFactory,
		})

		input := word.CreateWordInput{
			Text:          "Hello",
			Transcription: ptr("həˈloʊ"),
			Meanings: []word.CreateMeaningInput{
				{
					PartOfSpeech:  model.PartOfSpeechNoun,
					TranslationRu: "привет",
					Examples: []word.CreateExampleInput{
						{SentenceEn: "Hello, world!"},
					},
					Tags: []string{"greetings"},
				},
			},
		}

		result, err := svc.Create(ctx, input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Word.ID != 1 {
			t.Errorf("expected word ID=1, got %d", result.Word.ID)
		}
		if result.Word.Text != "hello" {
			t.Errorf("expected text='hello', got %q", result.Word.Text)
		}
		if len(result.Meanings) != 1 {
			t.Errorf("expected 1 meaning, got %d", len(result.Meanings))
		}
	})

	t.Run("empty text", func(t *testing.T) {
		svc := word.New(word.Deps{
			Words:      &mockWordRepository{},
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		input := word.CreateWordInput{
			Text:     "   ",
			Meanings: []word.CreateMeaningInput{{TranslationRu: "тест"}},
		}

		_, err := svc.Create(ctx, input)

		if !errors.Is(err, service.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("no meanings", func(t *testing.T) {
		svc := word.New(word.Deps{
			Words:      &mockWordRepository{},
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		input := word.CreateWordInput{
			Text:     "hello",
			Meanings: []word.CreateMeaningInput{},
		}

		_, err := svc.Create(ctx, input)

		if !errors.Is(err, service.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("word already exists", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				return model.Word{ID: 1, Text: "hello"}, nil
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		input := word.CreateWordInput{
			Text:     "hello",
			Meanings: []word.CreateMeaningInput{{TranslationRu: "привет"}},
		}

		_, err := svc.Create(ctx, input)

		if !errors.Is(err, service.ErrWordAlreadyExists) {
			t.Errorf("expected ErrWordAlreadyExists, got %v", err)
		}
	})

	t.Run("meaning without translation", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				return model.Word{}, database.ErrNotFound
			},
			CreateFunc: func(ctx context.Context, w *model.Word) error {
				w.ID = 1
				return nil
			},
		}

		meaningRepo := &mockMeaningRepository{}
		exampleRepo := &mockExampleRepository{}
		tagRepo := &mockTagRepository{}
		meaningTagRepo := &mockMeaningTagRepository{}

		txRunner := &mockTxRunner{}
		repoFactory := &mockRepositoryFactory{
			wordRepo:       wordRepo,
			meaningRepo:    meaningRepo,
			exampleRepo:    exampleRepo,
			tagRepo:        tagRepo,
			meaningTagRepo: meaningTagRepo,
		}

		svc := word.New(word.Deps{
			Words:       wordRepo,
			Meanings:    meaningRepo,
			Examples:    exampleRepo,
			Tags:        tagRepo,
			MeaningTag:  meaningTagRepo,
			TxRunner:    txRunner,
			RepoFactory: repoFactory,
		})

		input := word.CreateWordInput{
			Text: "hello",
			Meanings: []word.CreateMeaningInput{
				{PartOfSpeech: model.PartOfSpeechNoun, TranslationRu: ""},
			},
		}

		_, err := svc.Create(ctx, input)

		if !errors.Is(err, service.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func ptr(s string) *string {
	return &s
}
