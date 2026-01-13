package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/testutil"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v2"
)

func TestDictionaryRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		entry   *model.DictionaryEntry
		setup   func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "successful creation",
			entry: &model.DictionaryEntry{
				Text:           "Hello",
				TextNormalized: "hello",
			},
			setup: func(mock pgxmock.PgxPoolIface) {
				entryID := uuid.New()
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "text", "text_normalized", "created_at", "updated_at"}).
					AddRow(entryID, "Hello", "hello", now, now)
				mock.ExpectQuery(`INSERT INTO dictionary_entries`).
					WithArgs("Hello", "hello").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querier, mock := testutil.NewMockQuerier(t)
			repo := NewDictionaryRepository(querier)

			tt.setup(mock)

			ctx := context.Background()
			result, err := repo.Create(ctx, tt.entry)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("Create() returned nil result")
			}

			testutil.ExpectationsWereMet(t, mock)
		})
	}
}

func TestDictionaryRepository_CreateOrGet(t *testing.T) {
	tests := []struct {
		name    string
		entry   *model.DictionaryEntry
		setup   func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "returns existing entry on conflict",
			entry: &model.DictionaryEntry{
				Text:           "Hello",
				TextNormalized: "hello",
			},
			setup: func(mock pgxmock.PgxPoolIface) {
				entryID := uuid.New()
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "text", "text_normalized", "created_at", "updated_at"}).
					AddRow(entryID, "Hello", "hello", now, now)
				// Атомарный INSERT ... ON CONFLICT возвращает существующую запись
				mock.ExpectQuery(`INSERT INTO dictionary_entries`).
					WithArgs("Hello", "hello").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "creates new entry if not exists",
			entry: &model.DictionaryEntry{
				Text:           "World",
				TextNormalized: "world",
			},
			setup: func(mock pgxmock.PgxPoolIface) {
				entryID := uuid.New()
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "text", "text_normalized", "created_at", "updated_at"}).
					AddRow(entryID, "World", "world", now, now)
				// Атомарный INSERT ... ON CONFLICT создает новую запись
				mock.ExpectQuery(`INSERT INTO dictionary_entries`).
					WithArgs("World", "world").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querier, mock := testutil.NewMockQuerier(t)
			repo := NewDictionaryRepository(querier)

			tt.setup(mock)

			ctx := context.Background()
			result, err := repo.CreateOrGet(ctx, tt.entry)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("CreateOrGet() returned nil result")
			}

			testutil.ExpectationsWereMet(t, mock)
		})
	}
}

func TestDictionaryRepository_Update(t *testing.T) {
	entryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		id      uuid.UUID
		entry   *model.DictionaryEntry
		setup   func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "successful update",
			id:   entryID,
			entry: &model.DictionaryEntry{
				Text:           "Hello Updated",
				TextNormalized: "hello updated",
			},
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "text", "text_normalized", "created_at", "updated_at"}).
					AddRow(entryID, "Hello Updated", "hello updated", now, now)
				mock.ExpectQuery(`UPDATE dictionary_entries`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   entryID,
			entry: &model.DictionaryEntry{
				Text:           "Hello",
				TextNormalized: "hello",
			},
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`UPDATE dictionary_entries`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(pgx.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querier, mock := testutil.NewMockQuerier(t)
			repo := NewDictionaryRepository(querier)

			tt.setup(mock)

			ctx := context.Background()
			result, err := repo.Update(ctx, tt.id, tt.entry)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("Update() returned nil result")
			}

			testutil.ExpectationsWereMet(t, mock)
		})
	}
}

func TestDictionaryRepository_Delete(t *testing.T) {
	entryID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   entryID,
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`DELETE FROM dictionary_entries`).
					WithArgs(pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   entryID,
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`DELETE FROM dictionary_entries`).
					WithArgs(pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querier, mock := testutil.NewMockQuerier(t)
			repo := NewDictionaryRepository(querier)

			tt.setup(mock)

			ctx := context.Background()
			err := repo.Delete(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && err != nil {
				if err != database.ErrNotFound {
					t.Errorf("Delete() unexpected error = %v", err)
				}
			}

			testutil.ExpectationsWereMet(t, mock)
		})
	}
}

func TestDictionaryRepository_ExistsByNormalizedText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		setup   func(mock pgxmock.PgxPoolIface)
		want    bool
		wantErr bool
	}{
		{
			name: "exists",
			text: "hello",
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"1"}).AddRow(1)
				mock.ExpectQuery(`SELECT`).
					WithArgs("hello").
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not exists",
			text: "world",
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT`).
					WithArgs("world").
					WillReturnError(pgx.ErrNoRows)
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querier, mock := testutil.NewMockQuerier(t)
			repo := NewDictionaryRepository(querier)

			tt.setup(mock)

			ctx := context.Background()
			got, err := repo.ExistsByNormalizedText(ctx, tt.text)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExistsByNormalizedText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ExistsByNormalizedText() = %v, want %v", got, tt.want)
			}

			testutil.ExpectationsWereMet(t, mock)
		})
	}
}

func TestDictionaryRepository_FindByNormalizedText(t *testing.T) {
	entryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		text    string
		setup   func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "found",
			text: "hello",
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "text", "text_normalized", "created_at", "updated_at"}).
					AddRow(entryID, "Hello", "hello", now, now)
				mock.ExpectQuery(`SELECT`).
					WithArgs("hello").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "not found",
			text: "world",
			setup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT`).
					WithArgs("world").
					WillReturnError(pgx.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querier, mock := testutil.NewMockQuerier(t)
			repo := NewDictionaryRepository(querier)

			tt.setup(mock)

			ctx := context.Background()
			result, err := repo.FindByNormalizedText(ctx, tt.text)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByNormalizedText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("FindByNormalizedText() returned nil result")
			}

			if tt.wantErr && err != database.ErrNotFound {
				t.Errorf("FindByNormalizedText() expected ErrNotFound, got %v", err)
			}

			testutil.ExpectationsWereMet(t, mock)
		})
	}
}

func TestDictionaryRepository_CountTotal(t *testing.T) {
	tests := []struct {
		name    string
		filter  DictionaryFilter
		setup   func(mock pgxmock.PgxPoolIface)
		want    int64
		wantErr bool
	}{
		{
			name:   "count all",
			filter: DictionaryFilter{},
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(10))
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "count with search filter",
			filter: DictionaryFilter{
				Search: "hello",
			},
			setup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(5))
				// Fuzzy search генерирует SQL с оператором % и двумя аргументами
				mock.ExpectQuery(`SELECT COUNT`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(rows)
			},
			want:    5,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querier, mock := testutil.NewMockQuerier(t)
			repo := NewDictionaryRepository(querier)

			tt.setup(mock)

			ctx := context.Background()
			got, err := repo.CountTotal(ctx, tt.filter)

			if (err != nil) != tt.wantErr {
				t.Errorf("CountTotal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("CountTotal() = %v, want %v", got, tt.want)
			}

			testutil.ExpectationsWereMet(t, mock)
		})
	}
}
