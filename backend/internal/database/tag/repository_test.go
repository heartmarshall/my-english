package tag_test

import (
	"context"
	"testing"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/tag"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v2"
)

var tagColumns = []string{"id", "name"}

func TestRepo_Create(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := tag.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tg := &model.Tag{Name: "vocabulary"}

		mock.ExpectQuery(`INSERT INTO tags`).
			WithArgs("vocabulary").
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(ctx, tg)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tg.ID != 1 {
			t.Errorf("expected ID=1, got %d", tg.ID)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("nil tag", func(t *testing.T) {
		err := repo.Create(ctx, nil)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		tg := &model.Tag{Name: "   "}

		err := repo.Create(ctx, tg)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestRepo_GetByID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := tag.New(q)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		rows := pgxmock.NewRows(tagColumns).AddRow(1, "business")

		mock.ExpectQuery(`SELECT (.+) FROM tags WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		tg, err := repo.GetByID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tg.Name != "business" {
			t.Errorf("expected Name='business', got %q", tg.Name)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM tags WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetByID(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByName(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := tag.New(q)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		rows := pgxmock.NewRows(tagColumns).AddRow(1, "travel")

		mock.ExpectQuery(`SELECT (.+) FROM tags WHERE name = \$1`).
			WithArgs("travel").
			WillReturnRows(rows)

		tg, err := repo.GetByName(ctx, "travel")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tg.Name != "travel" {
			t.Errorf("expected Name='travel', got %q", tg.Name)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM tags WHERE name = \$1`).
			WithArgs("nonexistent").
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetByName(ctx, "nonexistent")

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetByNames(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := tag.New(q)
	ctx := context.Background()

	t.Run("found multiple", func(t *testing.T) {
		rows := pgxmock.NewRows(tagColumns).
			AddRow(1, "business").
			AddRow(2, "travel")

		mock.ExpectQuery(`SELECT (.+) FROM tags WHERE name IN \(\$1,\$2\)`).
			WithArgs("business", "travel").
			WillReturnRows(rows)

		tags, err := repo.GetByNames(ctx, []string{"business", "travel"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tags) != 2 {
			t.Errorf("expected 2 tags, got %d", len(tags))
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty input", func(t *testing.T) {
		tags, err := repo.GetByNames(ctx, []string{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tags == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(tags))
		}
	})
}

func TestRepo_List(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := tag.New(q)
	ctx := context.Background()

	t.Run("returns all", func(t *testing.T) {
		rows := pgxmock.NewRows(tagColumns).
			AddRow(1, "a-tag").
			AddRow(2, "b-tag")

		mock.ExpectQuery(`SELECT (.+) FROM tags ORDER BY name ASC`).
			WillReturnRows(rows)

		tags, err := repo.List(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tags) != 2 {
			t.Errorf("expected 2 tags, got %d", len(tags))
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_Delete(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := tag.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM tags WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM tags WHERE id = \$1`).
			WithArgs(int64(999)).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.Delete(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_GetOrCreate(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := tag.New(q)
	ctx := context.Background()

	t.Run("returns existing", func(t *testing.T) {
		rows := pgxmock.NewRows(tagColumns).AddRow(1, "existing")

		mock.ExpectQuery(`SELECT (.+) FROM tags WHERE name = \$1`).
			WithArgs("existing").
			WillReturnRows(rows)

		tg, err := repo.GetOrCreate(ctx, "existing")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tg.ID != 1 {
			t.Errorf("expected ID=1, got %d", tg.ID)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("creates new", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM tags WHERE name = \$1`).
			WithArgs("new-tag").
			WillReturnError(pgx.ErrNoRows)

		mock.ExpectQuery(`INSERT INTO tags`).
			WithArgs("new-tag").
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(5))

		tg, err := repo.GetOrCreate(ctx, "new-tag")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tg.ID != 5 {
			t.Errorf("expected ID=5, got %d", tg.ID)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := repo.GetOrCreate(ctx, "   ")

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}
