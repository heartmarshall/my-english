package inbox_test

import (
	"context"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
	"github.com/heartmarshall/my-english/internal/service/inbox"
)

// mockInboxRepository — мок для InboxRepository.
type mockInboxRepository struct {
	CreateFunc func(ctx context.Context, item *model.InboxItem) error
	GetByIDFunc func(ctx context.Context, id int64) (model.InboxItem, error)
	ListFunc    func(ctx context.Context) ([]model.InboxItem, error)
	DeleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockInboxRepository) Create(ctx context.Context, item *model.InboxItem) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, item)
	}
	return nil
}

func (m *mockInboxRepository) GetByID(ctx context.Context, id int64) (model.InboxItem, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return model.InboxItem{}, database.ErrNotFound
}

func (m *mockInboxRepository) List(ctx context.Context) ([]model.InboxItem, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []model.InboxItem{}, nil
}

func (m *mockInboxRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		sourceContext := "Harry Potter, page 50"
		repo := &mockInboxRepository{
			CreateFunc: func(ctx context.Context, item *model.InboxItem) error {
				if item.Text != "hello" {
					t.Errorf("expected Text='hello', got %q", item.Text)
				}
				if item.SourceContext == nil || *item.SourceContext != sourceContext {
					t.Errorf("expected SourceContext=%q, got %v", sourceContext, item.SourceContext)
				}
				item.ID = 1
				item.CreatedAt = time.Now()
				return nil
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		item, err := svc.Create(ctx, "hello", &sourceContext)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID != 1 {
			t.Errorf("expected ID=1, got %d", item.ID)
		}
		if item.Text != "hello" {
			t.Errorf("expected Text='hello', got %q", item.Text)
		}
	})

	t.Run("success without source context", func(t *testing.T) {
		repo := &mockInboxRepository{
			CreateFunc: func(ctx context.Context, item *model.InboxItem) error {
				if item.SourceContext != nil {
					t.Errorf("expected SourceContext=nil, got %v", item.SourceContext)
				}
				item.ID = 2
				item.CreatedAt = time.Now()
				return nil
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		item, err := svc.Create(ctx, "world", nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.Text != "world" {
			t.Errorf("expected Text='world', got %q", item.Text)
		}
	})

	t.Run("empty text", func(t *testing.T) {
		repo := &mockInboxRepository{}
		svc := inbox.New(inbox.Deps{Inbox: repo})

		_, err := svc.Create(ctx, "   ", nil)

		if err != service.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &mockInboxRepository{
			CreateFunc: func(ctx context.Context, item *model.InboxItem) error {
				return database.ErrDuplicate
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		_, err := svc.Create(ctx, "hello", nil)

		if err != database.ErrDuplicate {
			t.Errorf("expected ErrDuplicate, got %v", err)
		}
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedItem := model.InboxItem{
			ID:    1,
			Text:  "hello",
			CreatedAt: time.Now(),
		}

		repo := &mockInboxRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (model.InboxItem, error) {
				if id != 1 {
					t.Errorf("expected id=1, got %d", id)
				}
				return expectedItem, nil
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		item, err := svc.GetByID(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID != 1 {
			t.Errorf("expected ID=1, got %d", item.ID)
		}
		if item.Text != "hello" {
			t.Errorf("expected Text='hello', got %q", item.Text)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockInboxRepository{
			GetByIDFunc: func(ctx context.Context, id int64) (model.InboxItem, error) {
				return model.InboxItem{}, database.ErrNotFound
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		_, err := svc.GetByID(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("returns all items", func(t *testing.T) {
		expectedItems := []model.InboxItem{
			{ID: 1, Text: "hello", CreatedAt: time.Now()},
			{ID: 2, Text: "world", CreatedAt: time.Now().Add(time.Hour)},
		}

		repo := &mockInboxRepository{
			ListFunc: func(ctx context.Context) ([]model.InboxItem, error) {
				return expectedItems, nil
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		items, err := svc.List(ctx)

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
	})

	t.Run("empty result", func(t *testing.T) {
		repo := &mockInboxRepository{
			ListFunc: func(ctx context.Context) ([]model.InboxItem, error) {
				return []model.InboxItem{}, nil
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		items, err := svc.List(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if items == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(items) != 0 {
			t.Errorf("expected 0 items, got %d", len(items))
		}
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockInboxRepository{
			DeleteFunc: func(ctx context.Context, id int64) error {
				if id != 1 {
					t.Errorf("expected id=1, got %d", id)
				}
				return nil
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		err := svc.Delete(ctx, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockInboxRepository{
			DeleteFunc: func(ctx context.Context, id int64) error {
				return database.ErrNotFound
			},
		}

		svc := inbox.New(inbox.Deps{Inbox: repo})

		err := svc.Delete(ctx, 999)

		if err != database.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

