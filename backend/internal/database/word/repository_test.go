package word_test

import (
	"context"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	"github.com/heartmarshall/my-english/internal/database/word"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v2"
)

func TestRepo_Create(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	clock := testutil.NewMockClock()
	repo := word.New(q, word.WithClock(clock))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		w := &model.Word{
			Text:          "Hello",
			Transcription: ptr("həˈloʊ"),
			AudioURL:      ptr("https://example.com/hello.mp3"),
			FrequencyRank: intPtr(100),
		}

		transcription := "həˈloʊ"
		audioURL := "https://example.com/hello.mp3"
		freqRank := int64(100)
		rows := pgxmock.NewRows([]string{"id"}).AddRow(int64(1))
		mock.ExpectQuery(`INSERT INTO words`).
			WithArgs("hello", &transcription, &audioURL, &freqRank, clock.Now()).
			WillReturnRows(rows)

		err := repo.Create(ctx, w)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.ID != 1 {
			t.Errorf("expected ID=1, got %d", w.ID)
		}
		if w.Text != "hello" {
			t.Errorf("expected text='hello', got %q", w.Text)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("nil word", func(t *testing.T) {
		err := repo.Create(ctx, nil)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("empty text", func(t *testing.T) {
		w := &model.Word{Text: "   "}

		err := repo.Create(ctx, w)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		w := &model.Word{Text: "duplicate"}

		mock.ExpectQuery(`INSERT INTO words`).
			WithArgs("duplicate", nil, nil, nil, clock.Now()).
			WillReturnError(pgx.ErrNoRows) // simplified; real error would be postgres duplicate key

		err := repo.Create(ctx, w)

		if err == nil {
			t.Error("expected error, got nil")
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := word.New(q)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		now := time.Now()
		transcription := "həˈloʊ"
		audioURL := "https://example.com/hello.mp3"
		freqRank := int64(100)
		rows := pgxmock.NewRows([]string{"id", "text", "transcription", "audio_url", "frequency_rank", "created_at"}).
			AddRow(int64(1), "hello", &transcription, &audioURL, &freqRank, &now)

		mock.ExpectQuery(`SELECT (.+) FROM words WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		w, err := repo.GetByID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.ID != 1 {
			t.Errorf("expected ID=1, got %d", w.ID)
		}
		if w.Text != "hello" {
			t.Errorf("expected text='hello', got %q", w.Text)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM words WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetByID(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByText(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := word.New(q)
	ctx := context.Background()

	t.Run("found with trimming and lowercasing", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "text", "transcription", "audio_url", "frequency_rank", "created_at"}).
			AddRow(int64(1), "hello", nil, nil, nil, &now)

		mock.ExpectQuery(`SELECT (.+) FROM words WHERE text = \$1`).
			WithArgs("hello"). // trimmed and lowercased
			WillReturnRows(rows)

		w, err := repo.GetByText(ctx, "  HELLO  ")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.Text != "hello" {
			t.Errorf("expected text='hello', got %q", w.Text)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM words WHERE text = \$1`).
			WithArgs("nonexistent").
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetByText(ctx, "nonexistent")

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_List(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := word.New(q)
	ctx := context.Background()

	t.Run("without filter", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "text", "transcription", "audio_url", "frequency_rank", "created_at"}).
			AddRow(int64(1), "hello", nil, nil, nil, &now).
			AddRow(int64(2), "world", nil, nil, nil, &now)

		mock.ExpectQuery(`SELECT (.+) FROM words ORDER BY created_at DESC LIMIT 20 OFFSET 0`).
			WillReturnRows(rows)

		words, err := repo.List(ctx, nil, 0, 0) // default limit

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(words) != 2 {
			t.Errorf("expected 2 words, got %d", len(words))
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("with search filter", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "text", "transcription", "audio_url", "frequency_rank", "created_at"}).
			AddRow(int64(1), "hello", nil, nil, nil, &now)

		search := "hel"
		filter := &model.WordFilter{Search: &search}

		mock.ExpectQuery(`SELECT (.+) FROM words WHERE text ILIKE \$1 ORDER BY created_at DESC LIMIT 10 OFFSET 0`).
			WithArgs("%hel%").
			WillReturnRows(rows)

		words, err := repo.List(ctx, filter, 10, 0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(words) != 1 {
			t.Errorf("expected 1 word, got %d", len(words))
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "text", "transcription", "audio_url", "frequency_rank", "created_at"})

		mock.ExpectQuery(`SELECT (.+) FROM words`).
			WillReturnRows(rows)

		words, err := repo.List(ctx, nil, 20, 0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if words == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(words) != 0 {
			t.Errorf("expected 0 words, got %d", len(words))
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_Count(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := word.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(42))

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM words`).
			WillReturnRows(rows)

		count, err := repo.Count(ctx, nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 42 {
			t.Errorf("expected count=42, got %d", count)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_Update(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := word.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		w := &model.Word{
			ID:   1,
			Text: "Updated",
		}

		mock.ExpectExec(`UPDATE words SET`).
			WithArgs("updated", nil, nil, nil, int64(1)).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(ctx, w)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.Text != "updated" {
			t.Errorf("expected text='updated', got %q", w.Text)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		w := &model.Word{ID: 999, Text: "test"}

		mock.ExpectExec(`UPDATE words SET`).
			WithArgs("test", nil, nil, nil, int64(999)).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := repo.Update(ctx, w)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("nil word", func(t *testing.T) {
		err := repo.Update(ctx, nil)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestRepo_Delete(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := word.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM words WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM words WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.Delete(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_Exists(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := word.New(q)
	ctx := context.Background()

	t.Run("exists", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"1"}).AddRow(int64(1))

		mock.ExpectQuery(`SELECT 1 FROM words WHERE id = \$1 LIMIT 1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		exists, err := repo.Exists(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !exists {
			t.Error("expected exists=true")
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT 1 FROM words WHERE id = \$1 LIMIT 1`).
			WithArgs(int64(999)).
			WillReturnError(pgx.ErrNoRows)

		exists, err := repo.Exists(ctx, 999)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exists {
			t.Error("expected exists=false")
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

// Helper functions
func ptr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
