package example_test

import (
	"context"
	"testing"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/example"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v2"
)

var exampleColumns = []string{"id", "meaning_id", "sentence_en", "sentence_ru", "source_name"}

func TestRepo_Create(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := example.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		src := model.ExampleSourceFilm
		ex := &model.Example{
			MeaningID:  1,
			SentenceEn: "Hello, world!",
			SentenceRu: ptr("Привет, мир!"),
			SourceName: &src,
		}

		mock.ExpectQuery(`INSERT INTO examples`).
			WithArgs(int64(1), "Hello, world!", "Привет, мир!", &src).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(ctx, ex)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ex.ID != 1 {
			t.Errorf("expected ID=1, got %d", ex.ID)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("nil example", func(t *testing.T) {
		err := repo.Create(ctx, nil)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("missing sentence", func(t *testing.T) {
		ex := &model.Example{MeaningID: 1}

		err := repo.Create(ctx, ex)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestRepo_GetByID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := example.New(q)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		rows := pgxmock.NewRows(exampleColumns).
			AddRow(1, 1, "Hello!", "Привет!", "film")

		mock.ExpectQuery(`SELECT (.+) FROM examples WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		ex, err := repo.GetByID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ex.SentenceEn != "Hello!" {
			t.Errorf("expected SentenceEn='Hello!', got %q", ex.SentenceEn)
		}
		if ex.SourceName == nil || *ex.SourceName != model.ExampleSourceFilm {
			t.Error("expected SourceName=film")
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM examples WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetByID(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByMeaningID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := example.New(q)
	ctx := context.Background()

	t.Run("found multiple", func(t *testing.T) {
		rows := pgxmock.NewRows(exampleColumns).
			AddRow(1, 1, "First", nil, nil).
			AddRow(2, 1, "Second", nil, "book")

		mock.ExpectQuery(`SELECT (.+) FROM examples WHERE meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		examples, err := repo.GetByMeaningID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(examples) != 2 {
			t.Errorf("expected 2 examples, got %d", len(examples))
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := pgxmock.NewRows(exampleColumns)

		mock.ExpectQuery(`SELECT (.+) FROM examples WHERE meaning_id = \$1`).
			WithArgs(int64(999)).
			WillReturnRows(rows)

		examples, err := repo.GetByMeaningID(ctx, 999)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if examples == nil {
			t.Error("expected empty slice, got nil")
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_Delete(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := example.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM examples WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM examples WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.Delete(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_DeleteByMeaningID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := example.New(q)
	ctx := context.Background()

	t.Run("deletes multiple", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM examples WHERE meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 3))

		count, err := repo.DeleteByMeaningID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 3 {
			t.Errorf("expected count=3, got %d", count)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func ptr(s string) *string {
	return &s
}
