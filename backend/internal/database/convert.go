package database

import (
	"database/sql"
	"time"
)

// Этот файл содержит функции конвертации между Go-типами и SQL-типами.
// Экспортированы для использования в подпакетах (word, meaning, etc.)

// --- Null converters (Go -> SQL) ---
// Для pgx эти функции просто возвращают указатели напрямую.

// NullString возвращает указатель на строку (для pgx).
// Если s == nil, возвращает nil.
func NullString(s *string) any {
	return s
}

// NullInt возвращает указатель на int64 (для pgx).
// Если i == nil, возвращает nil.
func NullInt(i *int) any {
	if i == nil {
		return nil
	}
	val := int64(*i)
	return &val
}

// NullFloat возвращает указатель на float64 (для pgx).
// Если f == nil, возвращает nil.
func NullFloat(f *float64) any {
	return f
}

// NullTime возвращает указатель на time.Time (для pgx).
// Если t == nil, возвращает nil.
func NullTime(t *time.Time) any {
	return t
}

// --- Ptr converters (SQL -> Go) ---

// PtrString конвертирует sql.NullString в *string.
func PtrString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// PtrInt конвертирует sql.NullInt64 в *int.
func PtrInt(ni sql.NullInt64) *int {
	if !ni.Valid {
		return nil
	}
	i := int(ni.Int64)
	return &i
}

// PtrFloat конвертирует sql.NullFloat64 в *float64.
func PtrFloat(nf sql.NullFloat64) *float64 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float64
}

// PtrTime конвертирует sql.NullTime в *time.Time.
func PtrTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// --- Scanner interface ---

// Scanner — интерфейс для сканирования одной строки.
// Реализуется как *sql.Row, так и *sql.Rows.
type Scanner interface {
	Scan(dest ...any) error
}
