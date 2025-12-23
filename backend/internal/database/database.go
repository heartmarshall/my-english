// Package database предоставляет слой доступа к данным.
// Содержит репозитории для работы с PostgreSQL.
package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
)

// Pagination defaults — константы для пагинации.
const (
	DefaultLimit    = 20  // Лимит по умолчанию для списков
	MaxLimit        = 100 // Максимальный лимит
	DefaultSRSLimit = 10  // Лимит по умолчанию для SRS очереди
)

// Querier — абстракция над *sql.DB и *sql.Tx.
// Позволяет использовать репозитории как с прямым подключением,
// так и в рамках транзакции.
type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Compile-time проверки, что *sql.DB и *sql.Tx реализуют Querier.
var (
	_ Querier = (*sql.DB)(nil)
	_ Querier = (*sql.Tx)(nil)
)

// Builder — squirrel builder с PostgreSQL плейсхолдерами ($1, $2, ...).
var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// Clock — интерфейс для получения текущего времени.
// Позволяет инжектить time.Now() для тестирования.
type Clock interface {
	Now() time.Time
}

// RealClock — реальная реализация Clock.
type RealClock struct{}

// Now возвращает текущее время.
func (RealClock) Now() time.Time {
	return time.Now()
}

// DefaultClock — clock по умолчанию для продакшена.
var DefaultClock Clock = RealClock{}

// NormalizePagination нормализует параметры пагинации.
// limit: [1, MaxLimit], offset: >= 0
func NormalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

// NormalizeLimit нормализует только limit с указанным значением по умолчанию.
func NormalizeLimit(limit, defaultVal int) int {
	if limit <= 0 {
		limit = defaultVal
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return limit
}

// TxFunc — функция, выполняемая в рамках транзакции.
type TxFunc func(ctx context.Context, q Querier) error

// WithTx выполняет функцию fn в рамках транзакции.
// Если fn возвращает ошибку или возникает паника — транзакция откатывается.
// Иначе — коммитится.
func WithTx(ctx context.Context, db *sql.DB, fn TxFunc) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
