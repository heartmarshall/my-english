package meaning_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	"github.com/heartmarshall/my-english/internal/model"
)

func TestRepo_GetDueForReview(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	clock := testutil.NewMockClock()
	repo := meaning.New(db, meaning.WithClock(clock))
	ctx := context.Background()

	t.Run("returns due meanings", func(t *testing.T) {
		pastTime := clock.Now().Add(-1 * time.Hour)
		rows := sqlmock.NewRows(meaningColumns).
			AddRow(1, 1, "noun", nil, "тест", nil, nil, "review", pastTime, 7, 2.5, 3, pastTime, pastTime)

		mock.ExpectQuery(`SELECT (.+) FROM meanings WHERE next_review_at < \$1 ORDER BY next_review_at ASC LIMIT 10`).
			WithArgs(clock.Now()).
			WillReturnRows(rows)

		meanings, err := repo.GetDueForReview(ctx, 0) // default limit

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(meanings) != 1 {
			t.Errorf("expected 1 meaning, got %d", len(meanings))
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByStatus(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaning.New(db)
	ctx := context.Background()

	t.Run("returns meanings with status", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows(meaningColumns).
			AddRow(1, 1, "noun", nil, "тест", nil, nil, "new", nil, nil, nil, nil, now, now)

		mock.ExpectQuery(`SELECT (.+) FROM meanings WHERE learning_status = \$1 ORDER BY created_at ASC LIMIT 10`).
			WithArgs(model.LearningStatusNew).
			WillReturnRows(rows)

		meanings, err := repo.GetByStatus(ctx, model.LearningStatusNew, 0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(meanings) != 1 {
			t.Errorf("expected 1 meaning, got %d", len(meanings))
		}
		if meanings[0].LearningStatus != model.LearningStatusNew {
			t.Errorf("expected status=new, got %s", meanings[0].LearningStatus)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetStudyQueue(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	clock := testutil.NewMockClock()
	repo := meaning.New(db, meaning.WithClock(clock))
	ctx := context.Background()

	t.Run("returns new and due meanings", func(t *testing.T) {
		now := clock.Now()
		pastTime := now.Add(-1 * time.Hour)
		rows := sqlmock.NewRows(meaningColumns).
			AddRow(1, 1, "noun", nil, "новое", nil, nil, "new", nil, nil, nil, nil, now, now).
			AddRow(2, 2, "verb", nil, "на повторение", nil, nil, "review", pastTime, 7, 2.5, 3, now, now)

		mock.ExpectQuery(`SELECT (.+) FROM meanings WHERE \(learning_status = \$1 OR next_review_at < \$2\)`).
			WithArgs(model.LearningStatusNew, now).
			WillReturnRows(rows)

		meanings, err := repo.GetStudyQueue(ctx, 10)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(meanings) != 2 {
			t.Errorf("expected 2 meanings, got %d", len(meanings))
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetStats(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	clock := testutil.NewMockClock()
	repo := meaning.New(db, meaning.WithClock(clock))
	ctx := context.Background()

	t.Run("returns stats", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"total", "mastered", "learning", "due"}).
			AddRow(100, 50, 30, 20)

		mock.ExpectQuery(`SELECT (.+) FROM meanings`).
			WithArgs(
				model.LearningStatusMastered,
				model.LearningStatusLearning,
				clock.Now(),
				model.LearningStatusNew,
			).
			WillReturnRows(rows)

		stats, err := repo.GetStats(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.TotalWords != 100 {
			t.Errorf("expected TotalWords=100, got %d", stats.TotalWords)
		}
		if stats.MasteredCount != 50 {
			t.Errorf("expected MasteredCount=50, got %d", stats.MasteredCount)
		}
		if stats.LearningCount != 30 {
			t.Errorf("expected LearningCount=30, got %d", stats.LearningCount)
		}
		if stats.DueForReviewCount != 20 {
			t.Errorf("expected DueForReviewCount=20, got %d", stats.DueForReviewCount)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_UpdateSRS(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	clock := testutil.NewMockClock()
	repo := meaning.New(db, meaning.WithClock(clock))
	ctx := context.Background()

	t.Run("success with all fields", func(t *testing.T) {
		nextReview := clock.Now().Add(24 * time.Hour)
		interval := 7
		easeFactor := 2.5
		reviewCount := 5

		srs := &meaning.SRSUpdate{
			LearningStatus: model.LearningStatusLearning,
			NextReviewAt:   &nextReview,
			Interval:       &interval,
			EaseFactor:     &easeFactor,
			ReviewCount:    &reviewCount,
		}

		// Порядок: SET поля (learning_status, updated_at, next_review_at, interval, ease_factor, review_count), затем WHERE id
		mock.ExpectExec(`UPDATE meanings SET`).
			WithArgs(
				model.LearningStatusLearning,
				clock.Now(),
				nextReview,
				interval,
				easeFactor,
				reviewCount,
				int64(1), // id в WHERE в конце
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateSRS(ctx, 1, srs)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("success with minimal fields", func(t *testing.T) {
		srs := &meaning.SRSUpdate{
			LearningStatus: model.LearningStatusMastered,
		}

		// Только learning_status, updated_at, затем id в WHERE
		mock.ExpectExec(`UPDATE meanings SET`).
			WithArgs(
				model.LearningStatusMastered,
				clock.Now(),
				int64(1),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateSRS(ctx, 1, srs)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("nil srs", func(t *testing.T) {
		err := repo.UpdateSRS(ctx, 1, nil)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		srs := &meaning.SRSUpdate{
			LearningStatus: model.LearningStatusNew,
		}

		mock.ExpectExec(`UPDATE meanings SET`).
			WithArgs(
				model.LearningStatusNew,
				clock.Now(),
				int64(999),
			).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateSRS(ctx, 999, srs)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}
