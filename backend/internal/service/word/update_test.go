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

func TestService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		existingWord := model.Word{
			ID:        1,
			Text:      "hello",
			CreatedAt: now,
		}

		var updatedWord *model.Word

		wordRepo := &mockWordRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (model.Word, error) {
				if id == 1 {
					return existingWord, nil
				}
				return model.Word{}, database.ErrNotFound
			},
			UpdateFunc: func(ctx context.Context, w *model.Word) error {
				updatedWord = w
				return nil
			},
		}

		meaningRepo := &mockMeaningRepository{
			DeleteByWordIDFunc: func(ctx context.Context, wordID int64) (int64, error) {
				return 1, nil
			},
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

		newTranscription := "həˈloʊ"
		input := word.UpdateWordInput{
			Text:          "hello",
			Transcription: &newTranscription,
			Meanings: []word.UpdateMeaningInput{
				{Translations: []string{"привет"}},
			},
		}

		result, err := svc.Update(ctx, 1, input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updatedWord == nil {
			t.Fatal("word was not updated")
		}
		if updatedWord.Transcription == nil || *updatedWord.Transcription != newTranscription {
			t.Errorf("transcription was not updated correctly")
		}
		if result.Word.Text != "hello" {
			t.Errorf("expected text='hello', got %q", result.Word.Text)
		}
	})

	t.Run("not found", func(t *testing.T) {
		wordRepo := &mockWordRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (model.Word, error) {
				return model.Word{}, database.ErrNotFound
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		input := word.UpdateWordInput{
			Text:     "hello",
			Meanings: []word.UpdateMeaningInput{{Translations: []string{"привет"}}},
		}
		_, err := svc.Update(ctx, 999, input)

		if !errors.Is(err, service.ErrWordNotFound) {
			t.Errorf("expected ErrWordNotFound, got %v", err)
		}
	})

	t.Run("clear transcription", func(t *testing.T) {
		transcription := "həˈloʊ"
		existingWord := model.Word{
			ID:            1,
			Text:          "hello",
			Transcription: &transcription,
		}

		var updatedWord *model.Word

		wordRepo := &mockWordRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (model.Word, error) {
				return existingWord, nil
			},
			UpdateFunc: func(ctx context.Context, w *model.Word) error {
				updatedWord = w
				return nil
			},
		}

		meaningRepo := &mockMeaningRepository{
			DeleteByWordIDFunc: func(ctx context.Context, wordID int64) (int64, error) {
				return 0, nil
			},
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

		// Передаём nil для очистки transcription
		input := word.UpdateWordInput{
			Text:          "hello",
			Transcription: nil,
			Meanings: []word.UpdateMeaningInput{
				{Translations: []string{"привет"}},
			},
		}

		_, err := svc.Update(ctx, 1, input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updatedWord.Transcription != nil {
			t.Errorf("expected transcription to be nil, got %v", *updatedWord.Transcription)
		}
	})

	t.Run("change text to already existing word", func(t *testing.T) {
		existingWord := model.Word{ID: 1, Text: "hello"}
		anotherWord := model.Word{ID: 2, Text: "world"}

		wordRepo := &mockWordRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (model.Word, error) {
				if id == 1 {
					return existingWord, nil
				}
				return model.Word{}, database.ErrNotFound
			},
			GetByTextFunc: func(ctx context.Context, text string) (model.Word, error) {
				if text == "world" {
					return anotherWord, nil
				}
				return model.Word{}, database.ErrNotFound
			},
		}

		svc := word.New(word.Deps{
			Words:      wordRepo,
			Meanings:   &mockMeaningRepository{},
			Examples:   &mockExampleRepository{},
			Tags:       &mockTagRepository{},
			MeaningTag: &mockMeaningTagRepository{},
		})

		input := word.UpdateWordInput{
			Text:     "world", // пытаемся изменить на уже существующее слово
			Meanings: []word.UpdateMeaningInput{{Translations: []string{"мир"}}},
		}

		_, err := svc.Update(ctx, 1, input)

		if !errors.Is(err, service.ErrWordAlreadyExists) {
			t.Errorf("expected ErrWordAlreadyExists, got %v", err)
		}
	})
}
