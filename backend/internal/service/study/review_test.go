package study_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
	"github.com/heartmarshall/my-english/internal/service/study"
)

func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }

func TestService_ReviewMeaning(t *testing.T) {
	ctx := context.Background()

	t.Run("first review correct", func(t *testing.T) {
		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		existingMeaning := &model.Meaning{
			ID:             1,
			WordID:         1,
			TranslationRu:  "привет",
			LearningStatus: model.LearningStatusNew,
			EaseFactor:     floatPtr(2.5),
			Interval:       nil,
			ReviewCount:    nil,
		}

		var capturedUpdate *study.SRSUpdate

		meaningRepo := &mockMeaningRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (*model.Meaning, error) {
				if id == 1 {
					return existingMeaning, nil
				}
				return nil, database.ErrNotFound
			},
		}

		srsRepo := &mockSRSRepository{
			UpdateSRSFunc: func(ctx context.Context, id int64, srs *study.SRSUpdate) error {
				capturedUpdate = srs
				return nil
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      srsRepo,
			Clock:    &mockClock{now: now},
		})

		result, err := svc.ReviewMeaning(ctx, 1, 4) // grade = 4 (correct, easy)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Проверяем результат
		if result.ReviewCount == nil || *result.ReviewCount != 1 {
			t.Errorf("expected ReviewCount=1, got %v", result.ReviewCount)
		}
		// Первый ответ нового слова - интервал = 1
		if result.Interval == nil || *result.Interval != 1 {
			t.Errorf("expected Interval=1, got %v", result.Interval)
		}

		// Проверяем update
		if capturedUpdate == nil {
			t.Fatal("UpdateSRS was not called")
		}
		if capturedUpdate.Interval == nil || *capturedUpdate.Interval != 1 {
			t.Errorf("expected update Interval=1, got %v", capturedUpdate.Interval)
		}
	})

	t.Run("incorrect answer resets interval", func(t *testing.T) {
		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		existingMeaning := &model.Meaning{
			ID:             1,
			WordID:         1,
			TranslationRu:  "привет",
			LearningStatus: model.LearningStatusLearning,
			EaseFactor:     floatPtr(2.5),
			Interval:       intPtr(6),
			ReviewCount:    intPtr(3),
		}

		var capturedUpdate *study.SRSUpdate

		meaningRepo := &mockMeaningRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (*model.Meaning, error) {
				return existingMeaning, nil
			},
		}

		srsRepo := &mockSRSRepository{
			UpdateSRSFunc: func(ctx context.Context, id int64, srs *study.SRSUpdate) error {
				capturedUpdate = srs
				return nil
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      srsRepo,
			Clock:    &mockClock{now: now},
		})

		result, err := svc.ReviewMeaning(ctx, 1, 1) // grade = 1 (incorrect)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Ответ неправильный - интервал сбрасывается
		if result.Interval == nil || *result.Interval != 1 {
			t.Errorf("expected Interval=1 after incorrect, got %v", result.Interval)
		}
		// ReviewCount увеличивается
		if result.ReviewCount == nil || *result.ReviewCount != 4 {
			t.Errorf("expected ReviewCount=4, got %v", result.ReviewCount)
		}

		// Проверяем, что Easiness уменьшился
		if capturedUpdate.EaseFactor == nil || *capturedUpdate.EaseFactor >= 2.5 {
			t.Errorf("expected EaseFactor to decrease, got %v", capturedUpdate.EaseFactor)
		}
	})

	t.Run("meaning not found", func(t *testing.T) {
		meaningRepo := &mockMeaningRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (*model.Meaning, error) {
				return nil, database.ErrNotFound
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      &mockSRSRepository{},
			Clock:    &mockClock{now: time.Now()},
		})

		_, err := svc.ReviewMeaning(ctx, 999, 4)

		if !errors.Is(err, service.ErrMeaningNotFound) {
			t.Errorf("expected ErrMeaningNotFound, got %v", err)
		}
	})

	t.Run("invalid grade too low", func(t *testing.T) {
		svc := study.New(study.Deps{
			Meanings: &mockMeaningRepository{},
			SRS:      &mockSRSRepository{},
			Clock:    &mockClock{now: time.Now()},
		})

		_, err := svc.ReviewMeaning(ctx, 1, 0)

		if !errors.Is(err, service.ErrInvalidGrade) {
			t.Errorf("expected ErrInvalidGrade, got %v", err)
		}
	})

	t.Run("invalid grade too high", func(t *testing.T) {
		svc := study.New(study.Deps{
			Meanings: &mockMeaningRepository{},
			SRS:      &mockSRSRepository{},
			Clock:    &mockClock{now: time.Now()},
		})

		_, err := svc.ReviewMeaning(ctx, 1, 6)

		if !errors.Is(err, service.ErrInvalidGrade) {
			t.Errorf("expected ErrInvalidGrade, got %v", err)
		}
	})

	t.Run("becomes mastered after long interval", func(t *testing.T) {
		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Meaning уже в статусе Review с хорошим интервалом
		existingMeaning := &model.Meaning{
			ID:             1,
			WordID:         1,
			TranslationRu:  "привет",
			LearningStatus: model.LearningStatusReview,
			EaseFactor:     floatPtr(2.5),
			Interval:       intPtr(15), // интервал 15 дней
			ReviewCount:    intPtr(5),
		}

		meaningRepo := &mockMeaningRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (*model.Meaning, error) {
				return existingMeaning, nil
			},
		}

		srsRepo := &mockSRSRepository{
			UpdateSRSFunc: func(ctx context.Context, id int64, srs *study.SRSUpdate) error {
				return nil
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      srsRepo,
			Clock:    &mockClock{now: now},
		})

		result, err := svc.ReviewMeaning(ctx, 1, 5) // perfect answer

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// После успешного ответа с большим интервалом (15*2.5=37) должен стать Mastered
		if result.LearningStatus != model.LearningStatusMastered {
			t.Errorf("expected status Mastered, got %v", result.LearningStatus)
		}
	})

	t.Run("correct answer increases interval", func(t *testing.T) {
		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		existingMeaning := &model.Meaning{
			ID:             1,
			WordID:         1,
			TranslationRu:  "привет",
			LearningStatus: model.LearningStatusLearning,
			EaseFactor:     floatPtr(2.5),
			Interval:       intPtr(1),
			ReviewCount:    intPtr(1),
		}

		meaningRepo := &mockMeaningRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (*model.Meaning, error) {
				return existingMeaning, nil
			},
		}

		srsRepo := &mockSRSRepository{
			UpdateSRSFunc: func(ctx context.Context, id int64, srs *study.SRSUpdate) error {
				return nil
			},
		}

		svc := study.New(study.Deps{
			Meanings: meaningRepo,
			SRS:      srsRepo,
			Clock:    &mockClock{now: now},
		})

		result, err := svc.ReviewMeaning(ctx, 1, 4)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// После второго правильного ответа интервал должен быть 6
		if result.Interval == nil || *result.Interval != 6 {
			t.Errorf("expected Interval=6, got %v", result.Interval)
		}
	})
}
