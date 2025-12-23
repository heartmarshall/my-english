package database

import (
	"errors"
	"strings"
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

// IsDuplicateError проверяет, является ли ошибка нарушением UNIQUE constraint.
// Работает с github.com/lib/pq и github.com/jackc/pgx.
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Проверка для lib/pq: "pq: duplicate key value violates unique constraint"
	if strings.Contains(errStr, "duplicate key value violates unique constraint") {
		return true
	}

	// Проверка для pgx: содержит код ошибки 23505
	if strings.Contains(errStr, pgUniqueViolation) {
		return true
	}

	// Проверка для pgx v5: "SQLSTATE 23505"
	if strings.Contains(errStr, "SQLSTATE "+pgUniqueViolation) {
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
