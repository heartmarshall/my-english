package meaningtag_test

import (
	"context"
	"testing"

	"github.com/heartmarshall/my-english/internal/database/meaningtag"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	pgxmock "github.com/pashagolub/pgxmock/v2"
)

func TestRepo_AttachTag(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaningtag.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO meanings_tags`).
			WithArgs(int64(1), int64(2)).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.AttachTag(ctx, 1, 2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_AttachTags(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaningtag.New(q)
	ctx := context.Background()

	t.Run("attaches multiple", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO meanings_tags`).
			WithArgs(int64(1), int64(2), int64(1), int64(3), int64(1), int64(4)).
			WillReturnResult(pgxmock.NewResult("INSERT", 3))

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
	q, mock := testutil.NewMockQuerier(t)
	repo := meaningtag.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings_tags WHERE meaning_id = \$1 AND tag_id = \$2`).
			WithArgs(int64(1), int64(2)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.DetachTag(ctx, 1, 2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_DetachAllFromMeaning(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaningtag.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM meanings_tags WHERE meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 3))

		err := repo.DetachAllFromMeaning(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetTagIDsByMeaningID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := meaningtag.New(q)
	ctx := context.Background()

	t.Run("returns tag ids", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"tag_id"}).
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
		rows := pgxmock.NewRows([]string{"tag_id"})

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
	q, mock := testutil.NewMockQuerier(t)
	repo := meaningtag.New(q)
	ctx := context.Background()

	t.Run("returns all relations", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"meaning_id", "tag_id"}).
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
	q, mock := testutil.NewMockQuerier(t)
	repo := meaningtag.New(q)
	ctx := context.Background()

	t.Run("syncs tags - adds new and removes old", func(t *testing.T) {
		// Получаем текущие теги
		rows := pgxmock.NewRows([]string{"tag_id"}).
			AddRow(int64(2)).
			AddRow(int64(3))
		mock.ExpectQuery(`SELECT (.+) FROM meanings_tags WHERE meanings_tags\.meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		// Удаляем тег 2 (есть в текущих, но нет в новых)
		// Squirrel добавляет скобки вокруг WHERE условий
		mock.ExpectExec(`DELETE FROM meanings_tags WHERE.*meanings_tags\.meaning_id = \$1.*meanings_tags\.tag_id = \$2`).
			WithArgs(int64(1), int64(2)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		// Удаляем тег 3 (есть в текущих, но нет в новых)
		mock.ExpectExec(`DELETE FROM meanings_tags WHERE.*meanings_tags\.meaning_id = \$1.*meanings_tags\.tag_id = \$2`).
			WithArgs(int64(1), int64(3)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		// Добавляем теги 5 и 6 (есть в новых, но нет в текущих)
		mock.ExpectExec(`INSERT INTO meanings_tags`).
			WithArgs(int64(1), int64(5), int64(1), int64(6)).
			WillReturnResult(pgxmock.NewResult("INSERT", 2))

		err := repo.SyncTags(ctx, 1, []int64{5, 6})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("syncs tags - no changes needed", func(t *testing.T) {
		// Получаем текущие теги
		rows := pgxmock.NewRows([]string{"tag_id"}).
			AddRow(int64(5)).
			AddRow(int64(6))
		mock.ExpectQuery(`SELECT (.+) FROM meanings_tags WHERE meanings_tags\.meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		// Ничего не нужно добавлять или удалять
		err := repo.SyncTags(ctx, 1, []int64{5, 6})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("syncs tags - adds only", func(t *testing.T) {
		// Получаем текущие теги (пусто)
		rows := pgxmock.NewRows([]string{"tag_id"})
		mock.ExpectQuery(`SELECT (.+) FROM meanings_tags WHERE meanings_tags\.meaning_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		// Добавляем новые теги
		mock.ExpectExec(`INSERT INTO meanings_tags`).
			WithArgs(int64(1), int64(5), int64(1), int64(6)).
			WillReturnResult(pgxmock.NewResult("INSERT", 2))

		err := repo.SyncTags(ctx, 1, []int64{5, 6})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}
