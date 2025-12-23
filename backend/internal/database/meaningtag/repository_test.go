package meaningtag_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/heartmarshall/my-english/internal/database/meaningtag"
	"github.com/heartmarshall/my-english/internal/database/testutil"
)

func TestRepo_AttachTag(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaningtag.New(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO meanings_tags`).
			WithArgs(int64(1), int64(2)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.AttachTag(ctx, 1, 2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_AttachTags(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaningtag.New(db)
	ctx := context.Background()

	t.Run("attaches multiple", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO meanings_tags`).
			WithArgs(int64(1), int64(2), int64(1), int64(3), int64(1), int64(4)).
			WillReturnResult(sqlmock.NewResult(0, 3))

		err := repo.AttachTags(ctx, 1, []int64{2, 3, 4})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty list", func(t *testing.T) {
		err := repo.AttachTags(ctx, 1, []int64{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestRepo_DetachTag(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaningtag.New(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings_tags WHERE meaning_id = \$1 AND tag_id = \$2`).
			WithArgs(int64(1), int64(2)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DetachTag(ctx, 1, 2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_DetachAllFromMeaning(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaningtag.New(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings_tags WHERE meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 3))

		err := repo.DetachAllFromMeaning(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetTagIDsByMeaningID(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaningtag.New(db)
	ctx := context.Background()

	t.Run("returns tag ids", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"tag_id"}).
			AddRow(1).
			AddRow(2).
			AddRow(3)

		mock.ExpectQuery(`SELECT tag_id FROM meanings_tags WHERE meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		ids, err := repo.GetTagIDsByMeaningID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ids) != 3 {
			t.Errorf("expected 3 ids, got %d", len(ids))
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"tag_id"})

		mock.ExpectQuery(`SELECT tag_id FROM meanings_tags WHERE meaning_id = \$1`).
			WithArgs(int64(999)).
			WillReturnRows(rows)

		ids, err := repo.GetTagIDsByMeaningID(ctx, 999)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ids == nil {
			t.Error("expected empty slice, got nil")
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByMeaningIDs(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaningtag.New(db)
	ctx := context.Background()

	t.Run("returns all relations", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"meaning_id", "tag_id"}).
			AddRow(1, 10).
			AddRow(1, 20).
			AddRow(2, 10)

		mock.ExpectQuery(`SELECT meaning_id, tag_id FROM meanings_tags WHERE meaning_id IN \(\$1,\$2\)`).
			WithArgs(int64(1), int64(2)).
			WillReturnRows(rows)

		relations, err := repo.GetByMeaningIDs(ctx, []int64{1, 2})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(relations) != 3 {
			t.Errorf("expected 3 relations, got %d", len(relations))
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty input", func(t *testing.T) {
		relations, err := repo.GetByMeaningIDs(ctx, []int64{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if relations == nil {
			t.Error("expected empty slice, got nil")
		}
	})
}

func TestRepo_SyncTags(t *testing.T) {
	db, mock := testutil.NewMockDB(t)
	repo := meaningtag.New(db)
	ctx := context.Background()

	t.Run("syncs tags", func(t *testing.T) {
		// Сначала удаляет старые
		mock.ExpectExec(`DELETE FROM meanings_tags WHERE meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 2))

		// Затем добавляет новые
		mock.ExpectExec(`INSERT INTO meanings_tags`).
			WithArgs(int64(1), int64(5), int64(1), int64(6)).
			WillReturnResult(sqlmock.NewResult(0, 2))

		err := repo.SyncTags(ctx, 1, []int64{5, 6})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}
