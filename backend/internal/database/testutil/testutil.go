// Package testutil содержит утилиты для тестирования репозиториев.
package testutil

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/heartmarshall/my-english/internal/database"
)

// MockClock — мок для database.Clock.
type MockClock struct {
	Time time.Time
}

// Now возвращает заранее заданное время.
func (m *MockClock) Now() time.Time {
	return m.Time
}

// FixedTime возвращает фиксированное время для тестов.
func FixedTime() time.Time {
	return time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
}

// NewMockClock создаёт MockClock с фиксированным временем.
func NewMockClock() *MockClock {
	return &MockClock{Time: FixedTime()}
}

// Compile-time проверка.
var _ database.Clock = (*MockClock)(nil)

// NewMockDB создаёт мок базы данных для тестов.
func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db, mock
}

// ExpectationsWereMet проверяет, что все ожидания sqlmock выполнены.
func ExpectationsWereMet(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
