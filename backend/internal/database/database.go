// Package database предоставляет слой доступа к данным.
// Содержит репозитории для работы с PostgreSQL.
package database

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pagination defaults — константы для пагинации.
const (
	DefaultLimit    = 20  // Лимит по умолчанию для списков
	MaxLimit        = 100 // Максимальный лимит
	DefaultSRSLimit = 10  // Лимит по умолчанию для SRS очереди
)

// Querier — абстракция над pgxpool.Pool и pgx.Tx.
// Позволяет использовать репозитории как с прямым подключением,
// так и в рамках транзакции.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Compile-time проверки, что pgxpool.Pool и pgx.Tx реализуют Querier.
var (
	_ Querier = (*pgxpool.Pool)(nil)
	_ Querier = (pgx.Tx)(nil)
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
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn TxFunc) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
