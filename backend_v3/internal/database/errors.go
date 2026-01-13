package database

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Кастомные ошибки слоя данных.
// Позволяют бизнес-логике не зависеть от sql пакета.
var (
	// ErrNotFound возвращается, когда запись не найдена.
	ErrNotFound = errors.New("record not found")

	// ErrDuplicate возвращается при попытке создать дубликат
	// (нарушение UNIQUE constraint).
	ErrDuplicate = errors.New("record already exists")

	// ErrInvalidInput возвращается при невалидных входных данных.
	ErrInvalidInput = errors.New("invalid input")
)

// PostgreSQL error codes
// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	// pgUniqueViolation — код ошибки нарушения уникальности в PostgreSQL.
	pgUniqueViolation = "23505"
)

// IsNotFoundError проверяет, является ли ошибка "запись не найдена".
// Проверяет как кастомный ErrNotFound, так и оригинальный pgx.ErrNoRows.
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, pgx.ErrNoRows)
}

// IsDuplicateError проверяет, является ли ошибка нарушением UNIQUE constraint.
// Работает с github.com/jackc/pgx/v5.
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	// Проверка для pgx v5 через pgconn.PgError
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgUniqueViolation
	}

	// Fallback: проверка по строке (для совместимости)
	errStr := err.Error()
	if strings.Contains(errStr, pgUniqueViolation) ||
		strings.Contains(errStr, "duplicate key value violates unique constraint") {
		return true
	}

	return false
}

// WrapDBError оборачивает ошибку базы данных в domain-specific ошибку.
func WrapDBError(err error) error {
	if err == nil {
		return nil
	}

	if IsDuplicateError(err) {
		return ErrDuplicate
	}

	return err
}
