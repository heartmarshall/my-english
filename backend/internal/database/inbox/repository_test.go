package inbox_test

import (
	"context"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/inbox"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v2"
)

var inboxItemColumns = []string{"id", "text", "source_context", "created_at"}

func TestRepo_Create(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	clock := testutil.NewMockClock()
	repo := inbox.New(q, inbox.WithClock(clock))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		sourceContext := "Harry Potter, page 50"
		item := &model.InboxItem{
			Text:         "hello",
			SourceContext: &sourceContext,
		}

		mock.ExpectQuery(`INSERT INTO inbox_items`).
			WithArgs("hello", &sourceContext, clock.Time).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))

		err := repo.Create(ctx, item)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID != 1 {
			t.Errorf("expected ID=1, got %d", item.ID)
		}
		if item.CreatedAt != clock.Time {
			t.Errorf("expected CreatedAt=%v, got %v", clock.Time, item.CreatedAt)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("success without source context", func(t *testing.T) {
		item := &model.InboxItem{
			Text: "world",
		}

		mock.ExpectQuery(`INSERT INTO inbox_items`).
			WithArgs("world", (*string)(nil), clock.Time).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(2)))

		err := repo.Create(ctx, item)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID != 2 {
			t.Errorf("expected ID=2, got %d", item.ID)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("nil item", func(t *testing.T) {
		err := repo.Create(ctx, nil)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("empty text", func(t *testing.T) {
		item := &model.InboxItem{Text: "   "}

		err := repo.Create(ctx, item)

		if err != database.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestRepo_GetByID(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := inbox.New(q)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		sourceContext := "Book, page 10"
		now := time.Now()
		rows := pgxmock.NewRows(inboxItemColumns).
			AddRow(int64(1), "hello", &sourceContext, now)

		mock.ExpectQuery(`SELECT (.+) FROM inbox_items WHERE inbox_items\.id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		item, err := repo.GetByID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.Text != "hello" {
			t.Errorf("expected Text='hello', got %q", item.Text)
		}
		if item.SourceContext == nil || *item.SourceContext != sourceContext {
			t.Errorf("expected SourceContext=%q, got %v", sourceContext, item.SourceContext)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM inbox_items WHERE inbox_items\.id = \$1`).
			WithArgs(int64(999)).
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetByID(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_List(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := inbox.New(q)
	ctx := context.Background()

	t.Run("returns all items", func(t *testing.T) {
		now := time.Now()
		sourceContext1 := "Book 1"
		sourceContext2 := "Book 2"
		rows := pgxmock.NewRows(inboxItemColumns).
			AddRow(int64(1), "hello", &sourceContext1, now).
			AddRow(int64(2), "world", &sourceContext2, now.Add(time.Hour))

		mock.ExpectQuery(`SELECT (.+) FROM inbox_items ORDER BY inbox_items\.created_at DESC`).
			WillReturnRows(rows)

		items, err := repo.List(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("expected 2 items, got %d", len(items))
		}
		if items[0].Text != "hello" {
			t.Errorf("expected first item Text='hello', got %q", items[0].Text)
		}
		if items[1].Text != "world" {
			t.Errorf("expected second item Text='world', got %q", items[1].Text)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := pgxmock.NewRows(inboxItemColumns)

		mock.ExpectQuery(`SELECT (.+) FROM inbox_items ORDER BY inbox_items\.created_at DESC`).
			WillReturnRows(rows)

		items, err := repo.List(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// database.Select может вернуть nil для пустого результата
		if items != nil && len(items) != 0 {
			t.Errorf("expected 0 items, got %d", len(items))
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

func TestRepo_Delete(t *testing.T) {
	q, mock := testutil.NewMockQuerier(t)
	repo := inbox.New(q)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM inbox_items WHERE inbox_items\.id = \$1`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM inbox_items WHERE inbox_items\.id = \$1`).
			WithArgs(int64(999)).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.Delete(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		testutil.ExpectationsWereMet(t, mock)
	})
}

