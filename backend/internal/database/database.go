package database

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pagination defaults
const (
	DefaultLimit    = 20
	MaxLimit        = 100
	DefaultSRSLimit = 10
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var (
	_ Querier = (*pgxpool.Pool)(nil)
	_ Querier = (pgx.Tx)(nil)
)

var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

var DefaultClock Clock = RealClock{}

type TxFunc func(ctx context.Context, q Querier) error

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

// --- Scany Helpers ---

// GetOne сканирует одну структуру.
func GetOne[T any](ctx context.Context, q Querier, sql string, args ...any) (*T, error) {
	var dest T
	err := pgxscan.Get(ctx, q, &dest, sql, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &dest, nil
}

// Select сканирует список структур. Возвращает []*T.
func Select[T any](ctx context.Context, q Querier, sql string, args ...any) ([]T, error) {
	var dest []T // Слайс значений
	err := pgxscan.Select(ctx, q, &dest, sql, args...)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// GetScalar сканирует одно скалярное значение (int, string, bool).
func GetScalar[T any](ctx context.Context, q Querier, sql string, args ...any) (T, error) {
	var dest T
	err := pgxscan.Get(ctx, q, &dest, sql, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return dest, nil
		}
		return dest, err
	}
	return dest, nil
}

// SelectScalars сканирует список скалярных значений (например, []int64).
func SelectScalars[T any](ctx context.Context, q Querier, sql string, args ...any) ([]T, error) {
	var dest []T
	err := pgxscan.Select(ctx, q, &dest, sql, args...)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// CheckExists проверяет наличие строки.
func CheckExists(ctx context.Context, q Querier, sql string, args ...any) (bool, error) {
	var dummy int
	err := pgxscan.Get(ctx, q, &dummy, sql, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

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

func NormalizeLimit(limit, defaultVal int) int {
	if limit <= 0 {
		limit = defaultVal
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return limit
}
