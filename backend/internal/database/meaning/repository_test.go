package meaning_test

import (
	"context"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v2"
)

var meaningColumns = []string{
	"id", "word_id", "part_of_speech", "definition_en", "translation_ru",
	"cefr_level", "image_url", "learning_status", "next_review_at",
	"interval", "ease_factor", "review_count", "created_at", "updated_at",
}

func TestRepo_Create(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	clock := testutil.NewMockClock()
	repo := meaning.New(q, meaning.WithClock(clock))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		m := &model.Meaning{
			WordID:        1,
			PartOfSpeech:  model.PartOfSpeechNoun,
			TranslationRu: "привет",
			DefinitionEn:  ptr("a greeting"),
		}

		mock.ExpectQuery(`INSERT INTO meanings`).
			WithArgs(
				int64(1),                // word_id
				model.PartOfSpeechNoun,  // part_of_speech
				"a greeting",            // definition_en
				"привет",                // translation_ru
				nil,                     // cefr_level
				nil,                     // image_url
				model.LearningStatusNew, // learning_status (default)
				nil,                     // next_review_at
				nil,                     // interval
				nil,                     // ease_factor
				nil,                     // review_count
				clock.Now(),             // created_at
				clock.Now(),             // updated_at
			).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(ctx, m)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.ID != 1 {
			t.Errorf("expected ID=1, got %d", m.ID)
		}
		if m.LearningStatus != model.LearningStatusNew {
			t.Errorf("expected status=new, got %s", m.LearningStatus)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("nil meaning", func(t *testing.T) {
		err := repo.Create(ctx, nil)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("missing word_id", func(t *testing.T) {
		m := &model.Meaning{TranslationRu: "тест"}

		err := repo.Create(ctx, m)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("missing translation", func(t *testing.T) {
		m := &model.Meaning{WordID: 1}

		err := repo.Create(ctx, m)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("default part_of_speech", func(t *testing.T) {
		m := &model.Meaning{
			WordID:        1,
			TranslationRu: "тест",
		}

		mock.ExpectQuery(`INSERT INTO meanings`).
			WithArgs(
				int64(1),
				model.PartOfSpeechOther, // default
				nil,
				"тест",
				nil, nil, model.LearningStatusNew, nil, nil, nil, nil,
				clock.Now(), clock.Now(),
			).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(ctx, m)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.PartOfSpeech != model.PartOfSpeechOther {
			t.Errorf("expected part_of_speech=other, got %s", m.PartOfSpeech)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaning.New(q)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows(meaningColumns).
			AddRow(1, 1, "noun", "a greeting", "привет", "A1", nil, "new", nil, nil, 2.5, 0, now, now)

		mock.ExpectQuery(`SELECT (.+) FROM meanings WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		m, err := repo.GetByID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.ID != 1 {
			t.Errorf("expected ID=1, got %d", m.ID)
		}
		if m.TranslationRu != "привет" {
			t.Errorf("expected translation='привет', got %q", m.TranslationRu)
		}
		if m.PartOfSpeech != model.PartOfSpeechNoun {
			t.Errorf("expected part_of_speech=noun, got %s", m.PartOfSpeech)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM meanings WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetByID(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByWordID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaning.New(q)
	ctx := context.Background()

	t.Run("found multiple", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows(meaningColumns).
			AddRow(1, 1, "noun", nil, "привет", nil, nil, "new", nil, nil, nil, nil, now, now).
			AddRow(2, 1, "verb", nil, "приветствовать", nil, nil, "learning", nil, nil, nil, nil, now, now)

		mock.ExpectQuery(`SELECT (.+) FROM meanings WHERE word_id = \$1 ORDER BY created_at ASC`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		meanings, err := repo.GetByWordID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(meanings) != 2 {
			t.Errorf("expected 2 meanings, got %d", len(meanings))
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := pgxmock.NewRows(meaningColumns)

		mock.ExpectQuery(`SELECT (.+) FROM meanings WHERE word_id = \$1`).
			WithArgs(int64(999)).
			WillReturnRows(rows)

		meanings, err := repo.GetByWordID(ctx, 999)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if meanings == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(meanings) != 0 {
			t.Errorf("expected 0 meanings, got %d", len(meanings))
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_Delete(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaning.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.Delete(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_DeleteByWordID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaning.New(q)
	ctx := context.Background()

	t.Run("deletes multiple", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings WHERE word_id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 3))

		count, err := repo.DeleteByWordID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 3 {
			t.Errorf("expected count=3, got %d", count)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("no rows affected", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings WHERE word_id = \$1`).
			WithArgs(int64(999)).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		count, err := repo.DeleteByWordID(ctx, 999)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected count=0, got %d", count)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

// Helper function
func ptr(s string) *string {
	return &s
}
