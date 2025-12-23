package study_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/study"
)

func TestService_GetStats(t *testing.T) {
	ctx := context.Background()

	t.Run("returns stats", func(t *testing.T) {
		expectedStats := &model.Stats{
			TotalWords:        100,
			MasteredCount:     50,
			LearningCount:     30,
			DueForReviewCount: 20,
		}

		meaningRepo := &mockMeaningRepository{
			GetStatsFunc: func(ctx context.Context) (*model.Stats, error) {
				return expectedStats, nil
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      &mockSRSRepository{},
			Clock:    &mockClock{now: time.Now()},
		})

		stats, err := svc.GetStats(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.TotalWords != 100 {
			t.Errorf("expected TotalWords=100, got %d", stats.TotalWords)
		}
		if stats.MasteredCount != 50 {
			t.Errorf("expected MasteredCount=50, got %d", stats.MasteredCount)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		meaningRepo := &mockMeaningRepository{
			GetStatsFunc: func(ctx context.Context) (*model.Stats, error) {
				return nil, expectedErr
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      &mockSRSRepository{},
			Clock:    &mockClock{now: time.Now()},
		})

		_, err := svc.GetStats(ctx)

		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error to be propagated, got %v", err)
		}
	})
}
