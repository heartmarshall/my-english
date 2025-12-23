package database

import (
	"database/sql"
	"time"
)

// Этот файл содержит функции конвертации между Go-типами и SQL-типами.
// Экспортированы для использования в подпакетах (word, meaning, etc.)

// --- Null converters (Go -> SQL) ---

// NullString конвертирует *string в sql.NullString.
func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// NullInt конвертирует *int в sql.NullInt64.
func NullInt(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

// NullFloat конвертирует *float64 в sql.NullFloat64.
func NullFloat(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

// NullTime конвертирует *time.Time в sql.NullTime.
func NullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
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
