package study_test

import (
	"context"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/study"
)

func TestService_GetStudyQueue(t *testing.T) {
	ctx := context.Background()

	t.Run("returns queue", func(t *testing.T) {
		now := time.Now()

		meaningRepo := &mockMeaningRepository{
			GetStudyQueueFunc: func(ctx context.Context, limit int) ([]model.Meaning, error) {
				return []model.Meaning{
					{ID: 1, WordID: 1, TranslationRu: "привет", LearningStatus: model.LearningStatusNew},
					{ID: 2, WordID: 1, TranslationRu: "здравствуй", LearningStatus: model.LearningStatusLearning},
				}, nil
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      &mockSRSRepository{},
			Clock:    &mockClock{now: now},
		})

		queue, err := svc.GetStudyQueue(ctx, 10)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(queue) != 2 {
			t.Errorf("expected 2 meanings, got %d", len(queue))
		}
	})

	t.Run("respects limit", func(t *testing.T) {
		var capturedLimit int

		meaningRepo := &mockMeaningRepository{
			GetStudyQueueFunc: func(ctx context.Context, limit int) ([]model.Meaning, error) {
				capturedLimit = limit
				return []model.Meaning{}, nil
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      &mockSRSRepository{},
			Clock:    &mockClock{now: time.Now()},
		})

		_, err := svc.GetStudyQueue(ctx, 5)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedLimit != 5 {
			t.Errorf("expected limit=5, got %d", capturedLimit)
		}
	})
}
